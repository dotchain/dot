// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/changes/run"
	"github.com/dotchain/dot/changes/types"
	"reflect"
	"testing"
)

type S = types.S8
type A = types.A

func TestRunRevert(t *testing.T) {
	r := run.Run{0, 5, nil}
	if r.Revert() != nil {
		t.Error("Unexpected revert", r.Revert())
	}

	r = run.Run{0, 5, changes.Move{1, 2, 2}}
	expected := run.Run{0, 5, r.Change.Revert()}
	if r.Revert() != expected {
		t.Error("Unexpected revert", r.Revert())
	}
}

func TestRunMergeNils(t *testing.T) {
	r := run.Run{0, 5, changes.Move{5, 5, 1}}
	ox, rx := r.Merge(nil)
	if ox != nil || rx != r {
		t.Error("Unexpected merge", ox, rx)
	}

	ox, rx = r.ReverseMerge(nil)
	if ox != nil || rx != r {
		t.Error("Unexpected merge", ox, rx)
	}

	r.Change = nil
	ox, rx = r.Merge(changes.Move{5, 5, 1})
	if ox != (changes.Move{5, 5, 1}) || rx != nil {
		t.Error("Unexpected merge", ox, rx)
	}

	ox, rx = r.ReverseMerge(changes.Move{5, 5, 1})
	if ox != (changes.Move{5, 5, 1}) || rx != nil {
		t.Error("Unexpected merge", ox, rx)
	}
}

var first = A{S("a"), S("b")}
var second = A{S("c"), S("d")}
var third = A{S("e"), S("f")}
var initial = A{first, second, third}

func TestRunMergeReplace(t *testing.T) {
	l := run.Run{0, 2, changes.Move{0, 1, 1}}
	r := changes.Replace{Before: initial, After: S("OK")}
	validateMerge(t, l, r)
	validateMerge(t, l, changes.Replace{initial, changes.Nil})
}

func TestRunMergeSplice(t *testing.T) {
	l := run.Run{1, 1, changes.Move{0, 1, 1}}
	r := changes.Splice{Offset: 0, Before: A{}, After: A{S("OK")}}
	validateMerge(t, l, r)
	r = changes.Splice{Offset: 0, Before: A{first}, After: A{}}
	validateMerge(t, l, r)
	r = changes.Splice{Offset: 2, Before: A{}, After: A{S("OK")}}
	validateMerge(t, l, r)
	r = changes.Splice{Offset: 1, Before: A{second}, After: A{S("OK")}}
	validateMerge(t, l, r)
	r = changes.Splice{Offset: 0, Before: initial, After: A{S("OK")}}
	validateMerge(t, l, r)

	l = run.Run{0, 2, changes.Move{0, 1, 1}}
	r = changes.Splice{1, A{second, third}, A{S("OK")}}
	validateMerge(t, l, r)
	l = run.Run{1, 2, changes.Move{0, 1, 1}}
	r = changes.Splice{0, A{first, second}, A{S("OK")}}
	validateMerge(t, l, r)

	l = run.Run{0, 3, changes.Move{0, 1, 1}}
	r = changes.Splice{1, A{second}, A{S("OK")}}
	validateMerge(t, l, r)
}

func TestRunMergeMove(t *testing.T) {
	for count := 1; count < 3; count++ {
		for offset := 0; offset <= 3-count; offset++ {
			for dest := 0; dest <= 3; dest++ {
				if dest >= offset && dest <= offset+count {
					continue
				}
				r := changes.Move{offset, count, dest - offset - count}
				if dest < offset {
					r.Distance = dest - offset
				}
				l := run.Run{1, 1, changes.Move{0, 1, 1}}
				validateMerge(t, l, r)
				l = run.Run{0, 2, changes.Move{0, 1, 1}}
				validateMerge(t, l, r)
				l = run.Run{0, 3, changes.Move{0, 1, 1}}
				validateMerge(t, l, r)
			}
		}
	}
}

func TestRunMergeRun(t *testing.T) {
	ForEachRun(changes.Splice{0, A{}, A{S("Left")}}, func(l run.Run) {
		ForEachRun(changes.Splice{0, A{}, A{S("Right")}}, func(r run.Run) {
			validateMerge(t, l, r)
		})
	})
}

