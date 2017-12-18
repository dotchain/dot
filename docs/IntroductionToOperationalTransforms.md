# Introduction to Operational Transforms

Operational Transformation is an approach to conflict-free distributed
synchronization of data.  It allows multiple clients, starting from the same 
initial state, to independently make local changes to their state and eventually
converge.

There is a lot of information about this on the web and elsewhere. There are also
a variety of ways to implement this technique.

The documentation here covers the specifics of the technique used here from the 
point of view of a beginner who may not know much about operational transforms. 
In particular, the terms used here may not accurately match terms used elsewhere.

A lot of effort has been made here to avoid jargon or overly technical terms and
contributions or questions relating this section are not only welcome but also
appreciated.

## The general approach

Central to the approach here is clients representing changes to their state via
small mutations.

Consider the example of a collaborative editor where multiple clients can edit
the same text.  Lets start with an initial text of:

```
	The fox jumped over the fence
```

Now when one user edits the sentence to read `The red fox jumped over the fence`,
its client would represent this change as something like this:

```golang
	ot.Splice{Offset: 4, Before: "", After: "red "}
```

This change represents an **insertion** of the string `"red "`.  Representing all
edits as a sequence of changes is powerful in itself -- it allows recreating any
state from only the list of changes as well as potentially recreating the state at
a particular moment in the past.

But it is also a crucial element of the approach of synchronization as this is
done by exchanging the sequence of changes.

**Terminology**  The documentation sometime refers to changes as mutations. An
**Operation** is a different beast and is covered a little later. The unit of
synchronization is referred as a **Model** 

Now consider a second user editing the same text in parallel, changing the sentence
to read `The fox jumped over the beautiful fence`.

This second change would be represented as:

```golang
	ot.Splice{Offset: 24, Before: "", After: "beautiful "}
```

In particular, note that the offset 24 refers to the position where the text `"beautiful "`
is inserted into the sentence.  It currently does not account of the other client having
inserted `"red "` earlier.

This situation can be captured with a simple diagram.

