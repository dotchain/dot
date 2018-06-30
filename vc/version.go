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
//    version := vc.New(...)
//
// The version metadata is an immutable structure that keeps track of
// version derivations -- so when updates are made in a structured
// way, it can track them and transform them in a sensible way to
// derive the latest value.
//
// Consider the concrete example of a versioned string:
//
//     initial := "hello"
//     ver := vc.New(initial)
//     s := vc.String{Version: ver, Value: initial}
//
// A string created in a fashion as above can be treated as an
// immutable object providing simple mutations -- with the additional
// caveate that all the mutations are considered to belong to the same
// underlying storage. So, the updates are carried over to the
// underlying storate with the updated underlying storage available
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
//      collection := []interface{}{1, 2, 3}
//      ver := vc.New(collection)
//      managedCollection := vc.Slice{Version: ver, Value: collection}
//
//
// Garbage collection
//
// The awkward structure of the API such as separating the control
// methods for updates from the actual value is primarily to reduce
// the garbage collection overhead of old values. Basically, the only
// strong references the system  holds is to the latest  value being
// tracked and any explicit values being held onto by the application.
// Internally, the package holds references to the metadata associated
// with parents and such for any older values but the actual values
// themselves (which could be large) are not maintained.  In
// particular, this approach for versioning guarantees that a long
// chain of sequential changes will not keep leaking memory by
// holding onto the older references
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
	"strconv"
	"sync"
)

// Version implements the methods to manage a specific version
type Version interface {
	UpdateSync(changes []dot.Change) Version
	UpdateAsync(changes []dot.Change) Version
	Latest() (interface{}, Version)
	LatestAt(start, end *int) (interface{}, Version, *int, *int)
	Child(key string) Version
	ChildAt(inde int) Version
}

// basis is simply to identify changes
type basis struct{}

// version holds state about a particular version of a data-structure
// It does not hold a reference to the data-structure or any event in
// the past to make sure there are no memory leaks
type version struct {
	sync.Mutex
	// the basis pointer tags this version so changes can be based
	// and merged properly
	*basis
	// own is there just to provide space for the basis pointer above
	own basis
	// parent refers the the container this version is a part of
	parent
}

func (v *version) UpdateSync(changes []dot.Change) Version {
	changes = append([]dot.Change(nil), changes...)
	result := &version{parent: v.parent}
	result.basis = &result.own
	result.Lock()
	defer result.Unlock()
	v.Lock()
	defer v.Unlock()
	v.parent.Bubble(v.basis, result.basis, changes)
	return result
}

func (v *version) UpdateAsync(changes []dot.Change) Version {
	changes = append([]dot.Change(nil), changes...)
	result := &version{parent: v.parent}
	result.basis = &result.own
	result.Lock()
	go func() {
		defer result.Unlock()
		v.Lock()
		defer v.Unlock()
		v.parent.Bubble(v.basis, result.basis, changes)
	}()
	return result
}

func (v *version) Child(key string) Version {
	// TODO: use a local map and memoize version per key so that
	// it is properly synchronized?
	return &version{parent: &dictitem{key, v.parent}, basis: v.basis}
}

func (v *version) ChildAt(index int) Version {
	// TODO: use a local map and memoize version per key so that
	// it is properly synchronized?
	key := strconv.Itoa(index)
	item := &arrayitem{key: key, index: index, array: v.parent}
	return &version{parent: item, basis: v.basis}
}

func (v *version) Latest() (interface{}, Version) {
	val, parent, _, b := v.parent.Latest(nil, v.basis)
	if parent == nil {
		return nil, nil
	}
	val = unwrap(val)

	return val, &version{parent: parent, basis: b}
}

func (v *version) LatestAt(startp, endp *int) (interface{}, Version, *int, *int) {
	var nilpath *dot.RefPath
	val, parent, _, b := v.parent.Latest(nil, v.basis)
	if parent == nil {
		return nil, nil, nil, nil
	}
	val = unwrap(val)

	if startp != nil {
		s := &dot.RefIndex{Index: *startp, Type: dot.RefIndexStart}
		key := v.parent.MapPath(nilpath.Append("", s), v.basis, b)[0]
		start := dot.NewRefIndex(key).Index
		startp = &start
	}

	if endp != nil {
		e := &dot.RefIndex{Index: *endp, Type: dot.RefIndexEnd}
		key := v.parent.MapPath(nilpath.Append("", e), v.basis, b)[0]
		end := dot.NewRefIndex(key).Index
		endp = &end
	}

	return val, &version{parent: parent, basis: b}, startp, endp
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

// New returns a new version
func New(initial interface{}) Version {
	r := &root{v: initial}
	return &version{parent: r, basis: &r.own}
}

var x = dot.Transformer{}
var utils = dot.Utils(x)
