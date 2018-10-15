// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"errors"
	"github.com/dotchain/dot/ops"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// Fake unreliable store
type unreliable struct {
	err   error
	ops   []ops.Op
	count int
}

func (u *unreliable) Append(ctx context.Context, ops []ops.Op) error {
	err := u.err
	if err == nil {
		u.ops = append(u.ops, ops...)
	}
	u.count++
	return err
}

func (u *unreliable) Poll(ctx context.Context, version int) error {
	u.count++
	return u.err
}

func (u *unreliable) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	return nil, nil
}

func (u *unreliable) Close() {
}

func TestReliableAppend(t *testing.T) {
	u := &unreliable{err: errors.New("something")}
	r := ops.ReliableStore(u, rand.Float64, time.Millisecond, 10*time.Millisecond)
	go func() {
		time.Sleep(50 * time.Millisecond)
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
	if u.count < 10 || !reflect.DeepEqual(u.ops, expected) {
		t.Error("Unexpected state", u.count, u.ops)
	}
}

func TestReliablePoll(t *testing.T) {
	u := &unreliable{err: errors.New("something")}
	r := ops.ReliableStore(u, rand.Float64, time.Millisecond, 10*time.Millisecond)
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
