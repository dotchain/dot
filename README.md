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
    4. [Code generation](#code-generation)
    5. [Toggling Complete](#toggling-complete)
    6. [Changing description](#changing-description)
    7. [Adding Todos](#adding-todos)
    8. [Client connection](#client-connection)
    9. [Running the demo](#running-the-demo)
4. [How it all works](#how-it-all-works)
    1. [Applying changes](#applying-changes)
    2. [Applying changes with streams](#applying-changes-with-streams)
    3. [Composition of changes](#composition-of-changes)
    4. [Convergence](#convergence)
    5. [Convergence using streams](#convergence-using-streams)
    6. [Revert](#revert)
    7. [References](#references)
    8. [Branching of streams](#branching-of-streams)
    9. [Network synchronization](#network-synchronization)
5. [Backend storage implementations](#backend-storage-implementations)
6. [Undo log, folding and extras](#undo-log-folding-and-extras)
7. [Not yet implemented](#not-yet-implemented)
8. [Contributing](#contributing)

## TODO Example

The standard TODO-MVC example demonstrates the features of
collaborative (eventually consistent) distributed data structures.

### Server

The DOT backend is essentially a simple log store.  All mutations to
the application state are represented as a **sequence of operations**
and written in append-only fashion onto the log.  The following
snippet shows how to start a web server (though it does not include
authentication or CORs for example).

```go example.global

func Server() {
	// import net/http
	// import github.com/dotchain/dot/ops/nw
	// import github.com/dotchain/dot/ops/bolt

        // uses a local-file backed bolt DB backend
	store, _ := bolt.New("file.bolt", "instance", nil)
	store = nw.MemPoller(store)
        defer store.Close()
	http.Handle("/api/", &nw.Handler{Store: store})
        http.ListenAndServe(":8080", nil)
}
```

The example above uses the
[Bolt](http://godoc.org/github.com/dotchain/dot/ops/bolt)
for the actual storage of the operations.  There is also a
[Postgres](http://godoc.org/github.com/dotchain/dot/ops/pg) backend
available.

Note that the server above has no real reference to any application
logic: it simply accepts operations and writes them out in a
guaranteed order broadcasting these to all the clients.

### Types

A TODO MVC app consists of only two core types: `Todo` and `TodoList`:

```go example.global

// Todo tracks a single todo item
type Todo struct {
	Complete bool
        Description string
}

// TodoList tracks a collection of todo items
type TodoList []Todo

```

### Type registration

To use the types across the network, they have to be registered with
the codec (which will be `gob` in this example)

```go example.global
// import encoding/gob

func init() {
	gob.Register(Todo{})
        gob.Register(TodoList{})
}
```

### Code generation

For use with **DOT**, these types need to be augmented with standard
methods of the [Value](https://godoc.org/github.com/dotchain/dot/changes#Value)
interface (or in the case of lists like `TodoList`, also implement the
[Collection](https://godoc.org/github.com/dotchain/dot/changes#Collection)
interface).

These interfaces are essentially the ability to take changes of the
form **replace a sub field** or **replace items in the array** and
calculate the result of applying them.  They are mostly boilerplate
and so can be autogenerated easily via the
[dotc](https://godoc.org/github.com/dotchain/dot/x/dotc) package. See
[code generation](codegen.md) for augmenting the above type
information.

The code generation not only implements these two interfaces, it also
produces a new **Stream** type for **Todo** and **TodoList**.  A
stream type is like a linked list with the `Value` field being the
underlying value and **Next()** returning the next entry in the stream
(in case the value was modified).  And **Latest** returns the
last entry in the stream at that point.  Also, each stream type
implements mutation methods to easily modify the value associated with
a stream.

What makes the streams interesting is that two different modifications
from the same state cause both **Latest** of both to be the same with
the effect of both *merged*.  (This is done using the magic of
operational transformations)

### Toggling Complete

The code to toggle the `Complete` field of a particular todo item
looks like the following:

```go example.global
func Toggle(t *TodoListStream, index int) {
	// TodoListStream.Item() is generated code. It returns
        // a stream of the n'th element of the slice so that
        // particular stream can be modified. When that stream is
        // modified, the effect is automatically merged into the
        // parent (and available via .Next of the parent stream)
	todoStream := t.Item(index) 

	// TodoStream.Complete is generated code. It returns a stream
        // for the Todo.Complete field so that it can be modified. As
        // with slices above, mutations on the field's stream are
        // reflected on the struct stream (via .Next or .Latest())
        completeStream := todoStream.Complete()

	// completeStream is of type streams.Bool. All streams
        // implement the simple Update(newValue) method that replaces
        // the current value with a new value.
        completeStream.Update(!completeStream.Value)
}
```

Note that the function does not return any value here but the updates
can be fetched by calling `.Latest()` on any of the corresponding
streams. If a single stream instance has multiple edits, the
`Latest()` value is the merged value of all those edits. 

### Changing description

The code for Changing description is similar.  The string
`Description` field in `Todo` maps to a `streams.S16` stream. This
implements an `Update()` method like all streams.

But to make things interesting, lets look at **splicing** rather
than replacing the whole string. Splicing is taking a subsequence of
the string at a particular position and replacing it with the provided
value.  It captures insert,  delete and replace in one operation.

This probably better mimics what text editors do and a benefit of such
high granularity edits is that when two users edit the same text, so
long as they don't directly touch the same characters, the edits will
merge quite cleanly.

```go example.global
func SpliceDescription(t *TodoListStream, index, offset, count int, replacement string) {
	// TodoListStream.Item() is generated code. It returns
        // a stream of the n'th element of the slice so that
        // particular stream can be modified. When that stream is
        // modified, the effect is automatically merged into the
        // parent (and available via .Next of the parent stream)
	todoStream := t.Item(index) 

	// TodoStream.Description is generated code. It returns a
        // stream for the Todo.Description field so that it can be
        // modified. As with slices above, mutations on the field's
        // stream are reflected on the struct stream (via .Next or
        // .Latest()) 
	// TodoStream.Description() returns streams.S16 type
        descStream := todoStream.Description()

	// streams.S16 implements Splice(offset, removeCount, replacement)
        descStream.Splice(offset, count, replacement)
}
```

### Adding Todos

Adding a Todo is relatively simple as well:

```go example.global
func AddTodo(t *TodoListStream, todo Todo) {
	// All slice streams implement Splice(offset, removeCount, replacement)
	t.Splice(len(t.Value), 0, todo)
}
```

The use of `Splice` in this example should hint that (just like
strings) slices support insertion/deletion at arbitrary points within
via the Splice method. In addition to supporting this, streams also
support the `Move(offset, count, distance)` method to move some items
around within the slice

### Client connection

Setting up the client requires connecting to the URL where the server
is hosted.  In addition, the code below illustrations how sessions
could be saved and restarted if needed.

```go example.global
// import github.com/dotchain/dot/ops/nw
// import github.com/dotchain/dot/ops
// import math/rand

func Client(stop chan struct{}, render func(*TodoListStream)) {
	version, pending, todos := SavedSession()

	store := &nw.Client{URL: "http://localhost:8080/api/"}
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
seconds.

## How it all works

There are values, changes and streams.

1. **Values** implement the
[Value](https://godoc.org/github.com/dotchain/dot/changes#Value)
interface. If the value represents a collection, it also implements
the
[Collection](https://godoc.org/github.com/dotchain/dot/changes#Collection)
interface.
2. **Changes** represent mutations to values that can be *merged*. If
two independent changes are made to the same value, they can be merged
so that the `A + merged(B) = B + merged(A)`.  This is represented by
the [Change](https://godoc.org/github.com/dotchain/dot/changes#Change)
interface. The
[changes](https://godoc.org/github.com/dotchain/dot/changes) package
implements the core changes with composition that allow richer changes
to be implemented.
3. **Streams** represent a sequence of changes to a value, except it
is **convergent** -- if multiple writers modify a value, they each get
a separate stream instance that only reflects their local change but
following the *Next* chain will guarantee that both end up with the
same value.

### Applying changes

The following example illustrates how to edit a string with values and
changes

```golang dot_test.Example_applying_changes
	// import fmt
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types

	// S8 is DOT-compatible string type with UTF8 string indices
	initial := types.S8("hello")

        append := changes.Splice{
        	Offset: len("hello"), // end of "hello"
                Before: types.S8(""), // nothing to remove
                After: types.S8(" world"), // insert " world"
        }

        // apply the change
        updated := initial.Apply(nil, append)

	fmt.Println(updated)
        // Output: hello world
```

### Applying changes with streams

A less verbose *stream* based version (preferred) would look like so:

```golang dot_test.Example_apply_stream
	// import fmt
        // import github.com/dotchain/dot/streams

        initial := &streams.S8{Stream: streams.New(), Value: "hello"}
        updated := initial.Splice(5, 0, " world")

	fmt.Println(updated.Value)
        // Output: hello world
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

```golang dot_test.Example_changeset_composition
	// import fmt
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types

	initial := types.S8("hello")

        // append " world" => "hello world"
        append1 := changes.Splice{
        	Offset: len("hello"),
                Before: types.S8(""),
                After: types.S8(" world"),
        }

        // append "." => "hello world."
        append2 := changes.Splice{
        	Offset: len("hello world"),
                Before: types.S8(""),
                After: types.S8("."),
        }
        
        // now combine the two appends and apply
        both := changes.ChangeSet{append1, append2}
        updated := initial.Apply(nil, both)
        fmt.Println(updated)

	// Output: hello world.
```

Another form of composition is modifying a sub-element such as an
array element or a dictionary path:

```golang dot_test.Example_path_composition
	// import fmt
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types

        // types.A is a generic array type and types.M is a map type
        initial := types.A{types.M{"hello": types.S8("world")}}

        // replace "world" with "world!"
        replace := changes.Replace{Before: types.S8("world"), After: types.S8("world!")}

        // replace "world" with "world!" of initial[0]["hello"]
        path := []interface{}{0, "hello"}
        c := changes.PathChange{Path: path, Change: replace}
        updated := initial.Apply(nil, c)
        fmt.Println(updated)

	// Output: [map[hello:world!]]        
```

### Convergence

The core property of all changes is the ability to guarantee
*convergence* when two mutations are attempted on the same state:

```golang dot_test.Example_convergence
	// import fmt
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types

	initial := types.S8("hello")

	// two changes: append " world" and delete "lo"
	insert := changes.Splice{Offset: 5, Before: types.S8(""), After: types.S8(" world")}
	remove := changes.Splice{Offset: 3, Before: types.S8("lo"), After: types.S8("")}

	// two versions derived from initial
        inserted := initial.Apply(nil, insert)
        removed := initial.Apply(nil, remove)

        // merge the changes
        removex, insertx := insert.Merge(remove)

        // converge by applying the above
        final1 := inserted.Apply(nil, removex)
        final2 := removed.Apply(nil, insertx)

        fmt.Println(final1, final1 == final2)
        // Output: hel world true
```

### Convergence using streams

The same convergence example is a lot easier to read with streams:

```golang dot_test.Example_convergence_streams
	// import fmt
        // import github.com/dotchain/dot/streams

	initial := streams.S8{Stream:  streams.New(), Value: "hello"}

	// two changes: append " world" and delete "lo"
        s1 := initial.Splice(5, 0, " world")
	s2 := initial.Splice(3, len("lo"), "")

	// streams automatically merge because they are both
        // based on initial
        s1 = s1.Latest()
        s2 = s2.Latest()

        fmt.Println(s1.Value, s1.Value == s2.Value)
        // Output: hel world true
```

The ability to *merge* two independent changes done to the same
initial state is the basis for the eventual convergence of the data
structures.  The
[changes](http://godoc.org/github.com/dotchain/dot/changes) package 
has fairly intensive tests to cover the change types defined there,
both individually and in composition. 

### Revert

In addition to convergence, the set of change types are chosen
carefully to make it easy to implement *Revert()* (undo of the
change). This allows the ability to build a generic
[undo stack](https://godoc.org/github.com/dotchain/dot/streams/undo) as well
as somewhat fancy features like
[folding](https://godoc.org/github.com/dotchain/dot/x/fold).

    Example TODO

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

### Branching of streams

Streams can also be branched *a la* Git. Changes made in branches do not
affect the master or vice-versa -- until one of Pull or Push are
called.

```golang dot_test.Example_branching
	// import fmt
        // import github.com/dotchain/dot/streams
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types        
        
        // here a generic change stream is created
        master := streams.New()
        local := streams.Branch(master)

        // changes will not be reflected on master yet
        c := changes.Replace{Before: changes.Nil, After: types.S8("hello")}
        local = local.Append(c)

	if x, c1 := master.Next(); x != nil || c1 != nil {
        	fmt.Println("Master unexepectedly changed")
        }

	// push local changes up to master now
        streams.Push(local)
	if x, c2 := master.Next(); x == nil || c2 != c {
        	fmt.Println("Master changed but unexpectedly", x, c2)
        }

	// Output:
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

```golang dot_test.skip

// github.com/dotchain/dot/ops/nw"
// github.com/dotchain/dot/ops"
// github.com/dotchain/dot/streams"
// github.com/dotchain/dot/x/idgen"       

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

```golang dot_test.skip
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

