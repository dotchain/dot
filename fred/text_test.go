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

	x = s.ApplyCollection(nil, changes.Splice{Offset: 7, Before: s.Slice(7, 2), After: fred.Text("-")})
	if x != fred.Text("hello, -🌂") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.Apply(nil, changes.Splice{Offset: 11, Before: fred.Text(""), After: fred.Text("-")})
	if x != fred.Text("hello, 🌂🌂-") {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = s.ApplyCollection(nil, changes.Move{Offset: 7, Count: 2, Distance: -1})
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

func TestTextConcatCall(t *testing.T) {
	base := fred.Fixed(fred.Text("hello "))
	concat := fred.Field(base, fred.Fixed(fred.Text("concat")))
	suffix := fred.Fixed(fred.Text("world"))
	expr := fred.Call(concat, suffix)
	if x := expr.Eval(env); x != fred.Text("hello world") {
		t.Error("Unexpected", x)
	}
}

func TestTextSlice1Call(t *testing.T) {
	base := fred.Fixed(fred.Text("hello"))
	slice := fred.Field(base, fred.Fixed(fred.Text("slice")))
	expr := fred.Call(slice, fred.Fixed(fred.Num("3")))
	if x := expr.Eval(env); x != fred.Text("lo") {
		t.Error("Unexpected", x)
	}
}

func TestTextSlice2Call(t *testing.T) {
	base := fred.Fixed(fred.Text("hello"))
	slice := fred.Field(base, fred.Fixed(fred.Text("slice")))
	expr := fred.Call(
		slice,
		fred.Fixed(fred.Num("0")),
		fred.Fixed(fred.Num("4")),
	)
	if x := expr.Eval(env); x != fred.Text("hell") {
		t.Error("Unexpected", x)
	}
}

func TestTextSliceErrors(t *testing.T) {
	base := fred.Fixed(fred.Text("hello"))
	slice := fred.Field(base, fred.Fixed(fred.Text("slice")))
	wat := fred.Error("wat")
	boo := fred.Error("math/big: cannot unmarshal \"boo\" into a *big.Rat")
	cases := map[string][]fred.Val{
		"err0":  {fred.ErrInvalidArgs, fred.Num("1/2"), fred.Num("2")},
		"err1":  {fred.ErrInvalidArgs},
		"err2":  {fred.ErrInvalidArgs, fred.Num("1"), fred.Num("2"), fred.Num("3")},
		"err3":  {fred.ErrNotNumber, fred.Text("wat"), fred.Num("0")},
		"err4":  {fred.ErrNotNumber, fred.Num("0"), fred.Text("wat")},
		"err5":  {wat, wat, fred.Num("5")},
		"err6":  {boo, fred.Num("boo"), fred.Num("5")},
		"err7":  {boo, fred.Num("5"), fred.Num("boo")},
		"err8":  {fred.ErrInvalidArgs, fred.Num("-1"), fred.Num("0")},
		"err9":  {fred.ErrInvalidArgs, fred.Num("0"), fred.Num("-1")},
		"err10": {fred.ErrInvalidArgs, fred.Num("0"), fred.Num("100")},
		"err11": {fred.ErrInvalidArgs, fred.Num("100"), fred.Num("0")},
	}

	for name, vals := range cases {
		t.Run(name, func(t *testing.T) {
			expected := vals[0]
			args := []fred.Def{}
			for _, arg := range vals[1:] {
				args = append(args, fred.Fixed(arg))
			}

			x := fred.Call(slice, args...).Eval(env)
			if x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}

func TestTextSpliceCall(t *testing.T) {
	base := fred.Fixed(fred.Text("hello"))
	splice := fred.Field(base, fred.Fixed(fred.Text("splice")))
	expr := fred.Call(
		splice,
		fred.Fixed(fred.Num("3")),
		fred.Fixed(fred.Num("2")),
		fred.Fixed(fred.Text("la")),
	)
	if x := expr.Eval(env); x != fred.Text("hella") {
		t.Error("Unexpected", x)
	}
}

func TestTextSpliceErrors(t *testing.T) {
	base := fred.Fixed(fred.Text("hello"))
	slice := fred.Field(base, fred.Fixed(fred.Text("splice")))
	wat := fred.Error("wat")
	boo := fred.Error("math/big: cannot unmarshal \"boo\" into a *big.Rat")
	goo := fred.Text("goo")
	cases := map[string][]fred.Val{
		"err1":  {fred.ErrInvalidArgs},
		"err2":  {fred.ErrInvalidArgs, fred.Num("1"), fred.Num("2")},
		"err3":  {fred.ErrNotNumber, fred.Text("wat"), fred.Num("0"), goo},
		"err4":  {fred.ErrNotNumber, fred.Num("0"), fred.Text("wat"), goo},
		"err5":  {wat, wat, fred.Num("5"), goo},
		"err6":  {boo, fred.Num("boo"), fred.Num("5"), goo},
		"err7":  {boo, fred.Num("5"), fred.Num("boo"), goo},
		"err8":  {fred.ErrInvalidArgs, fred.Num("-1"), fred.Num("0"), goo},
		"err9":  {fred.ErrInvalidArgs, fred.Num("0"), fred.Num("-1"), goo},
		"err10": {fred.ErrInvalidArgs, fred.Num("0"), fred.Num("100"), goo},
		"err11": {fred.ErrInvalidArgs, fred.Num("100"), fred.Num("0"), goo},
	}

	for name, vals := range cases {
		t.Run(name, func(t *testing.T) {
			expected := vals[0]
			args := []fred.Def{}
			for _, arg := range vals[1:] {
				args = append(args, fred.Fixed(arg))
			}

			x := fred.Call(slice, args...).Eval(env)
			if x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}

func TestTextLengthCall(t *testing.T) {
	base := fred.Fixed(fred.Text("hello"))
	length := fred.Field(base, fred.Fixed(fred.Text("length")))
	if x := length.Eval(env); x != fred.Num("5") {
		t.Error("Unexpected", x)
	}
}

func TestTextUnexpectedField(t *testing.T) {
	base := fred.Fixed(fred.Text("hello "))
	expr := fred.Field(base, fred.Fixed(fred.Text("booya")))
	if x := expr.Eval(env); x != fred.ErrNoSuchField {
		t.Error("Unexpected", x)
	}
}