![Fig1](https://cdn.rawgit.com/dotchain/dot/master/assets/fig1.svg)

Now lets assume both clients send their mutations to the server which then proceeds to
share this back with the other clients.

**Note #1**  Not all operations transforms approaches require a server.  But in practice
almost all of them use some form of a server atleast to maintain a persistent record.
The approach used here relies on a **Journal Service** to provide both persistence of
these mutations but also to provide a consistent order of changes which helps get all
clients to converge and agree on the same final state.

Lets further assume that the **Journal Service** orders the `"red "` change before the
`"beautiful "` change.

This can be reprsented visually as follows:

![Fig2](https://cdn.rawgit.com/dotchain/dot/master/assets/fig2.svg)

When the second client receives the `"red "` change, *if it applies this change 
directly*, it will end up with a text of `The red fox jumped over the beautiful fence`
which matches the desired outcome.

But when the first client receives the `"beautiful "` change, *if it applies this change
directly*, it will end up with a test of `The red fox jumpled overbeautiful the fence`.

Why?  Because the offset in the second change of `24` was made on the initial sentence
and does not account for the change made by the first client meanwhile.

This is where the heart of operation transforms comes in.  It provides an algorithm that
takes a *new change* (such as `insert "beautiful " at offset 24`) and a *previous change*
(`insert "red " at offset 4`) and transforms the new change against the old so that
we have a *rebased new change* and a *merged old change*

In this particular case, the algorithm would yield:

```golang

	Rebased Change: ot.Splice{Offset: 28, Before: "", After: "beautiful "}
	Merged Change: ot.Splice{Offset: 4, Before: "", After: "red "}
```

If the first client applies the **rebased change** (i.e change that has been rebased to
the latest journal entry) to its local state, it would get the expected outcome:
`The red fox jumped over the beautiful fence`.

If the second client applies the **merged change** (i.e a change that happend before but
has not been merged into the current state) to its local state, it would get
the same out come as well. 

In other words, both clients then end up with the same state -- convergence!

A graphic way to represent this is as follows from the point of view of the first client.

![Fig3](https://cdn.rawgit.com/dotchain/dot/master/assets/fig3.svg)

Another way to look at this act of transformation is to consider the **rebased change**
as a change which has the `intention` of the *local change* but which has been altered to be
safe to apply after the *remote change*.  And the **merged change** can correspondingly
be thought of as having captured the `intention` of the *remote change* but which has been
adapted to factor in the already applied *local change*

**Note 2** In the example above, notice that the second client didn't have to modify the
change and it worked right away.  Why does this work for second client where it ends up 
with the correct outcome?
This just happens to be an accidental outcome because of the specific changes we considered
here.  There are other edits which can cause both clients to be wrong.  
Note that the **rebased change** for the transformation is
the same as the original change -- this is why it worked for the second client before.

### Single transformation

When a change A is being transformed against change B (where both changes were based on
the same initial state):

1. The **Rebased** change A can be safely applied by any client which has applied B
2. The **Merged** change B can be safely applied by any client which has applied A
3. Both of these clients will end up with the same result.

```
	A + Merged(B) = B + Rebased (A)
```

![Fig4](https://cdn.rawgit.com/dotchain/dot/master/assets/fig4.svg)

**Note 3** Notice the terms `rebase` and `merge` are reminiscent of `git` terms.  For those
who are familiar with `git`, there is a strong analogy between operational transforms and
git and the terms are referring to similar processes.

**Note 4** If we carefully choose the format and structure of the changes, it is possible
to define this transformation function without reference to the state of the client at all.
In other words, it is possible to define a transformation function that only takes as input
the two changes.  In fact, this is the approach taken here.  For example, the function to
transform two splices is `mergeSpliceSpliceSamePath` which only looks at the two splices.

### Multiple transformations

The next wrinkle to consider is what if each of the clients had done mulitple changes?

When there are two sequences of changes A1, A2, ...An and B1, B2, ...Bn and both sequences
are based on the same initial state, it is possible to use the single transformation procedure
over and over until sequence A is transformed against sequence B to yield a **Rebased sequence A**
and **Merged sequence B** such that:

1. The **Rebased sequence A** can be safely applied by any client which has applied B
2. The **Merged sequnce B** can be safely applied by any client which has applied A
3. Both of these clients are guaranteed to end up with the same state.

How does this procedure work?  The exact function which does it is 
[MergeChanges](https://github.com/dotchain/dot/blob/master/transformer.go#L78) and it is
defined recursively.  

Here is a pictorial representation...

![Fig 5](https://cdn.rawgit.com/dotchain/dot/master/assets/fig5.svg)

The following induction method considers how to calculate the rebase and merge of a
sequence of changes (of size n + 1) against a single change B -- based on the recursively
calculated rebase and merge of the subsequence of size n against B.

```
	Consider transform(A[1:n], B) = RebasedA[1:n] and MergedB
	where: A[1:n] + MergedB = B + RebasedA[1:n]

	Now consider a change that happens after A[n] = A[n+1]
	Note that this is in parallel to MergedB (i.e. A[n+1] and MergedB are both based
	on the same state:  A[1:n]).  So, we can transform A[n+1] against MergedB.

	transform(A[n+1], MergedB) = RebasedA[n+1], FullyMergedB 
	where: A[n+1] + FullyMergedB = MergedB + Rebased[n+1]
	
	Now lets add A[n+1] before both: 
	A[1:n] + A[n+1] + FullyMergedB = A[1:n] + MergedB + Rebased[n+1]

	Now considering the equality A[1:n] + MergedB = B + RebasedA[1:n],
	we get:

	A[1:n] + A[n+1] + FullyMergedB = B + RebasedA[1:n] + Rebased[n+1]
	
	In other words:
	A[1:n+1] + FullyMergedB = B + Rebase[1:n+1]

	This translates into the following recursive algorithm:

	function transform(A, b) {
		 if A.length === 1 {
		    // use single transformation code
		 }
		 let {rebased, merged} = transform(A.slice(0, A.length - 1), B)
		 let {rebased2, merged2} = transform(A[A.length - 1], merged)

		 return {
		 	rebased: rebased.concat(rebased2),
			merged: merged2
		 }
	}
```

The same mechanism can be extended to deal with transformation of a sequence of changes
against a sequence of changes (instead of just a single change).

## Operations

The approach detailed here consists of a small number of mutation types.  It is often
important to clients to group together these changes for semantic purposes (because any
intermediate state may not satisfy the application's integrity constraints) as well as 
for performance reasons.

An *operation* is a sequence of changes (they are meant to be applied in strict order)
annotated with a few extra properties.

## Journal Service

The journal service provides for an append-only persistent storage of operations as 
well as notifications to clients when operations have been appended.  The journal service
does not modify operations -- in particularly, it does not transform operations. The
journal service also guarantees that all clients will see the operations in the same
order.

The sequence of operations as ordered by the journal service is also loosely referred
to as the **journal**

**Note 5** Some OT systems have a server component that does more work but the design
chosen here has a very minimal server side component as a deliberate choice as it makes
maintenance and scale of this service relatively cheap.

## Log Service

A client that wishes to use the journal service needs to implement the reconcilation
algoirthm described here which is a bit intense.  The log service provides the
reconciliation effort -- clients can submit raw operations and get back compensating
actions to fix up their state.

The "rebased" sequence of operations in the Journal (i.e. each operations transformed
so a client can apply them to rebuild the model) is referred to as a **log**.

## Stale operations

The transformation logic worked out above performs well when both sequence of changes
had the same initial state.  But how do clients communicate which version their
changes were based on and how do we deal with transformations when the initial versions
behind the changes do not match?

Since we have a journal service which guarantees order of operations, we can use this
as a basis to refer to the version of the state.  That is, the ID of the operation in
the **journal** can stand for the version of the state that is obtained by applying 
all the operations in the journal up to that operation.  

The **basis ID** of an operation is the ID of the last operation in the **journal**
which has been factored into the client state (the client being referred to here is
the client which originated the operation). We also loosely refer to the
basis of a particular operation to mean the state produced by the act of applying all
operations up to and including the operation that matches the basis ID (in the journal).

Consider a client which sends an operation in flight and follows up immediately with
another operation.  Since there were no server changes factored into the client in
between these, these two will have the same **basis ID**.  To indicate that the
last operation is meant to be applied on top of the previous one, the client would 
need to tack on a **parent ID** to the last operation and set it to the **ID** of the
previous in-flight operation.

The **parent ID** of an operation is the previous in-flight operation from the same
client on top of which the current operation is meant to be applied.  Since mulitple
clients may be active, by the time the  **parent** operation shows up in the **journal**, 
there may be more operations (from other clients) between the **basis** and the **parent** 
and similarly more operations may end up appearing between the **parent** and 
the current operation**.

**Note 6** Some OT systems do not allow multiple operations in flight.  They insist
on forcing the client to wait for the first operation to appear in the **Journal**
before they send another operation.  Even with this approach, it is preferable to 
maintain the separation between **basis** and **parent** as it will allow the clients
to keep their operations untransformed.

With this understanding, we can look at the **journal** slightly differently.  It
is a sequence of operations, where a particular operation may have its **basis ID**
and **parent ID** pointing at any prior operation.

**Note 7** Note that mulitple operations can have the same **basis ID** because multiple
clients may have made parallel changes on top of this.  But mulitple operations cannot
have the same **parent ID** -- a client is not allowed to fork off a branch.  While this
can theoretically be supported, the current approach does not handle this.  So, well
behaved clients should ensure this. 

## Rebased Journal

The sequence of operations in the **journal** cannot be directly applied to a model
because an operation is not guaranteed to be based on the previous operation in 
the **journal**. But using the transformation algorithms already defined, it is 
posible to build a procedure for doing this.

For each operation in the **Journal**, we define its **rebased operation** and **merge chain**
as follows:

* The **rebased journal operation**, when applied on top of all previous rebased journal operations
  captures the state of journal up to this operation.  In other words, the rebased operations are
  all based sequentially and can be applied on top of each other to recreate the state up to that point.
* The **merge chain** for a given operation, when applied on top of the state of the client
  which sent the operation (i.e. the state after the client had applied this operation) will
  converge to the same state as above.  That is, clients can use the **merge chain** to sync up.
  Note that the **merge chain** is a sequence of operations. In fact,
  since the basis ID of the current operation may point to several operations before it
  in the **journal**, the **merge
  chain** of operations is expected to have one entry for each operation in between.

The algorithm for calculating the **rebased operation** and **merge chain** for an operation
in the **journal** is recursively defined based on the **rebased operation** and **merged chain**
of its predecessors!

Consider a *journal* which has already calculated all the **rebased operation** and **merge chain**

Now consider a new operation which lands on this journal. If this operation only has a basis ID
and no parent, things are a bit simple.  The lack of a parent ID implies that the operation was
executed directly on the basis of the operation -- i.e. one can consider all rebased operations
after the basis as being parallel to the new operation.  So, the procedure is to simply transform
the new operation against the **rebased journal operations**, 
but only those that appear after the basis in the **journal**. The transformation yields both
a **rebased operation** and a **merge chain** both of which are the **rebased journal operation**
and the **merge chain** needed for the new operation.

What if the new operation landing on the journal has a parent? Consider the **merge chain** of the
parent operation concatenated with all the **rebased journal operations** that appear after the
**parent** in the **journal**.
The client which submitted the parent operation can immediately apply this sequence to get up to date
with the journal. But the client has applied the current new operation in parallel.  So, the new 
operation is transformed against this sequence to obtain the new **rebased journal operation**
and the new **merge chain**

A minor wrinkle: the new operation may have a basis that is **later** than the previous operation.
Consider the example of a client that submits a couple of operations (A1 with basis ID b1, A2 with
basis ID b1 as well but parent ID = A1). Then it receives a new operation in the journal. It uses
an approach where it acts as if A1 and A2 are at end of the of the journal and figures out the 
fake **merge chain** for A2 (it is fake because A2 has not appeared on the **journal** yet and
there may be more operations that appear before A1) to get synced up with the new basis.  Now, if
this client makes another change and sends this operation up to the server, the new operation will
have a **later** basis than its parent.

So, we amend the procedure very slightly for calculating the **rebased operation** an **merge chain**
of a new operation whose basis is later than that of its parent.  The *unmodified* transformation
procedure above uses a concatenation of the **merge chain** of the parent operation and the
**rebased operations**
that follow the parent operation.  The *modification* is to skip all the leading operations from
this sequence which are earlier than the basis of the new operation. The logic for that is that
all these operations have already been factored into the new operation and transforming against
them is not appropriate.

With this change, the algorithm for calculating the rebased operation and merge chain looks
like [TransformAgainstLag](https://github.com/dotchain/dot/blob/master/transformer.go#L157))

## Reconciliaton

It is possible to extend the mechanism above to take care of an additional set of 
local operations.  In this case, a new remote operation will not only have to be
rebased against the known **journal** (as described above), the **local** operations 
would also need to be transformed against the **rebased** incoming remote operation.

The actual implementation of such a process would yield a transformed chain of
operations, which when applied to the local model would have the same effect as
if the remote operation were applied before the pending **local** operations.

This is implemented [here]((https://github.com/dotchain/dot/blob/master/sync/client_log.go)
