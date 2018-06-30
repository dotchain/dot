// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package vc implements  versioned datastructures.
//
// Versioned datastructures are similar to GIT-files -- changes can be
// done in a multiple-writer fashion with automatic merging.
//
// The basic idea is to consider all versions of a value as consiting
// of the actual raw value and a control structure which keeps track
// of changes.  So, if multiple changes are done in parallel on top of
// the same version, the control structure can identify as having
// originated from the same version (and so the effects would have to
// be merged in the "master" version).
//
// In this package, the control data struture is represented by  the
// Control interface.  In practical terms, this provides a very low
// level interface for book-keeping.  Instead, most callers will use
// the List/Map or String implementations which keep track of both the
// value as well as the control associated with a version. These
// concrete types also provide more natural mechanisms for modifying.
//
// Example
//
// It is possible to create a versioned type out of a value simply by
// calling New on that value:
//
//    ctl := vc.New(initialValue)
//
// The returned control is an immutable structure that keeps track of
// version derivations as all mutations derived from this particular
// value  are expected to pass through this ctl instance.
//
// Consider the concrete example of a versioned string:
//
//  initial := "hello"
//  ctl := vc.New(initial)
//  s := vc.String{Control: ctl, Value: initial}
//
// A string created like so can be treated as an immutable object
// with the basic operation of Splice() to modify the
// string. The interesting effect is what happens when the initial
// string is Spliced two different times. While each of the return
// values will reflect the individual splice operations, the two
// operations are also merged together and the merge value can be
// obtained at any time using the Latest() call (on any version):
//
//  s1 = s.Splice(5, 0, " world") // this will return "hello world"
//  s2 = s.Splice(0, 1, "H") // this will return "Hello"
//  s3 = s.Latest() // this will return "Hello world" merging both
//
// Without the call to Latest(), the string type acts like a regular
// immutable string in all respects with each operation have no
// visible effect on the others.
//
// Branching and merging
//
// The default behavior of mutations is to have them show up on the
// Latest immediately.  It is possible to create Branches where the
// default behavior is to act like a git-branch -- all changes made on
// the branch are reflected on the branch but not propagated to the
// parent.  Creating a branch allows the caller to control when the
// branch can be pushed up to the main line (if at all):
//
//   b, s1 := s.Branch()
//   s1.Splice(5, 0, " world")
//   // s1.Latest() == "hello world" but s.Latest() == "hello" still
//   b.Push()
//   // now s.Latest() is also "hello world"
//
// Thread safety and concurrency
//
// All the methods are threadsafe. There is limited locking at this
// point though some of it can be removed. When multiple concurrent
// changes are made or multiple changes are made on the same version,
// there are limited guarantees made: that the merge process will not
// break logical constraints (so if one splices  "hello", the merge
// process may move where the insert happened but not have other
// parallel changes be inserted within "hello" or change things in
// such a way that characters that were before the splice point appear
// after hello etc).  In addition, the non-Async methods guarantee
// that the effect of the method will get reflected in an immediate
// call to Latest whereas even that guarantee is not provided by the
// Async variations.  In all cases, basic causality is maintained --
// if a version is derived from another, the parent change is applied
// before the child.
//
// Composition
//
// It is possible to use composition of types and create richer
// types:
//
//      value := []interface{}{map[string]interface{}{"x": 5}}
//      ctl := vc.New(value)
//      collection := vc.Slice{Control: ctl, Value: value}
//      mapCtl := collection.ChildAt(0) // get control for first elt
//      map := vc.Map{Control: mapCtl, Value: value[0]}
//
// When an object is composed, mutations on the inner elements are
// propagated and the outer elements are modified but only reflected
// if using the "Latest"
//
//      map.SetKey("x", 5)
//      // collection.Value["x"] is not  modified
//      latest, _ := collection.Latest()
//      // latest.Value["x"] is now 5
//
// Note that the consumer of the "inner" elements have no way of
// knowing their outer containers or if they even exist.  Calling
// "Latest()" on the inner elements will be narrowed down to their
// view only.
//
// Caveat: It is possible that a mutation higher up in the container
// hierarchy can completely remove an element out. This will cause the
// Latest on a stale inner version to fail.  Hence the boolean second
// return value from Latest().
//
// Separation of value and control
//
// There is a fair degree of awkwardness because of the separation of
// the value from the control for making changes on the value.  This
// separation provides for a shared implementation of the control
// behavior irrespective of the "strong type" of the structure
// allowing  for custom types to be implemented.  This is still quite
// awkward to do in Golang due to the lack of generics but there is a
// plan in the works to use code generation and reflection to make
// this work better.  At this point, the Control object should be used
// to get the control of a child so that mutations on the children
// can be propagated to the parent.
//
//
// Garbage collection
//
// The structure of the codebase has been deliberate to avoid leaking
// memory. In particular, if one makes a sequence of changes, all the
// intermediate values are not maintained. In particular, if one calls
// Latest() and ensures no reference exists to prior versions, the
// overhead induced by the book-keeping is minimal and fixed to the
// number of objects used.
//
// But the overhead is not trivial.  Benchmarks are in the works but
// it is expected that the memory and CPU overhead, while being
// acceptable, could be improved significantly.  In particular, for
// large deeply nested structures, there is a fair amount of overhead
// in both recalculating and in the basic collection implementations
// that can be optimized (without actually changing the interfaces
// presented by the package).
//
// Issues
//
// The merging and transformation uses OT which guarantees
// conflict-free convergence but if there are application level data
// constraints that are not captured by the datastructure itself,
// concurrent edits can lead to voilations of such.  It is currently
// not posssible to define custom changes
//
// There are also a few gotchas with how Slice.Slice windows and in
// general how path invalidation happens (when a mutation higher up in
// the container hierarchy either deletes the current value or edits
// it in ways that make its "slice window" change)
//
// References
//
// There is no explicit support for references yet but without built
// in support for references,  it is extremely difficult to implement
// them outside as they need to be transformed by the operations.  For
// example, if the overall structure is not a tree but a DAG with a
// pointer to another element, the pointer needs to be represented as
// a path that gets evaluated dynamically.  This requires coordination
// with the version control mechanism.
//
// Distributing changes
//
// There is no ability to detect local changes, no notifications for
// such etc.  It is relatively easy to add this and when added, one
// can implement a OT client without any caller being aware of the
// existence of this infrastructure at all.
package vc

