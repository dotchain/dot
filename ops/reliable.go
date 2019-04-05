// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"time"

	"github.com/dotchain/dot/log"
)

// Reliable takes a store that can fail and converts it to a
// reliable store. All Append() calls return success immediately with
// background attempts to deliver/retry.
//
// Note that GetSince is not modified by this -- it can still be
// unreliable.
//
// Poll is modified to retry up to the specified timeout.
func Reliable(s Store, rand func() float64, initial, max time.Duration, l log.Log) Store {
	i, m := float64(initial), float64(max)
	ctx, cancel := context.WithCancel(context.Background())
	r := &reliable{s, nil, make(chan func(), 10000), rand, i, m, ctx, cancel, l}
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
	Store
	pending      []Op
	jobs         chan func()
	rand         func() float64
	initial, max float64

	deliverCtx    context.Context
	cancelDeliver func()

	log.Log
}

func (r *reliable) Close() {
	r.cancelDeliver()
	r.Store.Close()
}

func (r *reliable) Append(ctx context.Context, ops []Op) error {
	r.jobs <- func() {
		wasPending := len(r.pending) > 0
		r.pending = append(r.pending, ops...)
		if size := len(r.pending); !wasPending && size > 0 {
			go r.deliver(r.pending[:size:size])
		}
	}
	return nil
}

func (r *reliable) deliver(pending []Op) {
	current := r.initial

	for {
		err := r.Store.Append(context.Background(), pending)
		if err == nil {
			r.jobs <- func() {
				r.pending = r.pending[len(pending):]
				if size := len(r.pending); size > 0 {
					go r.deliver(r.pending[:size:size])
				}
			}
			return
		}

		delta := 0.5 * current
		min := current - delta
		max := current + delta
		next := min + r.rand()*(max-min+1)

		r.Log.Println("Retrying delivery after", time.Duration(next), err)
		timer := time.NewTimer(time.Duration(next))
		select {
		case <-r.deliverCtx.Done():
			timer.Stop()
			return
		case <-timer.C:
			timer.Stop()
		}

		current *= 1.5
		if current > r.max {
			current = r.max
		}
	}
}

func (r *reliable) Poll(ctx context.Context, version int) error {
	fn := func() error {
		if _, ok := ctx.Deadline(); !ok {
			ctx2, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()
			ctx = ctx2
		}

		err := r.Store.Poll(ctx, version)
		if ctx.Err() != nil {
			// if canceled due to timeouts, try GetSince again
			return nil
		}
		return err
	}
	return r.retry(ctx, fn)
}

func (r *reliable) retry(ctx context.Context, fn func() error) error {
	current := r.initial

	for {
		err := fn()
		if err == nil || err == ctx.Err() {
			return err
		}

		delta := 0.5 * current
		min := current - delta
		max := current + delta
		next := min + r.rand()*(max-min+1)
		timer := time.NewTimer(time.Duration(next))

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			timer.Stop()
		}

		current *= 1.5
		if current > r.max {
			current = r.max
		}
	}
}
