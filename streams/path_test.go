// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestChildOf_ModifyChild(t *testing.T) {
	base := streams.New()
	child := streams.ChildOf(base, 5, 2)
	move := changes.Move{2, 2, 2}
	child = child.Append(move)
	cx, base := base.Next()
	expected := changes.PathChange{[]interface{}{5, 2}, move}
	if !reflect.DeepEqual(cx, expected) {
		t.Error("unexpected change", cx)
	}

	child.ReverseAppend(move)
	cx, _ = base.Next()
	if !reflect.DeepEqual(cx, expected) {
		t.Error("unexpected change", cx)
	}
	cx, child2 := child.Next()
	if cx != move {
		t.Error("unexpeced change", cx)
	}
	if _, next := child2.Next(); next != nil {
		t.Error("Unexpected next", next)
	}

	count := 0
	child.Nextf("key", func(cx changes.Change, _ streams.Stream) {
		count++
		if cx != move {
			t.Error("Unexpected Nextf", cx)
		}
	})
	if count != 1 {
		t.Error("Unexpected callback count", count)
	}
}

func TestChildOf_InvalidRef(t *testing.T) {
	base := streams.New()
	child := streams.ChildOf(base, 5, 2)
	base = base.Append(changes.Replace{types.S8("OK"), changes.Nil})
	base.Append(changes.PathChange{[]interface{}{5, 2}, changes.Move{2, 2, 2}})

	if c, s := child.Next(); c != nil || s != nil {
		t.Error("Unexpected next value", c, s)
	}

	child.Nextf("key", func(c changes.Change, s streams.Stream) {
		t.Fatal("Unexpected callback", c, s)
	})
	child.Nextf("key", nil)
}

func TestChildOf_Scheduler(t *testing.T) {
	async := &streams.AsyncScheduler{}
	base := streams.New()
	child := streams.ChildOf(base, 5, 2).WithScheduler(async)
	if child.Scheduler() != async {
		t.Error("Could not update scheduler")
	}
}

func TestFilterPath(t *testing.T) {
	base := streams.New()
	child := streams.FilterPath(base, "hello", "world")
	pc := func(c changes.Change, keys ...interface{}) changes.Change {
		return changes.PathChange{keys, c}
	}

	base = base.Append(pc(changes.Move{2, 2, 2}, "bloomy"))
	cx, child := child.Next()
	if cx != nil {
		t.Error("Unexpected filter failure", cx)
	}

	change := pc(changes.Move{2, 2, 2}, "hello", "world", "ok")
	base = base.Append(change)
	if cx, _ := child.Next(); !reflect.DeepEqual(cx, change) {
		t.Error("Unexpected next change", cx)
	}

	change = pc(changes.Move{3, 3, 3}, "goop")
	child.Append(change)
	if cx, _ := base.Next(); !reflect.DeepEqual(cx, change) {
		t.Error("Unexpected next change", cx)
	}
}

func TestFilterOutPath(t *testing.T) {
	base := streams.New()
	child := streams.FilterOutPath(base, "hello", "world")
	pc := func(c changes.Change, keys ...interface{}) changes.Change {
		return changes.PathChange{keys, c}
	}

	base = base.Append(pc(changes.Move{2, 2, 2}, "hello", "world", "ok"))
	cx, child := child.Next()
	if cx != nil {
		t.Error("Unexpected filter failure", cx)
	}

	change := pc(changes.Move{2, 2, 2}, "boop")
	base = base.Append(change)
	if cx, _ := child.Next(); !reflect.DeepEqual(cx, change) {
		t.Error("Unexpected next change", cx)
	}

	change = pc(changes.Move{3, 3, 3}, "goop")
	child.Append(change)
	if cx, _ := base.Next(); !reflect.DeepEqual(cx, change) {
		t.Error("Unexpected next change", cx)
	}
}
