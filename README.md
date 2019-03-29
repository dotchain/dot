# Operational Transforms Package

[![Status](https://travis-ci.com/dotchain/dot.svg?branch=master)](https://travis-ci.com/dotchain/dot?branch=master)
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

## Status

Most of the code here is quite stable at this point but the project is
not yet ready for production:

1. Comprehensive end-to-end **stress** tests are missing.
2. The ops/nw package and how it integrates with streams.Async is a
bit wonky.

## Features

1. Small, well tested mutations that compose for rich JSON-like values
2. Immutable, Persistent types for ease of use
3. Strong references support that are automatically updated with changes
4. Streams and **Git-like** branching, merging support
5. Simple network support (Gob serialization)
6. Rich builtin undo support
7. Folding (committed changes on top of uncommitted changes)
8. Customizable rich types for values and changes

## Contents
1. [Status](#status)
2. [Features](#features)
3. [TODO Example](#todo-example)
    1. [Server](#server)
    2. [Types](#types)
    3. [Type registration](#type-registration)
    4. [Toggling Complete](#toggling-complete)
    5. [Changing description](#changing-description)
    6. [Adding Todos](#adding-todos)
    7. [Client connection](#client-connection)
4. [Walkthrough of the project](#walkthrough-of-the-project)
    1. [Composition of changes](#composition-of-changes)
    2. [Convergence](#convergence)
    3. [References](#references)
5. [Streams](#streams)
    1. [Branching of streams](#branching-of-streams)
    2. [Network synchronization](#network-synchronization)
6. [Backend storage implementations](#backend-storage-implementations)
7. [Undo log, folding and extras](#undo-log-folding-and-extras)
8. [Not yet implemented](#not-yet-implemented)
9. [Contributing](#contributing)

## TODO Example

The following walkthrough demonstrates the project by means of the
standard TODO-MVC example except that in this case, the application is
collaborative -- multiple clients can modify the same data and all
client UIs are expected to converge to the same data/visuals.

### Server

The DOT backend is essentially a simple log store with only append
operations and no modifications.  This is irrespective of whatever
types are used for the application state itself:

```go global

func Server() {
	// import net/http
	// import github.com/dotchain/dot/ops/nw
	// import github.com/dotchain/dot/ops/bolt

        // uses a local-file backed bolt DB backend
	store, _ := bolt.New("file.bolt", "instance", nil)
        defer store.Close()
	http.Handle("/api/", &nw.Handler{Store: store})
        http.ListenAndServe(":8080", nil)
}
```

The example above uses the
[Bolt](http://godoc.org/github.com/dotchain/dot/ops/bolt)
implementation of the store.  There is also a
[Postgres](http://godoc.org/github.com/dotchain/dot/ops/pg) backend
available.

### Types

A TODO MVC app consists of only two core types: `Todo` and `TodoList`:

```go global

// Todo tracks a single todo item
type Todo struct {
	Complete bool
        Description string
}

// TodoList tracks a collection of todo items
type TodoList []Todo

```

These types are incomplete as far as DOT is concerned because they do
not specify how to change them.  All `values` in DOT are expected to
be immutable and support the
[Value](https://godoc.org/github.com/dotchain/dot/changes#Value)
interface (or in the case of lists like `TodoList`, also implement the
[Collection](https://godoc.org/github.com/dotchain/dot/changes#Collection)
interface).  This allows for structured convergent mutations.

For example, such an implementation would indicate that `Complete` can
be modified independently and insertions into `TodoList` can happen
along with deletions etc.

For the most part, these implementations are routine for structs,
unions and sllices and so, they can be [code generated](codegen.md).

### Type registration

To use the types across the network, they have to be registered with
the codec (which will be `gob` in this example)

```go global
// import encoding/gob

func init() {
	gob.Register(Todo{})
        gob.Register(TodoList{})
}
```

### Toggling Complete

The code generation in DOT produces not only values, but also the
associated streams which allows standard types of mutations:

```go global
func Toggle(t *TodoListStream, index int) {
	// TodoListStream.Item() is implemented in the generated
        // code and returns *TodoStream
	itemStream := t.Item(index) 

	// Complete() is also implemented in the generated code.
        completeStream := itemStream.Complete()

	// Update() here refers to streams.Bool.Update
        completeStream.Update(!completeStream.Value)
}
```

Note that the function does not return any value here but the updates
can be fetched by calling `.Latest()` on any of the corresponding
streams. If a single stream instance has multiple edits, the
`Latest()` value is the merged value of all those edits.  If any
substreams (such as those produced by `t.Item` or
`itemStream.Complete`), then updates to the substreams gets reflected
on their parent streams appropriately.

Each stream also exposes the underlying type (such as `Todo` or
`TodoList`) via the `Value` field.

### Changing description

The code for Changing description is similar.  The string
`Description` field in `Todo` maps to a `streams.S16` stream in
`TodoStream` which allows `Update()` to modify the value.

But to make things interesting, lets look at **splicing** rather
than replacing the whole string. Splicing is taking a subsequence of
the string at a particular position and replacing it with the provided
value.  It captures insert,  delete and replace in one operation.

This probably better mimics what text editors do and a benefit of such
high granularity edits is that when two users edit the same text, so
long as they don't directly touch the same characters, the edits will
merge quite cleanly.

The `Splice` method is already implemented on the underlying string
stream and so the code here looks quite similar.

```go global
func SpliceDescription(t *TodoListStream, index, offset, count int, replacement string) {
	// TodoListStream.Item() is implemented in the generated
        // code and returns *TodoStream
	itemStream := t.Item(index) 

	// Description() is also implemented in the generated code.
        descStream := itemStream.Description()

	// Splice() here refers to streams.S16.Splice
        descStream.Splice(offset, count, replacement)
}
```

### Adding Todos

Adding a Todo is relatively simple as well:

```go global
func AddTodo(t *TodoListStream, todo Todo) {
	t.Splice(len(t.Value), 0, todo)
}
```

The use of `Splice` in this example should give away the fact that
items can be inserted into a specific position or deleted the same
way.

In addition to `Splice`, the oher basic operation available on
collections is `Move` (though this is not implemented in the automatic
code generation yet).

### Client connection

Setting up the client requires connecting to the URL where the server
is hosted. In the example below, a render function is expected as a
parameter. 

```go global
// import github.com/dotchain/dot/ops/nw
// import github.com/dotchain/dot/ops
// import math/rand

func Client(stop chan struct{}, url string, render func(*TodoListStream)) {
	version, pending, todos := SavedSession()

	store := &nw.Client{URL: url}
        defer store.Close()
        client := ops.NewConnector(version, pending, ops.Transformed(store), rand.Float64)
	stream := &TodoListStream{Stream: client.Stream, Value: todos}

	// start the network processing
	client.Connect()

        // save session before shutdown
	defer func() {
        	SaveSession(client.Version, client.Pending, stream.Latest().Value)
        }()
        defer client.Disconnect()

        client.Stream.Nextf("key", func() {
        	stream = stream.Latest()
        	render(stream)
        })
        render(stream)
	defer func() {
        	client.Stream.Nextf("key", nil)
        }()

	<- stop
}


func SaveSession(version int, pending []ops.Op, todos TodoList) {
	// this is not yet implemented. if it were, then
        // this value should be persisted locally and returned
        // by the call to savedSession
}

func SavedSession() (version int, pending []ops.Op, todos TodoList) {
	// this is not yet implemented. return default values
        return -1, nil, nil
}

```

### Running the demo

The TODO MVC demo is in the
[example](https://github.com/dotchain/dot/tree/master/example)
folder.

The snippets in this markdown file can be used to generate the
**todo.go** file and then auto-generate the "generated.go" file:

```sh
$ go get github.com/tvastar/test/cmd/testmd
$ testmd -pkg example -o examples/todo.go README.md
$ testmd -pkg main codegen.md > examples/generated.go
```

The server can then be started by:

```sh
$ go run server.go
```

The client can then be started by:

```sh
$ go run client.go
```

The provide client.go stub file simply appends a task every 10
seconds.  A real browser-based UX app is in the works.

## Design details of the project

The DOT project is based on *immutable* or *persistent* **values** and
**changes**. For example, inserting a character into a string would
look like this:

```golang skip
        // import github.com/dotchain/changes/types.S8
        // S8 is DOT-compatible string type with UTF8 string indices
        initial := types.S8("hello")

        // replace "" with " world" at offset = 5 (i.e. end)
        append := changes.Splice{5, types.S8(""), types.S8(" world")}

        // actually apply the change
        updated := initial.Apply(append)

        // now updated == "hello world"
```

A less verbose *stream* based version (preferred) would look like so:

```golang skip
        // import github.com/dotchain/streams

        initial := &streams.S8{Stream: streams.New(), Value: "hello"}
        updated := initial.Splice(5, 0, " world")

        // now updated.Value == "hello world"
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

```golang skip
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

```golang skip
        // types.A is an array type and types.M is a map type
        initial := types.A{types.M{"hello": types.S8("world")}}

        // replace "world" with "world!"
        replace := changes.Replace{types.S8("world"), types.S8("world!")}
        path := []interface{}{0, "hello"}

        // replace initial[0]["hello"]
        updated := initial.Apply(changes.PathChange{path, replace})
```


The [types](https://godoc.org/github.com/dotchain/dot/changes/types) package
implements standard value types (strings, arrays and maps) with which
arbitrary json-like value can be created.

### Convergence

The core property of all changes is the ability to guarantee
*convergence* when two mutations are attempted on the same state:

```golang skip
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
[undo stack](https://godoc.org/github.com/dotchain/dot/streams/undo) as well
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

```golang skip
       // import github.com/dotchain/streams

       // create a text stream
       initial := streams.S8{Stream: streams.New(), Value: "hello"}

       // two changes: append " world" and delete "lo"
       inserted := initial.Splice(5, 0, " world")
       removed := initial.Splice(3, len("lo"), "")

       // inserted.Value == "hello world" && removed.Value == "hel"

       // the converged value can be obtained from both:
       final1 := inserted.Latest()
       final2 := removed.Latest()

       // final1.Value == final2.Value && final1.Value == "hel world"

       // or even from the initial value
       final3 := initial.Latest()
       // final3.Value == "hel world"
```

The example above uses **streams.S8** which is a strongly typed
stream implemented on the weakly typed **streams.Stream**.  The
default stream implemention provided via **streams.New** only tracks
the stream changes and guarantees the correct seqeunce of changes for
convergence. The strongly typed **streams.S8** is built on top of that
to also track the current value.  Any custom type can similarly be
defined quite easily.

Those familiar with [ReactiveX](http://reactivex.io/) will find the
streams approach quite similar.   Except:

1. Streams here inherently expect multiple writers in different
context.
2. Each stream instance is immutable -- appending changes produces new
stream instances connected to the current one.
3. Multiple changes on the same stream instance cause convergence:
i.e. the result of any individual edit is consistent with the current
change but all other changes show as "futures" -- calling Next()
sequentially gets all the other changes made and Latest() is
guaranteed to be the same for all callers.
4. Streams are encouraged to be strongly typed.  So, methods like
filter or map etc are not inherently provided.  But code generation
for strongly typed streams is available at
[dotc](https://godoc.org/github.com/dotchain/dot/x/dotc)

### Branching of streams

Streams can also be branched a la Git. Changes made in branches do not
affect the master or vice-versa -- until one of Pull or Push are
called.

```golang skip

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

```golang skip

import (
       "github.com/dotchain/dot/ops/nw"
       "github.com/dotchain/dot/ops"
       "github.com/dotchain/dot/streams"
       "github.com/dotchain/dot/x/idgen"       
)

func connect() streams.Straem {
    c := nw.Client{URL: ...}
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

## Backend storage implementations

There are two storage implementations: [a local filesystem based
solution (using
BoltDB)](https://github.com/dotchain/dot/tree/master/ops/bolt) and a
[Postgres](https://github.com/dotchain/dot/tree/master/ops/pg)
solution.

A simple HTTP server can be created using the bolt/pg store implementations:

```golang skip
        // import github.com/dotchain/dot/ops/bolt
        // import github.com/dotchain/dot/ops/nw       

        store, _ := bolt.New("file.bolt", "instance", nil)
        defer  store.Close()
        handler := &nw.Handler{Store: store}
        h := func(w http.ResponseWriter, req  *http.Request) {
                // Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if req.Method == "OPTIONS" {
			return
		}
		handler.ServeHTTP(w, req)
	}
        http.HandleFunc("/api/", h)
        http.ListenAndServe()
```

## Undo log, folding and extras

The streams abstraction provides the basis for implementing
system-wide
[undo](https://godoc.org/github.com/dotchain/dot/streams/undo).

More interestingly, there is the ability to implement **Folding**. A
client can have a set of temporary changes (such as config or view
etc) which is not committed.  And then more changes can be made on top
of it which **are committed**.  These types of shenanigans is possible
with the use of a small fixed set of well-behaved changes.

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

