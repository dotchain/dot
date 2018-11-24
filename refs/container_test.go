// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/refs"
	"reflect"
	"testing"
)

func TestUpdateMerge(t *testing.T) {
	initial := refs.NewContainer(
		types.M{
			"ok":    changes.Nil,
			"bok":   changes.Nil,
			"loki":  changes.Nil,
			"array": types.A{changes.Nil, changes.Nil, changes.Nil, changes.Nil, changes.Nil},
		},
		nil,
	)
	initial, _ = initial.UpdateRef("goo", refs.Path{"Value", "ok"})
	initial, _ = initial.UpdateRef("array", refs.Path{"Value", "array", 1})

	pairs := [][2]changes.Change{
		{
			refs.Update{"goo", refs.Path{"Value", "ok"}, refs.Path{"Value", "array"}},
			nil,
		},
		{
			refs.Update{"array", refs.Path{"Value", "array", 1}, refs.Path{"Value", "array"}},
			changes.PathChange{refs.Path{"Value"}, changes.Replace{initial.Value, types.S8("ok")}},
		},
		{
			refs.Update{"goo", refs.Path{"Value", "ok"}, refs.Path{"Value", "bok"}},
			refs.Update{"goo", refs.Path{"Value", "ok"}, refs.Path{"Value", "loki"}},
		},
		{
			refs.Update{"goo", refs.Path{"Value", "ok"}, refs.Path{"Value", "bok"}},
			refs.Update{"goo", refs.Path{"Value", "ok"}, nil},
		},
		{
			refs.Update{"goop", nil, refs.Path{"Value", "bok"}},
			refs.Update{"goop", nil, refs.Path{"Value", "loki"}},
		},
		{
			refs.Update{"gool", nil, refs.Path{"Value", "bok"}},
			refs.Update{"good", nil, refs.Path{"Value", "loki"}},
		},
		{
			refs.Update{"gool", nil, refs.Path{"Value", "bok"}},
			changes.Replace{initial, types.S8("hello")},
		},
		{
			refs.Update{"gool", nil, refs.Path{"Value", "array", 5}},
			changes.PathChange{refs.Path{"Value", "array"}, changes.Move{0, 4, 1}},
		},
		{
			refs.Update{"gool", nil, refs.Path{"Value", "array", 5}},
			changes.PathChange{refs.Path{"Value", "array"}, changes.Move{0, 4, 1}},
		},
	}

	for _, pair := range pairs {
		c1, c2 := pair[0], pair[1]
		c1x, c2x := c1.Merge(c2)
		final1 := initial.Apply(c1).Apply(c1x)
		final2 := initial.Apply(c2).Apply(c2x)
		if !reflect.DeepEqual(final1, final2) {
			t.Error("Failed to merge", pair)
		}

		c1y, c2y := c1.Merge(changes.PathChange{nil, changes.ChangeSet{c2}})
		final1y := initial.Apply(c1).Apply(c1y)
		final2y := initial.Apply(c2).Apply(c2y)
		if !reflect.DeepEqual(final1y, final2y) || !reflect.DeepEqual(final1, final1y) {
			t.Error("Failed to merge", pair)
		}

		if custom, ok := c1.(changes.Custom); ok {
			c1z, c2z := custom.ReverseMerge(changes.PathChange{nil, changes.ChangeSet{c2}})
			final1z := initial.Apply(c1).Apply(c1z)
			final2z := initial.Apply(c2).Apply(c2z)
			if !reflect.DeepEqual(final1z, final2z) {
				t.Error("Failed to merge", pair)
			}
		}

		if custom, ok := c2.(changes.Custom); ok {
			c2x, c1x = custom.ReverseMerge(c1)
			final1 = initial.Apply(c1).Apply(c1x)
			final2 = initial.Apply(c2).Apply(c2x)
			if !reflect.DeepEqual(final1, final2) {
				t.Error("Failed to reverse merge", pair)
			}
		}
	}
}

func TestUpdateUnknownMerge(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Failed to panic")
		}
	}()
	u := refs.Update{"goo", refs.Path{"ok"}, refs.Path{"q"}}
	u.Merge(changes.Move{5, 2, 2})
}

func TestUpdateMiscReverseMerge(t *testing.T) {
	u := refs.Update{"goo", refs.Path{"ok"}, refs.Path{"q"}}
	x, y := u.ReverseMerge(nil)
	if !reflect.DeepEqual(y, u) || x != nil {
		t.Error("nil reverse merge failed", x, y)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Failed to panic")
		}
	}()
	u.ReverseMerge(changes.Move{5, 2, 2})
}

func TestUpdateApplyTo(t *testing.T) {
	initial := refs.Container{Value: types.S8("")}
	u := refs.Update{"goo", nil, refs.Caret{refs.Path{"Value"}, 5, false}}
	updated := initial.Apply(u)
	alt := u.ApplyTo(initial)
	if !reflect.DeepEqual(updated, alt) {
		t.Error("Unexpected ApplyTo", updated, alt)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Failed to panic")
		}
	}()
	u = refs.Update{"goo", refs.Path{"ok"}, refs.Path{"q"}}
	u.ApplyTo(types.S8(""))
}

func TestUpdateRevert(t *testing.T) {
	initial := refs.Container{
		Value: types.M{
			"ok":    changes.Nil,
			"bok":   changes.Nil,
			"loki":  changes.Nil,
			"array": types.A{changes.Nil, changes.Nil, changes.Nil, changes.Nil, changes.Nil},
		},
	}
	initial, _ = initial.UpdateRef("goo", refs.Path{"Value", "ok"})
	initial, _ = initial.UpdateRef("array", refs.Path{"Value", "array", 1})

	changes := []refs.Update{
		{"goo", refs.Path{"Value", "ok"}, nil},
		{"goo", refs.Path{"Value", "ok"}, refs.Path{"Value", "array", 1}},
		{"boo", nil, refs.Path{"Value", "ok"}},
	}
	for _, ch := range changes {
		reverted := initial.Apply(ch).Apply(ch.Revert())
		if !reflect.DeepEqual(initial, reverted) {
			t.Error("Failed to revert", ch)
		}
	}
}

func TestContainer(t *testing.T) {
	initial := refs.Container{Value: types.S8("OK")}
	x, c := initial.UpdateRef("boo", nil)
	if !reflect.DeepEqual(x, initial) || c != nil {
		t.Error("nil update", x, c)
	}

	v, _ := initial.UpdateRef("boo", refs.Path{"Value"})
	p := v.GetRef("boo")
	if !reflect.DeepEqual(p, refs.Path{"Value"}) {
		t.Error("GetRef failed")
	}

	vx := initial.Apply(changes.ChangeSet{refs.Update{"boo", nil, refs.Path{"Value"}}})
	if !reflect.DeepEqual(vx, v) {
		t.Error("Apply failed", vx, v)
	}

	v, _ = v.UpdateRef("boo", nil)
	if !reflect.DeepEqual(v.Refs(), initial.Refs()) {
		t.Error("Removing refs failed", v)
	}
}

func TestContainerPanics(t *testing.T) {
	mustPanic := func(msg string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Failed to panic", msg)
			}
		}()
		fn()
	}

	con := refs.Container{Value: types.S8("hello")}
	mustPanic("bad apply", func() {
		con.Apply(changes.Move{2, 2, 2})
	})
	mustPanic("bad path apply", func() {
		con.Apply(changes.PathChange{refs.Path{"zoo"}, changes.Move{2, 2, 2}})
	})
}
