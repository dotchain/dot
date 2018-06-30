// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import "github.com/dotchain/dot"

// Slice represents a slice of value of any type
type Slice struct {
	// Control is the version control metadata
	Control

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
	result := Slice{Control: l.Control, Value: value}
	if start != 0 || l.Start != nil {
		result.Start = l.add(l.Start, start)
	}

	if end != len(l.Value) || l.End != nil {
		result.End = l.add(l.Start, end)
	}

	return result
}

// Splice synchronously splices and returns the new slice.  If
// there were other changes done on the Slice before this operation,
// that will not be reflected in the output but it will be guaranteed
// to be reflected in the next call  to Latest.
func (l Slice) Splice(offset, removeCount int, replacement []interface{}) Slice {
	o := *l.add(l.Start, offset)
	before := l.Value[offset : offset+removeCount]
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	value := append(append(l.Value[:offset:offset], replacement...), l.Value[offset+removeCount:]...)
	version := l.Control.UpdateSync([]dot.Change{c})
	start, end := l.spliceOffsets(offset, removeCount, len(replacement))
	return Slice{Control: version, Value: value, Start: start, End: end}
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
	version := l.Control.UpdateAsync([]dot.Change{c})
	start, end := l.spliceOffsets(offset, removeCount, len(replacement))
	return Slice{Control: version, Value: value, Start: start, End: end}
}

// Move synchronously moves the specified elements (from offset to
// offset + count) by the specified distance (positive moves to the
// right, negative moves to the left).  It returns a new slice which
// reflects the effect of the move and the effect of the move is
// also guaranteed to be reflected in the next Latest call.
func (l Slice) Move(offset, count, distance int) Slice {
	move := &dot.MoveInfo{Offset: offset, Count: count, Distance: distance}
	changes := []dot.Change{{Move: move}}
	value, _ := unwrap(utils.Apply(l.Value, changes)).([]interface{})
	version := l.Control.UpdateSync(changes)
	s, e := l.Start, l.End
	return Slice{Control: version, Value: value, Start: s, End: e}
}

// MoveAsync asynchronously moves the specified elements (from offset to
// offset + count) by the specified distance (positive moves to the
// right, negative moves to the left).  It returns a new slice which
// reflects the effect of the move but the effect is not guaranteed to
//  be reflected into the next Latest() call (as this executes in a
//  separate go routine). The only guarantee is that caussality is
//  preserved with any operations on the output Slice being applied
//  after the current operation is applied.
func (l Slice) MoveAsync(offset, count, distance int) Slice {
	move := &dot.MoveInfo{Offset: offset, Count: count, Distance: distance}
	changes := []dot.Change{{Move: move}}
	value, _ := unwrap(utils.Apply(l.Value, changes)).([]interface{})
	version := l.Control.UpdateAsync(changes)
	s, e := l.Start, l.End
	return Slice{Control: version, Value: value, Start: s, End: e}
}

// Latest returns the latest value. The current object may have been
// deleted, in which case it returns the zero value and sets the bool
// to false.
func (l Slice) Latest() (Slice, bool) {
	v, ver, start, end := l.Control.LatestAt(l.Start, l.End)
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

	return Slice{Value: value[s:e:e], Start: start, End: end, Control: ver}, true
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
