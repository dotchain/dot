// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
)

// Substream creates a child stream
//
// The sequence of keys are used as paths into the logical value and
// work for both array indices and object-like keys.
//
// For instance: `streams.Substream(parent, 5, "count")` refers to the
// "count" field of the 5th element and any changes to it.
//
// Note that the path provided will be kept up-to-date. In the
// previous example, if 10 items were inserted into the root at index
// 0, the path would be internally updated to [15, "count"] at that
// point.  This guarantees that any updates to the substream get
// reflected at the right index of the parent stream.
//
// Substreams may "terminate" if a parent or some higher node is
// simply deleted.  Note that deleting the element referred to by the
// path itself does not cause the stream to dry up -- some element
// higher up needs to be replaced.  Dried up streams do not hold
// references to anything and so will not cause garbage collection
// issues. The only operation that such streams still permit would be
// the deletion of callbacks.
func Substream(parent Stream, key ...interface{}) Stream {
	return &substream{parent, refs.Path(key)}
}

type substream struct {
	parent Stream
	ref    refs.Ref
}

func (s *substream) Next() (Stream, changes.Change) {
	if s.ref == refs.InvalidRef {
		return nil, nil
	}

	next, nextc := s.parent.Next()
	if next == nil {
		return nil, nil
	}

	r, c := s.ref.Merge(nextc)
	if r == refs.InvalidRef {
		next = nil
	}

	return &substream{next, r}, c
}

func (s *substream) Nextf(key interface{}, fn func()) {
	if key == nil || s.ref != refs.InvalidRef {
		s.parent.Nextf(key, fn)
	}
}

func (s *substream) Append(c changes.Change) Stream {
	return s.apply(c, false)
}

func (s *substream) ReverseAppend(c changes.Change) Stream {
	return s.apply(c, true)
}

func (s *substream) apply(c changes.Change, reverse bool) Stream {
	if s.ref == refs.InvalidRef {
		return s
	}

	c = changes.PathChange{Path: []interface{}(s.ref.(refs.Path)), Change: c}
	if reverse {
		return &substream{s.parent.ReverseAppend(c), s.ref}
	}
	return &substream{s.parent.Append(c), s.ref}
}
