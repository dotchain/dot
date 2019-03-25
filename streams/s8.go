// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// S8 implements an UFT8 string stream
type S8 struct {
	Stream Stream
	Value  string
}

// Next returns the next if there is one.
func (s *S8) Next() (*S8, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	v, ok := types.S8(s.Value).Apply(nil, nextc).(types.S8)
	if !ok {
		return &S8{Value: s.Value}, nil
	}
	return &S8{Stream: next, Value: string(v)}, nextc
}

// Latest returns the latest non-nil entry in the stream
func (s *S8) Latest() *S8 {
	for next, _ := s.Next(); next != nil; next, _ = s.Next() {
		s = next
	}
	return s
}

// Update replaces the current value with the new value
func (s *S8) Update(val string) *S8 {
	before, after := types.S8(s.Value), types.S8(val)

	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: before, After: after})
		s = &S8{Stream: nexts, Value: val}
	}
	return s
}

// Splice replaces s[offset:offset+count] with the provided insert
// string value.
func (s *S8) Splice(offset, count int, insert string) *S8 {
	val := types.S8(s.Value)
	before := types.S8(s.Value[offset : offset+count])
	c := changes.Splice{Offset: offset, Before: before, After: types.S8(insert)}
	v := string(val.Apply(nil, c).(types.S8))
	return &S8{Stream: s.Stream.Append(c), Value: v}
}

// Move moves[offset:offset+count] by the provided distance to the
// right (or if distance is negative, to the left)
func (s *S8) Move(offset, count, distance int) *S8 {
	val := types.S8(s.Value)
	cx := changes.Move{Offset: offset, Count: count, Distance: distance}
	v := string(val.Apply(nil, cx).(types.S8))
	return &S8{Stream: s.Stream.Append(cx), Value: v}
}
