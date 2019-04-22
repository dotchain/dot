// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"log"
	"math/rand"
	"os"
	gosync "sync"
	"time"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/ops/sync"
	"github.com/dotchain/dot/streams"
)

// Session represents a client session
type Session struct {
	meta    streams.Stream
	close   func()
	version int
	pending []ops.Op
	x       map[int]ops.Op
	merge   map[int][]ops.Op

	gosync.Mutex
}

// Close closes the session
//
// The returned version and pending maybe reused to Reconnect from
// that state.
func (s *Session) Close() (version int, pending []ops.Op) {
	s.close()
	return s.version, s.pending
}

// Load implements the ops.Cache load interface
func (s *Session) Load(ver int) (ops.Op, []ops.Op) {
	return s.x[ver], s.merge[ver]
}

func (s *Session) makeChange(before, after interface{}, path ...interface{}) changes.Change {
	c := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: after}}
	if before != nil {
		c.Before = changes.Atomic{Value: before}
	}

	return changes.PathChange{Path: path, Change: c}
}

// Store implements the ops.Cache store interface
func (s *Session) Store(ver int, op ops.Op, merge []ops.Op) {
	s.Lock()
	defer s.Unlock()

	s.meta = s.meta.Append(changes.ChangeSet{
		s.makeChange(nil, op, "TransformedOp", ver),
		s.makeChange(nil, merge, "MergeOps", ver),
	})
	s.x[ver], s.merge[ver] = op, merge
}

// UpdateVersion updates the version/pending info
func (s *Session) UpdateVersion(version int, pending []ops.Op) {
	s.Lock()
	defer s.Unlock()

	s.meta = s.meta.Append(changes.ChangeSet{
		s.makeChange(s.version, version, "Version"),
		s.makeChange(s.pending, pending, "Pending"),
	})
	s.version, s.pending = version, pending
}

// Connect creates a fresh session to the provided URL
func Connect(url string) (*Session, streams.Stream) {
	session, updates, _ := Reconnect(url, -1, nil)
	return session, updates
}

// Reconnect creates a session using saved state from a prior session
//
// It returns a Session, the updates stream and the state stream.
//
// The Session must be closed when done at which time the current
// version and pending will be returned. That can be used to reconnect
// and create a new session.
//
// The updates streaem contains all the updates to the core stream
// while the meta stream contains info about the progress of the sync
// process itself.
//
// The meta stream can be thought of as happening on the following
// struct:
//
//    type SessionMeta struct {
//        Version int
//        Pending []ops.Op
//        TransformedOp map[int]ops.Op
//        MergeOps map[int][]ops.Op
//    }
//
// See x/meta for an example of how to use this.
func Reconnect(url string, version int, pending []ops.Op) (session *Session, updates, meta streams.Stream) {
	meta = streams.New()
	session = &Session{meta, nil, version, pending, map[int]ops.Op{}, map[int][]ops.Op{}, gosync.Mutex{}}
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	store := &nw.Client{URL: url, Log: logger, ContentType: "application/x-sjson"}
	stream, closefn := sync.Stream(
		store,
		sync.WithNotify(session.UpdateVersion),
		sync.WithSession(session.version, session.pending),
		sync.WithLog(log.New(os.Stderr, "C", log.Lshortfile|log.LstdFlags)),
		sync.WithBackoff(rand.Float64, time.Second, time.Minute),
		sync.WithAutoTransform(session),
	)
	session.close = func() {
		closefn()
		store.Close()
	}

	return session, stream, meta
}
