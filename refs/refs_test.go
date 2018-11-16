// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/changes/types"

	"reflect"
	"testing"
)

func TestInvalidRef(t *testing.T) {
	c := changes.Move{2, 2, 2}

	r, cx := refs.InvalidRef.Merge(c)
	if r != refs.InvalidRef || cx != nil {
		t.Error("InvalidRef failed merge", r, cx)
	}

	if !refs.InvalidRef.Equal(refs.InvalidRef) {
		t.Error("InvalidRef does not equal itself")
	}

	if refs.InvalidRef.Equal(refs.Path{5}) {
		t.Error("InvalidRef equals something else")
	}
}

func TestMerge(t *testing.T) {
	move1 := changes.Move{0, 1, 10}
	move2 := changes.Move{10, 1, 10}
	c := changes.ChangeSet{move1, move2}
	r := refs.Merge([]interface{}{4}, c)
	expected := &refs.MergeResult{P: []interface{}{3}, Unaffected: c}
	if r == nil || !reflect.DeepEqual(r, expected) {
		t.Error("Merge failed", r)
	}

	replace := changes.Replace{changes.Nil, types.S8("ok")}
	replace1 := changes.PathChange{[]interface{}{4, "key"}, replace}
	replace2 := changes.PathChange{[]interface{}{3}, replace}
	c = changes.ChangeSet{replace1, replace2}
	r = refs.Merge([]interface{}{4}, c)
	expected = &refs.MergeResult{
		P:          []interface{}{4},
		Scoped:     changes.PathChange{[]interface{}{"key"}, replace},
		Affected:   replace1,
		Unaffected: replace2,
	}
	if r == nil || !reflect.DeepEqual(r, expected) {
		t.Error("Merge failed", r.Affected)
	}

	move3 := changes.Move{2, 2, 2}
	move4 := changes.Move{1, 1, 1}
	c = changes.ChangeSet{
		changes.PathChange{[]interface{}{4}, move1},
		changes.ChangeSet{
			changes.PathChange{[]interface{}{4}, move2},
			changes.PathChange{[]interface{}{4}, move3},
		},
		changes.ChangeSet{
			changes.PathChange{[]interface{}{4}, move2},
			changes.PathChange{[]interface{}{4}, move3},
		},
		changes.PathChange{[]interface{}{4}, move4},
	}
	r = refs.Merge([]interface{}{4}, c)
	cx := changes.ChangeSet{
		changes.PathChange{[]interface{}{}, move1},
		changes.PathChange{[]interface{}{}, move2},
		changes.PathChange{[]interface{}{}, move3},
		changes.PathChange{[]interface{}{}, move2},
		changes.PathChange{[]interface{}{}, move3},
		changes.PathChange{[]interface{}{}, move4},
	}
	expected = &refs.MergeResult{
		P:      []interface{}{4},
		Scoped: cx,
		Affected: changes.ChangeSet{
			changes.PathChange{[]interface{}{4}, move1},
			changes.PathChange{[]interface{}{4}, move2},
			changes.PathChange{[]interface{}{4}, move3},
			changes.PathChange{[]interface{}{4}, move2},
			changes.PathChange{[]interface{}{4}, move3},
			changes.PathChange{[]interface{}{4}, move4},
		},
	}

	if r == nil || !reflect.DeepEqual(r, expected) {
		t.Error("Merge failed", r)
	}
}
