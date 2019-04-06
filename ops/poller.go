// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"sync"
)

// Polled implements a low-latency in memory poller. Any append will
// effectively trigger any pending polls to return.
//
// It can be used to wrap a store that does not support long polling
//
// Note: closing the wrapped store closes the original store as well
func Polled(s Store) Store {
	return &poller{store: s, waiters: map[chan error]bool{}}
}

type poller struct {
	sync.Mutex
	store   Store
	waiters map[chan error]bool
}

func (p *poller) Append(ctx context.Context, ops []Op) error {
	if err := p.store.Append(ctx, ops); err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()
	for ch := range p.waiters {
		ch <- nil
	}
	p.waiters = map[chan error]bool{}
	return nil
}

func (p *poller) GetSince(ctx context.Context, version, limit int) ([]Op, error) {
	result, err := p.store.GetSince(ctx, version, limit)
	if _, ok := ctx.Deadline(); ok && err == nil && len(result) == 0 {
		p.poll(ctx, version)
		if ctx.Err() != nil {
			return nil, nil
		}
		result, err = p.store.GetSince(ctx, version, limit)
	}
	return result, err
}

func (p *poller) poll(ctx context.Context, version int) {
	done := make(chan error, 1)

	p.Lock()
	p.waiters[done] = true
	p.Unlock()

	defer func() {
		p.Lock()
		defer p.Unlock()
		delete(p.waiters, done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}

func (p *poller) Close() {
	p.store.Close()
	p.Lock()
	defer p.Unlock()
	for ch := range p.waiters {
		ch <- nil
	}
	p.waiters = map[chan error]bool{}
}
