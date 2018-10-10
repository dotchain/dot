// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package text implements editable text streams
package text

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/x/types"
	"golang.org/x/text/unicode/norm"
)

var own = &struct{}{}

// Editable implements text editing functionality.  The main state
// maintained by Editable is the actual Text, the current location of
// the cursor and a set of selections that can be maintained with the
// text.
//
// Editable is an immutable type.  All mutations return a
// change.Change and the updated value
//
// There are two positions for each index: left or right. This is
// relevant when considering text that has wrapped around. The
// index in the text where wrapping occurs has two different positions
// on the screen: at the end of the line before wrapping and at the
// start of the line after wrapping.  The top position is considered
// "left" and the bottom line position is considered "right".
//
// There is another consideration: when a remote change causes an
// insertion at exactly the index of the cursor/caret, the caret can
// either be left alone or the caret can be pushed to the right by the
// inserted text.  The "left" position and "right" position match the
// two behaviors (respectively)
type Editable struct {
	Text   string
	Cursor refs.Range
	Refs   map[interface{}]refs.Ref
	Use16  bool

	// atomic is not used, just there to provide the Count/Slice methods
	changes.Atomic
}

var p = refs.Path{"Value"}

func (e *Editable) prioritize(c changes.Change, l refs.List) (changes.Change, refs.List) {
	return prioritized{c}, l
}

// SetSelection sets the selection range for text.
func (e *Editable) SetSelection(start, end int, left bool) (changes.Change, *Editable) {
	start, end = e.toValueOffset(start), e.toValueOffset(end)
	startx := refs.Caret{p, start, start > end || start == end && left}
	endx := refs.Caret{p, end, start < end || start == end && left}
	c, l := e.prioritize(e.toList().Update(own, refs.Range{startx, endx}))
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
	cx, lx := e.prioritize(l.Update(own, refs.Range{caret, caret}))
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
		idx := e.fromValueOffset(offset)
		idx -= e.PrevCharWidth(idx)
		caret.Index = e.toValueOffset(idx)
		before = e.stringToValue(e.Text).Slice(caret.Index, offset-caret.Index)
		offset = caret.Index
	}

	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := e.toList()
	cx, lx := e.prioritize(l.Update(own, refs.Range{caret, caret}))
	lx = lx.Apply(splice).(refs.List)
	return changes.ChangeSet{cx, splice}, e.fromList(lx)
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
	cx, lx := e.prioritize(l.Update(own, refs.Range{start, end}))
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
	_, l = l.Add(own, e.cursor())
	return l
}

func (e *Editable) fromList(l refs.List) *Editable {
	text := e.valueToString(l.V)
	cursor := l.R[own].(refs.Range)
	delete(l.R, own)
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

// NextCharWidth returns the width of a user-perceived character.  This
// takes care of combining characters and such.
func (e *Editable) NextCharWidth(idx int) int {
	return norm.NFC.NextBoundaryInString(e.Text[idx:], true)
}

// PrevCharWidth returns the width of a user-perceived character
// before the provided index.  This takes care of combining characters
// and such.
func (e *Editable) PrevCharWidth(idx int) int {
	offset := norm.NFC.LastBoundary([]byte(e.Text[:idx]))
	if offset < 0 {
		return 0
	}
	return idx - offset
}

// prioritized inverts the priority of operations
type prioritized struct {
	changes.Change
}

func (p prioritized) Merge(other changes.Change) (otherx, cx changes.Change) {
	own := p.Change
	if other != nil {
		own, other = other.Merge(own)
	}

	return other, own
}

func (p prioritized) Revert() changes.Change {
	return prioritized{p.Change.Revert()}
}

func (p prioritized) ReverseMerge(other changes.Change) (otherx, cx changes.Change) {
	own := p.Change
	if other != nil {
		other, own = own.Merge(other)
	}

	return other, own
}

func (p prioritized) ApplyTo(v changes.Value) changes.Value {
	return v.Apply(p.Change)
}

func (p prioritized) MergePath(path []interface{}) *refs.MergeResult {
	return refs.Merge(path, p.Change)
}