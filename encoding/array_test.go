// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import (
	"fmt"
	"github.com/dotchain/dot/encoding"
	"strconv"
	"testing"
)

type Array struct {
	initial, empty, insert, other interface{}
	objectInsert                  interface{}
	offset                        string
}

func ArrayTest() Array {
	return Array{
		initial: []interface{}{0, 1, 2, 3, 4, "hello", 6},
		empty:   []interface{}{},
		insert:  []interface{}{13, 14},
		other: map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  []interface{}{5, "q"},
		},
		objectInsert: "hello",
		offset:       "2",
	}
}

func (s Array) TestAll(t *testing.T) {
	if s.other == nil {
		// set other to be simple array
		s.other = []interface{}{"q", "p", 32, 44}
	}

	t.Run("Empty", s.TestEmpty)
	t.Run("UnmarshalMarshal", s.TestUnmarshalMarshal)
	t.Run("ObjectBehavior", s.TestObjectBehavior)
	t.Run("Start", s.TestStart)
	t.Run("Middle", s.TestMiddle)
	t.Run("End", s.TestEnd)
	t.Run("StringInserts", s.TestStringInserts)
	t.Run("TestDictionaryInserts", s.TestDictionaryInserts)
	t.Run("OtherInserts", s.TestOtherInserts)
}

func (s Array) TestEmpty(t *testing.T) {
	x := encoding.Get(s.empty)
	ensureEqual(t, x.IsArray(), true)
	ensureEqual(t, encoding.IsString(x), false)
	ensureEqual(t, encoding.IsString(s.empty), false)

	if x.Count() != 0 {
		t.Error("Count is not zero", x.Count())
	}
	if x.Slice(0, 0).Count() != 0 {
		t.Error("Count is not zero", x.Slice(0, 0))
	}

	ensureEqual(t, x.Splice(0, nil, nil), x)

	s1 := x.Splice(0, s.empty, s.insert)
	s2 := x.Splice(0, nil, encoding.Get(s.insert))
	expected := encoding.Get(s.insert)

	ensureEqual(t, s1, expected)
	ensureEqual(t, s2, expected)

	s3 := encoding.Get(s1.Splice(0, s.insert, s.empty))
	s4 := encoding.Get(s1.Splice(0, encoding.Get(s.insert), nil))

	ensureEqual(t, s3, x)
	ensureEqual(t, s4, x)

	// error cases
	shouldPanic(t, "out of bounds1", func() { x.Slice(1, 0) })
	shouldPanic(t, "out of bounds2", func() { x.Slice(0, 1) })
	shouldPanic(t, "out of bounds3", func() { x.Splice(1, nil, nil) })
	shouldPanic(t, "out of bounds4", func() { x.Splice(0, "a", nil) })
}

func (s Array) TestUnmarshalMarshal(t *testing.T) {

	initial := encoding.Get(s.initial)
	dupe := encoding.Get(unmarshal(marshal(initial)))

	ensureEqual(t, dupe, initial)

	initial = encoding.Get(s.empty)
	dupe = encoding.Get(unmarshal(marshal(initial)))

	ensureEqual(t, dupe, initial)
}

func (s Array) TestOtherInserts(t *testing.T) {
	// validate that inserting "other" into an empty yields same
	// results
	result := encoding.Get(s.empty).Splice(0, nil, s.other)
	expected := encoding.Get(s.other)

	if expected.Count() != result.Count() {
		t.Error("Counts differ", expected.Count(), result.Count())
	}

	expected.ForEach(func(kk int, val interface{}) {
		actual := encoding.Get(result).Get(fmt.Sprintf("%d", kk))
		ensureEqual(t, actual, val)
	})
}

func (s Array) TestStart(t *testing.T) {
	s.testAtOffset(t, 0)
}

func (s Array) TestMiddle(t *testing.T) {
	s.testAtOffset(t, 1)
}

func (s Array) TestEnd(t *testing.T) {
	s.testAtOffset(t, encoding.Get(s.initial).Count()-encoding.Get(s.insert).Count())
}

