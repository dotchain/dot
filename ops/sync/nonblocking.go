// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"context"
	"sync"

	"github.com/dotchain/dot/ops"
)

// NonBlocking converts a regular store into a non-blocking store
//
// This modifies the GetSince call to return immediately and fetch
// results asynchronously.
func NonBlocking(s ops.Store) ops.Store {
	return &nonblocking{Store: s, cache: map[int][]ops.Op{}, progress: map[int]bool{}}
}

type nonblocking struct {
	sync.Mutex
	ops.Store
	cache    map[int][]ops.Op
	progress map[int]bool
}

func (n *nonblocking) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	n.Lock()
	defer n.Unlock()

	if n.progress[version] {
		return nil, nil
	}

	if ops, ok := n.cache[version]; ok {
		delete(n.cache, version)
		return ops, nil
	}

	n.progress[version] = true
	go n.fetch(ctx, version, limit)
	return nil, nil
}

func (n *nonblocking) fetch(ctx context.Context, version, limit int) {
	ops, err := n.Store.GetSince(ctx, version, limit)
	n.Lock()
	defer n.Unlock()

	delete(n.progress, version)
	if err == nil {
		n.cache[version] = ops
	}
}
