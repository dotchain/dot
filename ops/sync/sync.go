// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
)

// Stream converts an ops Store into a stream.
//
// Calling Append() on the Stream is asynchronously appended to the
// Store and remote changes are also automatically merged into the
// returned stream.
//
// Sync works with a Transformed store.  A raw store can be auto
// transformed via the WithAutoTransform option.
//
// Additional sync options can be used to configure the behavior.
//
// The returned close function can be used to shutdown the
// synchronization but the underlying store still needs to be
// separately released.
func Stream(store ops.Store, opts ...Option) streams.Stream {
	notify := func(version int, pending, merge []ops.Op) {}
	c := &Config{Store: store, Log: log.Default(), Notify: notify, Version: -1}

	for _, opt := range opts {
		opt(c)
	}
	if c.AutoTransform {
		c.Store = ops.Transformed(c.Store, c.Cache)
	}
	c.Store = Reliable(c.Store, opts...)
	if c.NonBlocking {
		c.Store = NonBlocking(c.Store)
	}

	out := append([]ops.Op(nil), c.Pending...)
	session := &session{config: c, stream: streams.New(), out: out}

	return stream{session, session.stream}
}

type stream struct {
	*session
	stream streams.Stream
}

func (s stream) Append(c changes.Change) streams.Stream {
	return stream{s.session, s.stream.Append(c)}
}

func (s stream) ReverseAppend(c changes.Change) streams.Stream {
	return stream{s.session, s.stream.ReverseAppend(c)}
}

func (s stream) Push() error {
	return s.session.push()
}

func (s stream) Pull() error {
	return s.session.pull()
}

func (s stream) Undo() {
	s.stream.Undo()
}

func (s stream) Redo() {
	s.stream.Redo()
}

func (s stream) Next() (streams.Stream, changes.Change) {
	next, c := s.stream.Next()
	if next == nil {
		return nil, nil
	}
	return stream{s.session, next}, c
}
