// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"
	"unicode/utf16"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"

	"github.com/dotchain/dot/fred"
)

func TestTextSlice(t *testing.T) {
	s := fred.Text("hello, 🌂🌂")
	if x := s.Slice(3, 0); x != fred.Text("") {
		t.Error("Unexpected Slice(3, 0)", x)
	}
	if x := s.Slice(7, 2); x != fred.Text("🌂") {
		t.Error("Unexpected Slice()", x)
	}
}

func TestTextCount(t *testing.T) {
	if x := fred.Text("🌂").Count(); x != len(utf16.Encode([]rune("🌂"))) {
		t.Error("Unexpected Count()", x)
	}
}

func TestTextApply(t *testing.T) {
	s := fred.Text("hello, 🌂🌂")

	x := s.Apply(nil, nil)
	if x != s {
		t.Error("Unexpected Apply.nil", x)
	}

	x = s.Apply(nil, changes.Replace{Before: s, After: changes.Nil})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = s.Apply(nil, changes.Replace{Before: s, After: types.S8("OK")})
	if x != types.S8("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	x = s.Apply(nil, changes.Splice{Offset: 7, Before: s.Slice(7, 2), After: fred.Text("-")})
	if x != fred.Text("hello, -🌂") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.Apply(nil, changes.Splice{Offset: 11, Before: fred.Text(""), After: fred.Text("-")})
	if x != fred.Text("hello, 🌂🌂-") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.Apply(nil, changes.Move{Offset: 7, Count: 2, Distance: -1})
	if x != fred.Text("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}

	x = s.Apply(nil, changes.ChangeSet{changes.Move{Offset: 7, Count: 2, Distance: -1}})
	if x != fred.Text("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}

	x = s.Apply(nil, changes.PathChange{Change: changes.Move{Offset: 7, Count: 2, Distance: -1}})
	if x != fred.Text("hello,🌂 🌂") {
		t.Error("Unexpected Apply.Move", x)
	}
}
