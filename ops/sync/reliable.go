// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"time"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
)

// Reliable takes a store that can fail and converts it to a
// reliable store. All Append() calls return success immediately with
// background attempts to deliver/retry.
//
// Poll and GetSince are modified to retry up to the specified timeout.
//
// Note: the only options that affects Reliable() are WithBackoff()
// and WithLog()
func Reliable(s ops.Store, opts ...Option) ops.Store {
	c := &Config{Store: s, Log: log.Default()}
	c.Backoff.Rand = func() float64 { return 1.0 }
	c.Backoff.Initial = time.Second
	c.Backoff.Max = time.Minute
	for _, opt := range opts {
		opt(c)
	}

	return newReliable(c)
}

func newReliable(c *Config) ops.Store {
	ctx, cancel := context.WithCancel(context.Background())
	r := &reliable{c, nil, make(chan func(), 10000), ctx, cancel}
	go func() {
		for {
			select {
			case fn := <-r.jobs:
				fn()
			case <-ctx.Done():
				return
			}
		}
	}()
	return r
}

type reliable struct {
	*Config

	pending       []ops.Op
	jobs          chan func()
	deliverCtx    context.Context
	cancelDeliver func()
}

func (r *reliable) Close() {
	r.cancelDeliver()
	r.Store.Close()
}

func (r *reliable) Append(ctx context.Context, ops []ops.Op) error {
	r.jobs <- func() {
		wasPending := len(r.pending) > 0
		r.pending = append(r.pending, ops...)
		if size := len(r.pending); !wasPending && size > 0 {
			go r.deliver(r.pending[:size:size])
		}
	}
	return nil
}

func (r *reliable) deliver(pending []ops.Op) {
	err := r.retry(r.deliverCtx, func() error {
		return r.Store.Append(r.deliverCtx, pending)
	})

	if err == nil {
		r.jobs <- func() {
			r.pending = r.pending[len(pending):]
			if size := len(r.pending); size > 0 {
				go r.deliver(r.pending[:size:size])
			}
		}
	}
}

func (r *reliable) Poll(ctx context.Context, version int) error {
	if _, ok := ctx.Deadline(); !ok {
		ctx2, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		ctx = ctx2
	}

	fn := func() error {
		err := r.Store.Poll(ctx, version)
		if ctx.Err() != nil {
			return nil
		}
		return err
	}

	if err := r.retry(ctx, fn); err != nil {
		r.Log.Println("poll: ", err)
	}
	return ctx.Err()
}

func (r *reliable) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	var result []ops.Op
	fn := func() error {
		r, err := r.Store.GetSince(ctx, version, limit)
		if err == nil && len(r) == 0 {
			return errors.New("retrying on empty result")
		}
		result = r
		return err
	}

	if err := r.retry(ctx, fn); err != nil && err != ctx.Err() {
		r.Log.Println("GetSince: ", err)
	}
	return result, ctx.Err()
}

func (r *reliable) retry(ctx context.Context, fn func() error) error {
	current := float64(r.Backoff.Initial)

	for {
		err := fn()
		if err == nil || err == ctx.Err() {
			return err
		}

		delta := 0.5 * current
		min := current - delta
		max := current + delta
		next := min + r.Backoff.Rand()*(max-min+1)
		timer := time.NewTimer(time.Duration(next))

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			timer.Stop()
		}

		current *= 1.5
		if current > float64(r.Backoff.Max) {
			current = float64(r.Backoff.Max)
		}
	}
}
