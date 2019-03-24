// Copyright (C) 2018 rameshvk. All rights reserved.
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

func TestInvalidRef(t *testing.T) {
	c := changes.Move{Offset: 2, Count: 2, Distance: 2}

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
	move1 := changes.Move{Offset: 0, Count: 1, Distance: 10}
	move2 := changes.Move{Offset: 10, Count: 1, Distance: 10}
	c := changes.ChangeSet{move1, move2}
	r := refs.Merge([]interface{}{4}, c)
	expected := &refs.MergeResult{P: []interface{}{3}, Unaffected: c}
	if r == nil || !reflect.DeepEqual(r, expected) {
		t.Error("Merge failed", r)
	}

	replace := changes.Replace{Before: changes.Nil, After: types.S8("ok")}
	replace1 := changes.PathChange{Path: []interface{}{4, "key"}, Change: replace}
	replace2 := changes.PathChange{Path: []interface{}{3}, Change: replace}
	c = changes.ChangeSet{replace1, replace2}
	r = refs.Merge([]interface{}{4}, c)
	expected = &refs.MergeResult{
		P:          []interface{}{4},
		Scoped:     changes.PathChange{Path: []interface{}{"key"}, Change: replace},
		Affected:   replace1,
		Unaffected: replace2,
	}
	if r == nil || !reflect.DeepEqual(r, expected) {
		t.Error("Merge failed", r.Affected)
	}

	move3 := changes.Move{Offset: 2, Count: 2, Distance: 2}
	move4 := changes.Move{Offset: 1, Count: 1, Distance: 1}
	c = changes.ChangeSet{
		changes.PathChange{Path: []interface{}{4}, Change: move1},
		changes.ChangeSet{
			changes.PathChange{Path: []interface{}{4}, Change: move2},
			changes.PathChange{Path: []interface{}{4}, Change: move3},
		},
		changes.ChangeSet{
			changes.PathChange{Path: []interface{}{4}, Change: move2},
			changes.PathChange{Path: []interface{}{4}, Change: move3},
		},
		changes.PathChange{Path: []interface{}{4}, Change: move4},
	}
	r = refs.Merge([]interface{}{4}, c)
	cx := changes.ChangeSet{
		changes.PathChange{Path: []interface{}{}, Change: move1},
		changes.PathChange{Path: []interface{}{}, Change: move2},
		changes.PathChange{Path: []interface{}{}, Change: move3},
		changes.PathChange{Path: []interface{}{}, Change: move2},
		changes.PathChange{Path: []interface{}{}, Change: move3},
		changes.PathChange{Path: []interface{}{}, Change: move4},
	}
	expected = &refs.MergeResult{
		P:      []interface{}{4},
		Scoped: cx,
		Affected: changes.ChangeSet{
			changes.PathChange{Path: []interface{}{4}, Change: move1},
			changes.PathChange{Path: []interface{}{4}, Change: move2},
			changes.PathChange{Path: []interface{}{4}, Change: move3},
			changes.PathChange{Path: []interface{}{4}, Change: move2},
			changes.PathChange{Path: []interface{}{4}, Change: move3},
			changes.PathChange{Path: []interface{}{4}, Change: move4},
		},
	}

	if r == nil || !reflect.DeepEqual(r, expected) {
		t.Error("Merge failed", r)
	}
}
