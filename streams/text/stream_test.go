// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package text_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/text"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestStream(t *testing.T) {
	t.Run("Use16=false", streamSuite(false).Run)
	t.Run("Use16=true", streamSuite(true).Run)
}

type streamSuite bool

func (suite streamSuite) Run(t *testing.T) {
	t.Run("Append", suite.testAppend)
	t.Run("ReverseAppend", suite.testReverseAppend)
	t.Run("CollapsedSelection", suite.testCollapsedSelection)
	t.Run("NonCollapsedSelection", suite.testNonCollapsedSelection)
	t.Run("Paste", suite.testPaste)
	t.Run("Delete", suite.testDelete)
	t.Run("WithoutOwnCursor", suite.testWithoutOwnCursor)
	t.Run("CursorAdjustment", suite.testCursorAdjustment)
}

func (suite streamSuite) testWithoutOwnCursor(t *testing.T) {
	s := text.StreamFromString("Hello", false)
	sx := &streams.ValueStream{
		types.M{"Value": types.S8("Hello")},
		s.WithoutOwnCursor(),
	}

	s2 := s.Paste("Boo")
	s2.Paste("Hoo")

	s3, c1 := sx.Next()
	s3, c2 := s3.Next()

	if c1 == nil || c2 == nil {
		t.Fatal("Unexpected issue", c1, c2)
	}

	v := s3.(*streams.ValueStream).Value
	if !reflect.DeepEqual(v, types.M{"Value": types.S8("HooHello")}) {
		t.Error("Unexpected value", v)
	}
}

func (suite streamSuite) testCursorAdjustment(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	s.Paste("boo")
	s.SetSelection(2, 2, false)

	for v, _ := s.Next(); v != nil; v, _ = s.Next() {
		s = v.(*text.Stream)
	}

	start, _ := s.E.Start()
	end, _ := s.E.End()

	if end != 5 || start != 5 || s.E.Text != "booHello" {
		t.Error("Unexpected caret", start, end, s.E.Text)
	}
}

func (suite streamSuite) testPaste(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	sx := s.Paste("boo")
	suite.validate(t, s, sx)
	if sx.E.Text != "booHello" {
		t.Error("Unexpected text", sx.E.Text)
	}
	s = sx.Paste("Hoo")
	suite.validate(t, sx, s)
	if s.E.Text != "HooHello" {
		t.Error("Unexpected text", s.E.Text)
	}
}

func (suite streamSuite) testDelete(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	s = s.SetSelection(3, 5, false)

	sx := s.Delete()
	suite.validate(t, s, sx)
	start, _ := sx.E.Start()
	end, _ := sx.E.End()
	if sx.E.Text != "Hel" || start != 3 || end != 3 {
		t.Error("Unexpected text", sx.E.Text, start, end)
	}

	// the unicode chars below = a + agontek + acute. They should
	// be treated as one grapheme cluster
	s = text.StreamFromString("\u0061\u0328\u0301", bool(suite))
	s = s.SetSelection(len(s.E.Text), len(s.E.Text), false)

	sx = s.Delete()
	suite.validate(t, s, sx)
	if sx.E.Text != "" {
		t.Error("Unexpected text", s.E.Text)
	}
}

func (suite streamSuite) testAppend(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	change := changes.PathChange{[]interface{}{"Value"}, changes.Move{0, 1, 1}}
	after := s.Append(change)
	suite.validate(t, s, after.(*text.Stream))

	sx, _ := s.Next()
	if x, _ := sx.Next(); x != nil {
		t.Error("Unexpected non-nil next", x)
	}

	after = sx.Append(changes.Replace{s.E, types.S8("okok")})
	vs, ok := after.(*streams.ValueStream)
	if !ok || vs.Value != types.S8("okok") {
		t.Error("Unexpected replace result", after)
	}
	if x, _ := sx.Next(); !reflect.DeepEqual(after, x) {
		t.Error("Unexpected divergence", x)
	}
}

func (suite streamSuite) testReverseAppend(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	change := changes.PathChange{[]interface{}{"Value"}, changes.Move{0, 1, 1}}
	after := s.ReverseAppend(change)
	suite.validate(t, s, after.(*text.Stream))

	sx, _ := s.Next()
	if x, _ := sx.Next(); x != nil {
		t.Error("Unexpected non-nil next", x)
	}

	after = sx.ReverseAppend(changes.Replace{s.E, types.S8("okok")})
	vs, ok := after.(*streams.ValueStream)
	if !ok || vs.Value != types.S8("okok") {
		t.Error("Unexpected replace result", after)
	}
	if x, _ := sx.Next(); !reflect.DeepEqual(after, x) {
		t.Error("Unexpected divergence", x)
	}
}

func (suite streamSuite) testCollapsedSelection(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))

	// test caret
	after := s.SetSelection(3, 3, false)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 3 || left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 3 || left {
		t.Error("Unexpected end", idx, left)
	}

	s = after
	after = s.SetSelection(3, 3, true)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 3 || !left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 3 || !left {
		t.Error("Unexpected end", idx, left)
	}
}

func (suite streamSuite) testNonCollapsedSelection(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))

	after := s.SetSelection(3, 5, true)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 3 || left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 5 || !left {
		t.Error("Unexpected end", idx, left)
	}

	s = after
	after = s.SetSelection(5, 3, true)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 5 || !left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 3 || left {
		t.Error("Unexpected end", idx, left)
	}
}

func (suite streamSuite) validate(t *testing.T, before, after *text.Stream) {
	if next, _ := before.Next(); !reflect.DeepEqual(next, after) {
		t.Error("Divergent change", next.(*text.Stream).E, "x", after.E)
	}
	var next streams.Stream
	before.Nextf("validate", func() {
		before.Nextf("validate", nil)
		next, _ = before.Next()
	})
	if !reflect.DeepEqual(next, after) {
		t.Error("Divergent change")
	}
}
