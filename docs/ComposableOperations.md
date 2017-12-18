# Composable Operations

One of the distinguishing features of the approach to operational transformations
in DOT is the reliance on a small number of mutation types: Splice, Move, Set and
Range.

## Why so few mutation types?

There are three factors motivating the minimal approach:

1. WRiting transforms is hard.  The more complicated the mutation, the more involved
it is to get it right but in my practical experience, even simple mutations involve 
a lot of effort in getting the code right.  Take a look at the effort involved in
the move transform for example.

2. Changing the code for mutations requires complicated upgrades.  Most OT systems
rely on the ability to rebuild the model from previous mutations.  Changing the 
meaning would need a special upgrade step or potentially adding new versions as
new mutations (in which case both old and new mutations may need to be supported
side by side for a while).  Sticking to a small but minimal set of mutations
sidesteps this problem for most typical usage.


3. The test matrix for mutations is non-linear on the number of mutattion types.
The OT system has to guarantee that every mutation is safe against any other
mutation type that a different client might have attempted in parallel.  The 
complexity get unwieldy with more mutation types.  It is inevitable that some
mutation types will need to be versioned (new types added because of desire for
slight variations).  Having a large set of mutations makes the cost of adding
these important variations very high.  Keeping it minimal at the start allows
for a bit of wiggle room without unmanageable complexity.

An unintended benefit of a small number of mutations is the abilty to port
the mutations to a large number of platforms with low likelihood of portability
or other errors.

##  How does composition help?

There are two ways to get rich mutations:

1. Treat the low level mutation types used in OT as an assembly language two write
high level mutations. This leads to a lot of fairly powerful high level mutation
types that are immediately compatible in the system and do not require any transforms
to be written for them.  Also, the high level methods can be changed at any time
without worrying about rewriting history of models (though if the underlying
model schema has changed, an upgrade may very well be needed)

An example of this is Rich Text.  Please see [ImplementingRichText]
(ImplementingRichText.md)

2. Treat mutation types as generic.  For example, the Splice mutation type in this
library works on arrays (of arbitrary types) or strings.  This is achieved by the
OT library passing the buck to the developer in providing an "Array-Like" access 
interface (to be specific, see
[encoding](https://godoc.org/github.com/dotchain/dot/encoding))

In fact, this can be used to implement stacks, queues or any collection as they
all conform to array-like semantics -- so long as the basic mutation of the data 
structure are represented as a `Splice`.  '

An esoteric example is counters.  On the face of it, incrementing or decrementing
numbers does not have anything to do with sets or arrays.  But a virtual array
of numbers can effectively specify a number (the total of all elements) and also
allow increments (insert the increment) and decrements (insert the negative).
This might seem like it would grow the array indefinitely and be a burden to
store the array but the interesting aspect here is that the array itself never
needs to be stored.  Instead only the mutations need to be represented as if
there were a real array. When applying the mutation, the code can simply keep
track of the total and forget the actual underlying array at all.

This approach of creating rich types allows a very large number of data types
to be created and covers a very large and broad category of applications.

## Data integrity issues

There are two common data integrity problems with the minimal approach that 
generally do not happen with other "large mutation" approaches.

1. Micro-mutations can lead to intermediate states that are not valid in the
application. Consider the example of a table and an index -- every insert into
the table should be accompanied by an insert into the index and vice versa. But
the intermediate state exposed by a micro-operation can break this assumption.

This is somewhat easy to fix if one of the operation can be derived from the
other -- by simply making only the primary operation and leaving the secondary
operation to be done in response to the edit of the primary. The other remedy
is to mark consistency boundaries.  For example, DOT defines operations as a
set of changes and consistency can be maintained at the operation boundary but
not the individual change boundary.

2. Transformation integrity.  Consider an integrity constraint like the total
count of elements in two seperate arrays if fixed.  It is possible that two
simultaneous legal edits, when merged, can lead to breaking the constraint.
This is generally a very difficult problem for all Operational Transformation
systems but particularly harder with a minimal approach.

1. Allow loose consistency of the virtual document. In practice, this would
mean allow tables with rows not indexed or rows with tables not indexed (and
potentially with a higher level API that masks this).  Or allow things to 
grow larger than the size but use a different strategy to ignore the elements
past the size limits.

2. Use id-based objects rather than deep graphs and use garbage collection.
This is relatively useful anyway since the JSON virtual document structure
does not natively allow cycles anyway but it gets added significance when
one considers that integrity constraints (like no orphaned children etc)
are the same type of problem as garbage collection.  So, this approach can
"fix" any data inconsistencies periodically and "purify" the data.  There
is a question of how to treat the mutation resulting from a "fix".  Given
the distributed nature, any client can do it or a specific client hosted
by the service itself can take care of implementing the garbage collection
policy without having to do any mutations outside the OT framework.

3. Yet another approach is a custom server where the validation of state
is implemented server side when operations are received.  This is relatively
easy to implement but at the performance cost of needing to maintain a
choke point on the backend and potentially have performance and stability
issues.  This is discouraged for all but the most difficult situations.

## Automatic undos

One of the benefits of the micro mutations is that they have well defined
undo operations that are also relatively easy to implement without involving
higher level logic. So, implementing undo and redo can be done in a 
centralized location with a undo manager that does not even have to understand
the basic types involved.


