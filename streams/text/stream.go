// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package text

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// StreamFromString constructs a new text stream
func StreamFromString(initialText string, use16 bool) *Stream {
	return &Stream{&Editable{Text: initialText}, streams.New()}
}

// Stream implements the streams.Stream interface on top of Editable.
type Stream struct {
	E *Editable
	S streams.Stream
}

// Append implements streams.Stream:Append
func (s *Stream) Append(c changes.Change) streams.Stream {
	v := s.E.Apply(c)
	sx := s.S.Append(c)
	if e, ok := v.(*Editable); ok {
		return &Stream{e, sx}
	}
	return &streams.ValueStream{v, sx}
}

// ReverseAppend implements streams.Stream:ReverseAppend
func (s *Stream) ReverseAppend(c changes.Change) streams.Stream {
	v := s.E.Apply(c)
	sx := s.S.ReverseAppend(c)
	if e, ok := v.(*Editable); ok {
		return &Stream{e, sx}
	}
	return &streams.ValueStream{v, sx}
}

// Scheduler implements streams.Stream:Scheduler
func (s *Stream) Scheduler() streams.Scheduler {
	return s.S.Scheduler()
}

// WithScheduler implements streams.Stream:WithScheduler
func (s *Stream) WithScheduler(sch streams.Scheduler) streams.Stream {
	return &Stream{s.E, s.S.WithScheduler(sch)}
}

// Next implements streams.Stream.Next
func (s *Stream) Next() (changes.Change, streams.Stream) {
	return s.mapChangeValue(s.S.Next())
}

// Nextf implements streams.Stream.Nextf
func (s *Stream) Nextf(key interface{}, fn func(c changes.Change, str streams.Stream)) {
	if fn == nil {
		s.S.Nextf(key, fn)
	} else {
		s.S.Nextf(key, func(c changes.Change, str streams.Stream) {
			fn(s.mapChangeValue(c, str))
		})
	}
}

// SetSelection sets the selection range for text
func (s *Stream) SetSelection(start, end int, left bool) *Stream {
	c1, e1 := s.E.SetStart(start, start > end || start == end && left)
	c2, e2 := e1.SetEnd(end, start < end || start == end && left)
	return &Stream{e2, s.S.Append(changes.ChangeSet{c1, c2})}
}

func (s Stream) mapChangeValue(c changes.Change, str streams.Stream) (changes.Change, streams.Stream) {
	if str == nil {
		return c, str
	}

	v := s.E.Apply(c)
	if e, ok := v.(*Editable); ok {
		return c, &Stream{e, str}
	}
	return c, &streams.ValueStream{v, str}
}
