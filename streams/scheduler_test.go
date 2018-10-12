// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"reflect"
	"testing"
)

func TestAsync(t *testing.T) {
	async := &streams.Async{}
	s := async.Wrap(streams.New())
	cx := []changes.Change{}
	var latest streams.Stream = s
	s.Nextf(struct{}{}, func() {
		var c changes.Change
		c, latest = latest.Next()
		cx = append(cx, c)
	})

	s1 := s.Append(changes.Move{0, 1, 2})
	s2 := s1.Append(changes.Move{5, 6, 7})
	_ = s2.Append(changes.Move{3, 4, 5})
	if len(cx) != 0 {
		t.Fatal("Async scheduler unexpectedly flushed", cx)
	}

	if count := async.Run(1); count != 1 {
		t.Fatal("Async Run(1) return unexpected count", count)
	}

	if count := async.Run(-1); count != 2 {
		t.Fatal("Async Run(-1) did not flush", count)
	}

	expected := []changes.Change{
		changes.Move{0, 1, 2},
		changes.Move{5, 6, 7},
		changes.Move{3, 4, 5},
	}
	if !reflect.DeepEqual(cx, expected) {
		t.Fatal("Unexpected result", cx)
	}
}

func TestAsyncMerge(t *testing.T) {
	async := &streams.Async{}
	up := async.Wrap(streams.New())
	down := async.Wrap(streams.New())
	b := &streams.Branch{up, down, false}
	b.Connect()

	up = up.Append(changes.Move{0, 2, 2})
	down = down.Append(changes.Move{10, 2, 2})
	if cx, x := up.Next(); x != nil {
		t.Fatal("unexpected sync behavior", cx)
	}
	if cx, x := down.Next(); x != nil {
		t.Fatal("unexpected sync behavior", cx)
	}

	async.Run(-1)
	change1, _ := up.Next()
	change2, _ := down.Next()
	exp1 := changes.ChangeSet{nil, changes.Move{10, 2, 2}}
	exp2 := changes.ChangeSet{nil, changes.Move{0, 2, 2}}
	if !reflect.DeepEqual(change1, exp1) || !reflect.DeepEqual(change2, exp2) {
		t.Fatal("Unexpected changes", change1, change2)
	}
}
