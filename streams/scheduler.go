// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

// Scheduler is the interface that Streams use to schedule
// calls to the function provided as argument to Nextf.
type Scheduler interface {
	Schedule(fn func())
}

type syncScheduler struct{}

func (s *syncScheduler) Schedule(fn func()) {
	fn()
}

// SyncScheduler is the default scheduler used with streams. The
// function provided as argument to Stream.Nextf is called
// synchronously via this scheduler.
var SyncScheduler = &syncScheduler{}

// AsyncScheduler queues any callbacks without executing them
// synchronously. These queued callbacks can be executed with a call
// to Run()
type AsyncScheduler struct {
	Pending []func()
}

// Schedule queues the provided callback. This is not safe for
// concurrent access but it is rentrant (i.e. one of the callbacks can
// schedule more callbacks)
func (as *AsyncScheduler) Schedule(fn func()) {
	as.Pending = append(as.Pending, fn)
}

// Run executes pending callbacks.  The number of callbacks to execute
// is controlled by the count parameter. If it is negative, the
// process continues until the queue is flushed.
//
// The return value is the number of queued functions that were
// executed.
func (as *AsyncScheduler) Run(count int) int {
	run := 0
	for len(as.Pending) > 0 && count != 0 {
		count--
		run++
		as.Pending[0]()
		as.Pending = as.Pending[1:]
	}
	return run
}
