// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"context"
	"sync"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
)

// Stream converts an ops Store into a stream.
//
// Calling Append() on the Stream is asynchronously appended to the
// Store and remote changes are also automatically merged into the
// returned stream.
//
// Sync works with a Transformed store.  A raw store can be converted
// to a transformed store using ops.Transformed(rawStore).
//
// Additional sync options can be used to configure the behavior
//
// The returned close function can be used to shutdown the
// synchronization but the underlying store still needs to be
// separately released.
func Stream(store ops.Store, opts ...Option) (s streams.Stream, closefn func()) {
	notify := func(version int, pending []ops.Op) {}
	c := &Config{Store: store, Log: log.Default(), Notify: notify, Version: -1}

	for _, opt := range opts {
		opt(c)
	}
	c.Store = Reliable(store, opts...)

	session := &session{config: c, stream: streams.New()}
	session.id = session.newID()
	session.stream = &stream{streams.New(), map[interface{}]func(){}, &sync.Mutex{}}

	// add fake entries for each pending as an entry is expected
	// per pending request. See the ack behavior in session.read
	latest := session.stream
	for range c.Pending {
		latest = latest.Append(nil)
	}

	if len(c.Pending) > 0 {
		session.write(c.Pending)
	}

	ctx, cancel := context.WithCancel(context.Background())
	closed := make(chan struct{})

	latest.Nextf(session, session.onAppend)
	go func() {
		session.read(ctx)
		session.stream.Nextf(session, nil)
		close(closed)
	}()

	closefn = func() {
		cancel()
		<-closed
		// clean up the reliable store only
		// leaving the original store as is
		c.Store.(*reliable).cancelDeliver()
	}

	return latest, closefn
}
