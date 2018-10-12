// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
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
//     c1Next, s1Next := s1.Next()
//     c2Next, s2Next := s2.Next()
//     initialValue.Apply(c1).Apply(c1Next) == initialValue.Apply(c2).Apply(c2Next)
//
// Basically, just chasing the sequence of changes from a particular
// stream instance is guaranteed to end with the same value as any
// other stream in that family.
//
// A "family" is any stream derived from another in by means of any
// number of "Append" calls.
//
// A branch is just a derived stream with ability to Commit and Cancel
// local changes.
package streams

import "github.com/dotchain/dot/changes"

// Stream is an immutable type to manage a sequence of changes.
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
// This list can be traversed via the Next() method.  The Nextf method
// sets up a listener (or clears it) so that it can be used to listen
// for changes that have not been made yet.
//
// Concurrency
//
// Streams are generally not safe for concurrent access.
//
// Branching
//
// All changes made on a stream are propagated to the source.  It is
// possible to create git-like branches using the Branch type, where
// the changes are cached until an explicit call to Pull or Push to
// move the changes between two branches.
//
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
	Next() (changes.Change, Stream)

	// Nextf calls the provided callback whenever a next value
	// appears in the current stream. If the current stream
	// instance already has a next, the callback is called
	// immediately.
	//
	// If the fn is  nil, the listener is removed instead
	Nextf(key interface{}, fn func())
}

// New returns a new Stream
func New() Stream {
	fns := map[interface{}]func(){}
	return &stream{fns: fns}
}

type stream struct {
	c    changes.Change
	next *stream
	fns  map[interface{}]func()
}

func (s *stream) Next() (changes.Change, Stream) {
	if s.next == nil {
		return nil, nil
	}
	return s.c, s.next
}

func (s *stream) Nextf(key interface{}, fn func()) {
	if fn == nil {
		delete(s.fns, key)
	} else {
		s.fns[key] = fn
		for next := s.next; next != nil; next = next.next {
			if fn := s.fns[key]; fn != nil {
				fn()
			}
		}
	}
}

func (s *stream) Append(c changes.Change) Stream {
	return s.apply(c, false)
}

func (s *stream) ReverseAppend(c changes.Change) Stream {
	return s.apply(c, true)
}

func (s *stream) apply(c changes.Change, reverse bool) *stream {
	result := &stream{fns: s.fns}
	next := result
	for s.next != nil {
		c, next.c = s.merge(s.c, c, reverse)
		s = s.next
		next.next = &stream{fns: s.fns}
		next = next.next
	}
	s.c, s.next = c, next
	for _, fn := range s.fns {
		fn()
	}
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

// Branch represents the logical binding between a Master and Child
// streams.  changes.Changes made to the child stream are not normally
// reflected on the Master until a call to Push. Similarly, changes to
// the Master stream are not reflected on the child branch until a
// call to Pull. Merge is a combination of Push and Pull.
//
// Calling Connect() on a branch makes changes on one stream visible
// to the other immediately. Disconnect() stops the auto merging
//
// Note that unlike Stream, Branch is mutable.
//
// Concurrency
//
// Branch is not safe for concurrent access.
type Branch struct {
	Master, Local Stream
	Merging       bool
}

// Connect automerges changes between Master and Local immediately
// when they happen.  Explicit calls to Pull and Push are not
// needed. It is not safe to call Connect from within the Nextf
// callback of either Master or Local stream
func (b *Branch) Connect() {
	b.Master.Nextf(b, b.Merge)
	b.Local.Nextf(b, b.Merge)
}

// Disconnect removes the auto-emrge between Master and Local. All
// merges will have to be attempted manually after this.
func (b *Branch) Disconnect() {
	b.Master.Nextf(b, nil)
	b.Local.Nextf(b, nil)
}

func (b *Branch) merge(from, to Stream, reverse bool) (fromx, tox Stream) {
	if b.Merging {
		return from, to
	}

	b.Merging = true
	c, next := from.Next()
	for next != nil {
		if reverse {
			to = to.ReverseAppend(c)
		} else {
			to = to.Append(c)
		}
		from = next
		c, next = from.Next()
	}
	b.Merging = false
	return from, to
}

// Push updates all local changes on the Master branch
func (b *Branch) Push() {
	b.Local, b.Master = b.merge(b.Local, b.Master, false)
}

// Pull updates all master changes onto the local branch
func (b *Branch) Pull() {
	b.Master, b.Local = b.merge(b.Master, b.Local, true)
}

// Merge is shorthand for Push and Pull combined.
func (b *Branch) Merge() {
	b.Push()
	b.Pull()
}
