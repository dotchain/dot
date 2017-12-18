# Operation

Please see [dot package](https://godoc.org/github.com/dotchain/dot) for the 
golang representation of operations.  The following document describes the 
JSON representation.

A DOT operation is a JSON object like so:

```js
{
	"ID": <unique ID of the operation>,
	"Parents": <array of operation IDs>,
	"Changes": <array of Change structures>,
}
```


1. The **ID** is any arbitrary string (under 1024 UTF8 bytes) and  must be 
   unique to the **Model**.  For portability reasons, the ID is best served
   by using 7bit characters as it avoids unicode encoding issues with respect
   to uniqueness
2. The **Parents** array is used to hold one or two IDs.  The first ID is
   always refers to **Basis** of the operation which is the last operation
   in the Journal that was applied to the model before the current operation
   was applied. Please read this 
   [Introduction](IntroductionToOperationalTransforms.md) for a detailed
   explanation of **Basis**.  This can be an empty string for operations at
   the very start of the model.  The second optional element in the array
   refers to the **Parent** of the current operation which is the last
   operation the client created and applied but which may not have been
   incorporated yet into the journal (though it is required to have been
   sent to the journal service before sending the current operation)

## No Model ID specified

There is no model ID specified in an operation.  This is intentional because
operations are usually provided in arrays that all have the same model ID.
To avoid duplicating the structure, the model ID is typically provided at 
a different layer.

## Change structure

The **Change** structure is a JSON object like so:

```js
{
	"Path": [... string elements of path ... ],
	"Splice": <SpliceInfo structure>,
	"Move": <MoveInfo structure>,
	"Range": <RangeInfo structure>,
	"Set": <SetInfo structure>
}
```

The **Path** refers to the location in the virtual JSON model where the
change occurs.  The virtual JSON model is a JSON representation of the
actual model and is made up entirely of arrays and json objects.  The path
refers to the keys and indices involved in traversing this virtual json
model to arrive at the location of the change.  Note that keys and indices
are stringified in Path.

It is possible that the real model uses maps with non-string keys (such
as floating point numbers or bools) but these have to be encoded into
strings.

An empty or missing **Path** refers to the root of the model.

**Note** When a **Change** structure is nested (see **RangeInfo**), the
**Path** is relative to the array element and not relative to the model

A **Change** structure must have exactly only one field other than 
the **Path**.  For example, it is illegal to have both **Splice** and
**Move** specfied in the same change structure.  Instead, the changes should
be separated into two change structures so that the order of application
is explicit.

## Array like encodings

DOT uses the same mutation type to represent mutation in a family of array
like types.  For example, a Splice can be used against strings and arrays.

It is important that all clients agree on what the offsets mean.

1. Strings are treated as UTF16 arrays logically. Offsets refer to UTF16
offsets.  The wire representation of strings in SpliceInfo (Before and
After fields) remains as JSON strings though.

2. Regular JSON arrays can be used as such and empty fields represent
zero length arrays irrespective of the type.

3. Custom encodings are possible.  The DOT engine defines only one
custom encoding at this time: "SparseArray".  Custom encodings are 
actually represented as JSON objects like so:

```js
{
	"dot:encoding": <Name>,
	"dot:encoded": <JSON to be interpreted by implementation>
}
```

In particular, the Sparse array implementation uses the following
representation with run length encoding of values:

```js
{
	"dot:encoding": "SparseArray",
	"dot:encoded": <JSON array with all even elements being counts and odd elements being values>
}
```

## Dictionary like encodings

Similar to Arrays, dictionaries may also have alternate representations.

An example is compact sets, which is the only custom dictionary like encoding
defined by DOT at this moment:

```js
{
	"dot:encoding": "Set",
	"dot:encoded": <JSON array of keys in the set>
}
```

## SpliceInfo structure

The **SpliceInfo** structure is a JSON object like so:

```js
{
	"Offset": <Number>,
	"Before": <Any>,
	"After": <Any>
}
```

A Splice change represents mutating an array-like structure (such as
real arrays, sparse arrays, stacks, queues, strings) by replacing elements
from a particular location with an alternate sequence provided.

1. The **Offset** specifies the index where the replacement starts. This must
   be a non-negative number unlike some langauges where a negative index
   can be used to refer to offsets from the end of the sequence.  Zero is
   the start of the sequence.  An offset which is the size of the sequence
   can be used to insert elements at the end.  It is illegal to specify
   an offset larger than the size of the array-like structure.
2. The **Before** field specifies the sequence that is being replaced. The
   exact JSON representation can either be a string (in which case the
   offset refers to UTF16 offsets) or it can be a JSON array or it can be
   a custom array like encoding defined earlier.  A missing or null field
   indicates an empty array
3. The **After** field represents the replacement and has the same structure
   as the **Before** field.

A pure insertion change will have the Before field missing or empty while a pure
deletion change will have the After field missing.

## MoveInfo structure

The **MoveInfo** structure is a JSON object like so:

```js
{
	"Offset": <Number>,
	"Count": <Number>,
	"Distance": <Number>,
}
```

1. The **Offset** refers to the index of start of the sub-sequence being moved.
   Like with **SpliceInfo**, this must not be negative or larger than the size
   of the full sequence.  This change is allowed on strings in which case the
   offset refers to the UTF16 representation of the string.
2. The **Count** refers to the size of the sub-sequence that is being moved.
3. The **Distance** refers to the number of elements being skipped over. If it
   is positive, the move effectively takes the sub-sequence and shifts it
   to the right skipping over the number of elements specified in **Distance**.
   If the **Distance** is negative, the shifting happens to the left.
   Obviously, it is illegal for the Distance to be negative with a number
   larger than the Count.

## SetInfo structure

The **SetInfo** structure is a JSON-object like so:

```js
{
	"Key": <string field name>,
	"Before": <Any>,
	"After": <Any>
}
```

1. The **Key** field represents the key in the virtual JSON object whose value is
being replaced.  The actual type within the model may not be a dictionary but
could instead be a structure (or a class field). But it could also be a real
dictionary with non-string keys (such as arrays).  These are expected to be
serialized into strings.

2. The **Before** and **After** fields can be any JSON type and they can also
be custom encodings (or can contain custom encodings).

As with **SpliceInfo**, a pure insertion has a null or missing **Before** and a 
pure deletion has an empty or missing **After**.

## RangeInfo structure

The **RangeInfo** structure is a JSON object like so:

```js
{
	"Offset": <Number>,
	"Count": <Number>,
	"Change": <Change structure>
}
```

The **RangeInfo** structure represents a bulk modification of multiple elements
of a sequence.  The sub-sequence that is being modified is identified by the
**Offset** and **Count** fields and the actual modification is represented by
the nested **Change** field.  

Note that the **Path** within the **Change** structure is relative to the array
element i.e. an empty or missing **Path** field within the nested **Change**
structure would imply that the nested **Change** applies to the array element
directly.

## Encodings

As hinted at above, it is possible to encode objects when serializing to JSON
to save space or for performance reasons.  

At this point only two encodings are defined: one for sparse arrays (run length
encoding of array values) and one for sets.

```js
{
	"dot:encoding": "SparseArray",
	"dot:encoded": <JSON array with all even elements being counts and odd elements being values>
}
```

```js
{
	"dot:encoding": "Set",
	"dot:encoded": <JSON array of keys in the set>
}
```
