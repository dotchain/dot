// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package text_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/streams/text"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestText(t *testing.T) {
	t.Run("Use16=false", textSuite(false).Run)
	t.Run("Use16=true", textSuite(true).Run)
}

type textSuite bool

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
}

func (s textSuite) editable(txt string) *text.Editable {
	return &text.Editable{Text: txt, Use16: bool(s), Refs: map[interface{}]refs.Ref{}}
}

func (s textSuite) testTextCursors(t *testing.T) {
	e := s.editable("Hello")
	c, ex := e.SetSelection(3, 3, false)
	e = validate(t, c, e, ex)
	c, ex = e.SetSelection(3, 5, false)
	e = validate(t, c, e, ex)

	if x := e.Copy(); x != "lo" {
		t.Error("Unexpected copy failure", x)
	}

	// make start > end
	c, ex = e.SetSelection(5, 3, true)
	e = validate(t, c, e, ex)

	if x := e.Copy(); x != "lo" {
		t.Error("Unexpected copy failure", x)
	}
}

func (s textSuite) testCaretRemoteInsertion(t *testing.T) {
	e := s.editable("Hello")
	c, ex := e.SetSelection(3, 3, true)
	e = validate(t, c, e, ex)

	insert := changes.Splice{3, types.S8(""), types.S8("book")}
	if s {
		insert = changes.Splice{3, types.S16(""), types.S16("book")}
	}

	cx := changes.PathChange{[]interface{}{"Value"}, insert}
	ex = e.Apply(cx).(*text.Editable)
	e = validate(t, cx, e, ex)
	if start, _ := e.Start(); start != 3 {
		t.Error("Unexpected start", start)
	}

	_, e = e.SetSelection(3, 3, false)
	ex = e.Apply(cx).(*text.Editable)
	e = validate(t, cx, e, ex)
	if start, _ := e.Start(); start != 3+len("book") {
		t.Error("Unexpected start", start)
	}
}

func (s textSuite) testTextInsertCollapsed(t *testing.T) {
	e := s.editable("Hello")
	c, ex := e.SetSelection(3, 3, true)
	e = validate(t, c, e, ex)

	c, ex = e.Insert("<boo>")
	e = validate(t, c, e, ex)

	if e.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", e.Text)
	}

	if x, left := e.Start(); x != len("Hel<boo>") || left {
		t.Error("Unexpected start", x, left)
	}

	if x, left := e.End(); x != len("Hel<boo>") || left {
		t.Error("Unexpected end", x, left)
	}
}

func (s textSuite) testTextInsertNonCollapsed(t *testing.T) {
	e := s.editable("HelOKlo")
	c, ex := e.SetSelection(3, 5, false)
	e = validate(t, c, e, ex)

	c, ex = e.Insert("<boo>")
	e = validate(t, c, e, ex)

	if e.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", e.Text)
	}

	if x, left := e.Start(); x != len("Hel<boo>") || left {
		t.Error("Unexpected start", x, left)
	}

	if x, left := e.End(); x != len("Hel<boo>") || left {
		t.Error("Unexpected end", x, left)
	}
}

func (s textSuite) testTextPasteCollapsed(t *testing.T) {
	e := s.editable("Hello")
	c, ex := e.SetSelection(3, 3, false)
	e = validate(t, c, e, ex)

	c, ex = e.Paste("<boo>")
	e = validate(t, c, e, ex)

	if e.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", e.Text)
	}

	if x, left := e.Start(); x != len("Hel") || left {
		t.Error("Unexpected start", x, left)
	}

	if x, left := e.End(); x != len("Hel<boo>") || !left {
		t.Error("Unexpected end", x, left)
	}
}

