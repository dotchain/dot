// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rich_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func TestTextSetAttrRevert(t *testing.T) {
	s := rich.NewText("hello world")
	c := s.SetAttribute(3, 5, html.FontBold)
	s = s.Apply(nil, c).(rich.Text)
	if x := html.Format(s, nil); x != "hel<b>lo wo</b>rld" {
		t.Error("Unexpected", x)
	}

	s = s.Apply(nil, c.Revert()).(rich.Text)
	if x := html.Format(s, nil); x != "hello world" {
		t.Error("Unexpected", x)
	}

	s = s.Apply(nil, c.Revert().Revert()).(rich.Text)
	if x := html.Format(s, nil); x != "hel<b>lo wo</b>rld" {
		t.Error("Unexpected", x)
	}
}

func TestTextSetAttrMergeNil(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := s.SetAttribute(3, 5, html.FontBold)
	testMerge(t, s, c1, nil)
}

func TestTextSetAttrMergeReplace(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.Replace{Before: s, After: rich.NewText("boo hoo")}
	c2 := s.SetAttribute(3, 5, html.FontBold)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)
}

func TestTextSetAttrMergeSpliceNoConflict(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.Splice{
		Offset: 2,
		Before: rich.NewText("l"),
		After:  rich.NewText("---"),
	}
	c2 := s.SetAttribute(3, 5, html.FontBold)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = changes.Splice{
		Offset: 8,
		Before: rich.Text{},
		After:  rich.NewText("---"),
	}
	c2 = s.SetAttribute(3, 5, html.FontBold)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)
}

func TestTextSetAttrMergeSpliceWithin(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.Splice{
		Offset: 4,
		Before: rich.NewText("l"),
		After:  rich.NewText("---"),
	}
	c2 := s.SetAttribute(3, 5, html.FontBold)

	result := testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	if x := html.Format(result, nil); x != "hel<b>l--- wo</b>rld" {
		t.Error("Unexpected", x)
	}
}

func TestTextSetAttrMergeSpliceConflict(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.Splice{
		Offset: 4,
		Before: rich.NewText("o wor"),
		After:  rich.NewText("---"),
	}
	c2 := s.SetAttribute(2, 3, html.FontBold)

	result := testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	if x := html.Format(result, nil); x != "he<b>ll</b>---ld" {
		t.Error("Unexpected", x)
	}

	c2 = s.SetAttribute(7, 3, html.FontBold)

	result = testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	if x := html.Format(result, nil); x != "hell---<b>l</b>d" {
		t.Error("Unexpected", x)
	}
}

func TestTextSetAttrMergePathNoConflict(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.PathChange{
		Path: []interface{}{1, html.FontBold.Name()},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  html.FontBold,
		},
	}
	c2 := s.SetAttribute(3, 5, html.FontBold)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)
}

func TestTextSetAttrMergePathConflict(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.PathChange{
		Change: changes.PathChange{
			Path: []interface{}{4, html.FontWeight(200).Name()},
			Change: changes.Replace{
				Before: changes.Nil,
				After:  html.FontWeight(200),
			},
		},
	}
	c2 := s.SetAttribute(3, 5, html.FontBold)

	testMerge(t, s, c1, c2)
	testMerge(t, s, c2, c1)
	testReverseMerge(t, s, c1, c2)
}

func TestTextSetAttrMergeChangeSet(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := changes.ChangeSet{}
	c2 := s.SetAttribute(3, 5, html.FontBold)

	testMerge(t, s, c1, c2)
	testMerge(t, s, c2, c1)
	testReverseMerge(t, s, c1, c2)
}

func TestTextSetAttrMergeNoConflict(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := s.SetAttribute(3, 5, html.FontBold)
	c2 := s.SetAttribute(1, 2, html.FontBold)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = s.SetAttribute(1, 2, html.FontBold)
	c2 = s.SetAttribute(3, 5, html.FontBold)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = s.SetAttribute(3, 5, html.FontBold)
	c2 = s.SetAttribute(1, 5, html.FontStyleItalic)
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)
}

func TestTextSetAttrMergeConflict(t *testing.T) {
	s := rich.NewText("hello world")
	c1 := s.SetAttribute(3, 5, html.FontBold)
	c2 := s.SetAttribute(3, 5, html.FontWeight(100))
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = s.SetAttribute(3, 5, html.FontBold)
	c2 = s.SetAttribute(1, 4, html.FontWeight(100))
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = s.SetAttribute(3, 5, html.FontBold)
	c2 = s.SetAttribute(1, 6, html.FontWeight(100))
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = s.SetAttribute(3, 5, html.FontBold)
	c2 = s.SetAttribute(4, 6, html.FontWeight(100))
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)

	c1 = s.SetAttribute(3, 5, html.FontBold)
	c2 = s.SetAttribute(4, 2, html.FontWeight(100))
	testMerge(t, s, c1, c2)
	testReverseMerge(t, s, c1, c2)
}

func testMerge(t *testing.T, s rich.Text, c1, c2 changes.Change) rich.Text {
	c1x, c2x := c1.Merge(c2)
	s1 := s.Apply(nil, c1).Apply(nil, c1x).(rich.Text)
	s2 := s.Apply(nil, c2).Apply(nil, c2x).(rich.Text)
	if x1, x2 := html.Format(s1, nil), html.Format(s2, nil); x1 != x2 {
		t.Error("Diverged", x1, x2)
	}
	return s1
}

func testReverseMerge(t *testing.T, s rich.Text, c1, c2 changes.Change) {
	x := html.Format(testMerge(t, s, c1, c2), nil)
	c2x, c1x := c2.(changes.Custom).ReverseMerge(c1)
	s1 := s.Apply(nil, c1).Apply(nil, c1x).(rich.Text)
	s2 := s.Apply(nil, c2).Apply(nil, c2x).(rich.Text)
	if x1, x2 := html.Format(s1, nil), html.Format(s2, nil); x1 != x2 || x1 != x {
		t.Error("Diverged reverse merge", x1, x2, x)
	}

}
