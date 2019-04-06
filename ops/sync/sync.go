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

// Sync creates a synchronized stream out of a reliable transformed
// store.
//
// Updating the stream locally updates the store and vice-versa. A raw
// store can be made into a transformed store with Transformed() and
// made more reliable using Reliable().  Sync does not do either of
// these automatically
//
// Additional sync options can be used to configure the behavior
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
	}

	return latest, closefn
}
