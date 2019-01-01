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

func TestCaretNil(t *testing.T) {
	ref := refs.Caret{refs.Path{"OK"}, 5, false}
	refx, cx := ref.Merge(nil)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestCaretEqual(t *testing.T) {
	x := refs.Caret{refs.Path{2}, 5, false}
	if !x.Equal(x) {
		t.Error("x != x")
	}
	if x.Equal(x.Path) {
		t.Error("refs.Caret equals refs.Path")
	}
	if x.Equal(refs.Caret{refs.Path{5}, 5, false}) {
		t.Error("path not tested")
	}
	if x.Equal(refs.Caret{x.Path, 6, false}) {
		t.Error("index not tested")
	}
	if x.Equal(refs.Caret{x.Path, 5, true}) {
		t.Error("index not tested")
	}
}

func TestCaretReplace(t *testing.T) {
	replace := changes.Replace{types.S8("OK"), types.S8("goop")}

	ref := refs.Caret{Index: 5}
	refx, cx := ref.Merge(replace)
	if refx != refs.InvalidRef || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = ref.Merge(changes.PathChange{nil, replace})
	if refx != refs.InvalidRef || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	change := changes.PathChange{[]interface{}{"xyz"}, replace}
	refx, cx = ref.Merge(change)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestCaretSplice(t *testing.T) {
	splice := changes.Splice{2, types.S8("OK"), types.S8("Boo")}

	ref := refs.Caret{nil, 1, false}
	refx, cx := ref.Merge(splice)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, ref, cx)
	}

	ref = refs.Caret{refs.Path{2}, 5, false}
	refx, cx = ref.Merge(splice)
	if refx != refs.InvalidRef || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{refs.Path{}, 3, false}
	refx, cx = ref.Merge(splice)
	expected := refs.Caret{nil, 2, false}
	if !reflect.DeepEqual(refx, expected) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{refs.Path{}, 4, false}
	refx, cx = ref.Merge(splice)
	expected = refs.Caret{nil, 5, false}
	if !reflect.DeepEqual(refx, expected) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestCaretSpliceAtIndex(t *testing.T) {
	splice := changes.Splice{5, types.S8(""), types.S8("boo")}
	ref := refs.Caret{refs.Path{}, splice.Offset, false}
	refx, _ := ref.Merge(splice)
	expected := refs.Caret{nil, splice.Offset + splice.After.Count(), false}

	if !reflect.DeepEqual(refx, expected) {
		t.Error("Unexpected Merge", refx)
	}

	splice = changes.Splice{5, types.S8("a"), types.S8("boo")}
	ref = refs.Caret{nil, splice.Offset, false}
	refx, _ = ref.Merge(splice)
	expected = ref

	if !reflect.DeepEqual(refx, expected) {
		t.Error("Unexpected Merge", refx)
	}

	splice = changes.Splice{5, types.S8(""), types.S8("boo")}
	ref = refs.Caret{nil, splice.Offset, true}
	refx, _ = ref.Merge(splice)
	expected = ref

	if !reflect.DeepEqual(refx, expected) {
		t.Error("Unexpected Merge", refx)
	}
}

func TestCaretMoveRight(t *testing.T) {
	move := changes.Move{2, 2, 2}

	p := refs.Path(nil)
	ref := refs.Caret{p, 1, false}
	refx, cx := ref.Merge(move)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{p, 3, false}
	refx, cx = ref.Merge(move)
	if !reflect.DeepEqual(refx, refs.Caret{p, 5, false}) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{p, 4, false}
	refx, cx = ref.Merge(move)
	if !reflect.DeepEqual(refx, refs.Caret{p, 2, false}) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{p, 7, false}
	refx, cx = ref.Merge(move)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestCaretMoveLeft(t *testing.T) {
	move := changes.Move{2, 2, -1}
	p := refs.Path(nil)
	ref := refs.Caret{p, 1, false}
	refx, cx := ref.Merge(move)
	if !reflect.DeepEqual(refx, refs.Caret{p, 3, false}) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{p, 3, false}
	refx, cx = ref.Merge(move)
	if !reflect.DeepEqual(refx, refs.Caret{p, 2, false}) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{p, 4, false}
	refx, cx = ref.Merge(move)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	ref = refs.Caret{p, 0, false}
	refx, cx = ref.Merge(move)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestCaretMoveAtIndex(t *testing.T) {
	move := changes.Move{2, 2, 2}

	p := refs.Path(nil)
	ref := refs.Caret{p, 6, false}
	refx, _ := ref.Merge(move)
	if !reflect.DeepEqual(refx, ref) {
		t.Error("Unexpected Merge", refx)
	}

	ref = refs.Caret{p, 6, true}
	refx, _ = ref.Merge(move)
	expected := refs.Caret{p, 4, true}
	if !reflect.DeepEqual(refx, expected) {
		t.Error("Unexpected Merge", refx)
	}

	move = changes.Move{2, 2, -1}

	ref = refs.Caret{p, 1, false}
	refx, _ = ref.Merge(move)
	expected = refs.Caret{p, 3, false}
	if !reflect.DeepEqual(refx, expected) {
		t.Error("Unexpected Merge", refx)
	}

	ref = refs.Caret{p, 1, true}
	refx, _ = ref.Merge(move)
	if !reflect.DeepEqual(refx, ref) {
		t.Error("Unexpected Merge", refx)
	}
}

func TestCaretChangeSet(t *testing.T) {
	move1 := changes.Move{2, 2, 1}
	move2 := changes.Move{3, 2, 5}
	p := refs.Path(nil)
	ref := refs.Caret{p, 2, false}
	refx, cx := ref.Merge(changes.ChangeSet{move1, move2})
	if !reflect.DeepEqual(refx, refs.Caret{p, 8, false}) || cx != nil {
		t.Error("Unexpected merge", refx, cx)
	}
}

func TestCaretPathChange(t *testing.T) {
	splice := changes.Splice{2, types.S8("OK"), types.S8("Boo")}

	p := refs.Path{"hello"}
	ref := refs.Caret{p, 4, false}
	refx, cx := ref.Merge(changes.PathChange{p, splice})
	if !reflect.DeepEqual(refx, refs.Caret{p, 5, false}) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = ref.Merge(changes.PathChange{p, changes.PathChange{nil, splice}})
	if !reflect.DeepEqual(refx, refs.Caret{p, 5, false}) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = ref.Merge(changes.PathChange{[]interface{}{"hello", 4}, splice})
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = ref.Merge(changes.PathChange{[]interface{}{"goop"}, splice})
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestCaretInvalidRef(t *testing.T) {
	p := refs.Path{"xyz"}
	ref := refs.Caret{p, 4, false}
	replace1 := changes.Replace{Before: types.S8("OK"), After: types.S8("Boo")}
	replace2 := changes.Replace{Before: types.S8("Boo"), After: types.S8("Goo")}
	cset := changes.ChangeSet{replace1, replace2}
	c := changes.PathChange{[]interface{}{"xyz"}, cset}
	refx, cx := ref.Merge(c)
	if refx != refs.InvalidRef || cx != nil {
		t.Error("Unexpected merge", refx, cx)
	}
}

func TestCaretUnknownChange(t *testing.T) {
	ref := refs.Caret{nil, 5, false}
	refx, cx := ref.Merge(myChange{})
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected merge with unknown change", refx, cx)
	}
}

func TestCaretMerger(t *testing.T) {
	p := refs.Path{42}
	ref := refs.Caret{p, 22, false}
	refx, cx := ref.Merge(caretMerger{})
	expected := refs.Caret{p, 1029, false}
	if !reflect.DeepEqual(refx, expected) || cx != nil {
		t.Errorf("Unexpected Merge %#v %#v", refx, cx)
	}
}

type caretMerger struct {
	myChange
}

func (cm caretMerger) Merge(o changes.Change) (changes.Change, changes.Change) {
	return o, cm
}

func (cm caretMerger) ReverseMerge(o changes.Change) (changes.Change, changes.Change) {
	return o, cm
}

func (cm caretMerger) Revert() changes.Change {
	return cm
}

func (cm caretMerger) MergePath(p []interface{}) *refs.MergeResult {
	return &refs.MergeResult{P: p, Scoped: cm, Affected: cm}
}

func (cm caretMerger) MergeCaret(caret refs.Caret) refs.Ref {
	caret.Index = 1029
	return caret
}
