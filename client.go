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
}

// Close closes the session
//
// The returned version and pending maybe reused to Reconnect from
// that state.
func (s *Session) Close() (version int, pending []ops.Op) {
	s.close()
	return s.version, s.pending
}

// Connect creates a fresh session to the provided URL
func Connect(url string) (*Session, streams.Stream) {
	return Reconnect(url, -1, nil)
}

// Reconnect creates a session using saved state from a prior session
func Reconnect(url string, version int, pending []ops.Op) (*Session, streams.Stream) {
	session := &Session{version: version, pending: pending}
	store := ops.Transformed(&nw.Client{URL: url}, cache{})
	opts := []sync.Option{
		sync.WithNotify(func(version int, pending []ops.Op) {
			session.version = version
			session.pending = pending
		}),
		sync.WithSession(version, pending),
		sync.WithLog(log.New(os.Stdout, "C", log.Lshortfile|log.LstdFlags)),
		sync.WithBackoff(rand.Float64, time.Second, time.Minute),
	}
	stream, closefn := sync.Stream(store, opts...)
	session.close = func() {
		closefn()
		store.Close()
	}

	return session, stream
}

type cache map[int]interface{}

func (c cache) Load(key interface{}) (interface{}, bool) {
	v, ok := c[key.(int)]
	return v, ok
}

func (c cache) Store(key, value interface{}) {
	c[key.(int)] = value
}
