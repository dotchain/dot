// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package vc implements  versioned datastructures.
//
// A basic type is any immutable type or a simple composition using
// []interface{} or map[string]interface{} where the inner types can
// be any basic type.
//
// It is possible to create a versioned type out of this simply by
// calling New on this:
//
//    ctl := vc.New(...)
//
// The returned control is an immutable structure that keeps track of
// version derivations -- so when updates are made in a structured
// way, it can track them and transform them in a sensible way to
// derive the latest value.
//
// Consider the concrete example of a versioned string:
//
//     initial := "hello"
//     ctl := vc.New(initial)
//     s := vc.String{Control: ctl, Value: initial}
//
// A string created in a fashion as above can be treated as an
// immutable object providing simple mutations -- with the additional
// caveat that all the mutations are considered to belong to the same
// underlying storage. So, the updates are carried over to the
// underlying storage with the updated underlying storage available
// via "Latest()".  So, any consumer of the type can treat it as a
// simple immutable type and the type will be consistent with that
// usage (with all updates only reflecting that specific variation)
// but in addition, fetching the "latest" would effectively include a
// Git-like merge of all applied changes.
//
//      s1 = s.SpliceSync(5, 0, "world") // this will return "hello world"
//      s2 = s.SpliceSync(0, 1, "H") // this will return "Hello"
//      s3 = s.Latest() // this will return "Hello world" merging both
//
// Note that the immutable feel of the API extends to thread safety.
// Updates can happen on separate go routines and will cause no
// problems. The order of consolidation of operations is not tightly
// specified. In fact, all "Splice" operations provide weak ordering
// guarantees -- if there are concurrent splices, it is  possible for
// another item from another location to be inserted at the boundaries
// of where changes happen but it is not possible that a change
// happening concurrently will get inserted within the string being
// inserted.  Splices themselves also won't make an element which
// logically before another get merged in a fashion where it comes
// logical after though other items may get inserted (or
// deleted). When one uses the "Move" operation, this order guarantee
// does not hold ofcourse.
//
// All the operations have a Sync and Async version -- the Sync
// version guarantees that an immediate call to Latest will reflect the
// effect of the Sync operation (and indirectly those of all mutations
// this is dependent on).
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
// Garbage collection
//
// The awkward structure of the API such as separating the control
// methods for updates from the actual value is primarily to reduce
// the garbage collection overhead of maintaining the
// history. Basically, the only strong references the system  holds is
// to the latest value and any explicit values being used by the
// application itself -- with the infrastructure only maintaining
// references to the control structures of any values that are
// referenced by the application (and intermediate changes).  In
// particular, history changes and intermediate changes do not cause
// references to those corresponding intermediate values.
//
// Strongly typing
//
// The ability to support custom types (such as structs or slices to
// arbitrary types) is possible but is a bit convoluted due to the
// lack of generics in Go.  There is some support planned to make this
// task easy.
//
// Fork/Merge
//
// It is possible to implement forking and merging like with Git. This
// is planned for in the future.
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
// This is typically done to maintain the root level state of an
// application.  See individaul type examples for usage.
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
