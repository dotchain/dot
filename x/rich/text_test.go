// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rich_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func TestTextPlainText(t *testing.T) {
	s1 := rich.NewText("hello", html.FontBold)
	s2 := rich.NewText(" world")
	if x := s1.Concat(s2).PlainText(); x != "hello world" {
		t.Error("Unexpected PlainText()", x)
	}
}

func TestTextConcat(t *testing.T) {
	s1 := rich.NewText("hel", html.FontBold)
	s2 := rich.NewText("lo", html.FontBold)
	s3 := rich.NewText(" ", html.FontBold, html.FontStyleItalic)
	s4 := rich.NewText("wor")
	s5 := rich.NewText("ld")
	s := s1.Concat(s2).Concat(s3).Concat(s4).Concat(s5)
	if x := html.Format(s, nil); x != "<b>hello</b><i><b> </b></i>world" {
		t.Error("Unexpected", x)
	}
}

func TestTextSlice(t *testing.T) {
	s1 := rich.NewText("hello", html.FontBold)
	if x := s1.Slice(0, 5); !reflect.DeepEqual(x, s1) {
		t.Fatal("Slice full copy failed", x)
	}

	if x := s1.Slice(2, 0).(rich.Text); len(x) != 0 {
		t.Fatal("Slice zero failed", x)
	}

	s2 := rich.NewText(" world")
	s := s1.Concat(s2).Slice(3, 5).(rich.Text)
	if x := html.Format(s, nil); x != "<b>lo</b> wo" {
		t.Error("Unexpected", x)
	}
}

func TestTextSetAttr(t *testing.T) {
	s := rich.NewText("hello world")
	c := s.SetAttribute(3, 5, html.FontBold)
	s = s.Apply(nil, c).(rich.Text)
	if x := html.Format(s, nil); x != "hel<b>lo wo</b>rld" {
		t.Error("Unexpected", x)
	}

	s = s.Apply(nil, s.RemoveAttribute(5, 3, "FontWeight")).(rich.Text)
	if x := html.Format(s, nil); x != "hel<b>lo</b> world" {
		t.Error("Unexpected", x)
	}

	s = s.Apply(nil, s.SetAttribute(2, 7, html.FontStyleItalic)).(rich.Text)
	if x := html.Format(s, nil); x != "he<i>l</i><i><b>lo</b></i><i> wor</i>ld" {
		t.Error("Unexpected", x)
	}
}

func TestTextApplyCollection(t *testing.T) {
	s := rich.NewText("hello world")
	c := changes.Move{Offset: 4, Count: 1, Distance: -4}
	s = s.ApplyCollection(nil, c).(rich.Text)
	if x := html.Format(s, nil); x != "ohell world" {
		t.Error("Unexpected", x)
	}
}

func TestTextModifyStyleAtIndex(t *testing.T) {
	s := rich.NewText("hello world")
	c := changes.PathChange{
		Path: []interface{}{5, "FontWeight"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  html.FontBold,
		},
	}
	s = s.Apply(nil, c).(rich.Text)
	if x := html.Format(s, nil); x != "hello<b> </b>world" {
		t.Error("Unexpected", x)
	}
}
