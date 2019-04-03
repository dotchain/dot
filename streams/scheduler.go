// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"sync"

	"github.com/dotchain/dot/changes"
)

// Async queues any callbacks without executing them
// synchronously. These queued callbacks can be executed with a call
// to Run()
type Async struct {
	c chan func()
}

// NewAsync creates a new async scheduler
func NewAsync(size int) *Async {
	if size == 0 {
		size = 1000
	}
	return &Async{make(chan func(), size)}
}

// Wrap wraps a stream with an updated scheduler. Any calls to Nextf
// on the returned stream will have the callback scheduled.
func (as *Async) Wrap(s Stream) Stream {
	return &async{as, s, &sync.Mutex{}}
}

// Loop executes pending callbacks in a loop.
func (as *Async) Loop(count int) int {
	run := 0
	for len(as.c) > 0 && count != 0 {
		count--
		run++
		(<-as.c)()
	}
	return run
}

// LoopForever executes all pending callbacks.  Calling "Close" will
// release the goroutine
func (as *Async) LoopForever() {
	go func() {
		for fn := range as.c {
			fn()
		}
	}()
}

// Close releases the channel
func (as *Async) Close() {
	close(as.c)
}

type async struct {
	as *Async
	Stream
	*sync.Mutex
}

func (a *async) Append(c changes.Change) Stream {
	a.Lock()
	defer a.Unlock()
	return &async{a.as, a.Stream.Append(c), a.Mutex}
}

func (a *async) ReverseAppend(c changes.Change) Stream {
	a.Lock()
	defer a.Unlock()
	return &async{a.as, a.Stream.ReverseAppend(c), a.Mutex}
}

func (a *async) Nextf(key interface{}, fn func()) {
	if fn != nil {
		old := fn
		fn = func() {
			// TODO: check if old is still registered
			a.as.c <- old
		}
	}
	a.Lock()
	defer a.Unlock()
	a.Stream.Nextf(key, fn)
}

func (a *async) Next() (Stream, changes.Change) {
	a.Lock()
	defer a.Unlock()
	n, c := a.Stream.Next()
	if n != nil {
		n = &async{a.as, n, a.Mutex}
	}
	return n, c
}

func (a *async) Schedule(fn func()) {
	a.as.c <- fn
}
