// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"github.com/dotchain/dot/ops"
)

// MemStore returns an in-memory implementation of the Store interface
func MemStore(initial []ops.Op) ops.Store {
	m := &mem{
		seen:    map[interface{}]bool{},
		notify:  map[chan struct{}]bool{},
		control: make(chan func()),
		closed:  make(chan struct{}),
	}
	go func() {
		for {
			select {
			case fn := <-m.control:
				fn()
			case <-m.closed:
				return
			}
		}
	}()
	_ = m.Append(context.Background(), initial)
	return m
}

type mem struct {
	ops     []ops.Op
	seen    map[interface{}]bool
	notify  map[chan struct{}]bool
	control chan func()
	closed  chan struct{}
}

func (m *mem) Append(ctx context.Context, ops []ops.Op) error {
	fn := func() {
		for _, op := range ops {
			if id := op.ID(); !m.seen[id] {
				m.seen[id] = true
				op = op.WithVersion(len(m.ops))
				m.ops = append(m.ops, op)
			}
		}
		for ch := range m.notify {
			ch <- struct{}{}
		}
		m.notify = map[chan struct{}]bool{}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.control <- fn:
	case <-m.closed:
	}
	return nil
}

func (m *mem) GetSince(ctx context.Context, version, limit int) (ops []ops.Op, err error) {
	done := make(chan struct{}, 1)
	fn := func() {
		if version < len(m.ops) {
			ops = m.ops[version:]
			if len(ops) > limit {
				ops = ops[:limit]
			}
		}
		done <- struct{}{}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.closed:
		return nil, nil
	case m.control <- fn:
	}
	<-done
	return ops[:len(ops):len(ops)], err
}

func (m *mem) Poll(ctx context.Context, version int) error {
	done := make(chan struct{}, 1)
	fn := func() {
		if version >= len(m.ops) {
			m.notify[done] = true
		} else {
			done <- struct{}{}
		}
	}

	select {
	case <-ctx.Done():
	case <-m.closed:
	case m.control <- fn:
	}

	select {
	case <-ctx.Done():
	case <-m.closed:
	case <-done:
	}

	return ctx.Err()
}

func (m *mem) Close() {
	close(m.closed)

}
