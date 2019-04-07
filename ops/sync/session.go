// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
)

type session struct {
	sync.Mutex
	config *Config
	stream streams.Stream
	id     interface{}
}

func (s *session) read(ctx context.Context) {
	var err error
	var operations []ops.Op

	c, ver := s.config, s.config.Version
	for err == nil {
		operations, err = c.Store.GetSince(ctx, ver+1, 1000)
		if err != nil {
			break
		}

		if len(operations) > 0 {
			ver = s.onStoreOps(operations)
		}
	}

	if err == ctx.Err() {
		err = nil
	}

	s.must(err, "reliable GetSince failed")
}

func (s *session) onStoreOps(operations []ops.Op) int {
	s.Lock()
	defer s.Unlock()

	s.stream.Nextf(s, nil)
	defer s.stream.Nextf(s, s.onAppend)

	c := s.config
	before := c.Version
	for _, op := range operations {
		if op.Version() != c.Version+1 {
			s.must(verMismatchError{op.Version(), c.Version + 1}, "")
		}

		c.Version++
		if len(c.Pending) > 0 && c.Pending[0].ID() == op.ID() {
			c.Pending = c.Pending[1:]
			s.stream, _ = s.stream.Next()
		} else {
			s.stream = s.stream.ReverseAppend(op.Changes())
		}
	}

	if c.Version > before {
		c.Notify(c.Version, c.Pending)
	}
	return c.Version
}

func (s *session) onAppend() {
	s.Lock()
	defer s.Unlock()

	stream, c := s.stream, s.config
	for range c.Pending {
		stream, _ = stream.Next()
	}

	before := len(c.Pending)
	for next, ch := stream.Next(); next != nil; next, ch = next.Next() {
		op := ops.Operation{OpID: s.newID(), BasisID: c.Version, Change: ch}
		if len(c.Pending) > 0 {
			op.ParentID = c.Pending[len(c.Pending)-1].ID()
		}
		c.Pending = append(c.Pending, op)
	}

	if len(c.Pending) > before {
		c.Notify(c.Version, c.Pending)
		s.write(c.Pending[before:len(c.Pending)])
	}
}

func (s *session) write(pending []ops.Op) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := s.config.Store.Append(ctx, pending)
	s.must(err, "reliable append failed")
}

func (s *session) must(err error, msg string) {
	if err != nil {
		s.config.Log.Fatal(msg, err)
	}
}

type verMismatchError struct {
	got, expected int
}

func (v verMismatchError) Error() string {
	exp, got := strconv.Itoa(v.expected), strconv.Itoa(v.got)
	return "version mismatched: " + got + " (expected " + exp + ")"
}
