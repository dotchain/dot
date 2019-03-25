// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import "github.com/dotchain/dot/changes"

// Bool implements a bool stream.
type Bool struct {
	Stream Stream
	Value  bool
}

// Next returns the next if there is one.
func (s *Bool) Next() (*Bool, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	v := s.Value
	val, ok := (changes.Atomic{Value: v}).Apply(nil, nextc).(changes.Atomic)
	if ok {
		v, ok = val.Value.(bool)
	}
	if !ok {
		next = nil
		v = s.Value
		nextc = nil
	}
	return &Bool{Stream: next, Value: v}, nextc
}

// Latest returns the latest non-nil entry in the stream
func (s *Bool) Latest() *Bool {
	for next, _ := s.Next(); next != nil; next, _ = s.Next() {
		s = next
	}
	return s
}

// Update replaces the current value with the new value
func (s *Bool) Update(val bool) *Bool {
	before, after := changes.Atomic{Value: s.Value}, changes.Atomic{Value: val}

	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: before, After: after})
		s = &Bool{Stream: nexts, Value: val}
	}
	return s
}
