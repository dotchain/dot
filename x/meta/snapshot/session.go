// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package snapshot

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
	"github.com/dotchain/dot/x/meta"
)

type session struct {
	streams.Stream
	meta.Data
	gosync.Mutex
}

func (s *session) Load(ver int) (ops.Op, []ops.Op) {
	s.Lock()
	defer s.Unlock()

	return s.TransformedOp[ver], s.MergeOps[ver]
}

func (s *session) Store(ver int, op ops.Op, merge []ops.Op) {
	s.Lock()
	defer s.Unlock()

	s.Stream = s.Stream.Append(changes.ChangeSet{
		s.makeChange(nil, op, "TransformedOp", ver),
		s.makeChange(nil, merge, "MergeOps", ver),
	})
	s.TransformedOp[ver], s.MergeOps[ver] = op, merge
}

func (s *session) updateVersion(version int, pending []ops.Op) {
	s.Lock()
	defer s.Unlock()

	s.Stream = s.Stream.Append(changes.ChangeSet{
		s.makeChange(s.Version, version, "Version"),
		s.makeChange(s.Pending, pending, "Pending"),
	})
	s.Version, s.Pending = version, pending
}

func (s *session) makeChange(before, after interface{}, path ...interface{}) changes.Change {
	c := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: after}}
	if before != nil {
		c.Before = changes.Atomic{Value: before}
	}

	return changes.PathChange{Path: path, Change: c}
}

func reconnect(url string, m meta.Data) (closer func(), updates, metas streams.Stream) {
	x, merge := meta.CachedOp{}, meta.CachedOps{}
	for k, v := range m.TransformedOp {
		x[k] = v
	}
	m.TransformedOp = x
	for k, v := range m.MergeOps {
		merge[k] = v
	}
	m.MergeOps = merge

	metas = streams.New()
	s := &session{Stream: metas, Data: m}
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	store := &nw.Client{URL: url, Log: logger}
	stream, closefn := sync.Stream(
		store,
		sync.WithNotify(s.updateVersion),
		sync.WithSession(m.Version, m.Pending),
		sync.WithLog(log.New(os.Stderr, "C", log.Lshortfile|log.LstdFlags)),
		sync.WithBackoff(rand.Float64, time.Second, time.Minute),
		sync.WithAutoTransform(s),
	)

	closer = func() {
		closefn()
		store.Close()
	}

	return closer, stream, metas
}