func (s Array) testAtOffset(t *testing.T, offset int) {
	initial, insert, empty := encoding.Get(s.initial), encoding.Get(s.insert), encoding.Get(s.empty)

	ensureEqual(t, initial.IsArray(), true)
	ensureEqual(t, encoding.IsString(initial), false)
	ensureEqual(t, encoding.IsString(s.initial), false)

	// delete tests
	slice := encoding.Get(initial.Slice(offset, 2))
	before := encoding.Get(initial.Slice(0, offset))
	rest := encoding.Get(initial.Slice(offset+2, initial.Count()-offset-2))
	s1 := slice.Splice(2, nil, rest).Splice(0, nil, before)
	s2 := slice.Splice(2, empty, rest).Splice(0, empty, before)
	if slice.Count() != 2 {
		t.Error("Slice of 2 has different count", slice, 2)
	}
	if rest.Count() != initial.Count()-offset-2 {
		t.Error("Rest has invalid count", rest.Count(), initial.Count()-offset-2)
	}

	ensureEqual(t, s1, initial)
	ensureEqual(t, s2, initial)

	// insert tests
	inserted1 := encoding.Get(initial.Splice(offset, empty, insert))
	inserted2 := encoding.Get(initial.Splice(offset, nil, insert))
	resliced := encoding.Get(inserted1.Slice(offset, insert.Count()))

	ensureEqual(t, inserted1, inserted2)
	ensureEqual(t, resliced, insert)

	if inserted1.Count() != initial.Count()+insert.Count() {
		t.Error("Count does not match", initial.Count(), insert.Count(), inserted1.Count())
	}

	// replace tests
	before1 := initial.Slice(offset, insert.Count())
	replaced := initial.Splice(offset, before1, insert)
	after := replaced.Slice(offset, insert.Count())
	undone := replaced.Splice(offset, after, before1)

	ensureEqual(t, after, insert)
	ensureEqual(t, initial, undone)

	// failing out of bounds tests
	shouldPanic(t, "out of bounds1", func() { initial.Slice(0, -1) })
	shouldPanic(t, "out of bounds2", func() { initial.Slice(0, 1000) })
	shouldPanic(t, "out of bounds3", func() { initial.Slice(1000, 0) })
	shouldPanic(t, "out of bounds4", func() { initial.Splice(1000, nil, nil) })
	shouldPanic(t, "out of bounds5", func() { insert.Splice(0, initial, nil) })
}

func (s Array) TestObjectBehavior(t *testing.T) {
	offset := s.offset
	initial := encoding.Get(s.initial)
	shouldPanic(t, "get non integer on arrays", func() { initial.Get("something") })
	shouldPanic(t, "set(nil) on non-integer key", func() { initial.Set("something", nil) })
	shouldPanic(t, "set(non-nil) on string", func() { initial.Set("something", initial) })

	// the following shoud work
	insert := s.objectInsert
	old := initial.Get(offset)
	updated := initial.Set(offset, insert)

	ensureEqual(t, encoding.Get(updated).Get(offset), insert)
	ensureEqual(t, initial.Get(offset), old)
	undone := encoding.Get(updated).Set(offset, old)
	ensureEqual(t, undone, initial)

	count := 0
	initial.ForKeys(func(key string, val interface{}) {
		if strconv.Itoa(count) != key {
			t.Error("Unexpected key")
		}
		ensureEqual(t, initial.Get(strconv.Itoa(count)), val)
		count++
	})

	if count != initial.Count() {
		t.Error("Got unexpected keys")
	}
}

func (s Array) TestStringInserts(t *testing.T) {
	initial := encoding.Get(s.initial)

	shouldPanic(t, "inserting string into array", func() { initial.Splice(0, nil, encoding.Get("hello")) })
	shouldPanic(t, "deleting string from array", func() { initial.Splice(0, encoding.Get(""), nil) })
	shouldPanic(t, "inserting string into array", func() { initial.Splice(0, nil, "hello") })
	shouldPanic(t, "deleting string from array", func() { initial.Splice(0, "", nil) })
	shouldPanic(t, "inserting string into array", func() { initial.Splice(0, nil, encoding.NewString16("hello")) })
	shouldPanic(t, "deleting string from array", func() { initial.Splice(0, encoding.NewString16(""), nil) })
}

func (s Array) TestDictionaryInserts(t *testing.T) {
	initial := encoding.Get(s.initial)
	insert := map[string]interface{}{"hello": "world"}

	shouldPanic(t, "inserting dictionary into array", func() { initial.Splice(0, nil, insert) })
	shouldPanic(t, "inserting dictionary into array", func() { initial.Splice(0, nil, encoding.Get(insert)) })
}
