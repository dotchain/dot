// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package streams defines convergent streams of changes
//
// A stream is like an event emitter or source: it tracks a sequence
// of changes on a value. It is an immutable value that logically maps
// to a "Git commit".  Appending a change to an event is equivalent to
// creating a new commit based on the previous stream.
//
// Streams differ from event emitters or immutable values in a
// fundamental way: they are convergent.  All streams from the same
// family converge to the same value
//
// For example, consider two changes on the same initial stream value.
//
//     s := ...stream...
//     s1 := s.Append(change1)
//     s2 := s.Append(change2)
//
// The two output streams converge in the following sense:
//
//     s1Next, c1Next := s1.Next()
//     s2Next, c2Next := s2.Next()
//     initial.Apply(nil, c1).Apply(nil, c1Next) == initial.Apply(nil, c2).Apply(nil, c2Next)
//
// Basically, just chasing the sequence of changes from a particular
// stream instance is guaranteed to end with the same value as any
// other stream in that family.
//
// A "family" is any stream derived from another in by means of any
// number of "Append" calls.
//
//
// Branching
//
// Streams support Git-like branching with local changes not
// automatically appearing on the parent until a call to Push.
//
// Substream
//
// It is possible to create sub-streams for elements rooted below the
// current element.  For example, one can take a stream of elements
// and only focus on the substream of changes to the 5th element. In
// this case, if the parent stream has a change which splices in a
// few elements before 5, the sub-stream should correspondingly refer
// to the new indices.   And any changes on the sub-stream should
// refer to the correct index on the parent.   The Substream() method
// provides the implementation of this concept.
//
// Value Streams
//
// Streams inherently only track the actual changes and not the
// underlying values but most applications also need to track the
// current value. See Int, Bool, S16 or S8 for an example stream that
// tracks an underlying value backed by a Stream.
//
// Custom stream implementations
//
// The dotc package (https://godoc.org/github.com/dotchain/dot/x/dotc)
// defines a mechanism to automatically generate the Stream related
// types for structs, slices and unions.
package streams

import "github.com/dotchain/dot/changes"

// Stream is an immutable type to track a sequence of changes.
//
// A change can be "applied" to a stream instance via the Append
// method. This results in a new stream instance.  The old and the new
// stream instances can both be used for further changes  but they
// represent different states: a change applied on an earlier version
// of the stream will be transformed onto the latest when it is
// actually applied.
//
// Logically, every stream is a change made on top of another and so
// forms a tree. But each stream instance is careful to not store
// any references to previous changes as this would cause the memory
// to constantly grow.  Instead, each stream instance maintains a
// forward list --  a list of changes that will effectively get it to
// the same converged state as any other related stream instance.
//
// This list can be traversed via the Next() method.
//
// Branching
//
// All changes made on a stream are propagated to the source.  It is
// possible to create git-like branches using the Branch type, where
// the changes are cached until an explicit call to Pull or Push to
// move the changes between two branches.
type Stream interface {
	// Append adds a change on top of the current change.  If
	// the current change has a Next, this is merged with the next
	// and applied to the Next instead.  That way, the change is
	// propagated all the way and applied at the end of the
	// stream.
	//
	// A listener on the stream can expect to get a change that is
	// safe to apply on top of the last change emitted.
	Append(changes.Change) Stream

	// ReverseAppend is just like Append except ReverseMerge is
	// used instead of Merge.  ReverseAppend is used to when a
	// remote change is being appended -- with the newly appended
	// change actually taking precedence over all other changes
	// that have been applied on top of the current instance.
	ReverseAppend(changes.Change) Stream

	// Next returns the change and the next stream. If no further
	// changes exist, it returns nil for both. All related stream
	// instances are guaranateed to converge -- i.e. irrespective
	// of which instance one holds, iterating over all the Next
	// values and applying them will get them all to converge  to
	// the same value.
	Next() (Stream, changes.Change)

	// Push pushes all local change up to any remote stream.
	// Does nothing if not connected to a remote stream
	Push() error

	// Pull pulls all changes from a remote stream.
	// Does nothing if not connected to a remote stream
	Pull() error

	// Undo undoes the last change on this branch
	// Does nothing if not connected to a undo stack
	Undo()

	// Redo redoes the last change on this branch.
	// Does nothing if not connected to a undo stack
	Redo()
}

// New returns a new Stream
func New() Stream {
	return &stream{}
}

type stream struct {
	c    changes.Change
	next *stream
}

func (s *stream) Next() (Stream, changes.Change) {
	if s.next == nil {
		return nil, nil
	}
	return s.next, s.c
}

func (s *stream) Append(c changes.Change) Stream {
	return s.apply(c, false)
}

func (s *stream) ReverseAppend(c changes.Change) Stream {
	return s.apply(c, true)
}

func (s *stream) apply(c changes.Change, reverse bool) *stream {
	result := &stream{}
	next := result
	for s.next != nil {
		c, next.c = s.merge(s.c, c, reverse)
		s = s.next
		next.next = &stream{}
		next = next.next
	}
	s.c, s.next = c, next
	return result
}

func (s *stream) merge(left, right changes.Change, reverse bool) (lx, rx changes.Change) {
	if reverse {
		lx, rx = s.merge(right, left, false)
		return rx, lx
	}

	if left == nil {
		return right, left
	}
	return left.Merge(right)
}

func (s *stream) Push() error {
	return nil
}

func (s *stream) Pull() error {
	return nil
}

func (s *stream) Undo() {
}

func (s *stream) Redo() {
}

// Latest returns the latest stream instance and the set of changes
// that have taken place until then
func Latest(s Stream) (Stream, changes.Change) {
	cs := changes.ChangeSet(nil)
	sx := s
	for v, c := sx.Next(); v != nil; v, c = sx.Next() {
		sx = v
		if c != nil {
			cs = append(cs, c)
		}
	}
	return sx, cs.Simplify()
}
