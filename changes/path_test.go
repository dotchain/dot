// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"github.com/dotchain/dot/changes"
	"reflect"
	"testing"
)

type Path []interface{}

func TestPathChangeRevert(t *testing.T) {
	pc := changes.PathChange{Path{"a"}, changes.Move{1, 2, 3}}
	expected := changes.PathChange{Path{"a"}, pc.Change.Revert()}
	if !reflect.DeepEqual(pc.Revert(), expected) {
		t.Error("Unexpected revert", pc.Revert())
	}

	pc = changes.PathChange{Path{"a"}, nil}
	if pc.Revert() != nil {
		t.Error("Unexpected revert", pc.Revert())
	}
}

func TestPathChangeReverseMergeSimple(t *testing.T) {
	pc := changes.PathChange{nil, changes.Replace{Before: S("ab"), After: S("bc")}}
	o := changes.PathChange{nil, changes.Move{1, 1, -1}}

	lx, rx := pc.ReverseMerge(o)
	lx = simplify(lx)
	rx = simplify(rx)
	if lx != nil || rx != (changes.Replace{Before: S("ba"), After: S("bc")}) {
		t.Error("Unexpected merge output", lx, rx)
	}
}

func TestPathChangeDifferentPaths(t *testing.T) {
	left := changes.Replace{Before: S("a"), After: S("b")}
	right := changes.Replace{Before: S("c"), After: S("d")}
	l := changes.PathChange{Path{"q", 2, "a"}, left}
	r := changes.PathChange{Path{"q", 9, "a"}, right}
	validateMergeResults(t, l, r, r, l)
}

func TestPathChangeNil(t *testing.T) {
	left := changes.Replace{Before: S("a"), After: S("b")}
	l := changes.PathChange{Path{"q", 9, "aa"}, left}
	r := changes.PathChange{Path{"q", 9}, nil}
	validateMergeResults(t, l, r, nil, l)
}

func TestPathChangeMergeReplace(t *testing.T) {
	initial := A{A{S("a"), S("b")}, A{S("c"), S("d")}}
	left := changes.Replace{Before: initial, After: A{S("Boo")}}

	right := changes.PathChange{Path{0}, changes.Replace{Before: initial[0], After: A{S("")}}}
	lexpected := changes.Change(nil)
	rexpected := changes.Replace{Before: initial.Apply(nil, right), After: left.After}
	validateMergeResults(t, left, right, lexpected, rexpected)
}

func TestPathChangeMergeSpliceLeft(t *testing.T) {
	initial := A{A{S("a"), S("b")}, A{S("c"), S("d")}}
	left := changes.Splice{1, initial.Slice(1, 1), S("Boo")}

	right := changes.PathChange{Path{0}, changes.Replace{Before: initial[0], After: S("")}}

	validateMergeResults(t, left, right, right, left)
}

func TestPathChangeMergeSpliceMiddle(t *testing.T) {
	initial := A{A{S("a"), S("b")}, A{S("c"), S("d")}}
	left := changes.Splice{1, initial.Slice(1, 1), A{S("Boo")}}

	right := changes.PathChange{Path{1}, changes.Replace{Before: initial[1], After: S("zoog")}}

	lexpected := changes.Change(nil)
	rexpected := changes.Splice{1, initial.ApplyCollection(nil, right).Slice(1, 1), left.After}

	validateMergeResults(t, left, right, lexpected, rexpected)
}

func TestPathChangeMergeSpliceRight(t *testing.T) {
	initial := A{A{S("a"), S("b")}, A{S("c"), S("d")}}
	left := changes.Splice{0, initial.Slice(0, 1), A{S("Boo"), S("Boo2"), S("Boo3")}}

	right := changes.PathChange{Path{1, 0}, changes.Replace{Before: S("c"), After: S("zoog")}}

	lexpected := changes.PathChange{Path{3, 0}, changes.Replace{Before: S("c"), After: S("zoog")}}
	rexpected := left

	validateMergeResults(t, left, right, lexpected, rexpected)
}

