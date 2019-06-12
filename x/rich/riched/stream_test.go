// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
	"github.com/dotchain/dot/x/rich/html"
	"github.com/dotchain/dot/x/rich/riched"
)

func TestStreamEmptyNext(t *testing.T) {
	s := riched.NewStream(rich.NewText("Hello world"))
	if x := s.Next(); x != nil {
		t.Error("Unexpected next", x)
	}

	if x := s.ClearOverrides(); x != s {
		t.Error("Unexpected next", x)
	}
}

func TestStreamTextEditUpdatesSelection(t *testing.T) {
	s := riched.NewStream(rich.NewText("Hello world", data.FontBold))
	s = s.SetSelection([]interface{}{5}, []interface{}{5})
	s.Stream.Append(changes.PathChange{
		Path: []interface{}{"Text"},
		Change: changes.Splice{
			Offset: 0,
			Before: &rich.Text{},
			After:  rich.NewText("Hi! "),
		},
	})

	s = s.Next()

	if !reflect.DeepEqual(s.Focus, []interface{}{5 + len("Hi! ")}) {
		t.Error("Unexpected focus value", s.Focus)
	}

	if !reflect.DeepEqual(s.Anchor, []interface{}{5 + len("Hi! ")}) {
		t.Error("Unexpected anchor value", s.Anchor)
	}
}

func TestStreamInsertString(t *testing.T) {
	v := rich.Text(nil)
	s := riched.NewStream(&v)
	s = s.InsertString("hello world")
	if x := html.Format(s.Text); x != "hello world" {
		t.Error("Unexpected", x)
	}
}

func TestStreamOverrides(t *testing.T) {
	s := riched.NewStream(rich.NewText("Hello world", data.FontBold))
	s = s.SetSelection([]interface{}{5}, []interface{}{5})
	s = s.SetOverride(riched.NoAttribute(data.FontBold.Name()))
	s = s.SetOverride(data.FontStyleItalic)

	t.Run("SetOverride", func(t *testing.T) {
		s2 := s.InsertString(" beautiful")

		if x := html.Format(s2.Text); x != "<b>Hello</b><i> beautiful</i><b> world</b>" {
			t.Error("Unexpected SetOverride", x)
		}

		if x := s.Editor.SetOverride(data.FontStyleItalic); x != nil {
			t.Error("Unexpected SetOverride2", x)
		}
	})

	t.Run("SetOverride update", func(t *testing.T) {
		s2 := s.SetOverride(data.FontStyleOblique).InsertString(" beautiful")

		expected := "<b>Hello</b><span style=\"font-style: oblique\"> beautiful</span><b> world</b>"
		if x := html.Format(s2.Text); x != expected {
			t.Error("Unexpected SetOverride", x)
		}
	})

	t.Run("RemoveOverride", func(t *testing.T) {
		s2 := s.RemoveOverride(data.FontStyleItalic.Name()).InsertString(" beautiful")

		if x := html.Format(s2.Text); x != "<b>Hello</b> beautiful<b> world</b>" {
			t.Error("Unexpected RemoveOverride", x)
		}
		if x := s2.Editor.RemoveOverride(data.FontStyleItalic.Name()); x != nil {
			t.Error("Unexpected RemoveOverride2", x)
		}
	})

	t.Run("ClearOverrides", func(t *testing.T) {
		s2 := s.ClearOverrides().InsertString(" beautiful")

		if x := html.Format(s2.Text); x != "<b>Hello beautiful world</b>" {
			t.Error("Unexpected ClearOverrides", x)
		}
		if x := s2.Editor.ClearOverrides(); x != nil {
			t.Error("Unexpected ClearOverrides2", x)
		}
	})
}
