// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import (
	"github.com/dotchain/dot/encoding"
	"testing"
)

type String16 struct{}

func (s String16) TestAll(t *testing.T) {
	t.Run("Empty", s.TestEmpty)
	t.Run("UnmarshalMarshal", s.TestUnmarshalMarshal)
	t.Run("ObjectBehavior", s.TestObjectBehavior)
	t.Run("Start", s.TestStart)
	t.Run("Middle", s.TestMiddle)
	t.Run("End", s.TestEnd)
	t.Run("ArrayInserts", s.TestArrayInserts)
	t.Run("DictionaryInserts", s.TestDictionaryInserts)
}

func (String16) TestEmpty(t *testing.T) {
	ensureEqual(t, encoding.IsString(nil), true)
	ensureEqual(t, encoding.IsString(""), true)
	ensureEqual(t, encoding.IsString(encoding.Get("")), true)

	s := encoding.Get("")
	if s.Count() != 0 {
		t.Error("Count is not zero", s.Count())
	}
	if s.Slice(0, 0).Count() != 0 {
		t.Error("Count is not zero", s.Slice(0, 0))
	}
	if j := marshal(s); j != `""` {
		t.Error("JSON encoding is not an empty string", j)
	}
	ensureEqual(t, s.Splice(0, nil, nil), s)

	s1 := s.Splice(0, "", "hello")
	s2 := s.Splice(0, nil, "hello")
	expected := encoding.Get("hello")

	ensureEqual(t, s1, expected)
	ensureEqual(t, s2, expected)

	s3 := encoding.Get(s1.Splice(0, "hello", ""))
	s4 := encoding.Get(s1.Splice(0, "hello", nil))

	ensureEqual(t, s3, s)
	ensureEqual(t, s4, s)

	// error cases

	shouldPanic(t, "out of bounds1", func() { s.Slice(1, 0) })
	shouldPanic(t, "out of bounds2", func() { s.Slice(0, 1) })
	shouldPanic(t, "out of bounds3", func() { s.Splice(1, nil, nil) })
	shouldPanic(t, "out of bounds4", func() { s.Splice(0, "a", nil) })
}

func (String16) TestUnmarshalMarshal(t *testing.T) {
	initial := encoding.Get("initial")

	// marshal/unmarshal test
	dupe := encoding.Get(unmarshal(marshal(initial)))
	ensureEqual(t, dupe, initial)
}

func (s String16) TestObjectBehavior(t *testing.T) {
	s.testObjectBehavior(t, "2", encoding.Get("hello"))
}

func (s String16) TestStart(t *testing.T) {
	initial := encoding.Get("hello world")
	insert := encoding.Get("insert")
	empty := encoding.Get("")
	s.testAtOffset(t, 0, initial, insert, empty)
}

func (s String16) TestMiddle(t *testing.T) {
	initial := encoding.Get("hello world")
	insert := encoding.Get("insert")
	empty := encoding.Get("")
	s.testAtOffset(t, 1, initial, insert, empty)
}

func (s String16) TestEnd(t *testing.T) {
	initial := encoding.Get("hello world")
	insert := encoding.Get("insert")
	empty := encoding.Get("")
	s.testAtOffset(t, initial.Count()-insert.Count(), initial, insert, empty)
}

func (String16) testAtOffset(t *testing.T, offset int, initial, insert, empty encoding.UniversalEncoding) {
	ensureEqual(t, encoding.IsString("initial"), true)
	ensureEqual(t, encoding.IsString(encoding.Get("initial")), true)

	// delete
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

	// insert some
	inserted1 := encoding.Get(initial.Splice(offset, empty, insert))
	inserted2 := encoding.Get(initial.Splice(offset, nil, insert))
	resliced := encoding.Get(inserted1.Slice(offset, insert.Count()))

	ensureEqual(t, inserted1, inserted2)
	ensureEqual(t, resliced, insert)

	if inserted1.Count() != initial.Count()+insert.Count() {
		t.Error("Count does not match", initial.Count(), insert.Count(), inserted1.Count())
	}

	// replace
	before1 := initial.Slice(offset, insert.Count())
	replaced := initial.Splice(offset, before1, insert)
	after := replaced.Slice(offset, insert.Count())
	undone := replaced.Splice(offset, after, before1)

	ensureEqual(t, after, insert)
	ensureEqual(t, initial, undone)
}

func (String16) testObjectBehavior(t *testing.T, offset string, initial encoding.UniversalEncoding) {
	shouldPanic(t, "get on string", func() { initial.Get(offset) })
	shouldPanic(t, "set(nil) on string", func() { initial.Set(offset, nil) })
	shouldPanic(t, "set(non-nil) on string", func() { initial.Set(offset, initial) })
}

func (String16) TestArrayInserts(t *testing.T) {
	insert := []interface{}{"hello"}
	insertEmpty := []interface{}{}
	initial := encoding.Get("hello world")

	shouldPanic(t, "inserting array into string", func() { initial.Splice(0, nil, encoding.Get(insert)) })
	shouldPanic(t, "deleting array from string", func() { initial.Splice(0, encoding.Get(insertEmpty), nil) })
	shouldPanic(t, "inserting array into string", func() { initial.Splice(0, nil, insert) })
	shouldPanic(t, "deleting array from string", func() { initial.Splice(0, insertEmpty, nil) })
	shouldPanic(t, "inserting array into string", func() { initial.Splice(0, nil, encoding.NewArray(encoding.Default, insert)) })
	shouldPanic(t, "deleting array from string", func() { initial.Splice(0, encoding.NewArray(encoding.Default, insertEmpty), nil) })

}

func (String16) TestDictionaryInserts(t *testing.T) {
	insert := map[string]interface{}{"hello": "world"}
	initial := encoding.Get("hello world")
	shouldPanic(t, "inserting dictionary into array", func() { initial.Splice(0, nil, insert) })
	shouldPanic(t, "inserting dictionary into array", func() { initial.Splice(0, nil, encoding.Get(insert)) })
}
