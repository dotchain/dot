// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package undo

import (
	"sync"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

type cType int

const (
	local cType = iota
	upstream
	undo
	redo
)

type stack struct {
	sync.Mutex
	base    streams.Stream
	changes []changes.Change
	types   []cType
}

func (s *stack) pullChanges(t cType) {
	for base, c := s.base.Next(); base != nil; base, c = s.base.Next() {
		s.base = base
		if c != nil {
			s.changes = append(s.changes, c)
			s.types = append(s.types, t)
		}
	}
}

func (s *stack) withLock(fn func()) {
	s.Lock()
	defer s.Unlock()
	s.pullChanges(upstream)
	fn()
}

func (s *stack) Undo() {
	s.withLock(func() {
		if c, ok := s.getUndoChange(); ok {
			s.base.Append(c)
			s.pullChanges(undo)
		}
	})
}

func (s *stack) Redo() {
	s.withLock(func() {
		if c, ok := s.getRedoChange(); ok {
			s.base.Append(c)
			s.pullChanges(redo)
		}
	})
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
