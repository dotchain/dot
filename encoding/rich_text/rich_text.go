// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

import (
	"encoding/json"
	"fmt"
	"github.com/dotchain/dot/encoding"
	"github.com/pkg/errors"
	"reflect"
	"unicode/utf16"
)

func init() {
	ctor := func(c encoding.Catalog, m map[string]interface{}) encoding.UniversalEncoding {
		return NewArray(c, m)
	}

	key := "github.com/dotchain/dot/encoding/rich_text"
	encoding.Default.RegisterTypeConstructor("RichText", key, ctor)
}

// Array is a run length encoding of an attributed string
type Array struct {
	c encoding.Catalog
	m []map[string]string
}

// NewArray is the constructor to be used with encoding.Catalog
// Note that this type is automatically registered with the default
// encoding catalog
func NewArray(c encoding.Catalog, m map[string]interface{}) Array {
	bytes, err := json.Marshal(m["dot:encoded"])
	mustNotFail(err)
	result := Array{c: c}
	mustNotFail(json.Unmarshal(bytes, &result.m))
	return result
}

// MarshalJSON provides the right JSON layout
func (r Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"dot:encoding": "RichText",
		"dot:encoded":  r.m,
	})
}

// Count returns the total count of the logical array
func (r Array) Count() int {
	count := 0
	for _, attr := range r.m {
		count += len(utf16.Encode([]rune(attr["text"])))
	}
	return count
}

func (r Array) withoutText(attr map[string]string) map[string]string {
	result := map[string]string{}
	for key, val := range attr {
		if key != "text" {
			result[key] = val
		}
	}
	return result
}

func (r Array) slice(offset, count int) Array {
	if count < 0 || offset+count > r.Count() {
		panic(errors.Errorf("Out of bounds slice(%d, %d) on rich text of length %d", offset, count, r.Count()))
	}

	if count == 0 {
		return Array{c: r.c}
	}

	var attrs []map[string]string
	seen, have := 0, 0
	for _, attr := range r.m {
		s16 := utf16.Encode([]rune(attr["text"]))
		size := len(s16)
		start := seen
		seen += size

		if seen <= offset {
			continue
		}
		if offset > start {
			s16 = s16[offset-start:]
		}
		if len(s16) > count-have {
			s16 = s16[:count-have]
		}
		if len(s16) > 0 {
			have += len(s16)
			entry := r.withoutText(attr)
			entry["text"] = string(utf16.Decode(s16))
			attrs = append(attrs, entry)
			if have == count {
				break
			}
		}
	}
	return Array{r.c, attrs}
}

func (r Array) concat(r2 Array) Array {
	attrs := append([]map[string]string{}, r.m...)
	if len(r.m) > 0 && len(r2.m) > 0 {
		last := attrs[len(attrs)-1]
		if reflect.DeepEqual(r.withoutText(last), r.withoutText(r2.m[0])) {
			text := last["text"] + r2.m[0]["text"]
			last = r.withoutText(last)
			last["text"] = text
			attrs[len(attrs)-1] = last
			attrs = append(attrs, r2.m[1:]...)
			return Array{r.c, attrs}
		}
	}

	attrs = append(attrs, r2.m...)
	return Array{r.c, attrs}
}

func (r Array) append(attr interface{}) Array {
	b, err := json.Marshal(attr)
	mustNotFail(err)
	var a map[string]string
	mustNotFail(json.Unmarshal(b, &a))

	s16 := utf16.Encode([]rune(a["text"]))
	if len(s16) != 1 {
		panic(errors.Errorf("Invalid rich text append: %s", a["text"]))
	}

	m := append([]map[string]string{}, r.m...)
	coalesced := false
	if len(m) > 0 {
		last := m[len(m)-1]
		if reflect.DeepEqual(r.withoutText(last), r.withoutText(a)) {
			text := last["text"] + a["text"]
			m[len(m)-1] = r.withoutText(last)
			m[len(m)-1]["text"] = text
			coalesced = true
		}
	}

	if !coalesced {
		m = append(m, a)
	}
	return Array{r.c, m}
}

func (r Array) fromInterface(i interface{}) Array {
	if i == nil {
		return Array{c: r.c}
	}

	switch insert := r.c.Get(i).(type) {
	case Array:
		return insert
	default:
		// simple array of maps
		result := Array{c: r.c}
		insert.ForEach(func(idx int, val interface{}) {
			result = result.append(val)
		})
		return result
	}
}

// Slice implements slice on the logical array
func (r Array) Slice(offset, count int) encoding.ArrayLike {
	return r.slice(offset, count)
}

// Splice implements the splice operation, allowing any form of Array
func (r Array) Splice(offset int, before, after interface{}) encoding.ArrayLike {
	deleted, mid := r.fromInterface(before).Count(), r.fromInterface(after)
	left, right := r.slice(0, offset), r.slice(offset+deleted, r.Count()-offset-deleted)

	return left.concat(mid).concat(right)
}

// RangeApply iterates over the logical array, returning a new
// Array element
func (r Array) RangeApply(offset, count int, fn func(interface{}) interface{}) encoding.ArrayLike {
	insert := Array{c: r.c}
	r.Slice(offset, count).ForEach(func(idx int, input interface{}) {
		insert = insert.append(fn(input))
	})
	return r.Splice(offset, r.Slice(offset, count), insert)
}

// ForEach iterates of the logical array
func (r Array) ForEach(fn func(offset int, val interface{})) {
	seen := 0
	for _, attr := range r.m {
		text := utf16.Encode([]rune(attr["text"]))
		for _, t := range text {
			without := r.withoutText(attr)
			without["text"] = string(utf16.Decode([]uint16{t}))
			fn(seen, without)
			seen++
		}
	}
}

// ForKeys iterates over the keys which is same as ForEach
func (r Array) ForKeys(fn func(key string, val interface{})) {
	r.ForEach(func(offset int, val interface{}) {
		fn(fmt.Sprintf("%d", offset), val)
	})
}

// Get array entry via ObjectLike interface
func (r Array) Get(key string) interface{} {
	var offset int
	count, err := fmt.Sscanf(key, "%d", &offset)
	mustNotFail(err)
	if count != 1 {
		panic(errors.Errorf("Invalid offset specified for Get %s", key))
	}
	var result interface{}
	r.slice(offset, 1).ForEach(func(o int, val interface{}) {
		result = val
	})
	return result
}

// Set array entry via ObjectLike interface
func (r Array) Set(key string, val interface{}) encoding.ObjectLike {
	var offset int
	count, err := fmt.Sscanf(key, "%d", &offset)
	mustNotFail(err)
	if count != 1 {
		panic(errors.Errorf("Invalid offset specified for Get %s", key))
	}
	spliced := r.Splice(offset, r.slice(offset, 1), []interface{}{val})
	return spliced.(Array)
}

// IsArray implements UniversalEncoding.IsArray
func (r Array) IsArray() bool {
	return true
}

func mustNotFail(err error) {
	if err != nil {
		panic(errors.Errorf("Unexpected error %v", err))
	}
}
