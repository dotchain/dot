// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package undo

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// New returns a new stream based on another stream but with added
// ability to undo and redo actions that happened on the base stream.
//
// The return Stack internally manages a stack of changes applied to
// the stream.  Undo() and Redo() on the stack behave as expected.
//
// When using changes.Branch to merge streams, the undo stack should
// be used a little carefully.  It identifies upstream merges and
// ensures that those are not part of the undo/redo stack but this
// only works if the newly returned stream is used in the branch, like
// so:
//
//       original := streams.New()
//       upstream := streams.New()
//       downstream, stack := undo.New(original)
//       branch := changes.Branch{upstream, downstream}
//
//       ... now stack.Undo() and stack.Redo() will not undo/redo
//       ... operations made on the upstream branch but will use the
//       ... upstream changes for proper transforms
//
// The undo setup can be terminated by calling Close() on the returned
// stack. This will free up the resources associated with the stack
func New(base streams.Stream) (streams.Stream, Stack) {
	s := newStack(base)
	return stream{base, s}, s
}

type stream struct {
	base streams.Stream
	*stack
}

func (s stream) Append(c changes.Change) streams.Stream {
	return stream{s.base.Append(c), s.stack}
}

func (s stream) ReverseAppend(c changes.Change) streams.Stream {
	result := s
	s.stack.changeType(upstream, func() {
		result.base = s.base.Append(c)
	})
	return result
}

func (s stream) Next() (changes.Change, streams.Stream) {
	c, base := s.base.Next()
	if base != nil {
		base = stream{base, s.stack}
	}
	return c, base
}

func (s stream) Nextf(key interface{}, fn func()) {
	s.base.Nextf(key, fn)
}

func (s stream) Scheduler() streams.Scheduler {
	return s.base.Scheduler()
}

func (s stream) WithScheduler(sch streams.Scheduler) streams.Stream {
	return stream{s.base.WithScheduler(sch), s.stack}
}
