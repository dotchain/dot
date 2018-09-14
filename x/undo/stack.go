// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package undo

import "github.com/dotchain/dot/changes"

type cType int

const (
	local cType = iota
	upstream
	undo
	redo
)

// Stack provides Undo/Redo capability tracking all changes so this
// can be done correctly.
type Stack interface {
	Undo()
	Redo()
	Close()
}

type stack struct {
	currentType cType
	base        changes.Stream
	changes     []changes.Change
	types       []cType
}

var key = struct{}{}

func newStack(base changes.Stream) *stack {
	s := &stack{base: base}
	base.Nextf(key, func(c changes.Change, base changes.Stream) {
		s.base = base
		s.changes = append(s.changes, c)
		s.types = append(s.types, s.currentType)
	})
	return s
}

func (s *stack) changeType(newType cType, fn func()) {
	oldType := s.currentType
	s.currentType = newType
	fn()
	s.currentType = oldType
}

func (s *stack) Close() {
	s.base.Nextf(key, nil)
	s.changes = nil
	s.types = nil
}

func (s *stack) Undo() {
	if c, ok := s.getUndoChange(); ok {
		s.changeType(undo, func() { s.base.Append(c) })
	}
}

func (s *stack) Redo() {
	if c, ok := s.getRedoChange(); ok {
		s.changeType(redo, func() { s.base.Append(c) })
	}
}

func (s *stack) getUndoChange() (changes.Change, bool) {
	skipCount := 0
	l := len(s.changes) - 1
	for kk := range s.changes {
		switch s.types[l-kk] {
		case redo, local:
			if skipCount == 0 {
				return s.undoAt(l - kk), true
			}
			skipCount--
		case undo:
			skipCount++
		}
	}
	return nil, false
}

func (s *stack) getRedoChange() (changes.Change, bool) {
	skipCount := 0
	l := len(s.changes) - 1
	for kk := range s.changes {
		switch s.types[l-kk] {
		case undo:
			if skipCount == 0 {
				return s.undoAt(l - kk), true
			}
			skipCount--
		case redo:
			skipCount++
		case local:
			return nil, false
		}
	}
	return nil, false
}

func (s *stack) undoAt(offset int) changes.Change {
	c := s.changes[offset]
	if c != nil {
		c = c.Revert()
	}

	rest := s.simplify(s.changes[offset+1:], s.types[offset+1:])
	cx, _ := (changes.ChangeSet(rest)).Merge(c)
	return cx
}

// simplify remove undo/redo pairs from the sequence so as to not
// confuse the merge which is not great with some of these cases
func (s *stack) simplify(cx []changes.Change, types []cType) []changes.Change {
	var result []changes.Change
	var resultTypes []cType
	for kk, opType := range types {
		l := len(result)
		if l > 0 {
			lastOpType := resultTypes[l-1]
			cancel1 := (lastOpType == local || lastOpType == redo) && opType == undo
			cancel2 := lastOpType == undo && opType == redo
			if cancel1 || cancel2 {
				result = result[0 : l-1]
				resultTypes = resultTypes[0 : l-1]
				continue
			}
		}
		resultTypes = append(resultTypes, opType)
		result = append(result, cx[kk])
	}
	return result
}
