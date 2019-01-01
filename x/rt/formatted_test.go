// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rt_test

import (
	//	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rt"
	"testing"
)

func TestFormattedSliceSimple(t *testing.T) {
	f := rt.Formatted{Text: types.S8("hello")}
	slice := f.Slice(0, 4)
	if !slice.Equal(rt.Formatted{Text: types.S8("hell")}) {
		t.Error("Unexpected slice", slice)
	}
	slice = f.Slice(0, 0)
	if !slice.Equal(rt.Formatted{Text: types.S8("")}) {
		t.Error("Unexpected slice", slice)
	}
}

func TestFormattedSliceWithMultipleSegments(t *testing.T) {
	f := rt.Formatted{Text: types.S8("hello")}
	f = f.UpdateSlice(2, 2, func(inner rt.Formatted) rt.Formatted {
		return inner.UpdateInline(rt.Styles{"font": "x"}, nil)
	})

	slice1 := f.Slice(0, 3)
	slice2 := f.Slice(3, 2)
	joined := slice1.Concat(slice2)
	if !joined.Equal(f) {
		t.Error("Unexpected join", slice1, slice2, joined, f)
	}
}

func TestFormattedConcat(t *testing.T) {
	f := rt.Formatted{Text: types.S8("hello")}
	f = f.UpdateSlice(0, 2, func(inner rt.Formatted) rt.Formatted {
		return inner.UpdateInline(rt.Styles{"font": "x"}, nil)
	})
	f = f.UpdateSlice(2, 2, func(inner rt.Formatted) rt.Formatted {
		return inner.UpdateInline(rt.Styles{"font": "y"}, nil)
	})

	slice0 := f.Slice(0, 0)
	slice1 := f.Slice(0, 1)
	slice2 := f.Slice(1, 2)
	slice3 := f.Slice(3, 1)
	slice4 := rt.Formatted{Text: types.S8("o")}
	slice5 := rt.Formatted{Text: types.S8("")}

	joined := slice0.Concat(slice1, slice2, slice3, slice4, slice5)
	if !joined.Equal(f) {
		t.Error("Unexpected concat", joined, f)
	}
}

func TestFormattedEqual(t *testing.T) {
	t1 := rt.Formatted{Text: types.S8("hello")}
	t2 := rt.Formatted{Text: types.S8("world")}
	if t1.Equal(t2) {
		t.Error("False equality", t1, t2)
	}

	t1x := t1.Slice(0, 5)
	if !t1.Equal(t1x) {
		t.Error("False inequality", t1, t1x)
	}

	t1x = t1.UpdateSlice(0, 2, func(inner rt.Formatted) rt.Formatted {
		return inner.UpdateInline(rt.Styles{"font": "x"}, nil)
	})
	t1y := t1.UpdateSlice(0, 2, func(inner rt.Formatted) rt.Formatted {
		return inner.UpdateInline(rt.Styles{"font": "y"}, nil)
	})

	if t1x.Equal(t1y) {
		t.Error("False equality", t1x, t1y)
	}

	if !t1x.Equal(t1x) {
		t.Error("false inequality", t1x)
	}

	if t1x.Equal(t1) {
		t.Error("false inequality", t1x, t1)
	}
}

func TestFormattedMove(t *testing.T) {
	affix := rt.Formatted{Text: types.S8("x")}
	t1 := (rt.Formatted{Text: types.S8("hello")}).
		UpdateInline(rt.Styles{"font": "x"}, nil).
		UpdateInline(rt.Styles{"color": "blue"}, nil)
	t2 := (rt.Formatted{Text: types.S8("world")}).
		UpdateInline(rt.Styles{"font": "y"}, nil).
		UpdateInline(rt.Styles{"color": "red"}, nil)
	t3 := (rt.Formatted{Text: types.S8("goopy")}).
		UpdateInline(rt.Styles{"font": "z"}, nil).
		UpdateInline(rt.Styles{"color": "green"}, nil)
	tx := affix.Concat(t1, t2, t3, affix)

	moved1 := tx.Move(1, 5, 5)
	moved2 := tx.Move(6, 5, -5)
	expected := affix.Concat(t2, t1, t3, affix)
	if !expected.Equal(moved1) {
		t.Error("Unexpected move result", moved1)
	}
	if !expected.Equal(moved2) {
		t.Error("Unexpected move result", moved2)
	}
}

func TestInsertUpdateBlock(t *testing.T) {
	v1 := rt.Formatted{Text: types.S8("hello")}
	v2 := rt.Formatted{Text: types.S8("world")}
	styles1 := rt.Styles{"font": "x", "size": 100}
	styles2 := rt.Styles{"font": "y", "size": 100}

	v1 = v1.InsertBlock(0, styles1)
	v2 = v2.InsertBlock(0, styles2)
	merged := v1.Concat(v2)

	updated := merged.UpdateBlock(0, nil, []string{"font"})
	expected := rt.Formatted{Text: types.S8("helloworld")}
	expected = expected.InsertBlock(0, rt.Styles{"size": 100})
	if !expected.Equal(updated) {
		t.Error("Unexpected InsertBlock failure", updated)
	}
}

func TestRemoveBlock(t *testing.T) {
	v := rt.Formatted{Text: types.S8("hello")}
	styles := rt.Styles{"font": "ax", "size": 100}

	vx := v.InsertBlock(0, styles).RemoveBlock(0)
	if !v.Equal(vx) {
		t.Error("Unexpected InsertBlock failure", vx)
	}
}