func (s textSuite) testTextPasteNonCollapsed(t *testing.T) {
	e := s.editable("HelOKlo")
	c, ex := e.SetSelection(3, 5, false)
	e = validate(t, c, e, ex)

	c, ex = e.Paste("<boo>")
	e = validate(t, c, e, ex)

	if e.Text != "Hel<boo>lo" {
		t.Error("Unexpected insert text", e.Text)
	}

	if x, left := e.Start(); x != len("Hel") || left {
		t.Error("Unexpected start", x, left)
	}

	if x, left := e.End(); x != len("Hel<boo>") || !left {
		t.Error("Unexpected end", x, left)
	}
}

func (s textSuite) testTextDeleteCollapsed(t *testing.T) {
	e := s.editable("HelOKlo")
	c, ex := e.SetSelection(3, 3, false)
	e = validate(t, c, e, ex)

	c, ex = e.Delete()
	e = validate(t, c, e, ex)

	if e.Text != "HeOKlo" {
		t.Error("Unexpected insert text", e.Text)
	}

	if x, left := e.Start(); x != len("He") || !left {
		t.Error("Unexpected start", x, left)
	}

	if x, left := e.End(); x != len("He") || !left {
		t.Error("Unexpected end", x, left)
	}
}

func (s textSuite) testTextDeleteNonCollapsed(t *testing.T) {
	e := s.editable("HelOKlo")
	c, ex := e.SetSelection(3, 5, false)
	e = validate(t, c, e, ex)

	c, ex = e.Delete()
	e = validate(t, c, e, ex)

	if e.Text != "Hello" {
		t.Error("Unexpected insert text", e.Text)
	}

	if x, left := e.Start(); x != len("Hel") || !left {
		t.Error("Unexpected start", x, left)
	}

	if x, left := e.End(); x != len("Hel") || !left {
		t.Error("Unexpected end", x, left)
	}
}

func (s textSuite) testEmptyDelete(t *testing.T) {
	e := s.editable("HelOKlo")
	c, ex := e.SetSelection(0, 0, false)
	e = validate(t, c, e, ex)

	c, ex = e.Delete()
	validate(t, c, e, ex)

	if ex != e || c != nil {
		t.Error("Unexpected delete behavior", ex, c)
	}
}

func (s textSuite) testReplace(t *testing.T) {
	e := s.editable("HelOKlo")
	result := e.Apply(changes.Replace{e, types.S8("boo")})
	if result != types.S8("boo") {
		t.Error("Unexpected Apply reult", result)
	}
}

func (s textSuite) testCharWidths(t *testing.T) {
	e := s.editable("bròwn")
	w := e.NextCharWidth(2)
	if e.Text[2:2+w] != "ò" {
		t.Error("NextCharWidth unexpected", w)
	}
	if x := e.PrevCharWidth(2 + w); x != w {
		t.Error("PrevCharWidth unexpected", x)
	}

	if e.PrevCharWidth(0) != 0 {
		t.Error("PrevCharWidth(0)", e.PrevCharWidth(0))
	}

	// lets test out some agontek magic: a + ogonek + acute
	e = s.editable("\u0061\u0328\u0301")
	w = e.NextCharWidth(0)
	if w != len(e.Text) {
		t.Error("Unexpected char width", w)
	}
	if x := e.PrevCharWidth(w); x != w {
		t.Error("PrevCharWidth unexpected", x)
	}
}

func validate(t *testing.T, c changes.Change, before, after *text.Editable) *text.Editable {
	if !reflect.DeepEqual(before.Apply(c), after) {
		t.Fatal("change diverged", c)
	}
	if c != nil {
		reverted := after.Apply(c.Revert()).(*text.Editable)
		if reverted.Text != before.Text {
			t.Fatal("revert diverged", reverted.Text, before.Text, c)
		}
		start, left := before.Start()
		rstart, rleft := reverted.Start()
		if start != rstart || left != rleft {
			t.Fatal("revert diverged", start, left, rstart, rleft, c)
		}

		end, left := before.End()
		rend, rleft := reverted.End()
		if end != rend || left != rleft {
			t.Fatal("revert diverged", end, left, rend, rleft, c)
		}
	}
	return after
}
