// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func TestListApply(t *testing.T) {
	s1, s2 := rich.NewText("list1"), rich.NewText("list2")
	l1 := html.List{Type: "circle", Text: &s1}
	l2 := html.List{Type: "square", Text: &s2}

	if x := l1.Apply(nil, nil); !reflect.DeepEqual(x, l1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: l2, After: l2}
	if x := l1.Apply(nil, replace); !reflect.DeepEqual(x, l2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path: []interface{}{"Type"},
		Change: changes.Replace{
			Before: types.S16("circle"),
			After:  types.S16("square"),
		},
	}
	if x := l1.Apply(nil, c).(html.List); x.Type != "square" {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"Text"},
		Change: changes.Replace{
			Before: s1,
			After:  s2,
		},
	}
	if x := l1.Apply(nil, c).(html.List); !reflect.DeepEqual(*x.Text, s2) {
		t.Error("Unexpected change", x.Text)
	}
}

func TestListFormatHTML(t *testing.T) {
	tests := map[[2]string]string{
		{"", "hello"}:        "<ul><li>hello</li></ul>",
		{"circle", "hello"}:  "<ul style=\"list-style-type: circle;\"><li>hello</li></ul>",
		{"disc", "hello"}:    "<ul style=\"list-style-type: disc;\"><li>hello</li></ul>",
		{"square", "hello"}:  "<ul style=\"list-style-type: square;\"><li>hello</li></ul>",
		{"decimal", "hello"}: "<ol style=\"list-style-type: decimal;\"><li>hello</li></ol>",
		{"", "hello\nworld"}: "<ul><li>hello</li><li>world</li></ul>",
	}

	for pair, expected := range tests {
		l := html.NewList(pair[0], rich.NewText(pair[1]))

		if x := html.Format(l, nil); x != expected {
			t.Error("Unexpected", x, expected)
		}
	}

	s := rich.NewText("hel").
		Concat(rich.NewText("lo\nwor", html.FontBold)).
		Concat(rich.NewText("ld"))
	l := html.NewList("", s)
	expected := "<ul><li>hel<b>lo</b></li><li><b>wor</b>ld</li></ul>"
	if x := html.Format(l, nil); x != expected {
		t.Error("Unexpected", x, expected)
	}
}
