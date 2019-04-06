// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/ops/sync"
	"github.com/dotchain/dot/streams"
)

// Session represents a client session
type Session struct {
	close   func()
	version int
	pending []ops.Op
	x       map[int]ops.Op
	merge   map[int][]ops.Op
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

// Store implements the ops.Cache store interface
func (s *Session) Store(ver int, op ops.Op, merge []ops.Op) {
	s.x[ver] = op
	s.merge[ver] = merge
}

// UpdateVersion updates the version/pending info
func (s *Session) UpdateVersion(version int, pending []ops.Op) {
	s.version, s.pending = version, pending
}

// Connect creates a fresh session to the provided URL
func Connect(url string) (*Session, streams.Stream) {
	return Reconnect(url, -1, nil)
}

// Reconnect creates a session using saved state from a prior session
func Reconnect(url string, version int, pending []ops.Op) (*Session, streams.Stream) {
	session := &Session{nil, version, pending, map[int]ops.Op{}, map[int][]ops.Op{}}
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	store := &nw.Client{URL: url, Log: logger}
	stream, closefn := sync.Stream(
		store,
		sync.WithNotify(session.UpdateVersion),
		sync.WithSession(version, pending),
		sync.WithLog(log.New(os.Stderr, "C", log.Lshortfile|log.LstdFlags)),
		sync.WithBackoff(rand.Float64, time.Second, time.Minute),
		sync.WithAutoTransform(session),
	)
	session.close = func() {
		closefn()
		store.Close()
	}

	return session, stream
}
