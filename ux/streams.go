// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

// BoolStream holds a boolean value and tracks changes of it.  The
// value is itself almost immutable.
//
// Changes can be listened to via the Notifier and actual change values
// can be tracked by chasing the Next value (or calling Latest())
type BoolStream struct {
	// Notifier provides On/Off/Notify support.  This is carried
	// forward on all elements of the stream
	*Notifier

	// Value represents the current value
	Value bool

	// Change represents the change that resulted in the next
	// value. It can be nil in case of a change that was not
	// represented in which case callers have to simply work with
	// it.
	Change

	// Next is the next value in sequence
	Next *BoolStream
}

// Latest returns the latest value of this stream
func (s *BoolStream) Latest() *BoolStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Update updates the stream with a new value and returns the latest
// in the stream. Callers are not notified -- use Notify() to notify
func (s *BoolStream) Update(c Change, b bool) *BoolStream {
	for s.Next != nil {
		// TODO: this can have more interesting merge logic
		s = s.Next
	}
	result := &BoolStream{s.Notifier, b, c, nil}
	s.Next = result
	return result
}

// TextStream holds a text value and tracks changes of it.  The
// value is itself almost immutable.
//
// Changes can be listened to via the Notifier and actual change values
// can be tracked by chasing the Next value (or calling Latest())
type TextStream struct {
	// Notifier provides On/Off/Notify support.  This is carried
	// forward on all elements of the stream
	*Notifier

	// Value represents the current value
	Value string

	// Change represents the change that resulted in the next
	// value. It can be nil in case of a change that was not
	// represented in which case callers have to simply work with
	// it.
	Change

	// Next is the next value in sequence
	Next *TextStream
}

// Latest returns the latest value of this stream
func (s *TextStream) Latest() *TextStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Update updates the stream with a new value and returns the latest
// in the stream. Callers are not notified -- use Notify() to notify
func (s *TextStream) Update(c Change, text string) *TextStream {
	for s.Next != nil {
		// TODO: this can have more interesting merge logic
		s = s.Next
	}
	result := &TextStream{s.Notifier, text, c, nil}
	s.Next = result
	return result
}
