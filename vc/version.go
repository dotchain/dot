// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package vc implements  versioned datastructures
package vc

import (
	"github.com/dotchain/dot"
	"strconv"
	"sync"
)

// Version implements the methods to manage a specific version
type Version interface {
	UpdateSync(changes []dot.Change) Version
	UpdateAsync(changes []dot.Change) Version
	Latest() (interface{}, Version)
	LatestAt(start, end int) (interface{}, Version, int, int)
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
	// parent refers the the container this version is a part of
	parent
}

func (v *version) UpdateSync(changes []dot.Change) Version {
	changes = append([]dot.Change(nil), changes...)
	result := &version{parent: v.parent, basis: &basis{}}
	result.Lock()
	defer result.Unlock()
	v.Lock()
	defer v.Unlock()
	v.parent.Bubble(v.basis, result.basis, changes)
	return result
}

func (v *version) UpdateAsync(changes []dot.Change) Version {
	changes = append([]dot.Change(nil), changes...)
	result := &version{parent: v.parent, basis: &basis{}}
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

	return val, &version{parent: parent, basis: b}
}

func (v *version) LatestAt(start, end int) (interface{}, Version, int, int) {
	panic("Not yet implemented")
	// TODO: map start/end up to the correct basis
}

// New returns a new version
func New(initial interface{}) Version {
	r := &root{v: initial}
	return &version{parent: r, basis: &r.own}
}

var x = dot.Transformer{}
var utils = dot.Utils(x)
