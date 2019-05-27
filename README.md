# DOT

[![Status](https://travis-ci.com/dotchain/dot.svg?branch=master)](https://travis-ci.com/dotchain/dot?branch=master)
[![GoDoc](https://godoc.org/github.com/dotchain/dot?status.svg)](https://godoc.org/github.com/dotchain/dot)
[![codecov](https://codecov.io/gh/dotchain/dot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotchain/dot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotchain/dot)](https://goreportcard.com/report/github.com/dotchain/dot)

The DOT project is a blend of [operational
transformation](https://en.wikipedia.org/wiki/Operational_transformation),
[persistent/immutable
datastructures](https://en.wikipedia.org/wiki/Persistent_data_structure)
and [reactive](https://en.wikipedia.org/wiki/Reactive_programming)
stream processing.

This is an implementation of distributed data synchronization of rich
custom data structures with conflict-free merging.

## Status

This is very close to v1 release.  The [ES6](https://github.com/dotchain/dotjs) version
interoperates well right now but outstanding short-term issues have
more to do with consistency of the API surface than features:

* ~The ES6 version has a simpler polling-based Network API that seems worth adopting here.~  ** Adopted **
* ~The ES6 branch/undo integration also feels a lot simpler.~ ** Adopted **
* The ES6 version prefers `replace()` instead of `update()`.
* Nullable value types (i.e typed Nil values vs change.Nil vs nil) seems confusing.

## Features

1. Small, well tested mutations and immutable persistent values
2. Support for rich user-defined types, not just collaborative text
3. Streams and **Git-like** branching, merging support
4. Simple network support (Gob serialization) and storage support
5. Strong references support that are automatically updated with changes
6. Rich builtin undo support for any type and mutation
7. Folding (committed changes on top of uncommitted changes)

An interoperable ES6 version is available on [dotchain/dotjs](https://github.com/dotchain/dotjs) with a TODO MVC demo of it [here](https://github.com/dotchain/demos)


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
    10. [In browser demo](#in-browser-demo)
4. [How it all works](#how-it-all-works)
    1. [Applying changes](#applying-changes)
    2. [Applying changes with streams](#applying-changes-with-streams)
    3. [Composition of changes](#composition-of-changes)
    4. [Convergence](#convergence)
    5. [Convergence using streams](#convergence-using-streams)
    6. [Revert and undo](#revert-and-undo)
    7. [Folding](#folding)
    8. [Branching of streams](#branching-of-streams)
    9. [References](#references)
    10. [Network synchronization and server](#network-synchronization-and-server)
5. [Broad Issues](#broad-issues)
6. [Contributing](#contributing)

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
	// import github.com/dotchain/dot

        // uses a local-file backed bolt DB backend
	http.Handle("/dot/", dot.BoltServer("file.bolt"))
        http.ListenAndServe(":8080", nil)
}
```

The example above uses the
[Bolt](http://godoc.org/github.com/dotchain/dot/ops/bolt)
backend for the actual storage of the operations.  There is also a
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
// import github.com/dotchain/dot/ops.nw

func init() {
	nw.Register(Todo{})
        nw.Register(TodoList{})
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
from the same state cause the **Latest** of both to be the same with
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
high granularity edits is that when two users edit the same text, the
edits will merge quite cleanly so
long as they don't directly touch the same characters.

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
// import time
// import sync
// import github.com/dotchain/dot

var Lock sync.Mutex
func Client(stop chan struct{}, render func(*TodoListStream)) {
	url := "http://localhost:8080/dot/"
        session, todos := SavedSession()
	s, store := session.NonBlockingStream(url, nil)
        defer store.Close()

	todosStream := &TodoListStream{Stream: s, Value: todos}

        ticker := time.NewTicker(500*time.Millisecond)
        changed := true
	for {
        	if changed {
			render(todosStream)
                }
        	select {
                case <- stop:
                	return
                case <- ticker.C:
                }

                Lock.Lock()
		s.Push()
                s.Pull()
                next := todosStream.Latest()
                changed = next != todosStream
                todosStream, s = next, next.Stream
                Lock.Unlock()
        }

       	SaveSession(session, todosStream.Value)
}


func SaveSession(s *dot.Session, todos TodoList) {
	// this is not yet implemented. if it were, then
        // this value should be persisted locally and returned
        // by the call to savedSession
}

func SavedSession() (s *dot.Session, todos TodoList) {
	// this is not yet implemented. return default values
        return dot.NewSession(), nil
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

### In browser demo

The [fuss](https://github.com/dotchain/fuss) project has demos of a
TODO-MVC app built on top of this framework using
[gopherjs](https://github.com/gopherjs/gopherjs).  In particular, the
[collab](https://github.com/dotchain/fuss/tree/master/todo/collab)
folder illustrates how simple the code is to make something work
collaboratively (the rest of the code base is not even aware of
whether things are collaborative).

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

```golang dot_test.Example_applyingChanges
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

```golang dot_test.Example_applyingChangesUsingStreams
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

```golang dot_test.Example_changesetComposition
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

```golang dot_test.Example_pathComposition
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

```golang dot_test.Example_convergenceUsingStreams
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

### Revert and undo

All the predefined types of changes in DOT (see
[changes](https://godoc.org/github.com/dotchain/dot/changes)) are
carefully designed so that every change can be inverted easily without
reference to the underlying value.  For example,
[changes.Replace](https://godoc.org/github.com/dotchain/dot/changes#Replace)
has both the **Before** and **After** fields instead of just keeping
the **After**.  This allows the reverse to be computed quite easily by
swapping the two fields.  This does generally incur additional storage
expenses but the tradeoff is that code gets much simpler  to work
with.

In particular, it is possible to build generic
[undo](https://godoc.org/github.com/dotchain/dot/streams/undo) support
quite easily and naturally.  The following example shows both **Undo**
and **Redo** being invoked from an undo stack.

```go dot_test.Example_undoStreams
	// import fmt
        // import github.com/dotchain/dot/streams
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types
        // import github.com/dotchain/dot/streams/undo

	// create master, undoable child and the undo stack itself
	master := &streams.S16{Stream: streams.New(), Value: "hello"}
        s := undo.New(master.Stream)
        undoableChild := &streams.S16{Stream: s, Value: master.Value}

	// change hello => Hello
	undoableChild = undoableChild.Splice(0, len("h"), "H")
	fmt.Println(undoableChild.Value)

	// for kicks, update master hello => hello$ as if it came
        // from the server
        master.Splice(len("hello"), 0, "$")

	// now undo this via the stack
        s.Undo()

	// now undoableChild should be hello$
        undoableChild = undoableChild.Latest()
        fmt.Println(undoableChild.Value)

	// now redo the last operation to get Hello$
        s.Redo()
        undoableChild = undoableChild.Latest()
        fmt.Println(undoableChild.Value)
        
	// Output:
        // Hello
        // hello$
        // Hello$
```

### Folding

In the case of editors, folding refers to a piece of text that has
been hidden away. The difficulty with implementing this in a
collaborative setting is that as external edits come in, the fold has
to be maintained.

The design of DOT allows for an elegant way to achieve this: consider
the "folding" as a local change (replacing the folded region with
say "..."). This local change is never meant to be sent out.  All
changes to the unfolded and folded versions can be proxied quite
nicely without much app involvement:

```go dot_test.Example_folding
	// import fmt
        // import github.com/dotchain/dot/streams
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types
        // import github.com/dotchain/dot/x/fold

	// create master, folded child and the folding itself
	master := &streams.S16{Stream: streams.New(), Value: "hello world!"}
        foldChange := changes.Splice{
        	Offset: len("hello"),
                Before: types.S16(" world"),
                After: types.S16("..."),
        }
        foldedStream := fold.New(foldChange, master.Stream)
        folded := &streams.S16{Stream: foldedStream, Value :"hello...!"}

        // folded:  hello...! => Hello...!!!
	folded = folded.Splice(0, len("h"), "H")
        folded = folded.Splice(len("Hello...!"), 0, "!!")
        fmt.Println(folded.Value)

	// master: hello world => hullo world
	master = master.Splice(len("h"), len("e"), "u")
        fmt.Println(master.Value)

        // now folded = Hullo...!!!
        fmt.Println(folded.Latest().Value)

        // master = Hullo world!!!
        fmt.Println(master.Latest().Value)

	// Output:
        // Hello...!!!
        // hullo world!
        // Hullo...!!!
        // Hullo world!!!
```

### Branching of streams

Streams in DOT can also be branched *a la* Git. Changes made in
branches do not affect the master or vice-versa -- until one of Pull
or Push are called.

```golang dot_test.Example_branching
	// import fmt
        // import github.com/dotchain/dot/streams
        // import github.com/dotchain/dot/changes
        // import github.com/dotchain/dot/changes/types        
        
        // local is a branch of master
        master := &streams.S16{Stream: streams.New(), Value: "hello"}
        local := &streams.S16{Stream: streams.Branch(master.Stream), Value: master.Value}

	// edit locally: hello => hallo
	local.Splice(len("h"), len("e"), "a")

	// changes will not be reflected on master yet
        fmt.Println(master.Latest().Value)

	// push local changes up to master now
        local.Stream.Push()

	// now master = hallo
	fmt.Println(master.Latest().Value)

        // Output:
        // hello
        // hallo
```

There are other neat benefits to the branching model: it provides a
fine grained control for pulling changes from the network on demand
and suspending it as well as providing a way for making local
changes.

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


### Network synchronization and server

DOT uses a fairly simple backend
[Store](https://godoc.org/github.com/dotchain/dot/ops#Store)
interface: an append-only dumb log. The
[Bolt](https://godoc.org/github.com/dotchain/dot/ops/bolt) and
[Postgres](https://godoc.org/github.com/dotchain/dot/ops/pg)
implementations are quite simple and other data backends can be
easily added.

See [Server](#server) and [Client connection](#client-connection) for
sample server and client applications.  Note that the journal approach
used implies that the journal size only increases and so clients will
eventually take a while to rebuild their state from the journal. The
client API allows snapshotting state to make the rebuilds faster.
There is no server support for snapshots though it is possible to
build one rather easily

## Broad Issues

1. changes.Context/changes.Meta are not fully integrated
2. ~gob-encoding makes it harder to deal with other languages but JSON
encodindg wont work with interfaces.~
   * Added `sjson encoding` as a portable (if verbose) format.
   * The [ES6 dotjs](https://github.com/dotchain/dotjs) package uses this as the native format.
3. Cross-object merging and persisted branches need more platform support
   * Snapshots are somewhat related to this as well.
4. Full rich-text support with collaborative cursors still needs work
with references and reference containers.
5. Code generation can infer types from regular go declarations
6. Snapshots and transient states need some sugar.

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

