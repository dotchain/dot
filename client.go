// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"log"
	"math/rand"
	"os"
	"time"

	dotlog "github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/ops/sync"
	"github.com/dotchain/dot/streams"
)

// Session represents a client session
type Session struct {
	Version        int
	Pending, Merge []ops.Op

	OpCache    map[int]ops.Op
	MergeCache map[int][]ops.Op
}

// NewSession creates an empty session
func NewSession() *Session {
	return &Session{
		Version:    -1,
		OpCache:    map[int]ops.Op{},
		MergeCache: map[int][]ops.Op{},
	}
}

// Stream returns the stream of changes for this session
//
// The returned store can be used to *close* the stream when needed
//
// Actual syncing of messages happens when Push and Pull are called on the stream
func (s *Session) Stream(url string, logger dotlog.Log) (streams.Stream, ops.Store) {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	}

	store := &nw.Client{
		URL:         url,
		Log:         logger,
		ContentType: "application/x-sjson",
	}

	stream := sync.Stream(
		store,
		sync.WithNotify(s.UpdateVersion),
		sync.WithSession(s.Version, s.Pending, s.Merge),
		sync.WithLog(logger),
		sync.WithBackoff(rand.Float64, time.Second, time.Minute),
		sync.WithAutoTransform(s),
	)
	return stream, store
}

// NonBlockingStream returns the stream of changes for this session
//
// The returned store can be used to *close* the stream when needed
//
// Actual syncing of messages happens when Push and Pull are called on
// the stream. Pull() does the server-fetch asynchronously, returning
// immediately if there is no server data available.
func (s *Session) NonBlockingStream(url string, logger dotlog.Log) (streams.Stream, ops.Store) {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	}

	store := &nw.Client{
		URL:         url,
		Log:         logger,
		ContentType: "application/x-sjson",
	}

	stream := sync.Stream(
		store,
		sync.WithNotify(s.UpdateVersion),
		sync.WithSession(s.Version, s.Pending, s.Merge),
		sync.WithLog(logger),
		sync.WithBackoff(rand.Float64, time.Second, time.Minute),
		sync.WithAutoTransform(s),
		sync.WithNonBlocking(true),
	)
	return stream, store
}

// Load implements the ops.Cache load interface
func (s *Session) Load(ver int) (ops.Op, []ops.Op) {
	return s.OpCache[ver], s.MergeCache[ver]
}

// Store implements the ops.Cache store interface
func (s *Session) Store(ver int, op ops.Op, merge []ops.Op) {
	s.OpCache[ver], s.MergeCache[ver] = op, merge
}

// UpdateVersion updates the version/pending info
func (s *Session) UpdateVersion(version int, pending, merge []ops.Op) {
	s.Version, s.Pending, s.Merge = version, pending, merge
}
