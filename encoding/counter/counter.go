// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package counter implements a simple counter encoding for DOT.
//
// Please see package ver for a client interface:
// https://godoc.org/github.com/dotchain/ver
package counter

import (
	"github.com/dotchain/dot"
	"github.com/dotchain/dot/encoding"
)

func init() {
	ctor := func(_ encoding.Catalog, m map[string]interface{}) encoding.UniversalEncoding {
		val, _ := m["dot:encoded"].([]interface{})
		var sum int64
		for _, v := range val {
			sum += int64(v.(float64))
		}
		return Counter(sum)
	}

	key := "github.com/dotchain/dot/encoding/counter"
	encoding.Default.RegisterTypeConstructor("Counter", key, ctor)
}

// Counter represents an int that can be incremented/decremented
type Counter int64

// NormalizeDOT returns a native golang object that can be marshaled as JSON
func (c Counter) NormalizeDOT() interface{} {
	return map[string]interface{}{
		"dot:encoding": "Counter",
		"dot:generic":  true,
		"dot:encoded":  []interface{}{float64(c)},
	}
}

// Increment returns the updated counter and the change structure
func (c Counter) Increment(by int64) (Counter, dot.Change) {
	splice := &dot.SpliceInfo{After: []interface{}{float64(by)}}
	return Counter(int64(c) + by), dot.Change{Splice: splice}
}

// Count implements encoding.UniversalEncoding.Count
func (c Counter) Count() int {
	return 1
}

// Slice implements encoding.UniversalEncoding.Slice
func (c Counter) Slice(offset, count int) encoding.ArrayLike {
	panic("unexpected Slice() called on counter")
}

// Splice implements encoding.UniversalEncoding.Splice
func (c Counter) Splice(offset int, before, after interface{}) encoding.ArrayLike {
	result := float64(int64(c))
	encoding.Get(after).ForEach(func(_ int, val interface{}) {
		result += val.(float64)
	})
	return Counter(int64(result))
}

// RangeApply implements encoding.UniversalEncoding.RangeApply
func (c Counter) RangeApply(offset, count int, fn func(interface{}) interface{}) encoding.ArrayLike {
	panic("Unexpected RangeApply called on counter")
}

// ForEach implements encoding.UniversalEncoding.ForEach
func (c Counter) ForEach(fn func(offset int, val interface{})) {
	fn(0, float64(int64(c)))
}

// ForKeys implements encoding.UniversalEncoding.ForKeys
func (c Counter) ForKeys(fn func(key string, val interface{})) {
	panic("Unexpected ForKeys called on counter")
}

// Get implements encoding.UniversalEncoding.Get
func (c Counter) Get(key string) interface{} {
	if key == "0" {
		return float64(int64(c))
	}
	panic("Unexpected Get called on counter")
}

// Set implements encoding.UniversalEncoding.Set
func (c Counter) Set(key string, val interface{}) encoding.ObjectLike {
	panic("Unexpected Set called on counter")
}

// IsArray implements UniversalEncoding.IsArray
func (c Counter) IsArray() bool {
	return true
}
