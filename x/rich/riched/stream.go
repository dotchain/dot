// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/rich"
)

// Stream wraps Editor with a stream.
//
// All the mutation methods on Editor are available directly on Stream
// but with the return value being the new stream instance.
//
// All the other methods of the underlying Editor are also exposed on
// Stream instances.
type Stream struct {
	streams.Stream
	*Editor
}

// NewStream creates a new stream out of a rich text value
func NewStream(r *rich.Text) *Stream {
	return &Stream{streams.New(), NewEditor(r)}
}

// Next returns the next instance in stream or nil if there isn't one
func (s *Stream) Next() *Stream {
	if n, c := s.Stream.Next(); n != nil {
		return &Stream{n, s.Editor.Apply(nil, c).(*Editor)}
	}
	return nil
}

// SetSelection updates selection state
func (s *Stream) SetSelection(focus, anchor []interface{}) *Stream {
	return s.append(s.Editor.SetSelection(focus, anchor))
}

// SetOverride update an override that is used for text insertion
func (s *Stream) SetOverride(attr rich.Attr) *Stream {
	return s.append(s.Editor.SetOverride(attr))
}

// RemoveOverride removes an override
func (s *Stream) RemoveOverride(name string) *Stream {
	return s.append(s.Editor.RemoveOverride(name))
}

// ClearOverrides removes any overrides if present
func (s *Stream) ClearOverrides() *Stream {
	return s.append(s.Editor.ClearOverrides())
}

// InsertString inserts a string at the current selection
func (s *Stream) InsertString(text string) *Stream {
	return s.append(s.Editor.InsertString(text))
}

func (s *Stream) append(c changes.Change) *Stream {
	if c == nil {
		return s
	}
	return &Stream{s.Stream.Append(c), s.Editor.Apply(nil, c).(*Editor)}
}
