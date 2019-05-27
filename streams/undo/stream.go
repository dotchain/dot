// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package undo

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// New returns a new stream with undo/redo capabilities
func New(base streams.Stream) streams.Stream {
	return stream{base, &stack{base: base}}
}

type stream struct {
	base streams.Stream
	*stack
}

func (s stream) Append(c changes.Change) streams.Stream {
	result := s
	s.stack.withLock(func() {
		result.base = s.base.Append(c)
		s.pullChanges(local)
	})
	return result
}

func (s stream) ReverseAppend(c changes.Change) streams.Stream {
	return stream{s.base.ReverseAppend(c), s.stack}
}

func (s stream) Next() (streams.Stream, changes.Change) {
	base, c := s.base.Next()
	if base != nil {
		base = stream{base, s.stack}
	}
	return base, c
}

func (s stream) Push() error {
	return s.base.Push()
}

func (s stream) Pull() error {
	return s.base.Pull()
}

func (s stream) Undo() {
	s.stack.Undo()
}

func (s stream) Redo() {
	s.stack.Redo()
}
