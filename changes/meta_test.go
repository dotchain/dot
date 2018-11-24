// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"github.com/dotchain/dot/changes"
	"reflect"
	"testing"
)

func TestMetaRevert(t *testing.T) {
	m := changes.Meta{"hello", changes.Move{1, 2, 3}}
	expected := changes.Meta{"hello", m.Change.Revert()}
	if !reflect.DeepEqual(m.Revert(), expected) {
		t.Error("Unexpected revert", m.Revert())
	}
}

func TestConvergenceMeta(t *testing.T) {
	validate := func(initial changes.Value, left, right changes.Change) {
		leftx, rightx := left.Merge(right)
		if custom, ok := right.(changes.Custom); ok && custom != nil {
			revr, revx := custom.ReverseMerge(left)
			if !reflect.DeepEqual(revr, rightx) {
				t.Error("ReverseMerge", revr, rightx)
			}
			if !reflect.DeepEqual(revx, simplify(leftx)) {
				t.Errorf("ReverseMerge %#v %#v\n", revx, leftx)
			}
		}

		lval := initial.Apply(nil, changes.ChangeSet{left, leftx})
		rval := initial.Apply(nil, changes.ChangeSet{right, rightx})
		if lval != rval {
			t.Error("Diverged", lval, rval, left, right)
		}
	}

	ForEachChange(S("xyz"), func(initial changes.Value, left changes.Change) {
		ForEachChange(S("ab"), func(_ changes.Value, right changes.Change) {
			validate(initial, changes.Meta{"ok", left}, right)
			validate(initial, changes.Meta{"ok", left}, changes.Meta{"boo", right})
			validate(initial, right, changes.Meta{"boo", left})
		})
	})
}
