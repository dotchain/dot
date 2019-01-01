// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package collab implements a collaborative text control
package collab

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/streams"
	"golang.org/x/text/unicode/norm"
)

// Text is an immutable collaborative text edit control.
type Text struct {
	// Text is the raw text
	Text string

	// SessionID uniquely defines the current "session"
	SessionID interface{}

	// Refs contain collaborative and own cursors. Each reference
	// in this map is a refs.Range type and the key is the
	// SessionID.  A new session may not have any cursor until an
	// edit or cursor movement.
	Refs map[interface{}]refs.Ref

	// Stream is used to track changes. A trivial streams.New()
	// call can be used to create one.
	Stream streams.Stream
}

var p = refs.Path{"Value"}

// SetSelection sets the selection range for the current session
func (t Text) SetSelection(start, end int, left bool) (Text, changes.Change) {
	startc := refs.Caret{p, start, start > end || start == end && left}
	endc := refs.Caret{p, end, start < end || start == end && left}
	l, c := t.toList().UpdateRef(t.SessionID, refs.Range{startc, endc})
	return t.fromList(l, c), c
}

// Insert inserts strings at the current cursor position.  If the
// cursor is not collapsed, it replaces the selection and collapses
// the cursor)
func (t Text) Insert(s string) (Text, changes.Change) {
	offset, before := t.selection()
	after := types.S8(s)
	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := t.toList().Apply(nil, splice).(refs.Container)
	caret := refs.Caret{p, offset + after.Count(), false}
	lx, cx := l.UpdateRef(t.SessionID, refs.Range{caret, caret})
	cx = changes.ChangeSet{splice, cx}
	return t.fromList(lx, cx), cx
}

// Delete deletes the selection. In the case of a collapsed selection,
// it deletes the last character
func (t Text) Delete() (Text, changes.Change) {
	offset, before := t.selection()
	if offset == 0 && string(before) == "" {
		return t, nil
	}

	after := types.S8("")
	caret := refs.Caret{p, offset, true}

	if string(before) == "" {
		caret.Index = offset - t.PrevCharWidth(offset)
		before = types.S8(t.Text[caret.Index:offset])
		offset = caret.Index
	}

	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := t.toList()
	lx, cx := l.UpdateRef(t.SessionID, refs.Range{caret, caret})
	lx = lx.Apply(nil, splice).(refs.Container)
	cx = changes.ChangeSet{cx, splice}
	return t.fromList(lx, cx), cx
}

// ArrowLeft implements left arrow key, taking care to properly account
// for unicode sequences.
func (t Text) ArrowLeft() (Text, changes.Change) {
	_, end := t.cursor()
	end -= t.PrevCharWidth(end)
	return t.SetSelection(end, end, true)
}

// ShiftArrowLeft implements shift left arrow key, taking care to
// properly account for unicode sequences.
func (t Text) ShiftArrowLeft() (Text, changes.Change) {
	start, end := t.cursor()
	return t.SetSelection(start, end-t.PrevCharWidth(end), true)
}

// ArrowRight implements right arrow key, taking care to properly account
// for unicode sequences.
func (t Text) ArrowRight() (Text, changes.Change) {
	_, end := t.cursor()
	end += t.NextCharWidth(end)
	return t.SetSelection(end, end, false)
}

// ShiftArrowRight implements shift right arrow key, taking care to
// properly account for unicode sequences.
func (t Text) ShiftArrowRight() (Text, changes.Change) {
	start, end := t.cursor()
	return t.SetSelection(start, end+t.NextCharWidth(end), false)
}

// Copy does not change Text.  It just returns the text currently
// selectet.
func (t Text) Copy() string {
	_, sel := t.selection()
	return string(sel)
}

// Paste is like insert except it keeps the cursor around the pasted
// string.
func (t Text) Paste(s string) (Text, changes.Change) {
	offset, before := t.selection()
	after := types.S8(s)
	splice := changes.PathChange{p, changes.Splice{offset, before, after}}
	l := t.toList().Apply(nil, splice).(refs.Container)
	start := refs.Caret{p, offset, after.Count() == 0}
	end := refs.Caret{p, offset + after.Count(), true}
	lx, cx := l.UpdateRef(t.SessionID, refs.Range{start, end})
	cx = changes.ChangeSet{splice, cx}
	return t.fromList(lx, cx), cx
}

// StartOf returns the cursor index of the specified session. If the
// session does not exist, it default to zero index.
func (t Text) StartOf(sessionID interface{}) (int, bool) {
	if r, ok := t.Refs[sessionID]; ok {
		caret := r.(refs.Range).Start
		return caret.Index, caret.IsLeft
	}
	return 0, true
}

// EndOf returns the cursor index of the specified session.
func (t Text) EndOf(sessionID interface{}) (int, bool) {
	if r, ok := t.Refs[sessionID]; ok {
		caret := r.(refs.Range).End
		return caret.Index, caret.IsLeft
	}
	return 0, true
}

// Next returns the next value of data as determined from the
// stream. The boolean param is set to false if there is no next value
// or if the next value is not a valit Text.
func (t Text) Next() (Text, bool) {
	if s, c := t.Stream.Next(); s != nil {
		if result, ok := t.Apply(nil, c).(Text); ok {
			result.Stream = s
			return result, true
		}
	}
	return t, false
}

// Latest returns the latest value of the data
func (t Text) Latest() Text {
	var ok bool
	for t, ok = t.Next(); ok; t, ok = t.Next() {
	}
	return t
}

// Apply implements changes.Value. Note that this does not update the
// stream and as such should be used with care.
func (t Text) Apply(ctx changes.Context, c changes.Change) changes.Value {
	result := t.toList().Apply(ctx, c)
	if l, ok := result.(refs.Container); ok {
		text := string(l.Value.(types.S8))
		return Text{text, t.SessionID, l.Refs(), nil}
	}
	return result
}

func (t Text) cursor() (int, int) {
	if r, ok := t.Refs[t.SessionID]; ok {
		return r.(refs.Range).Start.Index, r.(refs.Range).End.Index
	}
	return 0, 0
}

func (t Text) toList() refs.Container {
	return refs.NewContainer(types.S8(t.Text), t.Refs)
}

func (t Text) fromList(l refs.Container, c changes.Change) Text {
	text := string(l.Value.(types.S8))
	return Text{text, t.SessionID, l.Refs(), t.Stream.Append(c)}
}

func (t Text) selection() (int, types.S8) {
	start, end := t.cursor()
	if start > end {
		start, end = end, start
	}
	return start, types.S8(t.Text[start:end])
}

// NextCharWidth returns the width of a user-perceived character.  This
// takes care of combining characters and such.
func (t Text) NextCharWidth(idx int) int {
	return norm.NFC.NextBoundaryInString(t.Text[idx:], true)
}

// PrevCharWidth returns the width of a user-perceived character
// before the provided index.  This takes care of combining characters
// and such.
func (t Text) PrevCharWidth(idx int) int {
	text := []byte(t.Text)[:idx]

	offset := norm.NFC.LastBoundary(text)
	if offset < 0 {
		return 0
	}

	if offset < idx || idx == 0 {
		return idx - offset
	}

	// NFC.LastBoundary is quite buggy in some cases.
	// See: https://github.com/golang/go/issues/9055
	// The work around is to brute force it in those cases
	idx = len(text) - 100
	if idx < 0 {
		idx = 0
	}
	w := len(text) - idx
	for w > 1 && norm.NFC.NextBoundary(text[idx:], true) != w {
		idx++
		w--
	}
	return w
}
