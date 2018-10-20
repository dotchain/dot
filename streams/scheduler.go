// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"github.com/dotchain/dot/changes"
	"sync"
)

// Async queues any callbacks without executing them
// synchronously. These queued callbacks can be executed with a call
// to Run()
type Async struct {
	l       sync.Mutex
	pending []func()
}

// Wrap wraps a stream with an updated scheduler. Any calls to Nextf
// on the returned stream will have the callback scheduled.
func (as *Async) Wrap(s Stream) Stream {
	return &async{as, s}
}

// Loop executes pending callbacks.  The number of callbacks to
// executed is controlled by the count parameter. If it is negative,
// the process continues until the queue is flushed.
//
// The return value is the number of queued functions that were
// executed.
func (as *Async) Loop(count int) int {
	as.l.Lock()
	defer as.l.Unlock()

	run := 0
	for len(as.pending) > 0 && count != 0 {
		count--
		run++
		as.pending[0]()
		as.pending = as.pending[1:]
	}
	return run
}

// Run safely runs a callback.  No async stream calls are likely to be
// processed when this is happening.  It is not safe to call Schedule
// from within a stream callback -- it will deadlock
func (as *Async) Run(fn func()) {
	as.l.Lock()
	defer as.l.Unlock()
	fn()
}

type async struct {
	as *Async
	Stream
}

func (a *async) Append(c changes.Change) Stream {
	return &async{a.as, a.Stream.Append(c)}
}

func (a *async) ReverseAppend(c changes.Change) Stream {
	return &async{a.as, a.Stream.ReverseAppend(c)}
}

func (a *async) Nextf(key interface{}, fn func()) {
	if fn != nil {
		old := fn
		fn = func() {
			a.as.pending = append(a.as.pending, old)
		}
	}
	a.Stream.Nextf(key, fn)
}

func (a *async) Next() (Stream, changes.Change) {
	n, c := a.Stream.Next()
	if n != nil {
		n = &async{a.as, n}
	}
	return n, c
}
