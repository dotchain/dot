// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	mustPanic := func(message string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Failed to panic", message)
			}
		}()
		fn()
	}

	l := refs.List{V: types.A{types.S8("")}, R: map[interface{}]refs.Ref{}}
	c1, l2 := l.Add("first", refs.Path{0})
	if !reflect.DeepEqual(l2.R["first"], refs.Path{0}) {
		t.Error("Add failed", l2.R)
	}
	if !reflect.DeepEqual(l2, l.Apply(c1)) {
		t.Error("Incorrect change", c1)
	}

	mustPanic("duplicate add", func() {
		l2.Add("first", refs.Path{1})
	})

	c2, l3 := l2.Update("first", refs.Path{1})
	if !reflect.DeepEqual(l3, l2.Apply(c2)) {
		t.Error("Update failed", c2)
	}

	if !reflect.DeepEqual(refs.Path{1}, l3.R["first"]) {
		t.Error("Update failed", l3.R)
	}

	c3, l4 := l3.Remove("first")
	if !reflect.DeepEqual(l4, l3.Apply(c3)) {
		t.Error("Remove failed", c3)
	}

	if !reflect.DeepEqual(l, l4) {
		t.Error("Remove failed", l4.R)
	}

	mustPanic("non-existent removal", func() {
		l4.Remove("first")
	})
	mustPanic("non-existent update", func() {
		l4.Update("first", refs.Path{2})
	})
	mustPanic("count", func() {
		l4.Count()
	})
	mustPanic("slice", func() {
		l4.Slice(0, 0)
	})
}

func TestListApply(t *testing.T) {
	l := refs.List{V: types.A{types.S8("")}, R: map[interface{}]refs.Ref{}}
	_, l = l.Add("first", refs.Path{"Value", 0})

	// now insert an element at the first value
	insert := changes.Splice{0, types.A{}, types.A{types.S8("first")}}
	c := changes.PathChange{[]interface{}{"Value"}, insert}
	lx := l.Apply(c).(refs.List)

	// confirm that the ref is now refs.Path{"Value", 1}
	if !reflect.DeepEqual(lx.R["first"], refs.Path{"Value", 1}) {
		t.Error("Apply failed to modify ref", lx.R)
	}

	// try the change again but wrap the PathChange
	lx = l.Apply(changes.PathChange{nil, c}).(refs.List)
	if !reflect.DeepEqual(lx.R["first"], refs.Path{"Value", 1}) {
		t.Error("Apply failed to modify ref", lx.R)
	}
}

func TestListApplyReplace(t *testing.T) {
	l := refs.List{V: types.A{types.S8("")}}
	lx := l.Apply(changes.Replace{l, types.S8("hello")})
	if lx != types.S8("hello") {
		t.Error("Replace didn't do its thing", lx)
	}
}

func TestListApplyMisc(t *testing.T) {
	l := refs.List{V: types.A{types.S8("")}}
	if !reflect.DeepEqual(l, l.Apply(nil)) {
		t.Error("Unexpected nil failure")
	}

	lx := l.Apply(changes.ChangeSet{nil, nil})
	if !reflect.DeepEqual(lx, l) {
		t.Error("Unexpected nil failure", lx)
	}
}
