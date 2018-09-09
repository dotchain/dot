// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// NewStream returns a well constructed stream
func NewStream() *Stream {
	return &Stream{fns: map[interface{}]func(Change, *Stream){}}
}

// Stream represents a sequence of changes.  A change can be "applied"
// to a stream instance via the Apply method. This results in a new
// stream instance.  The old and the new stream instances can both be
// used for further changes  but they represent different states: a
// change applied on an earlier version of the stream will be
// transformed onto the latest when it is actually applied.
//
// Logically, every stream is a change made on top of another and so
// forms a open tree. But each stream instance is careful to not store
// any references to previous changes as this would cause the memory
// to constantly grow.  Instead, each stream instance maintains a
// forward list --  a list of changes that will effectively get it to
// the same converged state as any other related stream instance.
//
// Setting up an event listener via On will cause this forward list to
// be scanned.
//
// Concurrency
//
// Streams are not safe for concurrent access.
//
// Branching
//
// All changes made on a stream are immediately visible. It is
// possible to create git-like branches using the Branch type, where
// the changes are cached until an explicit call to Pull or Push to
// move the changes between two branches.
type Stream struct {
	c    Change
	next *Stream
	fns  map[interface{}]func(c Change, latest *Stream)
}

// On adds or updates a listener specified  by the key. Setting the
// callback function to nil will remove the callback
//
// The same stream instance is returned for convenience.
func (s *Stream) On(key interface{}, fn func(c Change, latest *Stream)) *Stream {
	if fn == nil {
		delete(s.fns, key)
	} else {
		s.fns[key] = fn
		for s.next != nil {
			fn(s.c, s.next)
			s = s.next
		}
	}
	return s
}

// Apply applies a change and creates a new instance. If another
// change is applied on the old instance, it is assumed that  the
// second change is not aware of the effect of the first and so has to
// be transformed when it is applied.
//
// The returned stream is a logical different "version" than the
// current stream instance.
func (s *Stream) Apply(c Change) *Stream {
	return s.apply(c, false)
}

func (s *Stream) apply(c Change, reverse bool) *Stream {
	result := &Stream{fns: s.fns}
	next := result
	for s.next != nil {
		c, next.c = s.merge(s.c, c, reverse)
		s = s.next
		next.next = &Stream{fns: s.fns}
		next = next.next
	}
	s.c, s.next = c, next
	for _, fn := range s.fns {
		fn(c, next)
	}
	return result
}

func (s *Stream) merge(left, right Change, reverse bool) (lx, rx Change) {
	if reverse {
		return swap(s.merge(right, left, false))
	}

	if left == nil {
		return right, left
	}
	return left.Merge(right)
}

// Branch represents the logical binding between a parent and child
// streams.  Changes made to the child stream are not immediately
// reflected on the Master until a call to Push. Similarly, changes to
// the Masater stream are not reflected on the child branch until a
// call to Pull.  Merge is a combination of Push and Pull
//
// Concurrency
//
// Streams are not safe for concurrent access.
type Branch struct {
	Master, Local *Stream
}

func (b *Branch) merge(from, to *Stream, reverse bool) (fromx, tox *Stream) {
	for from.next != nil {
		to = to.apply(from.c, reverse)
		from = from.next
	}
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
