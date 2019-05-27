// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/sync"
	"github.com/dotchain/dot/test/testops"
)

// Fake blocking store
type blocking struct {
	ops.Store
	ops []ops.Op
	ch  chan bool
}

func (b *blocking) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	b.ch <- false
	return b.ops, nil
}

func TestNonBlockingGetSince(t *testing.T) {
	b := &blocking{ch: make(chan bool)}
	r := sync.NonBlocking(b)

	if ops, err := r.GetSince(context.Background(), 10, 1000); ops != nil || err != nil {
		t.Error("Unexpected result", ops, err)
	}
	if ops, err := r.GetSince(context.Background(), 10, 1000); ops != nil || err != nil {
		t.Error("Unexpected result", ops, err)
	}
	b.ops = []ops.Op{ops.Operation{OpID: "boo"}}
	<-b.ch

	// TODO: better way to ensure that the other goroutine has been scheduled
	time.Sleep(time.Millisecond)
	ops, err := r.GetSince(context.Background(), 10, 1000)
	if !reflect.DeepEqual(ops, b.ops) || err != nil {
		t.Error("Unexpected result", ops, err)
	}

}

func TestNonblockingSync(t *testing.T) {
	store := ops.Polled(testops.MemStore(nil))
	xformed := ops.Transformed(store, testops.NullCache())
	stream := sync.Stream(xformed, sync.WithNonBlocking(true))

	if err := stream.Pull(); err != nil {
		t.Error("Unexpected error", err)
	}

	if s, _ := stream.Next(); s != nil {
		t.Error("Unexpected next", s)
	}
}
