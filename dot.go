// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dot is a container for operational transformations
//
// Please see https://github.com/dotchain/dot for a tutorial on
// how to use DOT.
//
// The core functionality is spread out between dot/changes,
// dot/refs, dot/streams and dot/x but this package exposes simple
// client and server implementations:
//
// Server example
//
//      import "encoding/gob"
//      import "net/http"
//      import "github.com/dotchain/dot"
//      ...
//      gob.Register(..) // register any non-standard OT types used
//      http.Handle("/dot/", dot.BoltServer("file.bolt", "instance", nil))
//      http.ListenAndServer(":8080", nil)
//
//
// Client example
//
//      import "encoding/gob"
//      import "net/http"
//      import "github.com/dotchain/dot"
//      ...
//      gob.Register(..) // register any non-standard OT types used
//      clientStream := dot.Client(-1, nil) // start from scratch
//
//
// Immutable values
//
// DOT uses immutable values. Every Value must implement the
// change.Value interface which is a single Apply method that returns
// the result of applying a mutation (while leaving the original value
// effectively unchanged).
//
// If the underlying type behaves like a collection (such as with
// Slices), the type must also implement some collection specific
// methods specified in the changes.Collection interface.
//
// Most actual types are likely to be structs or slices with
// boilerplate implementaations of the interfaces. The x/dotc package
// has a code generator which can emit such boilerplate
// implementations simplifying this task.
//
// Changes
//
// The changes package implements a set of simple changes (Replace,
// Splice and Move). Richer changes are expected to be built up by
// composition via changes.ChangeSet (which is a sequence of changes)
// and changes.PathChange (which modifies a value at a path).
//
// Changes are immutable too and generally are meant to not maintain
// any reference to the value they apply on.  While custom changes are
// possible (they have to implement the changes.Custom interface),
// they are expected to be rare as the default set of chnange types cover
// a vast variety of scenarios.
//
// The core logic of DOT is in the Merge methods of changes: they
// guaranteee that if two independent changes are done to a value, the
// deviation in the values can be converged.  The basic property of
// any two changes (on the same value) is that:
//
//      leftx, rightx := left.Merge(right)
//      initial.Apply(nil, left).Apply(nil, leftx) ==
//      initial.Apply(nil, right).Apply(nil, rightx)
//
// Care must be taken with custom changes to ensure that this property
// is preserved.
//
// Streams
//
// Streams represent the sequence of changes associated with a single
// value. Stream instances behave like they are immutable: when a
// change happens, a new stream instance captures the change.  Streams
// also support multiple-writers: it is possible for two independent
// changes to the same stream instance. In this case, the
// newly-created  stream instances only capture the respective
// changes but these both have a "Next" value that converges to the
// same value.  That is, the two separate streams implicitly have the
// changes from each other (but after transforming through the Merge)
// method.
//
// This allows streams to perform quite nicely as convergent data
// structures without much syntax overhead:
//
//    initial := streams.S8{Stream:  streams.New(), Value: "hello"}
//
//    // two changes: append " world" and delete "lo"
//    s1 := initial.Splice(5, 0, " world")
//    s2 := initial.Splice(3, len("lo"), "")
//
//    // streams automatically merge because they are both
//    // based on initial
//    s1 = s1.Latest()
//    s2 = s2.Latest()
//
//    fmt.Println(s1.Value, s1.Value == s2.Value)
//    // Output: hel world true
//
// Strongly typed streams
//
// The streams package provides a generic Stream implementation (via
// the New function) which implements the idea of a sequence of
// convergent changes. But much of the power of streams is in having
// strongly type streams where the stream is associated with a
// strongly typed value.  The streams package provides simple text
// streamss (S8 and S16) as well as Bool and Counter types.  Richer
// types like structs and slices can be converted to their stream
// equivalent rather mechanically and  this is done by the x/dotc
// package -- using code generation.
//
//    Some day, Golang would support generics and then the code
//    generation ugliness of x/dotc will no longer be needed.
//
// Substreams are streams that refer into a particular field of a
// parent stream.   For example, if the parent value is a struct with
// a "Done" field, it is  possible to treat the "Done stream" as the
// changes scoped this field. This allows code to be written much more
// cleanly.   See the https://github.com/dotchain/dot#toggling-complete
// section of the documentation for an example.
//
// Other features
//
// Streams support branching (a la Git) and folding.  See the examples!
//
// Streams also support references. A typical use case is maintaining
// the user cursor within a region of text.  When remote changes
// happen to the text, the cursor needs to be updated.  In fact, when
// one takes a substream of an element of an array, the array index
// needs to be automatically  managed (i.e. insertions into the array
// before the index should automatically update the index etc).  This
// is managed within streams using references.
//
//
// Server implementations
//
// A particular value can be reconstituted from the sequence of
// changes to that value. In DOT, only these changes are stored and
// that too in an append-only log.  This make the backend rather
// simple and generally agnostic of application types to a large
// extent.
//
// See https://github.com/dotchain/dot#server for example code.
package dot

//go:generate go get github.com/tvastar/test/cmd/testmd
//go:generate go get github.com/tvastar/toc
//go:generate testmd -pkg dot_test -o dot_test.go README.md
//go:generate toc -h Contents -o README.md README.md
//go:generate testmd -pkg example -o example/todo.go README.md
//go:generate testmd -pkg main codegen.md
