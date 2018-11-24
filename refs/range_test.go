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

func TestRangeNil(t *testing.T) {
	p := refs.Path(nil)
	ref := refs.Range{refs.Caret{p, 5, false}, refs.Caret{p, 10, false}}
	refx, cx := ref.Merge(nil)
	if !reflect.DeepEqual(refx, ref) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}

func TestRangeEqual(t *testing.T) {
	c := refs.Caret{refs.Path{2}, 5, false}
	x := refs.Range{c, c}
	if !x.Equal(x) {
		t.Error("x != x")
	}
	if x.Equal(c) {
		t.Error("refs.Caret equals refs.Path")
	}
	cx := refs.Caret{nil, 5, false}
	if x.Equal(refs.Range{c, cx}) {
		t.Error("End not tested")
	}
	if x.Equal(refs.Range{cx, c}) {
		t.Error("start not tested")
	}
}

func TestRangeReplace(t *testing.T) {
	replace := changes.Replace{types.S8("OK"), types.S8("goop")}

	ref := refs.Range{refs.Caret{nil, 5, false}, refs.Caret{nil, 6, false}}
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

func TestRangeSplice(t *testing.T) {
	splice := changes.Splice{5, types.S8("OK"), types.S8("Boo")}
	newRange := func(start, end int) refs.Range {
		sleft, eleft := false, true
		if start == end {
			eleft = false
		}
		return refs.Range{refs.Caret{nil, start, sleft}, refs.Caret{nil, end, eleft}}
	}

	refx, cx := newRange(1, 5).Merge(splice)
	if !reflect.DeepEqual(refx, newRange(1, 5)) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = newRange(7, 9).Merge(splice)
	if !reflect.DeepEqual(refx, newRange(8, 10)) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = newRange(5, 7).Merge(splice)
	if !reflect.DeepEqual(refx, newRange(5, 8)) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	refx, cx = newRange(6, 7).Merge(splice)
	if !reflect.DeepEqual(refx, newRange(5, 8)) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}

	splice = changes.Splice{5, types.S8("OK"), types.S8("")}
	refx, cx = newRange(5, 7).Merge(splice)
	if !reflect.DeepEqual(refx, newRange(5, 5)) || cx != nil {
		t.Error("Unexpected Merge", refx, cx)
	}
}
