// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import "github.com/dotchain/dot/changes"

// Int implements an Int stream.
type Int struct {
	Stream Stream
	Value  int
}

// Next returns the next if there is one.
func (s *Int) Next() (*Int, changes.Change) {
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
		v, ok = val.Value.(int)
	}
	if !ok {
		next = nil
		v = s.Value
	}
	return &Int{Stream: next, Value: v}, nil
}

// Latest returns the latest non-nil entry in the stream
func (s *Int) Latest() *Int {
	for next, _ := s.Next(); next != nil; next, _ = s.Next() {
		s = next
	}
	return s
}

// Update replaces the current value with the new value
func (s *Int) Update(val int) *Int {
	before, after := changes.Atomic{Value: s.Value}, changes.Atomic{Value: val}

	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: before, After: after})
		s = &Int{Stream: nexts, Value: val}
	}
	return s
}
