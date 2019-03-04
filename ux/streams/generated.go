// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
//

package streams

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// BoolStream is a stream of bool values.
type BoolStream struct {
	// Notifier provides On/Off/Notify support. New instances of
	// BoolStream created via the AppendLocal or AppendRemote
	// share the same Notifier value.
	*streams.Notifier

	// Value holds the current value. The latest value may be
	// fetched via the Latest() method.
	Value bool

	// Change tracks the change that leads to the next value.
	Change changes.Change

	// Next tracks the next value in the stream.
	Next *BoolStream
}

// NewBoolStream creates a new bool stream
func NewBoolStream(value bool) *BoolStream {
	return &BoolStream{&streams.Notifier{}, value, nil, nil}
}

// Latest returns the latest value in the stream
func (s *BoolStream) Latest() *BoolStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Append appends a local change. isLocal identifies if the caller is
// local or remote. It returns the updated stream whose value matches
// the provided value and whose Latest() converges to the latest of
// the stream.
func (s *BoolStream) Append(c changes.Change, value bool, isLocal bool) *BoolStream {
	if c == nil {
		c = changes.Replace{Before: s.wrapValue(s.Value), After: s.wrapValue(value)}
	}

	// return value: after is correctly set to provided value
	result := &BoolStream{Notifier: s.Notifier, Value: value}

	// before tracks s, after tracks result, v tracks latest value
	// of after chain
	before := s
	var v changes.Value = changes.Atomic{value}

	// walk the chain of Next and find corresponding values to
	// add to after so that both s annd after converge
	after := result
	for ; before.Next != nil; before = before.Next {
		var afterChange changes.Change

		if isLocal {
			c, afterChange = before.Change.Merge(c)
		} else {
			afterChange, c = c.Merge(before.Change)
		}

		if c == nil {
			// the convergence point is before.Next
			after.Change, after.Next = afterChange, before.Next
			return result
		}

		if afterChange == nil {
			continue
		}

		// append this to after and continue with that
		v = v.Apply(nil, afterChange)
		after.Change = afterChange
		after.Next = &BoolStream{Notifier: s.Notifier, Value: s.unwrapValue(v)}
		after = after.Next
	}

	// append the residual change (c) to converge to wherever
	// after has landed. Notify since s.Latest() has now changed
	before.Change, before.Next = c, after
	s.Notify()
	return result
}

func (s *BoolStream) wrapValue(i interface{}) changes.Value {
	if x, ok := i.(changes.Value); ok {
		return x
	}
	return changes.Atomic{i}
}

func (s *BoolStream) unwrapValue(v changes.Value) bool {
	if x, ok := v.(interface{}).(bool); ok {
		return x
	}
	return v.(changes.Atomic).Value.(bool)
}

// TextStream is a stream of string values.
type TextStream struct {
	// Notifier provides On/Off/Notify support. New instances of
	// TextStream created via the AppendLocal or AppendRemote
	// share the same Notifier value.
	*streams.Notifier

	// Value holds the current value. The latest value may be
	// fetched via the Latest() method.
	Value string

	// Change tracks the change that leads to the next value.
	Change changes.Change

	// Next tracks the next value in the stream.
	Next *TextStream
}

// NewTextStream creates a new string stream
func NewTextStream(value string) *TextStream {
	return &TextStream{&streams.Notifier{}, value, nil, nil}
}

// Latest returns the latest value in the stream
func (s *TextStream) Latest() *TextStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Append appends a local change. isLocal identifies if the caller is
// local or remote. It returns the updated stream whose value matches
// the provided value and whose Latest() converges to the latest of
// the stream.
func (s *TextStream) Append(c changes.Change, value string, isLocal bool) *TextStream {
	if c == nil {
		c = changes.Replace{Before: s.wrapValue(s.Value), After: s.wrapValue(value)}
	}

	// return value: after is correctly set to provided value
	result := &TextStream{Notifier: s.Notifier, Value: value}

	// before tracks s, after tracks result, v tracks latest value
	// of after chain
	before := s
	var v changes.Value = changes.Atomic{value}

	// walk the chain of Next and find corresponding values to
	// add to after so that both s annd after converge
	after := result
	for ; before.Next != nil; before = before.Next {
		var afterChange changes.Change

		if isLocal {
			c, afterChange = before.Change.Merge(c)
		} else {
			afterChange, c = c.Merge(before.Change)
		}

		if c == nil {
			// the convergence point is before.Next
			after.Change, after.Next = afterChange, before.Next
			return result
		}

		if afterChange == nil {
			continue
		}

		// append this to after and continue with that
		v = v.Apply(nil, afterChange)
		after.Change = afterChange
		after.Next = &TextStream{Notifier: s.Notifier, Value: s.unwrapValue(v)}
		after = after.Next
	}

	// append the residual change (c) to converge to wherever
	// after has landed. Notify since s.Latest() has now changed
	before.Change, before.Next = c, after
	s.Notify()
	return result
}

func (s *TextStream) wrapValue(i interface{}) changes.Value {
	if x, ok := i.(changes.Value); ok {
		return x
	}
	return changes.Atomic{i}
}

func (s *TextStream) unwrapValue(v changes.Value) string {
	if x, ok := v.(interface{}).(string); ok {
		return x
	}
	return v.(changes.Atomic).Value.(string)
}
