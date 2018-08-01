// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

import "unicode/utf16"

// String16 implements a UTF16 string.  This is important
// because Javascript clients use this representation by default
// and if the dot.Transformer did not use the same representation,
// all the offsets and counts interpreted by it could be different
// and lead to issues
type String16 []uint16

// NewString16 is the constructor as used by catalog.Catalog
func NewString16(s string) String16 {
	return String16(utf16.Encode([]rune(s)))
}

func (s String16) fromInterface(i interface{}) String16 {
	if i == nil {
		return String16([]uint16{})
	}

	switch i := i.(type) {
	case String16:
		return i
	case string:
		return NewString16(i)
	case enrichArray:
		return s.fromInterface(i.ArrayLike)
	}

	panic(errNotStringType)
}

// Count returns the size of the string in UTF16 characters
func (s String16) Count() int {
	return len(s)
}

// Slice works on UTF16 offsets in the string
func (s String16) Slice(offset, count int) ArrayLike {
	if offset+count > s.Count() {
		panic(errStringIndexOutOfBounds)
	}
	return s[offset : offset+count]
}

// Splice works on UTF16 offsets and only accepts UTF16 strings
// (i.e. it interprets before/after as UTF16 arrays). This will
// yield the right result for JS clients
func (s String16) Splice(offset int, before, after interface{}) ArrayLike {
	deleted, a := s.fromInterface(before).Count(), s.fromInterface(after)

	left, middle, right := s[:offset], a, s[offset+deleted:]
	return String16(append(append(append([]uint16{}, left...), middle...), right...))
}

// RangeApply should not really be called for strings
func (s String16) RangeApply(offset, count int, fn func(interface{}) interface{}) ArrayLike {
	panic(errMethodNotSupported)
}

// ForEach should not really be called for strings
func (s String16) ForEach(func(offset int, val interface{})) {
	panic(errMethodNotSupported)
}
