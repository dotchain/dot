// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"testing"

	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func TestFormatBold(t *testing.T) {
	s := rich.NewText("Hello ").
		Concat(rich.NewText("beautiful", html.FontBold)).
		Concat(rich.NewText(" world"))

	result := html.Format(s)
	if result != "Hello <b>beautiful</b> world" {
		t.Error("Unexpected", result, s)
	}
}

func TestFormatBoldAndItalic(t *testing.T) {
	s := rich.NewText("Hello ").
		Concat(rich.NewText("bold", html.FontBold)).
		Concat(rich.NewText("and", html.FontBold, html.FontStyleItalic)).
		Concat(rich.NewText("italic", html.FontStyleItalic)).
		Concat(rich.NewText(" world"))

	result := html.Format(s)
	if result != "Hello <b>bold</b><i><b>and</b></i><i>italic</i> world" {
		t.Error("Unexpected", result, s)
	}
}

func TestFormatBlockQuote(t *testing.T) {
	bq := html.NewBlockQuote(rich.NewText("hello", html.FontBold))

	if x := html.Format(bq); x != "<blockquote><b>hello</b></blockquote>" {
		t.Error("Unexpected", x)
	}
}

func TestFormatFontStyle(t *testing.T) {
	styles := []html.FontStyle{
		html.FontStyleNormal,
		html.FontStyleItalic,
		html.FontStyleOblique,
	}

	for _, style := range styles {
		t.Run(string(style), func(t *testing.T) {
			s := rich.NewText("hello", style)
			expected := fmt.Sprintf(
				"<span style=\"font-style: %s\">hello</span>",
				style,
			)
			if style == html.FontStyleItalic {
				expected = "<i>hello</i>"
			}
			if x := html.Format(s); x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}

func TestFormatFontWeight(t *testing.T) {
	weights := []html.FontWeight{
		html.FontThin,
		html.FontExtraLight,
		html.FontLight,
		html.FontNormal,
		html.FontMedium,
		html.FontSemibold,
		html.FontBold,
		html.FontExtraBold,
		html.FontBlack,
	}
	for _, weight := range weights {
		str := fmt.Sprintf("%d", weight)
		t.Run(str, func(t *testing.T) {
			s := rich.NewText("hello", weight)
			expected := fmt.Sprintf(
				"<span style=\"font-weight: %d\">hello</span>",
				weight,
			)
			if weight == html.FontBold {
				expected = "<b>hello</b>"
			}
			if x := html.Format(s); x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}

func TestFormatHeading(t *testing.T) {
	levels := []string{"h1", "h1", "h2", "h3", "h4", "h5", "h6", "h1"}
	for l, str := range levels {
		test := fmt.Sprintf("%s-%d", str, l)
		t.Run(test, func(t *testing.T) {
			h := html.NewHeading(l, rich.NewText("x", html.FontBold))
			expected := fmt.Sprintf("<%s><b>x</b></%s>", str, str)
			if x := html.Format(h); x != expected {
				t.Error("Unexpected", x, expected)
			}
		})
	}
}

func TestFormatImage(t *testing.T) {
	i := html.NewImage("quote\"d", "a < b")

	if x := html.Format(i); x != "<img src=\"quote&#34;d\" alt=\"a &lt; b\"></img>" {
		t.Error("Unexpected", x)
	}
}

func TestFormatLink(t *testing.T) {
	s := rich.NewText("a < b")
	l := html.NewLink("quote\"d", s)

	if x := html.Format(l); x != "<a href=\"quote&#34;d\">a &lt; b</a>" {
		t.Error("Unexpected", x)
	}
}

func TestFormatList(t *testing.T) {
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

		if x := html.Format(l); x != expected {
			t.Error("Unexpected", x, expected)
		}
	}

	s := rich.NewText("hel").
		Concat(rich.NewText("lo\nwor", html.FontBold)).
		Concat(rich.NewText("ld"))
	l := html.NewList("", s)
	expected := "<ul><li>hel<b>lo</b></li><li><b>wor</b>ld</li></ul>"
	if x := html.Format(l); x != expected {
		t.Error("Unexpected", x, expected)
	}
}
