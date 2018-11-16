// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"context"
	"github.com/dotchain/dot/ops"
)

// MemStore returns an in-memory implementation of the Store interface.
// The Poll method either returns immediately or blocks for the full
// specified delay.  Use the Poller to reduce latencies.
func MemStore(initial []ops.Op) ops.Store {
	q := queue{control: make(chan func()), closed: make(chan struct{})}
	m := &store{seen: map[interface{}]bool{}, q: q}
	go q.run()
	m.append(context.Background(), initial)
	return m
}

type store struct {
	ops  []ops.Op
	seen map[interface{}]bool
	q    queue
}

func (m *store) append(ctx context.Context, ops []ops.Op) {
	m.q.push(func() {
		for _, op := range ops {
			if id := op.ID(); !m.seen[id] {
				m.seen[id] = true
				op = op.WithVersion(len(m.ops))
				m.ops = append(m.ops, op)
			}
		}
	})
}

func (m *store) Append(ctx context.Context, ops []ops.Op) error {
	m.append(ctx, ops)
	return nil
}

func (m *store) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	result := []ops.Op(nil)
	done := make(chan struct{}, 1)
	queued := m.q.push(func() {
		if version < len(m.ops) {
			result = m.ops[version:]
			if len(result) > limit {
				result = result[:limit]
			}
		}
		done <- struct{}{}
	})

	if queued {
		<-done
	}

	return result, nil
}

func (m *store) Poll(ctx context.Context, version int) error {
	done := make(chan struct{}, 1)
	queued := m.q.push(func() {
		if version < len(m.ops) {
			done <- struct{}{}
		}
	})
	if queued {
		select {
		case <-ctx.Done():
		case <-m.q.closed:
		case <-done:
		}
	}

	return ctx.Err()
}

func (m *store) Close() {
	m.q.close()
}
