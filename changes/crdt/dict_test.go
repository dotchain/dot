// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt_test

import (
	"testing"

	"github.com/dotchain/dot/changes/crdt"
)

func TestDictSetUnset(t *testing.T) {
	d := crdt.Dict{}
	c1, d := d.Set("hello", "world")
	if _, v := d.Get("hello"); v != "world" {
		t.Fatal("Set failed", v)
	}

	c2, d2 := d.Set("hello", "world2")
	if _, v := d2.Get("hello"); v != "world2" {
		t.Error("Undo diverged", v)
	}

	d3 := d2.Apply(nil, c2.Revert()).(crdt.Dict)
	if _, v := d3.Get("hello"); v != "world" {
		t.Fatal("Unset failed", v)
	}

	d4 := d2.Apply(nil, c1.Revert()).(crdt.Dict)
	if _, v := d4.Get("hello"); v != "world2" {
		t.Error("Undo diverged", v)
	}

	d5 := d4.Apply(nil, c2.Revert()).(crdt.Dict)
	if r, v := d5.Get("hello"); r != nil || v != nil {
		t.Error("Undo diverged", r, v)
	}
}

func TestDictDelUndel(t *testing.T) {
	d := crdt.Dict{}
	_, d = d.Set("hello", "world")
	if _, v := d.Get("hello"); v != "world" {
		t.Fatal("Set failed", v)
	}

	c2, d2 := d.Delete("hello")
	if r, v := d2.Get("hello"); r != nil || v != nil {
		t.Fatal("Delete failed", r, v)
	}

	d3 := d2.Apply(nil, c2.Revert()).(crdt.Dict)
	if _, v := d3.Get("hello"); v != "world" {
		t.Fatal("Set failed", v)
	}
}

func TestDictUpdate(t *testing.T) {
	inner := crdt.Dict{}
	_, inner = inner.Set("hello", "world")

	d := crdt.Dict{}
	_, d = d.Set("inner", inner)

	c, _ := inner.Set("hello", "world2")
	cx, d2 := d.Update("inner", c)

	_, v := d2.Get("inner")
	if _, x := v.(crdt.Dict).Get("hello"); x != "world2" {
		t.Fatal("Mismatched inner", x)
	}

	d2 = d.Apply(nil, cx).(crdt.Dict)
	_, v = d2.Get("inner")
	if _, x := v.(crdt.Dict).Get("hello"); x != "world2" {
		t.Fatal("Mismatched inner", x)
	}
}
