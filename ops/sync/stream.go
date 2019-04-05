// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"sync"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

type stream struct {
	inner streams.Stream
	fns   map[interface{}]func()
	*sync.Mutex
}

func (s stream) Append(c changes.Change) streams.Stream {
	defer s.notify()()
	s.Lock()
	defer s.Unlock()
	s.inner = s.inner.Append(c)
	return s
}

func (s stream) ReverseAppend(c changes.Change) streams.Stream {
	defer s.notify()()
	s.Lock()
	defer s.Unlock()
	s.inner = s.inner.ReverseAppend(c)
	return s
}

func (s stream) Next() (streams.Stream, changes.Change) {
	s.Lock()
	defer s.Unlock()
	next, c := s.inner.Next()
	if next != nil {
		next = stream{next, s.fns, s.Mutex}
	}
	return next, c
}

func (s stream) Nextf(key interface{}, fn func()) {
	s.Lock()
	defer s.Unlock()
	if fn == nil {
		delete(s.fns, key)
	} else {
		s.fns[key] = fn
	}
}

func (s stream) notify() func() {
	fns := make([]func(), 0, len(s.fns))
	for _, fn := range s.fns {
		fns = append(fns, fn)
	}

	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}
