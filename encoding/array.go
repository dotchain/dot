// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

// Array implements a JSON array. Please take a look at encoding/sparse
// for a better example of how to implement custom array-like encodings
type Array struct {
	c Catalog
	v []interface{}
}

// NewArray is the internal constructor. Note that this is
// automatically registered.
func NewArray(c Catalog, f []interface{}) Array {
	return Array{c, f}
}

func (s Array) fromInterface(i interface{}) Array {
	if i == nil {
		return Array{s.c, []interface{}{}}
	}

	switch i := i.(type) {
	case Array:
		return i
	case []interface{}:
		return Array{s.c, i}
	case enrichArray:
		return s.fromInterface(i.ArrayLike)
	default:
		result := []interface{}{}
		s.c.Get(i).ForEach(func(_ int, v interface{}) {
			result = append(result, v)
		})
		return Array{s.c, result}
	}
}

// Count is the total number elemnts in the array
func (s Array) Count() int {
	return len(s.v)
}

// Slice returns a new array which represents a slice of the contents
func (s Array) Slice(offset, count int) ArrayLike {
	if offset+count > s.Count() {
		panic(errIndexOutOfBounds)
	}
	return Array{s.c, s.v[offset : offset+count]}
}

// Splice implements delete + replace
func (s Array) Splice(offset int, before, after interface{}) ArrayLike {
	deleted, a := s.fromInterface(before).Count(), s.fromInterface(after)

	left, middle, right := s.v[:offset], a.v, s.v[offset+deleted:]
	joined := append(append(append([]interface{}{}, left...), middle...), right...)
	return Array{s.c, joined}
}

// RangeApply applies a function to the range and replaces
// the current array with the modified values. It returns
// a new array and does not actually modify the input array
// in place.
func (s Array) RangeApply(offset, count int, fn func(interface{}) interface{}) ArrayLike {
	after := make([]interface{}, count)
	before := make([]interface{}, count)
	s.Slice(offset, count).ForEach(func(kk int, v interface{}) {
		before[kk] = v
		after[kk] = fn(v)
	})
	return s.Splice(offset, before, after)
}

// ForEach iterates over all the items of the array in
// sequential order
func (s Array) ForEach(fn func(offset int, val interface{})) {
	for offset := 0; offset < len(s.v); offset++ {
		fn(offset, s.v[offset])
	}
}
