# Operational Transforms Package

[![Status](https://travis-ci.org/dotchain/dot.svg?branch=master)](https://travis-ci.org/dotchain/dot?branch=master)
[![GoDoc](https://godoc.org/github.com/dotchain/dot?status.svg)](https://godoc.org/github.com/dotchain/dot)
[![codecov](https://codecov.io/gh/dotchain/dot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotchain/dot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotchain/dot)](https://goreportcard.com/report/github.com/dotchain/dot)

The DOT project is a blend of [Operational
Transformation](https://en.wikipedia.org/wiki/Operational_transformation),
[Persistent
Datastructures](https://en.wikipedia.org/wiki/Persistent_data_structure)
and [reactive](https://en.wikipedia.org/wiki/Reactive_programming)
stream processing.

The whole project is in the middle of a refactoring and as such isn't
stable yet.

## Features and demos

1. Small, well tested mutations that compose for rich JSON-like values
2. Immutable, Persistent types for ease of use
3. Rich builtin undo support
4. Folding (committed changes on top of uncommitted changes)
5. Strong references support that are automatically updated with changes
6. Streams and **Git-like** branching, merging support
7. Customizable rich types for values and changes
8. Simple network support (Gob serialization)

See [Demos](https://github.com/dotchain/demos). Currently, this
requires running the client from a command line but a browser-based
demo is in the works.

## Tutorial

The DOT project is based on *immutable* **values** and
**changes**. For example, inserting a character into a string would
look like this:

```golang
        initial := types.S8("hello")
        insert := changes.Splice{5, types.S8(""), types.S8(" world")}
        updated := initial.Apply(insert)
        // now updated == "hello world"
```

The [changes](https://godoc.org/github.com/dotchain/dot/changes)
package implements the core changes: **Splice**, **Move** and
**Replace**.  The logical model for these changes is to treat all
values as either being like *arrays* (in which case the first two
operations apply) or *map like*.  The **Replace** change replaces any
value with a new value and when combined with a **PathChange** it can
replace any value within a map.  In addition to **PathChange**,
changes can be combined with **ChangeSet**.  Custom changes can be
implemented as well for rich text, for instance).

The core property of all changes is the ability to guarantee
*convergence* when two mutations are attempted on the same state:

```golang
       initial := types.S8("hello")

       // two changes: append " world" and delete "lo"
       insert := changes.Splice{5, types.S8(""), types.S8(" world")}
       remove := changes.Splice{3, types.S8("lo"), types.S8("")}

       // two versions
       inserted := initial.Apply(insert)
       removed := initial.Apply(remove)

       // merge
       removex, insertx := insert.Merge(remove)

       // converge
       final1 := inserted.Apply(removex)
       final2 := removed.Apply(insertx)
       // now final1 == final2
```

The ability to *merge* two independent changes done to the same
initial state is central to operational transformation and all the
changes defined in the
[changes](http://godoc.org/github.com/dotchain/dot/changes) package
implement this with fairly intensive tests to cover them individually
and in composition.

In addition to convergence, the set of change types are chosen
carefully to make it easy to implement undo (inversion of the
change). This allows the ability to build
[folding](https://godoc.org/github.com/dotchain/dot/x/fold) as well as
a general purpose
[undo](https://godoc.org/github.com/dotchain/dot/x/undo) stack.

The [types](https://godoc.org/github.com/dotchain/dot/x/types) package
implements standard value types (strings, arrays and maps) with which
fairly rich types can be composed.

Real world applications often have the need to work with non-tree data
structures (such as graphs or with pointers). To enable working with
these, the [refs](https://godoc.org/github.com/dotchain/dot/refs)
package defines a few types of *references*: **Path**, **Caret** and
**Range**.  For example, a text editor would need to track the current
selection and if any remote change modifies the text, the selection
would have to be carefully updated.  The
[refs](https://godoc.org/github.com/dotchain/dot/refs) package provide
the types needed to manage this.  In addition, it defines the concept
of **List** of references: a value can maintain a reference to another
part of the value, like a pointer. This allows non-tree structures to
be represented.

The [streams](https://godoc.org/github.com/dotchain/dot/streams)
package provides the basic abstraction of a **stream**.  This is
similar to event emitters and such with a twist: it works with
immutable objects and can capture the idea of merging. Consider the
same example of merging strings modeled with streams:

```golang

        initial := streams.ValueStream{types.S8("hello"), streams.New()}

       // two changes: append " world" and delete "lo"
       insert := changes.Splice{5, types.S8(""), types.S8(" world")}
       remove := changes.Splice{3, types.S8("lo"), types.S8("")}

       // two versions
       inserted := initial.Append(insert)
       removed := initial.Append(remove)

       // latest can be obtained from inserted or removed
       final1 := inserted
       for _, next := latest.Next(); next != nil; _ next = latest.Next() {
           latest = next
       }

       // final2 iterates on removed and is guranteed to have the same result
```

Basically, a stream instance acts like an immutable object in that any
changes `Appended` to it leave the original stream alone and produce a
new instance.  So, two separate mutations on the same stream will not
see the effect of the other.  But unlike immutable objects, streams
provide the ability to "navigate" the not-yet-merged changes using
*Next()* and get to the converged state: All streams in the same
family (i.e created by a tree of *Append* calls) flow into the same
final destination.

The example above uses a **ValueStream** which has both the value and
tracks changes but it is possible to just track changes.  One benefit
of doing so is the ability to "scope" changes down.  For example, to
only listen for changes at a particular path, say "rows/id:99k", we
can do `streams.ChildOf(baseStream, "rows", "id:99k")`.  It is also
possible to filter out paths etc.

The *streams* approach also provides a few powerful tools: the ability
to branch and merge.  For instance, a version of a particular *value*
can be branched off and changes made to it and eventually merged in:

```golang

        branch := streams.Branch{streams.New(), streams.New()}
        master := streams.ValueStream{initial, branch.Master}
        local := streams.ValueStream{initial, branch.Local}


        // changes will not be reflected on master yet
        local = local.Append(insert)

        // merge will push local up to master and pull down
        // any changes from master
        branch.Merge()
```

This powerful functionality is quite similar to
[Git](https://en.wikipedia.org/wiki/Git) and other source control
systems except it is applied to **data structures**

There are other neat benefits to the branching model: it provides a
fine grained control for pulling changes from the network on demand
and suspending it as well as providing a way for making local
changes.

## Network synchronization

DOT uses a fairly simple backend
[Store](https://godoc.org/github.com/dotchain/dot/ops#Store)
interface: an append-only dumb log. Each operation that is appended in
the log gets an incrementing integer version (starting at zero). DOT
allows operation pipe-lining (i.e. it doesnt wait for acknowledgments
from the server before sending more operations) and to clarify the
exact sequence, every operation carries both the last server version
the client is aware of and any previous client operation.

The [ops](https://godoc.org/github.com/dotchain/dot/ops) package takes
these raw entries in the log and provides the synchronization
mechanism to connect it to a stream which allows much of the
client/app logic to be written agnostic of the network.

## Undo log, folding and extras

The streams abstraction provides the basis for implementing
system-wide
[undo](https://godoc.org/github.com/dotchain/dot/x/undo).

More interestingly, there is the ability to implement **Folding**. A
client can have a set of temporary changes (such as config or view
etc) which is not committed.  And then more changes can be made on top
of it which **are committed**.  These types of shenanigans is possible
with the use of a small fixed set of well-behaved changes.

## Not yet implemented

There is no native JS version.

The storage layer is an in-memory version though it is relatively
simple to build storage backends given the very trivial storage
interface. 

The network API does not include authentication and authorization.

There is no snapshot storage mechanism (for operations as well as full
values) which would require replaying the log each time.

It is also possible to implement cross-object merging (i.e. sharing a
sub-object between two instances by using the OT merge approach to the
shared instance).  This is not implemented here but 

## Reactive computation, scheduler

Streams in DOT also allow for change notifications for change that
have not yet happened. These notifications are the same as following
the changes from a particular stream that is guaranteed to converge.

But streams programming with synchronous notification can get hairy
with reentrancy and locking issues.  The
[AsyncScheduler](https://godoc.org/github.com/dotchain/dot/streams#AsyncScheduler)
provides for an elegant way to manage an event pump/loop and avoid
re-entrancy issues. 

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

