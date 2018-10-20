// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/idgen"
	"math/rand"
	"time"
)

// Connector helps connect a Store to a stream, taking local changes
// on the stream and writing them to the store and vice versa.
//
// The Version represents the version of the last operation received
// from the store.  Pending represents the operations that have been
// sent but not yet acknowledged by the store.
//
// The Stream can be used to make local changes as well as keep up to
// date with remote changes.  For concurrency control, the stream
// should be wrapped with streams.Async.  The convenience function
// NewConnector takes care of this book-keeping though if multiple
// stores are in play, the same Async object is recommended.
type Connector struct {
	Version int
	Pending []Op
	streams.Stream
	*streams.Async
	Store
	NewID func() interface{}
}

// NewConnector creates a new connection between the store and a
// stream. It creates an Async object as well as the  stream taking
// care to wrap the stream via Async.Wrap.
func NewConnector(version int, pending []Op, store Store) *Connector {
	async := &streams.Async{}
	s := async.Wrap(streams.New())
	store = ReliableStore(store, rand.Float64, time.Second/2, time.Minute)
	return &Connector{version, pending, s, async, store, idgen.New}
}

// Connect starts the connect process.
func (c *Connector) Connect(ctx context.Context) {
	must(c.Store.Append(context.Background(), c.Pending))

	c.Stream.Nextf(c, func() {
		var change changes.Change
		c.Stream, change = streams.Latest(c.Stream)
		op := Operation{OpID: c.NewID(), BasisID: c.Version, VerID: -1, Change: change}
		if len(c.Pending) > 0 {
			op.ParentID = c.Pending[0].ID()
		}
		c.Pending = append(c.Pending, op)
		must(c.Store.Append(context.Background(), []Op{op}))
	})
	go func() {
		c.readLoop(ctx)
		c.Stream.Nextf(c, nil)
	}()
}

func (c *Connector) readLoop(ctx context.Context) {
	limit := 1000

	for {
		ops, err := c.Store.GetSince(ctx, c.Version+1, limit)
		if ctx.Err() != nil {
			return
		}
		must(err)

		if len(ops) == 0 {
			must(c.Store.Poll(ctx, c.Version+1))
			continue
		}

		c.Async.Run(func() {
			for _, op := range ops {
				c.Version = op.Version()
				change := op.Changes()
				if len(c.Pending) > 0 && c.Pending[0].ID() == op.ID() {
					change = nil
				}
				c.Stream = c.Stream.ReverseAppend(change)
			}
		})
	}
}

func must(err error) {}
