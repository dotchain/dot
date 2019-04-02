// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"sync"
	"time"

	"github.com/dotchain/dot/streams"
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
	close func()
	sync.Mutex
}

// NewConnector creates a new connection between the store and a
// stream. It creates an Async object as well as the  stream taking
// care to wrap the stream via Async.Wrap.
func NewConnector(version int, pending []Op, store Store, rand func() float64) *Connector {
	async := streams.NewAsync(0)
	s := async.Wrap(streams.New())
	async.LoopForever()
	store = ReliableStore(store, rand, time.Second/2, time.Minute)
	return &Connector{Version: version, Pending: pending, Stream: s, Async: async, Store: store}
}

// Connect starts the synchronization process.
func (c *Connector) Connect() {
	ctx, cancel := context.WithCancel(context.Background())
	closed := make(chan struct{})
	c.close = func() {
		cancel()
		<-closed
	}

	must(c.Store.Append(ctx, c.Pending))
	c.Stream.Nextf(c, func() { c.write(ctx) })

	go func() {
		c.readLoop(ctx)
		c.Stream.Nextf(c, nil)
		close(closed)
	}()
}

// Disconnect stops the synchronization process.  The version and
// pending are updated to the latest values when the call returns
func (c *Connector) Disconnect() {
	c.close()
}

// write takes any unwritten changes from c.Stream and writes it out
// to the ops store. note that write does not update c.Stream as
// c.Stream tracks the last upstream version
func (c *Connector) write(ctx context.Context) {
	var idx int
	var ops []Op

	c.Lock()
	for next, ch := c.Stream.Next(); next != nil; next, ch = next.Next() {
		if idx >= len(c.Pending) {
			op := Operation{OpID: NewID(), BasisID: c.Version, Change: ch}
			if len(c.Pending) > 0 {
				op.ParentID = c.Pending[len(c.Pending)-1].ID()
			}
			c.Pending = append(c.Pending, op)
			ops = append(ops, op)
		}
		idx++
	}
	c.Unlock()
	must(c.Store.Append(ctx, ops))
}

// readLoop reads operations from the store and adds it to c.Stream
// taking care to handle acknowledgements: acknowledgements are
// expected to be in order, so the pending unacknowledge list is
// checked. Acknowledgments are not merged as they are already merged
// -- the stream is simply advanced.
func (c *Connector) readLoop(ctx context.Context) {
	limit := 1000

	for {
		c.Lock()
		version := c.Version + 1
		c.Unlock()
		ops, err := c.Store.GetSince(ctx, version, limit)
		if ctx.Err() != nil {
			return
		}
		must(err)

		if len(ops) == 0 {
			must(c.Store.Poll(ctx, version))
			continue
		}

		for _, op := range ops {
			c.Lock()
			c.Version = op.Version()
			change := op.Changes()

			ack := len(c.Pending) > 0 && c.Pending[0].ID() == op.ID()
			if ack {
				c.Pending = c.Pending[1:]
				c.Stream, _ = c.Stream.Next()
			} else {
				c.Stream = c.Stream.ReverseAppend(change)
			}
			c.Unlock()
		}
	}
}

func must(err error) {}
