// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"context"
	"github.com/dotchain/dot/ops"
)

// MemPoller implements a low-latency in memory poller. Any append will
// effectively trigger any pending polls to return.
func MemPoller(s ops.Store) ops.Store {
	q := queue{control: make(chan func()), closed: make(chan struct{})}
	go q.run()
	return &poller{s, q, map[chan error]bool{}}
}

type poller struct {
	store   ops.Store
	q       queue
	waiters map[chan error]bool
}

func (p *poller) Append(ctx context.Context, ops []ops.Op) error {
	if err := p.store.Append(ctx, ops); err != nil {
		return err
	}

	p.q.push(func() {
		for ch := range p.waiters {
			ch <- nil
		}
		p.waiters = map[chan error]bool{}
	})
	return nil
}

func (p *poller) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	return p.store.GetSince(ctx, version, limit)
}

func (p *poller) Poll(ctx context.Context, version int) error {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan error, 2)
	p.q.push(func() {
		p.waiters[done] = true
	})
	defer func() {
		p.q.push(func() {
			delete(p.waiters, done)
		})
	}()
	go func() {
		done <- p.store.Poll(ctx2, version)
	}()
	return <-done
}

func (p *poller) Close() {
	p.store.Close()
	p.q.close()
}
