// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

func TestAsync(t *testing.T) {
	async := streams.NewAsync(0)
	s := async.Wrap(streams.New())
	cx := []changes.Change{}
	var latest streams.Stream = s
	s.Nextf(struct{}{}, func() {
		var c changes.Change
		latest, c = latest.Next()
		cx = append(cx, c)
	})

	s1 := s.Append(changes.Move{Offset: 0, Count: 1, Distance: 2})
	s2 := s1.Append(changes.Move{Offset: 5, Count: 6, Distance: 7})
	_ = s2.Append(changes.Move{Offset: 3, Count: 4, Distance: 5})
	if len(cx) != 0 {
		t.Fatal("Async scheduler unexpectedly flushed", cx)
	}

	if count := async.Loop(1); count != 1 {
		t.Fatal("Async Loop(1) return unexpected count", count)
	}

	if count := async.Loop(-1); count != 2 {
		t.Fatal("Async Loop(-1) did not flush", count)
	}

	expected := []changes.Change{
		changes.Move{Offset: 0, Count: 1, Distance: 2},
		changes.Move{Offset: 5, Count: 6, Distance: 7},
		changes.Move{Offset: 3, Count: 4, Distance: 5},
	}
	if !reflect.DeepEqual(cx, expected) {
		t.Fatal("Unexpected result", cx)
	}
}

func TestAsyncMerge(t *testing.T) {
	async := streams.NewAsync(0)
	up := async.Wrap(streams.New())
	down := async.Wrap(streams.New())
	streams.Connect(up, down)

	up = up.Append(changes.Move{Offset: 0, Count: 2, Distance: 2})
	down = down.Append(changes.Move{Offset: 10, Count: 2, Distance: 2})
	if x, cx := up.Next(); x != nil {
		t.Fatal("unexpected sync behavior", cx)
	}
	if x, cx := down.Next(); x != nil {
		t.Fatal("unexpected sync behavior", cx)
	}

	async.Loop(-1)
	_, change1 := up.Next()
	_, change2 := down.Next()
	exp1 := changes.ChangeSet{nil, changes.Move{Offset: 10, Count: 2, Distance: 2}}
	exp2 := changes.ChangeSet{nil, changes.Move{Offset: 0, Count: 2, Distance: 2}}
	if !reflect.DeepEqual(change1, exp1) || !reflect.DeepEqual(change2, exp2) {
		t.Fatal("Unexpected changes", change1, change2)
	}
}

func TestAsyncForever(t *testing.T) {
	async := streams.NewAsync(0)
	async.LoopForever()
	s := async.Wrap(streams.New())
	var wg sync.WaitGroup
	wg.Add(1)
	s.Nextf("key", wg.Done)
	s.Append(changes.Move{Offset: 2, Count: 2, Distance: 2})
	wg.Wait()
	async.Close()
}
