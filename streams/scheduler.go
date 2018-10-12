// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

// Async queues any callbacks without executing them
// synchronously. These queued callbacks can be executed with a call
// to Run()
type Async struct {
	pending []func()
}

// Wrap wraps a stream with an updated scheduler. Any calls to Nextf
// on the returned stream will have the callback scheduled.
func (as *Async) Wrap(s Stream) Stream {
	return async{as, s}
}

// Run executes pending callbacks.  The number of callbacks to execute
// is controlled by the count parameter. If it is negative, the
// process continues until the queue is flushed.
//
// The return value is the number of queued functions that were
// executed.
func (as *Async) Run(count int) int {
	run := 0
	for len(as.pending) > 0 && count != 0 {
		count--
		run++
		as.pending[0]()
		as.pending = as.pending[1:]
	}
	return run
}

type async struct {
	as *Async
	Stream
}

func (a async) Nextf(key interface{}, fn func()) {
	if fn != nil {
		old := fn
		fn = func() {
			a.as.pending = append(a.as.pending, old)
		}
	}
	a.Stream.Nextf(key, fn)
}
