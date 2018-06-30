// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import (
	"github.com/dotchain/dot"
	"unicode/utf16"
)

// String represents an immutable string slice
type String struct {
	// Control is the version control metadata
	Control

	// Start. End refer to the window into the original slice.
	// If they are nil, they refer to the logical start/end
	Start, End *int

	// Value is the actual underlying value of the slice
	Value string
}

// String creates a new slice from the current string.  It does not
// mutate the underlying value, just creates  a new value with the
// provided window.   The start/end refer to boundaries within the
// string being sliced rather than the underlying storage.  Unlike
// native Go slices it is not possible to increase the window slice
// with a String operation.
//
// The version associated with the String call does not change but
// mutations from the parent slice will not be reflected here.
func (l String) String(start, end int) String {
	value := l.Value[start:end]
	offset := len(utf16.Encode([]rune(l.Value[:start])))
	count := len(utf16.Encode([]rune(value)))

	result := String{Control: l.Control, Value: value}
	if start != 0 || l.Start != nil {
		result.Start = l.add(l.Start, offset)
	}

	if end != len(l.Value) || l.End != nil {
		result.End = l.add(l.Start, offset+count)
	}

	return result
}

// Splice synchronously splices and returns the new slice.  If
// there were other changes done on the String before this operation,
// that will not be reflected in the output but it will be guaranteed
// to be reflected in the next call  to Latest.
func (l String) Splice(offset, removeCount int, replacement string) String {
	before := l.Value[offset : offset+removeCount]
	value := l.Value[:offset] + replacement + l.Value[offset+removeCount:]

	//  update offset and count to utf16
	offset = len(utf16.Encode([]rune(l.Value[:offset])))
	removeCount = len(utf16.Encode([]rune(before)))

	o := *l.add(l.Start, offset)
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	ctl := l.Control.UpdateSync([]dot.Change{c})
	start, end := l.spliceOffsets(offset, removeCount, len(utf16.Encode([]rune(replacement))))
	return String{Control: ctl, Value: value, Start: start, End: end}
}

// SpliceAsync asynchronously splices and returns the new slice.  If
// there were other changes done on the String before this operation,
// that will not be reflected in the output but will eventually get
// reflected in a call to Latest.
func (l String) SpliceAsync(offset, removeCount int, replacement string) String {
	before := l.Value[offset : offset+removeCount]
	value := l.Value[:offset] + replacement + l.Value[offset+removeCount:]

	//  update offset and count to utf16
	offset = len(utf16.Encode([]rune(l.Value[:offset])))
	removeCount = len(utf16.Encode([]rune(before)))

	o := *l.add(l.Start, offset)
	splice := &dot.SpliceInfo{Offset: o, Before: before, After: replacement}
	c := dot.Change{Splice: splice}
	ctl := l.Control.UpdateAsync([]dot.Change{c})
	start, end := l.spliceOffsets(offset, removeCount, len(utf16.Encode([]rune(replacement))))
	return String{Control: ctl, Value: value, Start: start, End: end}
}

// Latest returns the latest value. The current object may have been
// deleted, in which case it returns the zero value and sets the bool
// to false.
func (l String) Latest() (String, bool) {
	v, ctl, start, end := l.Control.LatestAt(l.Start, l.End)
	if ctl == nil {
		return String{}, false
	}

	value := utf16.Encode([]rune(v.(string)))
	s, e := 0, len(value)
	if start != nil {
		s = *start
	}
	if end != nil {
		e = *end
	}

	val := string(utf16.Decode(value[s:e:e]))
	return String{Value: val, Start: start, End: end, Control: ctl}, true
}

func (l String) spliceOffsets(offset, removeCount, replaceCount int) (*int, *int) {
	var startp, endp *int

	if l.Start != nil {
		startp = l.Start
	}

	if l.End != nil {
		endp = l.add(l.End, replaceCount-removeCount)
	}
	return startp, endp
}

func (l String) add(p *int, v int) *int {
	if p != nil {
		v += *p
	}
	return &v
}
