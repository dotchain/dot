// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package text implements editable text streams
package text

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/x/types"
)

// Editable implements text editing functionality.  The main state
// maintained by Editable is the actual Text, the current location of
// the cursor and a set of selections that can be maintained with the
// text.
//
// Editable is an immutable type.  All mutations return a
// change.Change and the updated value
type Editable struct {
	Text   string
	Cursor refs.Range
	Refs   map[interface{}]refs.Ref
	Use16  bool

	// atomic is not used, just there to provide the Count/Slice methods
	changes.Atomic
}

var p = refs.Path{"Value"}

// SetCaret sets the cursor to a specific index.
func (e *Editable) SetCaret(idx int) (changes.Change, *Editable) {
	idx = e.toValueOffset(idx)
	left := idx == 0
	r := refs.Range{refs.Caret{p, idx, left}, refs.Caret{p, idx, left}}
	c, l := e.toList().Update(e, r)
	return c, e.fromList(l)
}

// SetStart sets the start index.
func (e *Editable) SetStart(idx int, left bool) (changes.Change, *Editable) {
	idx = e.toValueOffset(idx)
	r := refs.Range{refs.Caret{p, idx, left}, e.cursor().End}
	c, l := e.toList().Update(e, r)
	return c, e.fromList(l)
}

// SetEnd sets the end index.
func (e *Editable) SetEnd(idx int, left bool) (changes.Change, *Editable) {
	idx = e.toValueOffset(idx)
	r := refs.Range{e.cursor().Start, refs.Caret{p, idx, left}}
	c, l := e.toList().Update(e, r)
	return c, e.fromList(l)
}

// Insert inserts strings at the current cursor position.  If the
// cursor is not collapsed, it collapses the cursor)
func (e *Editable) Insert(s string) (changes.Change, *Editable) {
	offset, before := e.selection()
	after := e.stringToValue(s)
	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := e.toList().Apply(splice).(refs.List)
	caret := refs.Caret{p, offset + after.Count(), false}
	cx, lx := l.Update(e, refs.Range{caret, caret})
	return changes.ChangeSet{splice, cx}, e.fromList(lx)
}

// Delete deletes the selection. In the case of a collapsed selection,
// it deletes the last character
func (e *Editable) Delete() (changes.Change, *Editable) {
	offset, before := e.selection()
	if offset == 0 && before.Count() == 0 {
		return nil, e
	}

	after := before.Slice(0, 0)
	caret := refs.Caret{p, offset, true}

	if before.Count() == 0 {
		// TODO: index-- is incorrect. Take care of UTF8
		// encoding shit and find the right size
		caret.Index--
		before = e.stringToValue(e.Text).Slice(caret.Index, caret.Index+1)
	}

	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := e.toList().Apply(splice).(refs.List)
	cx, lx := l.Update(e, refs.Range{caret, caret})
	return changes.ChangeSet{splice, cx}, e.fromList(lx)
}

// Copy does not change editable.  It just returns the text currently
// selected.
func (e *Editable) Copy() string {
	_, sel := e.selection()
	return e.valueToString(sel)
}

// Start returns the cursor index
func (e *Editable) Start() (int, bool) {
	caret := e.cursor().Start
	return e.fromValueOffset(caret.Index), caret.IsLeft
}

// End returns the cursor end
func (e *Editable) End() (int, bool) {
	caret := e.cursor().End
	return e.fromValueOffset(caret.Index), caret.IsLeft
}

// Paste is like insert except it keeps the cursor around the pasted
// string.
func (e *Editable) Paste(s string) (changes.Change, *Editable) {
	offset, before := e.selection()
	after := e.stringToValue(s)
	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := e.toList().Apply(splice).(refs.List)
	start := refs.Caret{p, offset, after.Count() == 0}
	end := refs.Caret{p, offset + after.Count(), true}
	cx, lx := l.Update(e, refs.Range{start, end})
	return changes.ChangeSet{splice, cx}, e.fromList(lx)
}

// Apply implements the changes.Value interface
func (e *Editable) Apply(c changes.Change) changes.Value {
	result := e.toList().Apply(c)
	l, ok := result.(refs.List)
	if !ok {
		return result
	}

	return e.fromList(l)
}

func (e *Editable) stringToValue(s string) changes.Value {
	if e.Use16 {
		return types.S16(s)
	}
	return types.S8(s)
}

func (e *Editable) valueToString(v changes.Value) string {
	if e.Use16 {
		return string(v.(types.S16))
	}
	return string(v.(types.S8))
}

func (e *Editable) cursor() refs.Range {
	c := e.Cursor
	c.Start.Path = p
	c.End.Path = p
	return c
}

func (e *Editable) toValueOffset(idx int) int {
	if e.Use16 {
		return types.S16(e.Text).ToUTF16(idx)
	}
	// validate that the offset works
	_ = e.Text[idx:]
	return idx
}

func (e *Editable) fromValueOffset(idx int) int {
	if e.Use16 {
		return types.S16(e.Text).FromUTF16(idx)
	}
	// validate that the offset works
	_ = e.Text[idx:]
	return idx
}

func (e *Editable) toList() refs.List {
	l := refs.List{e.stringToValue(e.Text), e.Refs}
	_, l = l.Add(e, e.cursor())
	return l
}

func (e *Editable) fromList(l refs.List) *Editable {
	text := e.valueToString(l.V)
	cursor := l.R[e].(refs.Range)
	delete(l.R, e)
	return &Editable{text, cursor, l.R, e.Use16, changes.Atomic{nil}}
}

func (e *Editable) selection() (int, changes.Value) {
	c := e.cursor()
	v := e.stringToValue(e.Text)
	start, end := c.Start.Index, c.End.Index
	diff := end - start
	if start > end {
		start, end = end, start
		diff = end - start
	}
	return start, v.Slice(start, diff)
}
