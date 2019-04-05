// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"sync"
	"time"

	"github.com/dotchain/dot/log"
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
// stores are in play, the same Async object is recommended for all
// stores.
type Connector struct {
	Version int
	Pending []Op
	streams.Stream
	*streams.Async
	Store
	log.Log
	close func()
	sync.Mutex
}

// NewConnector creates a new connection between the store and a
// stream. It creates an Async object as well as the stream taking
// care to wrap the stream via Async.Wrap.
//
// Connector works with a reliable store, so it creates one using
// Reliable(store) (using the provided rand function for binary
// exponential backoff).
//
// The provided store parameter is expected to be already transformed
// (via Transformed(store) for instance).
//
// The version refers to the version of the last operation received
// from the server before.  In case of a fresh start, this should be
// -1.
//
// Pending refers to any prior operations that were not acknowledged
// by the store.  These may have been sent before but sending them
// again is harmless as the store will drop duplicates.
func NewConnector(version int, pending []Op, store Store, rand func() float64) *Connector {
	c := &Connector{
		Version: version,
		Pending: pending,
		Log:     log.Default(),
		Async:   streams.NewAsync(0),
	}

	c.Stream = c.Async.Wrap(streams.New())

	// add fake entries for each pending as an entry is expected
	// per pending request. See the ack behavior in readLoop
	clientStream := c.Stream
	for range pending {
		clientStream = clientStream.Append(nil)
	}

	c.Store = Reliable(store, rand, time.Second*2, time.Minute, log.Default())
	return c
}

// Connect starts the synchronization process.
//
// If logging is desired, the Log field should be set before being
// called.
//
// All the fields of the Connector object are not safe to use until a
// Disconnenct call.
func (c *Connector) Connect() {
	c.Async.LoopForever()

	// create a context that is canceled on close
	// this signals all go routines to stop
	ctx, cancel := context.WithCancel(context.Background())

	closed := make(chan struct{})
	c.close = func() {
		cancel()
		<-closed // wait for read loop to finish
	}

	c.must(c.Store.Append(ctx, c.Pending))
	c.Stream.Nextf(c, func() { c.write(ctx) })

	go func() {
		c.readLoop(ctx)
		c.Stream.Nextf(c, nil)
		c.Async.Close()
		c.updatePending(ctx)
		close(closed)
	}()
}

// Disconnect stops the synchronization process.  The version and
// pending are updated to the latest values when the call returns
//
// The Async field is automatically closed but the stores are not.
func (c *Connector) Disconnect() {
	c.close()
}

// write takes any unwritten changes from c.Stream and writes it out
// to the ops store. note that write does not update c.Stream as
// c.Stream tracks the last upstream version
func (c *Connector) write(ctx context.Context) {
	ops := c.updatePending(ctx)
	if len(ops) > 0 {
		c.must(c.Store.Append(ctx, ops))
	}
}

func (c *Connector) updatePending(ctx context.Context) []Op {
	var idx int
	var ops []Op

	c.Lock()
	defer c.Unlock()
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
	return ops
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
		c.must(err)

		if len(ops) == 0 {
			c.must(c.Store.Poll(ctx, version))
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

func (c *Connector) must(err error) {
	if err != nil {
		c.Log.Println(err)
	}
}
