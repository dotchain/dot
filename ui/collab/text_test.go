// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package collab_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/ui/collab"
	"github.com/dotchain/dot/changes/types"
	"reflect"
	"testing"
)

func TestText(t *testing.T) {
	(textSuite{}).Run(t)
}

type textSuite struct{}

func (s textSuite) Run(t *testing.T) {
	t.Run("Cursors", s.testTextCursors)
	t.Run("CaretRemoteInsertion", s.testCaretRemoteInsertion)
	t.Run("InsertCollapsed", s.testTextInsertCollapsed)
	t.Run("InsertNonCollapsed", s.testTextInsertNonCollapsed)
	t.Run("PasteCollapsed", s.testTextPasteCollapsed)
	t.Run("PasteNonCollapsed", s.testTextPasteNonCollapsed)
	t.Run("DeleteCollapsed", s.testTextDeleteCollapsed)
	t.Run("DeleteNonCollapsed", s.testTextDeleteNonCollapsed)
	t.Run("EmptyDelete", s.testEmptyDelete)
	t.Run("Replace", s.testReplace)
	t.Run("CharWidths", s.testCharWidths)
	t.Run("ArrowRightLeft", s.testArrowRightLeft)
	t.Run("ShiftArrowLeft", s.testShiftArrowLeft)
	t.Run("ShiftArrowRight", s.testShiftArrowRight)
	t.Run("Latest", s.testLatest)
}

func (s textSuite) text(txt string) collab.Text {
	return collab.Text{txt, nil, map[interface{}]refs.Ref{}, streams.New()}
}

func (s textSuite) validate(t *testing.T, c changes.Change, before, after collab.Text) collab.Text {
	before.Stream = streams.New()
	before.Stream.Append(c)

	next, ok := before.Next()
	if !ok || !reflect.DeepEqual(next, after) {
		t.Fatal("change diverged", c)
	}

	next.Stream = streams.New()
	next.Stream.Append((changes.ChangeSet{c}).Revert())
	reverted, _ := next.Next()

	if reverted.Text != before.Text {
		t.Fatal("revert diverged", before.Text, reverted.Text)
	}

	start, left := before.StartOf(before.SessionID)
	rstart, rleft := reverted.StartOf(reverted.SessionID)
	if start != rstart || left != rleft {
		t.Fatal("revert diverged", start, left, rstart, rleft, c)
	}

	end, left := before.EndOf(before.SessionID)
	rend, rleft := reverted.EndOf(reverted.SessionID)
	if end != rend || left != rleft {
		t.Fatal("revert diverged", end, left, rend, rleft, c)
	}

	return after
}

func (s textSuite) testTextCursors(t *testing.T) {
	x := s.text("Hello")
	x2, c := x.SetSelection(3, 3, false)
	x = s.validate(t, c, x, x2)
	x2, c = x.SetSelection(3, 5, false)
	x = s.validate(t, c, x, x2)

	if sel := x.Copy(); sel != "lo" {
		t.Error("Unexpected copy failure", sel)
	}

	// make start > end
	x2, c = x.SetSelection(5, 3, true)
	x = s.validate(t, c, x, x2)

	if sel := x.Copy(); sel != "lo" {
		t.Error("Unexpected copy failure", sel)
	}
}

func (s textSuite) testCaretRemoteInsertion(t *testing.T) {
	insert := changes.Splice{3, types.S8(""), types.S8("book")}
	cx := changes.PathChange{[]interface{}{"Value"}, insert}

	for _, isLeft := range []bool{true, false} {
		x := s.text("Hello")
		x2, c := x.SetSelection(3, 3, isLeft)
		x = s.validate(t, c, x, x2)

		x.Stream.Append(cx)
		x2, _ = x.Next()
		x = s.validate(t, cx, x, x2)

		expected := 3
		if !isLeft {
			expected += insert.After.Count()
		}

		if start, _ := x.StartOf(x.SessionID); start != expected {
			t.Error("Unexpected start", start, expected)
		}
		if end, _ := x.EndOf(x.SessionID); end != expected {
			t.Error("Unexpected end", end, expected)
		}
	}
}

