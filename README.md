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

The project marries conflict-free merging with eventually convergent
persistent datastrutures.

## Features and demos

1. Small, well tested mutations that compose for rich JSON-like values
2. Immutable, Persistent types for ease of use
3. Strong references support that are automatically updated with changes
4. Streams and **Git-like** branching, merging support
5. Simple network support (Gob serialization)
6. Rich builtin undo support
7. Folding (committed changes on top of uncommitted changes)
8. Customizable rich types for values and changes


See [Demos](https://dotchain.github.io/demos/). Currently, this
requires running the client from a command line but a browser-based
demo is in the works.

## Walkthrough of the project

The DOT project is based on *immutable* or *persistent* **values** and
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
value with a new value.

Changes can be *composed*.  For example, the **PathChange** type
allows modifying the value at a specific path in the value:

```golang
        initial := types.A{types.S8("hello"), types.S8("world")}
        insert := changes.Splice{5, types.S8(""), types.S8(" world")}
        change := changes.PathChange{[]interface{}{0}, insert}
        updated := initial.Apply(change)
```

The other composition is combining a sequence of changes using
**ChangeSet**.

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
initial state is the basis for the eventual convergence of the data
structures.  The
[changes](http://godoc.org/github.com/dotchain/dot/changes) package 
has fairly intensive tests to cover the change types defined there,
both individually and in composition. 

In addition to convergence, the set of change types are chosen
carefully to make it easy to implement *Revert()* (undo of the
change). This allows the ability to build a generic
[undo stack](https://godoc.org/github.com/dotchain/dot/x/undo) as well
as somewhat fancy features like
[folding](https://godoc.org/github.com/dotchain/dot/x/fold).


The [types](https://godoc.org/github.com/dotchain/dot/x/types) package
implements standard value types (strings, arrays and maps) with which
arbitrary json-like value can be created.

### References

There are two broad cases where a JSON-like structure is not quite
enough.

1. Editors often need to track the cursor or selection which can be
thought of as offsets in the editor text.  When changes happen, these
need to be **merged**.
2. Objects often need to refer to other parts of the JSON-tree. When
changes happen, these would need to be updated

The [refs](https://godoc.org/github.com/dotchain/dot/refs) package
implements a set of types that help work with these.  In particular,
it defines a
[Container](https://godoc.org/github.com/dotchain/dot/refs#Container)
value that allows elements within to refer to other elements.

## Streams

The [streams](https://godoc.org/github.com/dotchain/dot/streams)
package defines the Stream type which is best thought of as a
"convergent persistent data structure".  It is persistent in the sense
that mutations simply return new values leaving the existing values as
is. It is convergent in the sense that all mutations from an initial
value are considered part of the same "family" and iterating on its
**Next()** values will converge all the values to an identical final
value: 

```golang

        initial := streams.ValueStream{types.S8("hello"), streams.New()}

       // two changes: append " world" and delete "lo"
       insert := changes.Splice{5, types.S8(""), types.S8(" world")}
       remove := changes.Splice{3, types.S8("lo"), types.S8("")}

       // two versions directly on top of the initial value
       inserted := initial.Append(insert)
       removed := initial.Append(remove)

       // like persistent types,
       //    inserted == "helloworld" and removed = "hel"

       // the converged value can be obtained from both:
       final1 := streams.Latest(inserted).(streams.ValueStream)
       final2 := streams.Latest(removed).(streams.ValueStream)

       // or even from the initial value
       final1 := streams.Latest(initial).(streams.ValueStream)

       // all three are: "helworld"
```

The example above uses a **ValueStream** which has both the value and
tracks changes but it is possible to just track changes.  One benefit
of doing so is the ability to "scope" changes down, say to a specific
path. This allows removing all unnecessary storage for parts one is
not interested in.

Those familiar with [ReactiveX](http://reactivex.io/) will find the
streams approach quite similar -- except that streams are guaranteed
to converge.

### Branching of streams

Streams can also be branched a la Git. Changes made in branches do not
affect the master or vice-versa -- until one of Pull or Push are
called.

```golang

        master := streams.New()
        local := streams.Branch(master)

        // changes will not be reflected on master yet
        local = local.Append(insert)

        // push local changes up to master
        streams.Push(local)
```

There are other neat benefits to the branching model: it provides a
fine grained control for pulling changes from the network on demand
and suspending it as well as providing a way for making local
changes.

### Network synchronization

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

```golang

import (
       "github.com/dotchain/dot/x/nw"
       "github.com/dotchain/dot/ops"
       "github.com/dotchain/dot/streams"
       "github.com/dotchain/dot/x/idgen"       
)

func connect() streams.Straem {
    c := nw.Client{URL: ...}`
    defer c.Close()

    stream := streams.New()
    sync := ops.NewSync(ops.Transformed(c), -1, stream, idgen.New)

    go func() {
       version := 0
       for {
            ctx := context.WithTimeout(context.Background(), time.Second*30)
            sync.Poll(ctx, version)
            sync.Fetch(ctx, version, 1000)
            version = sync.Version()
       }
    }()

    return stream
}
    
```

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

The network API requires a careful management of event loops. This
should be simplified.

There is no snapshot storage mechanism (for operations as well as full
values) which would require replaying the log each time.

It is also possible to implement cross-object merging (i.e. sharing a
sub-object between two instances by using the OT merge approach to the
shared instance).  This is not implemented here but 

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

