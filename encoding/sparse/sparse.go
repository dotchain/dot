// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sparse

import (
	"encoding/json"
	"github.com/dotchain/dot/encoding"
	"github.com/pkg/errors"
)

func init() {
	encoding.Default.RegisterConstructor("SparseArray", NewArray)
}

// Array is a run length encoding of an array.
type Array struct {
	c encoding.Catalog
	v []interface{}
}

// NewArray is the constructor to be used with encoding.Catalog.
// Note that this type is automatically registered with encoding.
func NewArray(c encoding.Catalog, m map[string]interface{}) Array {
	// the actual value is always in the 'dot:encoded' field.
	values, _ := m["dot:encoded"].([]interface{})
	return Array{c, values}
}

// MarshalJSON provides the right JSON layout
func (s Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"dot:encoding": "Array",
		"dot:encoded":  s.v,
	})
}

// Count returns the total count of the logical array
func (s Array) Count() int {
	count := 0
	for kk := 0; kk < len(s.v); kk += 2 {
		count += s.v[kk].(int)
	}
	return count
}

// Slice implements slice on the logical array
func (s Array) Slice(offset, count int) encoding.ArrayLike {
	return s.slice(offset, count)
}

func (s Array) slice(offset, count int) Array {
	if count < 0 {
		panic(errors.Errorf("Slice out of bounds: count is negative: %d", count))
	}

	// seen tracks total seen, have = size of result
	seen, have := 0, 0
	result := []interface{}{}
	for kk := 0; kk < len(s.v); kk += 2 {
		repeat := s.v[kk].(int)
		seen += repeat

		if seen <= offset {
			continue
		}
		consume := repeat
		if offset >= seen-repeat {
			consume = seen - offset
		}

		if consume > count-have {
			consume = count - have
		}

		if consume > 0 {
			result = append(result, consume, s.v[kk+1])
		}

		have += consume
		if have == count {
			return Array{s.c, result}
		}
	}
	if count == have && offset+count == seen {
		return Array{s.c, result}
	}
	panic(errors.Errorf("Slice out of bounds: offset, count, size = %d, %d, %d", offset, count, seen))
}

func (s Array) fromInterface(i interface{}) Array {
	if i == nil {
		return Array{s.c, nil}
	}

	insert := []interface{}{}
	lastCount := 0
	s.c.Get(i).ForEach(func(idx int, val interface{}) {
		if idx == 0 || insert[len(insert)-1] != val {
			insert = append(insert, 1, val)
			lastCount = 1
		} else {
			lastCount++
			insert[len(insert)-2] = lastCount
		}
	})
	return Array{s.c, insert}
}

// Splice implements the splice operation, allowing any element type
// but not fake arrays like strings for the "before" and "after"
func (s Array) Splice(offset int, before, after interface{}) encoding.ArrayLike {
	deleted, a := s.fromInterface(before).Count(), s.fromInterface(after)
	left, middle, right := s.slice(0, offset).v, a.v, s.slice(offset+deleted, s.Count()-offset-deleted).v

	left = append([]interface{}{}, left...)

	// coalesce
	if len(middle) > 0 && len(left) > 0 && middle[1] == left[len(left)-1] {
		left[len(left)-2] = left[len(left)-2].(int) + middle[0].(int)
		middle = middle[2:]
	}

	if len(middle) > 0 && len(right) > 0 && right[1] == middle[len(middle)-1] {
		middle = append([]interface{}{}, middle...)
		middle[len(middle)-2] = middle[len(middle)-2].(int) + right[0].(int)
		right = right[2:]
	}

	return Array{s.c, append(append(left, middle...), right...)}
}

// RangeApply iterates over the logical array, returning a new run length array
// encoded with values from the callback
func (s Array) RangeApply(offset, count int, fn func(interface{}) interface{}) encoding.ArrayLike {
	insert := []interface{}{}
	s.Slice(offset, count).ForEach(func(idx int, input interface{}) {
		val := fn(input)
		if idx == 0 || insert[len(insert)-1] != val {
			insert = append(insert, 1, val)
		} else {
			insert[len(insert)-2] = 1 + (insert[len(insert)-2].(int))
		}
	})
	return s.Splice(offset, s.Slice(offset, count), Array{s.c, insert})
}

// ForEach iterates of the logical array
func (s Array) ForEach(fn func(offset int, val interface{})) {
	seen := 0
	for kk := 0; kk < len(s.v); kk += 2 {
		for count := 0; count < s.v[kk].(int); count++ {
			fn(seen, s.v[kk+1])
			seen++
		}
	}
}