import (
	"encoding/json"
	"github.com/dotchain/dot"
	"github.com/dotchain/dot/encoding"
)

// New returns a new control for the provided initial value. This
// maintains a reference to the latest value (and any mutations that
// happened though the intermediate values themselves are not held).
//
// New is the equivalent of creating a new repository. The call itself
// isn't expensive but it does cache the latest value for the
// repository (exposed via Latest()).  The returned Control interface
// is not expected to be directly used but instead wrapped with Slice
// or Map (based on the structure of the provided initial value)
func New(initial interface{}) Control {
	r := &root{v: initial}
	return &control{parent: r, basis: &r.own}
}

// TODO: move this unwrap crap into Encoding
func unwrap(i interface{}) interface{} {
	if i == nil {
		return nil
	}

	ue, ok := utils.C.TryGet(i)
	if !ok {
		return i
	}

	if encoding.IsString(ue) {
		var result string
		b, _ := json.Marshal(ue)
		_ = json.Unmarshal(b, &result)
		return result
	}

	if ue.IsArray() {
		result := make([]interface{}, ue.Count())
		ue.ForEach(func(offset int, val interface{}) {
			result[offset] = unwrap(val)
		})
		return result
	}

	result := map[string]interface{}{}
	ue.ForKeys(func(key string, val interface{}) {
		result[key] = unwrap(val)
	})
	return result
}

var x = dot.Transformer{}
var utils = dot.Utils(x)