func TestRunMergePathChange(t *testing.T) {
	l := run.Run{1, 1, changes.Splice{0, A{}, A{S("Left")}}}
	r := changes.PathChange{nil, changes.Replace{initial, changes.Nil}}
	validateMerge(t, l, r)

	l = run.Run{1, 2, changes.Splice{0, A{}, A{S("Left")}}}
	for kk := 0; kk < 3; kk++ {
		replace := changes.Replace{initial[kk], changes.Nil}
		r := changes.PathChange{[]interface{}{kk}, replace}
		validateMerge(t, l, r)
		validateMerge(t, r, l)
	}
}

func TestRunMergePath(t *testing.T) {
	r := run.Run{5, 10, changes.Replace{S("OK"), S("boo")}}
	p := refs.Path(nil)
	px, cx := p.Merge(r)
	if !reflect.DeepEqual(px, p) || !reflect.DeepEqual(cx, r) {
		t.Fatal("Empty refs.Path merge failed", px, cx)
	}

	p = refs.Path{4}
	px, cx = p.Merge(r)
	if !reflect.DeepEqual(px, p) || cx != nil {
		t.Fatal("Unaffected refs.Path merge failed", px, cx)
	}

	p = refs.Path{15}
	px, cx = p.Merge(r)
	if !reflect.DeepEqual(px, p) || cx != nil {
		t.Fatal("Unaffected refs.Path merge failed", px, cx)
	}

	p = refs.Path{5, 0}
	px, cx = p.Merge(r)
	if px != refs.InvalidRef || cx != nil {
		t.Fatal("Affected refs.Path merge failed", px, cx)
	}

	p = refs.Path{14, "x"}
	move := changes.Move{2, 3, 4}
	r = run.Run{5, 10, changes.PathChange{[]interface{}{"x"}, move}}
	px, cx = p.Merge(r)
	if !reflect.DeepEqual(px, p) ||
		!reflect.DeepEqual(cx, changes.PathChange{[]interface{}{}, move}) {
		t.Fatal("Affected refs.Path merge failed", px, cx)
	}

	p = refs.Path{14, 2}
	r = run.Run{5, 10, move}
	px, cx = p.Merge(r)
	if !reflect.DeepEqual(px, refs.Path{14, 6}) || cx != nil {
		t.Fatal("Affected refs.Path merge failed", px, cx)
	}
}

func ForEachRun(change changes.Change, fn func(r run.Run)) {
	for count := 1; count <= 3; count++ {
		for offset := 0; offset <= 3-count; offset++ {
			fn(run.Run{offset, count, change})
		}
	}
}

func validateMerge(t *testing.T, l, r changes.Change) {
	validateMerge1(t, initial, l, r)
	validateMerge1(t, initial, l, changes.PathChange{nil, r})
	validateMerge1(t, initial, changes.PathChange{nil, l}, r)
	validateMerge1(t, initial, changes.ChangeSet{l}, r)
	validateMerge1(t, initial, l, changes.ChangeSet{r})
}

func validateMerge1(t *testing.T, initial changes.Value, l, r changes.Change) {
	lx, rx := l.Merge(r)
	lval := initial.Apply(l).Apply(lx)
	rval := initial.Apply(r).Apply(rx)
	if !reflect.DeepEqual(lval, rval) {
		t.Error("merge failed", lval, rval, "---", l, "---", r)
	}
	if rev, ok := r.(revMerge); ok {
		rx2, lx2 := rev.ReverseMerge(l)
		lx, rx, lx2, rx2 = simplify(lx), simplify(rx), simplify(lx2), simplify(rx2)
		if !reflect.DeepEqual(rx, rx2) || !reflect.DeepEqual(lx, lx2) {
			t.Error("reverse merge diverged from merge", lx, lx2, rx, rx2)
		}
	}
	if rev, ok := l.(revMerge); ok {
		lx, rx = rev.ReverseMerge(r)
		rx2, lx2 := r.Merge(l)

		lx, rx, lx2, rx2 = simplify(lx), simplify(rx), simplify(lx2), simplify(rx2)
		if !reflect.DeepEqual(rx, rx2) || !reflect.DeepEqual(lx, lx2) {
			t.Error("reverse merge diverged from merge", lx, lx2, rx, rx2)
		}
	}
}

type revMerge interface {
	ReverseMerge(changes.Change) (changes.Change, changes.Change)
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
