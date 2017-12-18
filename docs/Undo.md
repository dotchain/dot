# Undo

One of the benefits of operational transforms as described here is the 
ability to implement undo operations very naturally.

For example, all the change types in this package have sufficient information
for the corresponding undo mutation to be constructed.  So a client can
simply use that to construct the mutations needed to reverse the last
operation. In fact, the operation type supports just this via an 'Undo'
method.

Note: a key reason for being able to do this is the fact that the mutation
types are defined to hold enough information of the prior state to do this.
For example, the splice operation not only says how many elements to delete,
it actually holds the copy of the deleted items.  This is a conscious choice.
Other systems do not do this and in those systems, a current model is needed
to identify the equivalent inverse operation which breaks down when remote
operations are allowed to intervene.

## Intervening operations

In a collaborative system, a wrinkle on top of the undo/redo setup is the 
fact that remote operations may intervene local changes.  How does one
undo an operation if a remote operation happened after it?

The simple approach is to consider the undo operation as if it is being
performed by a separate session with the basis and parent of the operation
being undone.  Then this can be transformed against all further operations
in the journal.

But an actual implementation can be a bit simpler without having to rely
on creating a separate session or haing awareness of the jouranl at all.
Instead, this can be accomplished purely by tracking the mutations applied
to the local model.

Consider the sequence of mutations to a local model like:

```
    C1, C2, S1, S2...
    where C = local client mutatoin
    	  S = remote mutation due to collaboration
	  U = local undo mutation
	  R = local redo mutation
```

To undo a particular operation in this chain, we can simply "merge" the
inverse effect against all the operations that follow it.  This will yield
the effective inverse that can be applied at the end -- exactly what we
would like to accomplish.

## Redo

What about redo?  A redo is simply the inverse of the corresponding undo.

Consider:

```
    C1, C2, S1, S2, U2, ...
```

Now to compute the redo to apply, we just locate the last undo U2, invert
its effect (Undo(U2)) and merge that with the rest of the operations to figure
out what R2 would be.

What about a more complex case like this:

```
    C1, C2, S1, S2, U2, U1, S2,  R1, ...
```

One may notice that the correct undo operation to redo here is actually is 
actually U2.  So the algorithm for finding this involves walking the operations
from the most recent back to find the first undo with the caveat that any
redo operation will cause the corresponding undo to be ignored.

Introducing a redo also slightly complicates the algorithm to find which
operation to undo slightly.  Firstly, undo will have to find the last local
operation or the last redo operation which ever comes first -- because the
redo operation effectively captures the intent of the original local operation,
there is no reason to go back further.

Secondly, a undo operation also need to skip a redo if there is a matching undo
after it..

## Canceling undo/redo

The merge operations in this package guarantee convergence but they do not
guarantee that a operation when merged against an (operation, inverse) pair
will remain exactly the same.  This is effectively very difficult to do in
all the cases and the package does not make much of an attempt at this.

There are OT systems that do this well but this is generally a rare problem.

The package attempts a very simple mitigation of removing consequitive pairs
of undo/redo changes before merge/transformation.  This will ensure that 
undo/redo works very well in a large number of cases.

## Removing all undo/redo pairs

Is it possible to remove non-consequtive undo redo pairs?  That is, can we 
deal with intervening operations between the original change and the undo
operation for it?

A fairly expensive method exists for just this.

Consider:

```
     O, A, B, C, Transformed(Undo(O)) ....
```

If one merges Undo(O) against A, B, C, we are going ot get Transformed(Undo(O)).
In addition, we would get A1, B1, C1 such that:

```
    O + Undo(O) + A1, B1, C1 = O + A + B + C + Transformed(Undo(O))
```

Since O and Undo(O) cancel each other, one can simply replace that original
sequence with:

```
   A1, B1, C1
```

The above algorithm is relatively simple if we know the bookends of the inverse
pairs and for local client operations, we can track this easily enough.  But
what if the remote operations are the undo/redo pairs?

For those cases, we can use a brute force algoirthm trying to match an
undo against all future operations (and repeating the process until all such
pairs are identified canceled).  An interesting side-effect of this method
is that if an user manually implements undo/redo, this will get detected and
produce a better merge operation.

A simpler and more performant approach would be if operations tagged the other
operation they are inverting.

This package only implements the algoirthm to remove adjacent local undo/redo
pairs. The simplified "undo"/"redo" marker idea is worth considering for the
core "MergeOperations" API but at a later phase of the project.