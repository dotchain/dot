// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"testing"
	"unicode/utf16"
)

func TestS8Slice(t *testing.T) {
	s := types.S8("hello, 🌂🌂")
	if x := s.Slice(3, 0); x != types.S8("") {
		t.Error("Unexpected Slice(3, 0)", x)
	}
	if x := s.Slice(7, 4); x != types.S8("🌂") {
		t.Error("Unexpected Slice()", x)
	}
}

func TestS8Count(t *testing.T) {
	if x := types.S8("🌂").Count(); x != len("🌂") {
		t.Error("Unexpected Count()", x)
	}
}

func TestS8Apply(t *testing.T) {
	s := types.S8("hello, 🌂🌂")

	x := s.Apply(nil, nil)
	if x != s {
		t.Error("Unexpected Apply.nil", x)
	}

	x = s.Apply(nil, changes.Replace{s, changes.Nil})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = s.Apply(nil, changes.Replace{s, types.S16("OK")})
	if x != types.S16("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	x = s.Apply(nil, changes.Splice{7, s.Slice(7, 4), types.S8("-")})
	if x != types.S8("hello, -🌂") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.Apply(nil, changes.Move{7, 4, -1})
	if x != types.S8("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}

	x = s.Apply(nil, changes.ChangeSet{changes.Move{7, 4, -1}})
	if x != types.S8("hello,🌂 🌂") {
		t.Error("Unexpected Apply.ChangeSet", x)
	}

	x = s.Apply(nil, changes.PathChange{nil, changes.Move{7, 4, -1}})
	if x != types.S8("hello,🌂 🌂") {
		t.Error("Unexpected Apply.PathChange", x)
	}
}

func TestS16Slice(t *testing.T) {
	s := types.S16("hello, 🌂🌂")
	if x := s.Slice(3, 0); x != types.S16("") {
		t.Error("Unexpected Slice(3, 0)", x)
	}
	if x := s.Slice(7, 2); x != types.S16("🌂") {
		t.Error("Unexpected Slice()", x)
	}
}

func TestS16Count(t *testing.T) {
	if x := types.S16("🌂").Count(); x != len(utf16.Encode([]rune("🌂"))) {
		t.Error("Unexpected Count()", x)
	}
	if x := types.S16("hello").ToUTF16(1); x != 1 {
		t.Error("Unexpected idx calculation", x)
	}
}

func TestS16Apply(t *testing.T) {
	s := types.S16("hello, 🌂🌂")

	x := s.Apply(nil, nil)
	if x != s {
		t.Error("Unexpected Apply.nil", x)
	}

	x = s.Apply(nil, changes.Replace{s, changes.Nil})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = s.Apply(nil, changes.Replace{s, types.S8("OK")})
	if x != types.S8("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	x = s.Apply(nil, changes.Splice{7, s.Slice(7, 2), types.S16("-")})
	if x != types.S16("hello, -🌂") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.Apply(nil, changes.Splice{11, types.S16(""), types.S16("-")})
	if x != types.S16("hello, 🌂🌂-") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.Apply(nil, changes.Move{7, 2, -1})
	if x != types.S16("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}

	x = s.Apply(nil, changes.ChangeSet{changes.Move{7, 2, -1}})
	if x != types.S16("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}

	x = s.Apply(nil, changes.PathChange{nil, changes.Move{7, 2, -1}})
	if x != types.S16("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}
}

// this implements Change but not CustomChange
type poorlyDefinedChange struct{}

func (p poorlyDefinedChange) Merge(o changes.Change) (changes.Change, changes.Change) {
	return o, nil
}

func (p poorlyDefinedChange) Revert() changes.Change {
	return p
}

func TestStringPanics(t *testing.T) {
	mustPanic := func(fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Failed to panic")
			}
		}()
		fn()
	}

	mustPanic(func() {
		types.S8("hello").Apply(nil, poorlyDefinedChange{})
	})

	mustPanic(func() {
		types.S16("hello").Apply(nil, poorlyDefinedChange{})
	})

	mustPanic(func() {
		s := types.S16("hello, 🌂🌂")
		s.Apply(nil, changes.ChangeSet{changes.Move{7, 3, -1}})
	})

	mustPanic(func() {
		s := types.S16("hello, 🌂🌂")
		s.ToUTF16(10)
	})
}
