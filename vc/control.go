// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import (
	"github.com/dotchain/dot"
	"strconv"
	"sync"
)

// Control implements the methods to manage a specific
// version. Control typically does not hold any references to the
// underlying value itself but simply tracks changes and provides
// access to fetch the latest values.
type Control interface {
	UpdateSync(changes []dot.Change) Control
	UpdateAsync(changes []dot.Change) Control
	Latest() (interface{}, Control)
	LatestAt(start, end *int) (interface{}, Control, *int, *int)
	Child(key string) Control
	ChildAt(index int) Control
	Branch(value interface{}) (Branch, Control)
}

// basis is simply to identify changes
type basis struct{}

// control holds state about a particular control of a data-structure
// It does not hold a reference to the data-structure or any event in
// the past to make sure there are no memory leaks
type control struct {
	sync.Mutex
	// the basis pointer tags this version so changes can be based
	// and merged properly
	*basis
	// own is there just to provide space for the basis pointer above
	own basis
	// parent refers the the container this version is a part of
	parent
}

func (c *control) UpdateSync(changes []dot.Change) Control {
	changes = append([]dot.Change(nil), changes...)
	result := &control{parent: c.parent}
	result.basis = &result.own
	result.Lock()
	defer result.Unlock()
	c.Lock()
	defer c.Unlock()
	c.parent.Bubble(c.basis, result.basis, changes)
	return result
}

func (c *control) UpdateAsync(changes []dot.Change) Control {
	changes = append([]dot.Change(nil), changes...)
	result := &control{parent: c.parent}
	result.basis = &result.own
	result.Lock()
	go func() {
		defer result.Unlock()
		c.Lock()
		defer c.Unlock()
		c.parent.Bubble(c.basis, result.basis, changes)
	}()
	return result
}

func (c *control) Child(key string) Control {
	// TODO: use a local map and memoize version per key so that
	// it is properly synchronized?
	return &control{parent: &dictitem{key, c.parent}, basis: c.basis}
}

func (c *control) ChildAt(index int) Control {
	// TODO: use a local map and memoize version per key so that
	// it is properly synchronized?
	key := strconv.Itoa(index)
	item := &arrayitem{key: key, index: index, array: c.parent}
	return &control{parent: item, basis: c.basis}
}

func (c *control) Latest() (interface{}, Control) {
	val, parent, _, b := c.parent.Latest(nil, c.basis)
	if parent == nil {
		return nil, nil
	}
	val = unwrap(val)

	return val, &control{parent: parent, basis: b}
}

func (c *control) LatestAt(startp, endp *int) (interface{}, Control, *int, *int) {
	var nilpath *dot.RefPath
	val, parent, _, b := c.parent.Latest(nil, c.basis)
	if parent == nil {
		return nil, nil, nil, nil
	}
	val = unwrap(val)

	if startp != nil {
		s := &dot.RefIndex{Index: *startp, Type: dot.RefIndexStart}
		key := c.parent.MapPath(nilpath.Append("", s), c.basis, b)[0]
		start := dot.NewRefIndex(key).Index
		startp = &start
	}

	if endp != nil {
		e := &dot.RefIndex{Index: *endp, Type: dot.RefIndexEnd}
		key := c.parent.MapPath(nilpath.Append("", e), c.basis, b)[0]
		end := dot.NewRefIndex(key).Index
		endp = &end
	}

	return val, &control{parent: parent, basis: b}, startp, endp
}

func (c *control) Branch(value interface{}) (Branch, Control) {
	r := &root{v: value}
	result := &control{parent: r, basis: &r.own}
	return &branch{parent: c, root: r}, result
}
