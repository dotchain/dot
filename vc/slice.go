// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import "github.com/dotchain/dot"

// Slice represents a slice of value of any type
type Slice struct {
	// Version is the version control metadata
	Version

	// Start. End refer to the window into the original slice.
	// If they are nil, they refer to the logical start/end
	Start, End *int

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
	value := l.Value[start:end:end]
	result := Slice{Version: l.Version, Value: value}
	if start != 0 || l.Start != nil {
		result.Start = l.add(l.Start, start)
	}

	if end != len(l.Value) || l.End != nil {
		result.End = l.add(l.Start, end)
	}

	return result
}

// SpliceSync synchronously splices and returns the new slice.  If
// there were other changes done on the Slice before this operation,
// that will not be reflected in the output but it will be guaranteed
// to be reflected in the next call  to Latest.
func (l Slice) SpliceSync(offset, removeCount int, replacement []interface{}) Slice {
	o := *l.add(l.Start, offset)
	before := l.Value[offset : offset+removeCount]
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	value := append(append(l.Value[:offset:offset], replacement...), l.Value[offset+removeCount:]...)
	version := l.Version.UpdateSync([]dot.Change{c})
	start, end := l.spliceOffsets(offset, removeCount, len(replacement))
	return Slice{Version: version, Value: value, Start: start, End: end}
}

// SpliceAsync asynchronously splices and returns the new slice.  If
// there were other changes done on the Slice before this operation,
// that will not be reflected in the output but will eventually get
// reflected in a call to Latest.
func (l Slice) SpliceAsync(offset, removeCount int, replacement []interface{}) Slice {
	o := *l.add(l.Start, offset)
	before := l.Value[offset : offset+removeCount]
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	value := append(append(l.Value[:offset:offset], replacement...), l.Value[offset+removeCount:]...)
	version := l.Version.UpdateAsync([]dot.Change{c})
	start, end := l.spliceOffsets(offset, removeCount, len(replacement))
	return Slice{Version: version, Value: value, Start: start, End: end}
}

// Latest returns the latest value. The current object may have been
// deleted, in which case it returns the zero value and sets the bool
// to false.
func (l Slice) Latest() (Slice, bool) {
	v, ver, start, end := l.Version.LatestAt(l.Start, l.End)
	if ver == nil {
		return Slice{}, false
	}

	value := v.([]interface{})
	s, e := 0, len(value)
	if start != nil {
		s = *start
	}
	if end != nil {
		e = *end
	}

	return Slice{Value: value[s:e:e], Start: start, End: end, Version: ver}, true
}

func (l Slice) spliceOffsets(offset, removeCount, replaceCount int) (*int, *int) {
	var startp, endp *int

	if l.Start != nil {
		startp = l.Start
	}

	if l.End != nil {
		endp = l.add(l.End, replaceCount-removeCount)
	}
	return startp, endp
}

func (l Slice) add(p *int, v int) *int {
	if p != nil {
		v += *p
	}
	return &v
}
