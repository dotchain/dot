// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/test/testops"
)

func TestPolled(t *testing.T) {
	store := ops.Polled(testops.MemStore(nil))
	defer store.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := store.GetSince(ctx, 0, 1000)
	if err != nil || ctx.Err() == nil {
		t.Error("unexpected poll result", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// now kick off a go routine to wakeup
	go func() {
		time.Sleep(10 * time.Millisecond)
		op := ops.Operation{OpID: "ID1"}
		_ = store.Append(context.Background(), []ops.Op{op})
	}()

	operations, err := store.GetSince(ctx, 0, 1000)
	if err != nil || len(operations) != 1 {
		t.Error("unexpected poll result", err)
	}
}

func TestClosedPolledStore(t *testing.T) {
	store := ops.Polled(testops.MemStore(nil))

	go func() {
		time.Sleep(10 * time.Millisecond)
		store.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := store.GetSince(ctx, 0, 1000)
	if err != nil {
		t.Error("unexpected poll result", err)
	}
}

func TestPolledError(t *testing.T) {
	myerr := errors.New("something")
	store := ops.Polled(fakeStore{append: myerr, get: myerr})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := store.GetSince(ctx, 0, 1000)
	if err != myerr || ctx.Err() != nil {
		t.Error("unexpected poll result", err, ctx.Err())
	}

	// GetSince succeeds but an append shouldn't unblock it
	store = ops.Polled(fakeStore{append: myerr})
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	go func() {
		time.Sleep(time.Millisecond * 10)
		if x := store.Append(ctx, []ops.Op{ops.Operation{}}); x == nil {
			t.Error("internal failure, append should have failed")
		}
	}()

	_, err = store.GetSince(ctx, 0, 1000)
	if err != nil || ctx.Err() == nil {
		t.Error("unexpected poll result", err, ctx.Err())
	}

}

type fakeStore struct {
	append, get error
}

func (f fakeStore) Append(_ context.Context, opx []ops.Op) error {
	return f.append
}

func (f fakeStore) GetSince(_ context.Context, version, limit int) ([]ops.Op, error) {
	return nil, f.get
}

func (f fakeStore) Close() {
}
