// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package fold implements a simple scheme for folding.
//
// Folding is similar to branching in that a set of changes are made
// but not pushed up to the original stream.  But unlike branching,
// folds allow more *unfolded* changes to be made on top, correctly
// transferring these upstream (after transforming them against the
// folded changes).   Similarly,  upstream changes are pulled in, also
// correctly transforming them.
//
// TODO: It is a bit tricky to remove folded changes. In an ideal
// setup, removing folds should have no effect on any streams derived
// from the fold while having the right effect upstream.  This is not
// yet implemented.
package fold

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// New returns a new stream with a "folded" change. This change is not
// applied onto the base stream but held back.  Any further changes on
// the underlying stream or the returned stream are properly proxied
// back and forth.
//
// The folded change can be fetched back by calling Unfold on the
// returned stream or any stream derived from it.
func New(c changes.Change, base streams.Stream) streams.Stream {
	return stream{c, base}
}

type stream struct {
	fold changes.Change
	streams.Stream
}

func (s stream) Append(c changes.Change) streams.Stream {
	fold := s.fold
	if fold != nil {
		fold, c = c.Merge(fold.Revert())
		if fold != nil {
			fold = fold.Revert()
		}
	}

	if c == nil {
		return stream{fold, s.Stream}
	}

	return stream{fold, s.Stream.Append(c)}
}

func (s stream) ReverseAppend(c changes.Change) streams.Stream {
	panic("Folded streams do not support ReverseAppend")
}

func (s stream) Next() (streams.Stream, changes.Change) {
	base, c := s.Stream.Next()
	if base == nil {
		return nil, nil
	}
	fold, cx := c.Merge(s.fold)
	return &stream{fold, base}, cx
}

// Unfold takes any stream derived from a folded stream (created by
// New) and returns the current state of the "change" that is folded
// as well as the modified base stream.
//
// It panics if the provided stream is not derived from New().
func Unfold(s streams.Stream) (changes.Change, streams.Stream) {
	x := s.(stream)
	return x.fold, x.Stream
}