func (s textSuite) testTextInsertCollapsed(t *testing.T) {
	x := s.text("Hello")
	x2, c := x.SetSelection(3, 3, true)
	x = s.validate(t, c, x, x2)

	x2, c = x.Insert("<boo>")
	x = s.validate(t, c, x, x2)

	if x.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); start != len("Hel<boo>") || left {
		t.Error("Unexpected start", start, left)
	}

	if end, left := x.EndOf(x.SessionID); end != len("Hel<boo>") || left {
		t.Error("Unexpected end", end, left)
	}
}

func (s textSuite) testTextInsertNonCollapsed(t *testing.T) {
	x := s.text("HelOKlo")
	x2, c := x.SetSelection(3, 5, false)
	x = s.validate(t, c, x, x2)

	x2, c = x.Insert("<boo>")
	x = s.validate(t, c, x, x2)

	if x.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); start != len("Hel<boo>") || left {
		t.Error("Unexpected start", start, left)
	}

	if end, left := x.EndOf(x.SessionID); end != len("Hel<boo>") || left {
		t.Error("Unexpected end", end, left)
	}
}

func (s textSuite) testTextPasteCollapsed(t *testing.T) {
	x := s.text("Hello")
	x2, c := x.SetSelection(3, 3, false)
	x = s.validate(t, c, x, x2)

	x2, c = x.Paste("<boo>")
	x = s.validate(t, c, x, x2)

	if x.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); start != len("Hel") || left {
		t.Error("Unexpected start", start, left)
	}

	if end, left := x.EndOf(x.SessionID); end != len("Hel<boo>") || !left {
		t.Error("Unexpected end", end, left)
	}
}

func (s textSuite) testTextPasteNonCollapsed(t *testing.T) {
	x := s.text("HelOKlo")
	x2, c := x.SetSelection(3, 5, false)
	x = s.validate(t, c, x, x2)

	x2, c = x.Paste("<boo>")
	x = s.validate(t, c, x, x2)

	if x.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); start != len("Hel") || left {
		t.Error("Unexpected start", start, left)
	}

	if end, left := x.EndOf(x.SessionID); end != len("Hel<boo>") || !left {
		t.Error("Unexpected end", end, left)
	}
}

