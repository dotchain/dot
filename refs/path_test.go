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

func TestPathNil(t *testing.T) {
	p := refs.Path{"OK"}
	px, cx := p.Merge(nil)
	if !reflect.DeepEqual(px, p) || cx != nil {
		t.Error("Unexpected Merge", px, cx)
	}
}

func TestPathReplace(t *testing.T) {
	replace := changes.Replace{types.S8("OK"), types.S8("goop")}

	p := refs.Path(nil)
	p2, cx := p.Merge(replace)
	if p2 != refs.InvalidRef || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	change := changes.PathChange{[]interface{}{"xyz"}, replace}
	p2, cx = p.Merge(change)
	if !reflect.DeepEqual(p2, p) || !reflect.DeepEqual(cx, change) {
		t.Error("Unexpected Merge", p2, cx)
	}
}

func TestPathSplice(t *testing.T) {
	splice := changes.Splice{2, types.S8("OK"), types.S8("Boo")}

	p := refs.Path{1}
	p2, cx := p.Merge(splice)
	if !reflect.DeepEqual(p2, p) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{2}
	p2, cx = p.Merge(splice)
	if p2 != refs.InvalidRef || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{4}
	p2, cx = p.Merge(splice)
	if !reflect.DeepEqual(p2, refs.Path{5}) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{}
	p2, cx = p.Merge(splice)
	if !reflect.DeepEqual(p2, p) || cx != splice {
		t.Error("Unexpected Merge", p2, cx)
	}
}

func TestPathMoveRight(t *testing.T) {
	move := changes.Move{2, 2, 2}

	p := refs.Path{1}
	p2, cx := p.Merge(move)
	if !reflect.DeepEqual(p2, p) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{3}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, refs.Path{5}) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{4}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, refs.Path{2}) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{7}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, p) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}
}

func TestPathMoveLeft(t *testing.T) {
	move := changes.Move{2, 2, -1}
	p := refs.Path{1}
	p2, cx := p.Merge(move)
	if !reflect.DeepEqual(p2, refs.Path{3}) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{3}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, refs.Path{2}) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{4}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, p) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{0}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, p) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{}
	p2, cx = p.Merge(move)
	if !reflect.DeepEqual(p2, p) || cx != move {
		t.Error("Unexpected Merge", p2, cx)
	}
}

func TestPathChangeSet(t *testing.T) {
	move1 := changes.Move{2, 2, 1}
	move2 := changes.Move{3, 2, 5}
	p := refs.Path{2}
	px, cx := p.Merge(changes.ChangeSet{move1, move2})
	if !reflect.DeepEqual(px, refs.Path{8}) || cx != nil {
		t.Error("Unexpected merge", px, cx)
	}

	p = refs.Path{}
	moves := changes.ChangeSet{move1, move2}
	px, cx = p.Merge(moves)
	if !reflect.DeepEqual(px, p) || !reflect.DeepEqual(cx, moves) {
		t.Error("Unexpected merge", px, cx)
	}

	px, cx = p.Merge(changes.ChangeSet{move1})
	if !reflect.DeepEqual(px, p) || !reflect.DeepEqual(cx, move1) {
		t.Error("Unexpected merge", px, cx)
	}
}

func TestPathPathChange(t *testing.T) {
	splice := changes.Splice{2, types.S8("OK"), types.S8("Boo")}

	p := refs.Path{"hello", 4}
	p2, cx := p.Merge(changes.PathChange{[]interface{}{"hello"}, splice})
	if !reflect.DeepEqual(p2, refs.Path{"hello", 5}) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{"hello"}
	p2, cx = p.Merge(changes.PathChange{[]interface{}{"hello"}, splice})
	if !reflect.DeepEqual(p2, p) || cx != splice {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{"hello"}
	p2, cx = p.Merge(changes.PathChange{[]interface{}{"hello", "world"}, splice})
	expected := changes.PathChange{[]interface{}{"world"}, splice}
	if !reflect.DeepEqual(p2, p) || !reflect.DeepEqual(cx, expected) {
		t.Error("Unexpected Merge", p2, cx)
	}

	p = refs.Path{"goop"}
	p2, cx = p.Merge(changes.PathChange{[]interface{}{"hello"}, splice})
	if !reflect.DeepEqual(p2, p) || cx != nil {
		t.Error("Unexpected Merge", p2, cx)
	}
}

func TestPathInvalidRef(t *testing.T) {
	p := refs.Path{"xyz", 4}
	replace1 := changes.Replace{Before: types.S8("OK"), After: types.S8("Boo")}
	replace2 := changes.Replace{Before: types.S8("Boo"), After: types.S8("Goo")}
	cset := changes.ChangeSet{replace1, replace2}
	c := changes.PathChange{[]interface{}{"xyz"}, cset}
	px, cx := p.Merge(c)
	if px != refs.InvalidRef || cx != nil {
		t.Error("Unexpected merge", px, cx)
	}
}

func TestPathUnknownChange(t *testing.T) {
	p := refs.Path{}
	px, cx := p.Merge(myChange{})
	if !reflect.DeepEqual(px, p) || cx != nil {
		t.Error("Unexpected merge with unknown change", px, cx)
	}
}

func TestPathMerger(t *testing.T) {
	p := refs.Path{}
	pm := pathMerger{myChange{}}
	px, cx := p.Merge(pm)
	if !reflect.DeepEqual(px, refs.Path{"OK"}) || cx != pm {
		t.Error("Unexpected Merge", px, cx)
	}
}

type myChange struct{}

func (m myChange) Merge(o changes.Change) (changes.Change, changes.Change) {
	return o, m
}

func (m myChange) ReverseMerge(o changes.Change) (changes.Change, changes.Change) {
	return o, m
}

func (m myChange) Revert() changes.Change {
	return m
}

type pathMerger struct {
	myChange
}

func (p pathMerger) MergePath(path refs.Path) (refs.Ref, changes.Change) {
	return refs.Path{"OK"}, p
}
