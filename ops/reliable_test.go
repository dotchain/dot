// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
)

// Fake unreliable store
type unreliable struct {
	sync.Mutex
	err   error
	ops   []ops.Op
	count int
}

func (u *unreliable) Append(ctx context.Context, ops []ops.Op) error {
	u.Lock()
	defer u.Unlock()
	err := u.err
	if err == nil {
		u.ops = append(u.ops, ops...)
	}
	u.count++
	return err
}

func (u *unreliable) Poll(ctx context.Context, version int) error {
	u.Lock()
	defer u.Unlock()
	u.count++
	return u.err
}

func (u *unreliable) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	u.Lock()
	defer u.Unlock()
	err := u.err
	u.count++
	if err == nil {
		return u.ops, nil
	}
	return nil, err
}

func (u *unreliable) Close() {
}

func TestReliableAppend(t *testing.T) {
	u := &unreliable{err: errors.New("something")}
	r := ops.Reliable(u, rand.Float64, time.Millisecond, 10*time.Millisecond, log.Default())
	go func() {
		time.Sleep(50 * time.Millisecond)
		u.Lock()
		defer u.Unlock()
		u.err = nil
	}()
	opx := []ops.Op{ops.Operation{OpID: "one"}}
	if err := r.Append(context.Background(), opx); err != nil {
		t.Fatal("Reliable append failed", err)
	}
	if err := r.Append(context.Background(), opx); err != nil {
		t.Fatal("Reliable append failed", err)
	}
	time.Sleep(100 * time.Millisecond)
	expected := append(append([]ops.Op(nil), opx...), opx...)

	u.Lock()
	defer u.Unlock()
	if u.count < 10 || !reflect.DeepEqual(u.ops, expected) {
		t.Error("Unexpected state", u.count, u.ops)
	}
}

func TestReliablePoll(t *testing.T) {
	u := &unreliable{err: errors.New("something")}
	r := ops.Reliable(u, rand.Float64, time.Millisecond, 10*time.Millisecond, log.Default())
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err := r.Poll(ctx, 100)
	if err != ctx.Err() || u.count < 10 {
		t.Fatal("Unexpected err", u.count, err)
	}
	cancel()

	*u = unreliable{}
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = r.Poll(ctx, 100)
	if err != nil || u.count != 1 {
		t.Fatal("Unexpected err", u.count, err)
	}
	cancel()
}

func TestReliableGetSince(t *testing.T) {
	u := &unreliable{err: errors.New("something")}
	u.ops = []ops.Op{ops.Operation{OpID: "one"}}
	r := ops.Reliable(u, rand.Float64, time.Millisecond, 10*time.Millisecond, log.Default())
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	result, err := r.GetSince(ctx, 100, 102)
	if err != u.err || ctx.Err() != nil || u.count != 1 || len(result) != 0 {
		t.Fatal("Unexpected err", u.count, len(result), err, ctx.Err())
	}
	cancel()

	u.err = nil
	u.count = 0
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	result, err = r.GetSince(ctx, 100, 102)
	if err != nil || u.count != 1 || !reflect.DeepEqual(result, u.ops) {
		t.Fatal("Unexpected err", u.count, err, result)
	}
	cancel()
}
