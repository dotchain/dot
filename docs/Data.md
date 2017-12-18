# Modeling Shared Data in DOT

This is an overview of the philosophy of DOT with regards to shared
application data.

For the purpose of this discussion, lets consider an application like
a collaborative rich text editor (like Word 365 or Google docs). To
make things interesting, lets further consider that the editor can
embed graphs (such as charts) and allows selection of individual graph
elements.

## Preliminaries

A reminder: The unit of granularity of DOT synchronization is a
"Model" and models are logical trees (irrespective of actual physical
storage structure).  Mutations to a model are represented by a path to
the location in the logical tree and the actual change.

## Primary state

The first things to consider are the data that the user thinks (s)he
can modify.  Of these, we restrict things further to things that can
only change when some user edits that data (i.e. it cannot change in
response to an edit of a different data).

This is necessarily a vague definition as it does depend on the exact
data structures used but most engineers would probably get the gist of
the idea.  For the collaborative editor setup, the primary data would
be:

1. The rich text part including any formatting etc.
2. The embedded graph elements

### Single vs multiple models

It is possible to fit all the primary data here into one Model object,
but notice that changing the graph does not affect the text much at
all and changing the rich text does not affect the graph elements.

This leads to a natural decomposition of the two.  It would be ideal
if we could create a model for each embedded graph and model for the
overall text.  A side benefit is it would allow the app to be built in
such a way that the graph editor can work directly on its own model
and not even communicate with the app which edits the rich text.

The philsophy in DOT leans towards composition.  So, lets consider how
the rich text can hold "references" to the individual graph.

### References

In this simple case, the can simply be the Model ID and nothing more
is needed.  But consider the case where the application, for some
purpose, put all embedded graphs together in one model ID.  In that
case, it would be useful for the reference to include both the model
ID of the embedded graphs collection as well as the path into the
virtual JSON that holds the individual graph.

A reference now looks like:

```json
[ "modelID", index ]
```

Now we immediately run into a problem with synchronization: which
version of the embedded model does the reference point to? This is
relevant because another client may not have synced up to the latest
version of the embedded graphs and the index would no longer even be
accurate.

To get around this, we really have to store the version with the model
ID. It is better to actually structure the reference like so:

```json
{
	"ModelID": <modelID>,
	"Version": <LastOperationID>,
	"Path": [.. path into the model at that version]
}
```

Note that `Path` here is exactly the same type of `Path` as would have
been used with operations if that graph were to be edited.  This also
provides a easy mechanism to look at changes to the embedded model and
filter out changes that do not affect a particular reference.

But there are other problems: what happens if an edit to the graphs
model (say an insert of a new graph at the top of the list of embedded
graphs) breaks the path used?

Any client that watches all changes on the embeddings model can
actually reliably keep track of the "real path" by means of a simple
transformation function.  But the important fact is that the calculted
value depends only on the reference value and all changes made to the
referred model past that initial version.

In a sense, the raw reference data is static but the editor uses a
calculated value.  The former is "primary data" (data that isn't
affected by change elsewhere) and the calculated value becomes
"secondary data".

### Cursors, Selections

A similar situation occurs with cursors and selections.  These are
also effectively references.  A simple cursor would refer to a
specific location in the text in the main model of the app.  But note
that while the raw cursor can only be modified by a user, the
effective cursor position can be modified by changes to the rich text
being made by other users.  But the calculation to maintain the cursor
position of the current user looks exactly like the work to maintain
the calculated reference value.

There is a slight difference from paths used in operations
though. Cursor selections often need to maintain the logical "start"
or end within selections.  So an index like "0" can sometimes be
translated to "1" if an object was inserted into the array at position
"0" but cursor sometimes want to refer to the logical start or logical
end.  A case could also be made for logical medians etc but lets
sidestep that for the moment.

A simple modification to the path scheme takes care of it: within
array elements, `[start]` and `[end]` can be treated as logical
starts/ends and not get transformed at all.

Selections are a bit different.  They are typically collections of
references.

### Functional derivation for references

There is a functional way to look at references. One can define the
current reference by providing a static initial data and a `pure
function` of a model and its changes like so:

```golang

func CurrentReference(initial Ref, changes []dot.Change) (effective Ref) {
	... pure function ...
}
```


But this requires the full cumulative changes since the original
version. There is a slightly better functional way of looking at this
particular derivation:

```golang

func RefDelta(last Ref, change dot.Change) (delta []dot.Change) {
	... pure function which returns how to modify ...
}
```

This pure function takes the current computed value of a ref and the
set of changes that happened to the referenced model since the
computation -- and returns the set of changes to apply to modify the
current reference to the right new reference.

### Automatic derivation

This leads to a interesting proposition -- if one were to use a
reactive or incremental computing approach with DOT, it should be
possible for the client platform to maintain calculated references
with a very simple minimal functional computation and leave the
mechanics of how these are maintained to the DOT reactive engine.

The second form of functional derivation in particular is very
efficient as it only requires the last known value plus the changes
since then.

### Line numbers

But not all state problems work that way.  Consider the case of line
numbers in the editor.

Note that this is not a primary data but a secondary data -- similar
to the calculated value of references but different in that it is not
actually ever edited by users.  

Lets first consider how a collaborative editor could represent it.  An
editor that can deal with very large documents should ideally, at any
given point, the editor should ideally maintain its "visibility
window" -- which can be represented by a pair of references within the 
virtual text array that represents the full editor (assuming the
editor uses a flat virtual array). 

Within this window, it would probably need the rendered text + the
starting line number. The rest of the calculated values can be
calculated efficiently each time it renders that we can effectively
ignore it.

So, the naive functional representation is via a function like so:

```golang

func getRenderState(input RichText, start, end Ref) (startLineNum int, window RichText) {
     ... pure function ...
}
```

But with a large document this gets expensive.  An alternate model is
to actually maintain the "start" of the window as an offset and then
one can actually write a pure function like so:

```golang

type Window struct {
	offset int
	startLineNumber int
	slice RichText
}

func updateWindow(initial window, changes []dot.Change)  []dot.Change {
}
```

Here the `updateWindow` takes as input the `changes` to the full rich
text but it has sufficient context to be able to compute the output
changes to apply to the `Window` struct.

