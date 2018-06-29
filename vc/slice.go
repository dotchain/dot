// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import (
	"github.com/dotchain/dot"
)

// Slice represents a slice of value of any type
type Slice struct {
	// Version is the version control metadata
	Version

	// Offset refers to the actual offset in the backing
	// array. This is zero unless a new list is created from
	// another via a Slice.Slice call which maps the offsets as
	// needed.
	Offset int

	// Value is the actual underlying value of the slice
	Value []interface{}
}

// Slice creates a new slice from the current list.  It does not
// mutate the underlying value, just creates  a new value with the
// provided window.   The start/end refer to boundaries within the
// current slice rather than the underlying window, so it is not
// possible to increase the window slice with a Slice operation.
//
// The version associated with the Slice call does not change but
// mutations from the parent slice will not be reflected  here.
func (l Slice) Slice(start, end int) Slice {
	s, value := l.Offset+start, l.Value[start:end:end]
	return Slice{Version: l.Version, Offset: s, Value: value}
}

// SpliceSync synchronously splices and returns the new slice.  If
// there were other changes done on the Slice before this operation,
// that will not be reflected in the output but it will be guaranteed
// to be reflected in the next call  to Latest.
func (l Slice) SpliceSync(offset, removeCount int, replacement []interface{}) Slice {
	o := l.Offset + offset
	before := l.Value[offset : offset+removeCount]
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	value := append(append(l.Value[:offset:offset], replacement...), l.Value[offset+removeCount:]...)
	version := l.Version.UpdateSync([]dot.Change{c})
	return Slice{Version: version, Value: value, Offset: l.Offset}
}

// SpliceAsync asynchronously splices and returns the new slice.  If
// there were other changes done on the Slice before this operation,
// that will not be reflected in the output but will eventually get
// reflected in a call to Latest.
func (l Slice) SpliceAsync(offset, removeCount int, replacement []interface{}) Slice {
	o := l.Offset + offset
	before := l.Value[offset : offset+removeCount]
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	value := append(append(l.Value[:offset:offset], replacement...), l.Value[offset+removeCount:]...)
	version := l.Version.UpdateAsync([]dot.Change{c})
	return Slice{Version: version, Value: value, Offset: l.Offset}
}

// Latest returns the latest value. The current object may have been
// deleted, in which case it returns the zero value and sets the bool
// to false.
func (l Slice) Latest() (Slice, bool) {
	v, ver, s, e := l.Version.LatestAt(l.Offset, l.Offset+len(l.Value))
	if ver == nil {
		return Slice{}, false
	}

	value := (v.([]interface{}))[s:e]
	return Slice{Value: value, Offset: s, Version: ver}, true
}
