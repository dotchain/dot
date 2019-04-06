// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package testops

import (
	"context"
	"sync"

	"github.com/dotchain/dot/ops"
)

// MemStore returns an in-memory implementation of the Store interface
func MemStore(initial []ops.Op) ops.Store {
	m := &mem{seen: map[interface{}]bool{}}
	_ = m.Append(context.Background(), initial)
	return m
}

type mem struct {
	ops  []ops.Op
	seen map[interface{}]bool

	sync.Mutex
}

func (m *mem) Append(ctx context.Context, ops []ops.Op) error {
	m.Lock()
	defer m.Unlock()
	for _, op := range ops {
		if id := op.ID(); !m.seen[id] {
			m.seen[id] = true
			op = op.WithVersion(len(m.ops))
			m.ops = append(m.ops, op)
		}
	}
	return nil
}

func (m *mem) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	result := []ops.Op(nil)
	m.Lock()
	defer m.Unlock()
	if version < len(m.ops) {
		result = m.ops[version:]
		if len(result) > limit {
			result = result[:limit]
		}
	}
	return result, nil
}

func (m *mem) Close() {
}
