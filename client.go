// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"math/rand"
	"sync"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/streams"
)

// Session represents a client session
type Session struct {
	store ops.Store
	c     *ops.Connector
}

// Close closes the session
//
// The returned version and pending maybe reused to Reconnect from
// that state.
func (s *Session) Close() (version int, pending []ops.Op) {
	s.c.Disconnect()
	s.c.Async.Close()
	s.store.Close()
	return s.c.Version, s.c.Pending
}

// Connect creates a fresh session to the provided URL
func Connect(url string) (*Session, streams.Stream) {
	return Reconnect(url, -1, nil)
}

// Reconnect creates a session using saved state from a prior session
func Reconnect(url string, version int, pending []ops.Op) (*Session, streams.Stream) {
	store := &nw.Client{URL: url}
	cache := &sync.Map{}
	xformed := ops.TransformedWithCache(store, cache)
	connector := ops.NewConnector(version, pending, xformed, rand.Float64)
	stream := connector.Stream
	connector.Connect()
	return &Session{store, connector}, stream
}
