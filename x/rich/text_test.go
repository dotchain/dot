// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rich_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
	"github.com/dotchain/dot/x/rich/html"
)

func TestTextPlainText(t *testing.T) {
	s1 := rich.NewText("hello", data.FontBold)
	s2 := rich.NewText(" world")
	if x := s1.Concat(s2).PlainText(); x != "hello world" {
		t.Error("Unexpected PlainText()", x)
	}
}

func TestTextConcat(t *testing.T) {
	s1 := rich.NewText("hel", data.FontBold)
	s2 := rich.NewText("lo", data.FontBold)
	s3 := rich.NewText(" ", data.FontBold, data.FontStyleItalic)
	s4 := rich.NewText("wor")
	s5 := rich.NewText("ld")
	s := s1.Concat(s2).Concat(s3).Concat(s4).Concat(s5)
	if x := html.Format(s); x != "<b>hello</b><i><b> </b></i>world" {
		t.Error("Unexpected", x)
	}
}

func TestTextSlice(t *testing.T) {
	s1 := rich.NewText("hello", data.FontBold)
	if x := s1.Slice(0, 5); !reflect.DeepEqual(x, s1) {
		t.Fatal("Slice full copy failed", x)
	}

	if x := s1.Slice(2, 0).(*rich.Text); len(*x) != 0 {
		t.Fatal("Slice zero failed", x)
	}

	s2 := rich.NewText(" world")
	s := s1.Concat(s2).Slice(3, 5).(*rich.Text)
	if x := html.Format(s); x != "<b>lo</b> wo" {
		t.Error("Unexpected", x)
	}
}

func TestTextSetAttr(t *testing.T) {
	s := rich.NewText("hello world")
	c := s.SetAttribute(3, 5, data.FontBold)
	s = s.Apply(nil, c).(*rich.Text)
	if x := html.Format(s); x != "hel<b>lo wo</b>rld" {
		t.Error("Unexpected", x)
	}

	s = s.Apply(nil, s.RemoveAttribute(5, 3, "FontWeight")).(*rich.Text)
	if x := html.Format(s); x != "hel<b>lo</b> world" {
		t.Error("Unexpected", x)
	}

	s = s.Apply(nil, s.SetAttribute(2, 7, data.FontStyleItalic)).(*rich.Text)
	if x := html.Format(s); x != "he<i>l</i><i><b>lo</b></i><i> wor</i>ld" {
		t.Error("Unexpected", x)
	}
}

func TestTextApplyCollection(t *testing.T) {
	s := rich.NewText("hello world")
	c := changes.Move{Offset: 4, Count: 1, Distance: -4}
	s = s.ApplyCollection(nil, c).(*rich.Text)
	if x := html.Format(s); x != "ohell world" {
		t.Error("Unexpected", x)
	}
}

func TestTextModifyStyleAtIndex(t *testing.T) {
	s := rich.NewText("hello world")
	c := changes.PathChange{
		Path: []interface{}{5, "FontWeight"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  data.FontBold,
		},
	}
	s = s.Apply(nil, c).(*rich.Text)
	if x := html.Format(s); x != "hello<b> </b>world" {
		t.Error("Unexpected", x)
	}
}

func TestEmbed(t *testing.T) {
	embed := data.Link{URL: "hello", Value: types.S16("world")}
	s := rich.NewEmbed(embed)
	if s.PlainText() != " " {
		t.Error("Unexpected", s.PlainText())
	}
	if len(*s) != 1 || (*s)[0].Attrs["Embed"] != embed {
		t.Error("Unexpected rich text contents", s)
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Unexpected success")
		}
	}()
	rich.NewEmbed(data.FontBold)
}
