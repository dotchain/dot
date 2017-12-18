## Implementing Rich Text

For the purposes of this discussion, we will consider the following properties of rich text only:

1. Text can be bolded or italicized (inline styles)
2. Ordered and unordered lists (block styles) can nest. New lines are used to break line items.

## Virtual rich text state

An important point to be clarified here is that what we are talking about here is the virtual
JSON structure of rich text used for operations.  The actual storage has a lot more options
possible and these are considered later on.  So, please set aside concerns of performance or
memory consumption etc.

## Basic structure

The basic structure at the top level is actually an array with each element in the array
representing a character and it associated styles.

Each element of an array has three properties:

1. **char** is a string field which represents the actual UTF16 character. Now it might
seem best if we grouped
together *grapheme clusters* and styled them (since they effectively render as individual
characters) but this allows for more flexibility.  For instance, it is not uncommon for
editors to support modifying just the combining mark.

2. **inlineStyles** is a map of strings which represents the inline styles associated
with this character.  For this discussion, this map will have two entries: one for **bold**
and another for **italics**, both of these have boolean values.  (In practise, this will
probably hold CSS properties and their value associated with inline styles)

3. **blockPath** is an array that holds the sequence of blocks that encapsulate a
particular element.  For example, if a character is within a numbered list, its blockPath
would be an array with "numberedList" as the only item. If it was in a bulleted list
that itself was within a numbered list, the blockPath would be `["bulletedList", "numberedList"]`.
Note that it is intentional that the path lists the blocks of an element from the inside out.
This allows for run-length encoding of list styles later on.

How does this deal with indentation?  Indentation can be thought of as a block type by itself.
Note that when a user "indents" within a list, it should be thought of as creating an extra
"list" type (duplicating the innermost list type).

## Applying inline styles

Applying inline styles on on a contiguous region can be done with a simple Range mutation.
The Range mutation itself takes an inner mutation which applies to each element of the array.
The inner mutation here would be to set the corresponding field of the **inlineStyles** map.

A common thought experiment for collaborative rich text editors is how the system behaves
when one side marks a range of text as bold while another inserts text within.  The way the
transforms are written for splice and range, these should generally work gracefully.

## Applying block styles

Applying block styles generally also involves Range operations with the inner mutation being
a splice to insert the new block name into the `blockPath` array.  Often though, the operation
being attempted is more complicated and involves changing the block paths of various ranges
slightly differently.  In those cases, the user intention would require multiple operations.

For example, if the user selects a block of text with some of it in a bulleted list, some of it
in a numbered list and some of it in no list at all, some editors may choose to implement this
as if the user wrapped the whole selection within a list.  But it is also possible to consider
this as if the user added the numbered list items within a numbered list, the bulleted list
items within a bulleted list and the rest as within an "indentation" block.  The latter 
implementation would require more than one Range mutation.

In general, depending on the specifics of the block-level mutations, it is possible that some
parallel block level operations yield surprising results.  But given the rarity of these
events, this seems like an acceptable tradeoff given the positives of having such a simple
structure.

## Links

How would we choose to represent links? Since a particular character can only be part of one
link at a time, it is possible to simply create a **link** property with the value of the 
property set to the actual link.

## Embedded objects and images

It is common to embed special objects within rich text which have their own special rendering.
A common example would be images and such.  A simple approach is to add an "object" property
to the array elements with the value being whatever is the value of it.  A good choice for the 
**char** field of this object element is to set it to a special unicode character.

One of the effects of this choice for embedded objects is that there is support for these
objects participating in inline and block styles as well as links.

## Lists as embedded objects?

Can lists be implemented as embedded objects?  Sure. But most user operations on documents
tend to treat lists as text. For example, it is common to select text that starts in the
middle of a list and ends up at some other point within the text.  It is quite nice to be
able to represent this selection logically as indices within a single rich text buffer and
where that is a consideration, the separate implementation of lists makes sense.

## Rich text as embedded objects?

There are two different ways rich text can be embedded within.  One is the use of embedded
object to simply hold a reference to an id with the actual virtual JSON storage of the inner
rich text being stored elsewhere.

The second is the contents of the value of the `object` field containing the actual inner
rich text (probably as a sub-field).

Both of these suffer from the selectability issue -- can a selection range from within an
inner rich text to an outer rich text.  But this can be handled as a purely UX matter.

There is also a slightly more thorny issue of how the inline styles compose. If the outer
text has "bold" style, does it override the inner bold style?  If it does, things are good.
But if it instead simply sets the default style for the inner rich text, things can get
a bit complicated -- we would need to define a special state for the **bold** style which is
**default** (as in apply parent default).  But as hinted elsewhere, this may actually be
appropriate anyway if one uses CSS property values (which provides us with "inherit" as an
option that maps exactly to this situation).

## Storage implementation choices

The virtual JSON structure elaborated above is independent of actual storage choices.

A simple storage choice is to have a top level map of fields `char`, `inlineStyles`, `blockStyles`
and any other fields of the element (`link`, `object` etc) where 

* `char` stores the concatenation of the individual UTF8 strings (i.e. it is a simple string)
* the other fields store a run-length encoding of the individual elements.

This is quite compact while also making it easy to implement the basic mutations on it.
But it is quite slow to implement other operations that editors often need.

It is possible to build tree-based data structures for editors whose operations map closer
to it and then implement the basic operations on top of it.

## Rich text in operations

Even if the storage and performance cost is addressed via the implementation choice, how would
we ensure that the binary format used to transmit operations is not expensive?

Note that the library explicitly supports passing any arbitrary data types for the `Before` and `After`
fields of `Splice` and `Set` operations.  It uses the provided `ImmutableData` interface as
a bridge to work with these and so, applications can actually use any encodings for these
as they see fit (including exactly the implementation choice) and pass that through as
the bridge.