func (s textSuite) testTextDeleteCollapsed(t *testing.T) {
	x := s.text("HelOKlo")
	x2, c := x.SetSelection(3, 3, false)
	x = s.validate(t, c, x, x2)

	x2, c = x.Delete()
	x = s.validate(t, c, x, x2)

	if x.Text != "HeOKlo" {
		t.Error("Unexpected insert text", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); start != len("He") || !left {
		t.Error("Unexpected start", start, left)
	}

	if end, left := x.EndOf(x.SessionID); end != len("He") || !left {
		t.Error("Unexpected end", end, left)
	}
}

func (s textSuite) testTextDeleteNonCollapsed(t *testing.T) {
	x := s.text("HelOKlo")
	x2, c := x.SetSelection(3, 5, false)
	x = s.validate(t, c, x, x2)

	x2, c = x.Delete()
	x = s.validate(t, c, x, x2)

	if x.Text != "Hello" {
		t.Error("Unexpected insert text", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); start != len("Hel") || !left {
		t.Error("Unexpected start", start, left)
	}

	if end, left := x.EndOf(x.SessionID); end != len("Hel") || !left {
		t.Error("Unexpected end", end, left)
	}
}

func (s textSuite) testEmptyDelete(t *testing.T) {
	x := s.text("HelOKlo")
	x2, c := x.SetSelection(0, 0, false)
	x = s.validate(t, c, x, x2)

	x2, c = x.Delete()
	s.validate(t, c, x, x2)

	if !reflect.DeepEqual(x, x2) || c != nil {
		t.Error("Unexpected delete behavior", x2, c)
	}
}

func (s textSuite) testReplace(t *testing.T) {
	x := s.text("HelOKlo")
	result := x.Apply(changes.Replace{x, types.S8("boo")})
	if result != types.S8("boo") {
		t.Error("Unexpected Next", result)
	}

	x.Stream.Append(changes.Replace{x, types.S8("boo")})
	if _, ok := x.Next(); ok {
		t.Error("Unexpected Next value")
	}
}

func (s textSuite) testCharWidths(t *testing.T) {
	x := s.text("bròwn")
	w := x.NextCharWidth(2)
	if x.Text[2:2+w] != "ò" {
		t.Error("NextCharWidth unexpected", w)
	}
	if x := x.PrevCharWidth(2 + w); x != w {
		t.Error("PrevCharWidth unexpected", x)
	}

	if x.PrevCharWidth(0) != 0 {
		t.Error("PrevCharWidth(0)", x.PrevCharWidth(0))
	}

	// ensure that prev char width works with " ", "!" and "♔"
	// and agonek a = agonek + acute: "\u0061\u0328\u0301"
	for _, choice := range []string{" ", "!", "♔", "\u0061\u0328\u0301"} {
		x = s.text("x" + choice)
		w = x.NextCharWidth(1)
		if w != len(choice) {
			t.Fatal("NextCharWidth gave odd answer", choice, w)
		}
		if w2 := x.PrevCharWidth(1 + w); w2 != len(choice) {
			t.Error("PrevCharWidth gave odd answer", choice, w, w2)
		}
	}
}

func (s textSuite) testArrowRightLeft(t *testing.T) {
	// lets test out some agontek magic: a + ogonek + acute
	x := s.text("\u0061\u0328\u0301")
	after, c := x.ArrowRight()
	s.validate(t, c, x, after)
	if x.Text != after.Text {
		t.Error("Unexpected value change", after.Text)
	}

	if start, left := after.StartOf(x.SessionID); left || start != len(x.Text) {
		t.Error("Unexpected start value", left, start)
	}

	if end, left := after.EndOf(x.SessionID); left || end != len(x.Text) {
		t.Error("Unexpected end value", left, end)
	}

	x, c = after.ArrowLeft()
	s.validate(t, c, after, x)

	if x.Text != after.Text {
		t.Error("Unexpected value change", x.Text)
	}

	if start, left := x.StartOf(x.SessionID); !left || start != 0 {
		t.Error("Unexpected start value", left, start)
	}

	if end, left := x.EndOf(x.SessionID); !left || end != 0 {
		t.Error("Unexpected end value", left, end)
	}
}

func (s textSuite) testShiftArrowLeft(t *testing.T) {
	// lets test out some agontek magic: a + ogonek + acute
	x := s.text("\u0061\u0328\u0301")
	x, _ = x.SetSelection(len(x.Text), len(x.Text), false)
	after, c := x.ShiftArrowLeft()
	s.validate(t, c, x, after)
	if x.Text != after.Text {
		t.Error("Unexpected value change", after.Text)
	}

	if start, left := after.StartOf(x.SessionID); !left || start != len(x.Text) {
		t.Error("Unexpected start value", left, start)
	}

	if end, left := after.EndOf(x.SessionID); left || end != 0 {
		t.Error("Unexpected end value", left, end)
	}
}

func (s textSuite) testShiftArrowRight(t *testing.T) {
	// lets test out some agontek magic: a + ogonek + acute
	x := s.text("\u0061\u0328\u0301")
	after, c := x.ShiftArrowRight()
	s.validate(t, c, x, after)

	if x.Text != after.Text {
		t.Error("Unexpected value change", after.Text)
	}

	if start, left := after.StartOf(x.SessionID); left || start != 0 {
		t.Error("Unexpected start value", left, start)
	}

	if end, left := after.EndOf(x.SessionID); !left || end != len(x.Text) {
		t.Error("Unexpected end value", left, end)
	}
}

func (s textSuite) testLatest(t *testing.T) {
	p := func(c changes.Change) changes.Change {
		return changes.PathChange{[]interface{}{"Value"}, c}
	}
	x := s.text("hello")
	x.Stream.
		Append(p(changes.Splice{5, types.S8(""), types.S8(" ")})).
		Append(p(changes.Splice{6, types.S8(""), types.S8("world")}))

	latest := x.Latest()
	if latest.Text != "hello world" {
		t.Error("Latest unexpected", latest.Text)
	}
}
