// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// S16 implements an UFT16 string stream
type S16 struct {
	Stream Stream
	Value  string
}

// Next returns the next if there is one.
func (s *S16) Next() (*S16, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	v, ok := types.S16(s.Value).Apply(nil, nextc).(types.S16)
	if !ok {
		return &S16{Value: s.Value}, nil
	}
	return &S16{Stream: next, Value: string(v)}, nextc
}

// Latest returns the latest non-nil entry in the stream
func (s *S16) Latest() *S16 {
	for next, _ := s.Next(); next != nil; next, _ = s.Next() {
		s = next
	}
	return s
}

// Update replaces the current value with the new value
func (s *S16) Update(val string) *S16 {
	before, after := types.S16(s.Value), types.S16(val)

	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: before, After: after})
		s = &S16{Stream: nexts, Value: val}
	}
	return s
}

// Splice replaces s[offset:offset+count] with the provided insert
// string value.
func (s *S16) Splice(offset, count int, insert string) *S16 {
	val := types.S16(s.Value)
	o := val.ToUTF16(offset)
	before := types.S16(s.Value[offset : offset+count])
	c := changes.Splice{Offset: o, Before: before, After: types.S16(insert)}
	v := string(val.Apply(nil, c).(types.S16))
	return &S16{Stream: s.Stream.Append(c), Value: v}
}

// Move moves[offset:offset+count] by the provided distance to the
// right (or if distance is negative, to the left)
func (s *S16) Move(offset, count, distance int) *S16 {
	val := types.S16(s.Value)
	o, c, d := val.ToUTF16(offset), val.ToUTF16(count), val.ToUTF16(count)
	cx := changes.Move{Offset: o, Count: c, Distance: d}
	v := string(val.Apply(nil, cx).(types.S16))
	return &S16{Stream: s.Stream.Append(cx), Value: v}
}
