// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

type queue struct {
	control chan func()
	closed  chan struct{}
}

func (q queue) push(fn func()) bool {
	select {
	case q.control <- fn:
		return true
	case <-q.closed:
		return false
	}
}

func (q queue) close() {
	close(q.closed)
}

func (q queue) run() {
	for {
		select {
		case fn := <-q.control:
			fn()
		case <-q.closed:
			return
		}
	}
}
