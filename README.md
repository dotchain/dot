# Operational Transforms Package

[![Status](https://travis-ci.com/dotchain/dot.svg?branch=master)](https://travis-ci.org/dotchain/dot?branch=master)
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

See [Demos](https://dotchain.github.io/demos/).

## Walkthrough of the project

The DOT project is based on *immutable* or *persistent* **values** and
**changes**. For example, inserting a character into a string would
look like this:

```golang
        // import "github.com/dotchain/x/types.S8
        // S8 is DOT-compatible string type with UTF8 string indices
        initial := types.S8("hello")

        // replace "" with " world" at offset = 5 (i.e. end)
        append := changes.Splice{5, types.S8(""), types.S8(" world")}

        // actually apply the change
        updated := initial.Apply(append)

        // now updated == "hello world"
```

The [changes](https://godoc.org/github.com/dotchain/dot/changes)
package implements the core changes: **Splice**, **Move** and
**Replace**.  The logical model for these changes is to treat all
values as either being like *arrays* or like *maps*.  The actual
underlying datatype can be different as long as the array/map
semantics is implemented.

### Composition of changes

Changes can be *composed* together. A simple form of composition is
just a set of changes:

```golang
        initial := types.S8("hello")

        // append " world"
        append1 := changes.Splice{5, types.S8(""), types.S8(" world")}

        // append "."
        append2 := changes.Splice{11, types.S8(""), types.S8(".")}
        
        // now combine the two appends and apply
        both := changes.ChangeSet{append1, append2}
        updated := initial.Apply(both)
```

Another form of composition is modifying a sub-element such as an
array element or a dictionary path:

```golang
        // types.A is an array type and types.M is a map type
        initial := types.A{types.M{"hello": types.S8("world")}}

        // replace "world" with "world!"
        replace := changes.Replace{types.S8("world"), types.S8("world!")}
        path := []interface{}{0, "hello"}
        change := changes.PathChange{path, replace}

        // replace initial[0]["hello"]
        updated := initial.Apply(changes.PathChange{path, replace})
```


The [types](https://godoc.org/github.com/dotchain/dot/x/types) package
implements standard value types (strings, arrays and maps) with which
arbitrary json-like value can be created.

### Convergence

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

### References

There are two broad cases where a JSON-like structure is not quite
enough.

1. Editors often need to track the cursor or selection which can be
thought of as offsets in the editor text.  When changes happen to the
text, for example, the offset would need to be updated.
2. Objects often need to refer to other parts of the JSON-tree. For
example, one can represent a graph using the array, map primitives
with the addition of references. When changes happen, these too would
need to be updated.

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
       // import "github.com/dotchain/streams/text

       // create an UTF8 text stream
       useUTF16 := false
       initial := text.StreamFromString("hello", useUTF16)

       // two changes: append " world" and delete "lo"
       insert := changes.Splice{5, types.S8(""), types.S8(" world")}
       remove := changes.Splice{3, types.S8("lo"), types.S8("")}

       // two versions directly on top of the initial value
       inserted := initial.Append(insert).(*text.Stream)
       removed := initial.Append(remove).(*text.Stream)

       // like persistent types,
       //    inserted.Value() == "helloworld" and removed.Value() = "hel"

       // the converged value can be obtained from both:
       final1 := streams.Latest(inserted).(*text.Stream)
       final2 := streams.Latest(removed).(*text.Stream)

       // or even from the initial value
       final3 := streams.Latest(initial).(*text.Stream)

       // all three are: "helworld"
```

The example above uses **text.Stream** which tracks not just the
changes but the effective value along with the changes.  The
[streams](https://godoc.org/github.com/dotchain/dot/streams) package
defines a
[ValueStream](https://godoc.org/github.com/dotchain/dot/streams#ValueStream)
type that is similar but there is also the ability to work purely with
a change stream with no associated value. This is useful for pure
transformations (such as "scoping" changes to specific fields or array
indices which allows applications to only maintain the values needed
rather than track the whole state).

Those familiar with [ReactiveX](http://reactivex.io/) will find the
streams approach quite similar -- except that streams are guaranteed
to converge.

### Branching of streams

Streams can also be branched a la Git. Changes made in branches do not
affect the master or vice-versa -- until one of Pull or Push are
called.

```golang

        // here a generic change stream is created
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
the client is aware of and the ID of any previous client
*unacknowledged* operation.

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

    // the following two can be used to restart a session
    initialVersion := -1
    unacknowledgedOps := []ops.Op(nil)
    conn := ops.NewConnector(initialVersion, unacknowledgedOps, c, rand.Float64)
    stream := conn.Stream
    conn.Connect()
    defer conn.Disconnect()

    // ... now stream starts receiving updates from the network
    // ... and local changes can also be applied to  it
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

## Backend storage implementations

There are two storage implementations: [a local filesystem based
solution (using
BoltDB)](https://github.com/dotchain/dot/tree/master/ops/bolt) and a
[Postgres](https://github.com/dotchain/dot/tree/master/ops/pg)
solution.

## Not yet implemented

There is no native JS version though the [browser
demo](https://dotchain.github.io/demos/) uses a GopherJS
transpiled version.  At some point, this will be packaged for native
JS consumption.

The async scheduler and the way it interacts with ops Connector are
still a bit awkward to use.

It is also possible to implement cross-object merging (i.e. sharing a
sub-object between two instances by using the OT merge approach to the
shared instance).  This is not implemented here but 

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