func TestPathChangeMergeMoveRight(t *testing.T) {
	initial := A{A{S("A"), S("a")}, A{S("B"), S("b")}, A{S("C"), S("c")}, A{S("D"), S("d")}, A{S("E"), S("e")}, A{S("F"), S("f")}}
	left := changes.Move{1, 2, 1} // abcdef => adbcef
	for idx := 0; idx < initial.Count(); idx++ {
		before := S(("abcdef")[idx : idx+1])
		re := changes.Replace{Before: before, After: S("boo")}
		right := changes.PathChange{Path{idx, 1}, re}
		lexpected := right
		rexpected := left
		if idx >= 1 && idx < 3 {
			lexpected.Path = Path{idx + 1, 1}
		} else if idx == 3 {
			lexpected.Path = Path{1, 1}
		}

		validateMergeResults(t, left, right, lexpected, rexpected)
	}
}

func TestPathChangeMergeMoveLeft(t *testing.T) {
	initial := A{A{S("A"), S("a")}, A{S("B"), S("b")}, A{S("C"), S("c")}, A{S("D"), S("d")}, A{S("E"), S("e")}, A{S("F"), S("f")}}
	left := changes.Move{3, 1, -2} // abcdef => adbcef
	for idx := 0; idx < initial.Count(); idx++ {
		before := S(("abcdef")[idx : idx+1])
		re := changes.Replace{Before: before, After: S("boo")}
		right := changes.PathChange{Path{idx, 1}, re}
		lexpected := right
		rexpected := left
		if idx >= 1 && idx < 3 {
			lexpected.Path = Path{idx + 1, 1}
		} else if idx == 3 {
			lexpected.Path = Path{1, 1}
		}

		validateMergeResults(t, left, right, lexpected, rexpected)
	}
}

func validateMergeResults(t *testing.T, l, r, lexpected, rexpected changes.Change) {
	validateMergeResults1(t, l, r, lexpected, rexpected)
	validateMergeResults1(t, changes.PathChange{nil, l}, changes.PathChange{nil, r}, lexpected, rexpected)
	validateMergeResults1(t, changes.PathChange{nil, l}, r, lexpected, rexpected)
	validateMergeResults1(t, l, changes.PathChange{nil, r}, lexpected, rexpected)
	lx := simplify(changes.PathChange{Path{"hello"}, lexpected})
	rx := simplify(changes.PathChange{Path{"hello"}, rexpected})
	left := simplify(changes.PathChange{Path{"hello"}, l})
	right := simplify(changes.PathChange{Path{"hello"}, r})
	validateMergeResults1(t, left, right, lx, rx)
	validateMergeResults1(t, changes.PathChange{nil, left}, right, lx, rx)
	validateMergeResults1(t, changes.ChangeSet{left}, right, lx, rx)
}

func validateMergeResults1(t *testing.T, l, r, lexpected, rexpected changes.Change) {
	lx, rx := l.Merge(r)
	if !reflect.DeepEqual(simplify(lx), lexpected) || !reflect.DeepEqual(simplify(rx), rexpected) {
		t.Error("Unexpected l, r", lx, rx)
	}

	if r == nil {
		return
	}

	rx, lx = r.Merge(l)
	if !reflect.DeepEqual(simplify(lx), lexpected) || !reflect.DeepEqual(simplify(rx), rexpected) {
		t.Error("Unexpected l, r", lx, rx)
	}

	if rev, ok := r.(changes.Custom); ok {
		rx, lx := rev.ReverseMerge(l)
		if !reflect.DeepEqual(simplify(lx), lexpected) || !reflect.DeepEqual(simplify(rx), rexpected) {
			t.Error("Unexpected l, r", lx, rx)
		}
	}
}

func simplify(c changes.Change) changes.Change {
	switch c := c.(type) {
	case nil:
		return nil
	case changes.ChangeSet:
		if len(c) == 0 {
			return nil
		}
		if len(c) == 1 {
			return simplify(c[0])
		}
	case changes.PathChange:
		if cx := simplify(c.Change); cx == nil {
			return nil
		} else {
			c.Change = cx
		}

		if len(c.Path) == 0 {
			return c.Change
		}

		if pc, ok := c.Change.(changes.PathChange); ok {
			c.Path = append([]interface{}(nil), c.Path...)
			c.Path = append(c.Path, pc.Path...)
			c.Change = pc.Change
		}
		return c
	}
	return c
}
