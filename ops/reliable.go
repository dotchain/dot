// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"time"
)

// ReliableStore takes a store that can fail and converts it to a
// reliable store. All Append() calls return success immediately but
// attempt to deliver with retires in the background. Poll calls block
// for the full timeout specified in the context, retrying the
// underlying store (unless the underlying call succeeds)
//
// The suggested rand function is math.rand.Float64.
//
// Note that the GetSince is not modified by this -- it can still be
// unreliable.
//
// Poll is modified to retry up to the specified timeout
func ReliableStore(s Store, rand func() float64, initial, max time.Duration) Store {
	i, m := float64(initial), float64(max)
	r := &reliable{s, nil, make(chan func(), 10000), rand, i, m}
	go func() {
		for fn := range r.jobs {
			fn()
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
		time.Sleep(time.Duration(next))
		current = current * 1.5
		if current > r.max {
			current = r.max
		}
	}
}

func (r *reliable) Poll(ctx context.Context, version int) error {
	current := r.initial

	for {
		err := r.Store.Poll(ctx, version)
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
		current = current * 1.5
		if current > r.max {
			current = r.max
		}
	}
}
