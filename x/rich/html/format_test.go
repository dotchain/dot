// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
	"github.com/dotchain/dot/x/rich/html"
)

func TestFormatMultiple(t *testing.T) {
	s := rich.NewText("Hello ").
		Concat(rich.NewText("bold", data.FontBold)).
		Concat(rich.NewText("and", data.FontBold, data.FontStyleItalic)).
		Concat(rich.NewText("italic", data.FontStyleItalic)).
		Concat(rich.NewText(" world"))

	result := html.Format(s)
	if result != "Hello <b>bold</b><i><b>and</b></i><i>italic</i> world" {
		t.Error("Unexpected", result, s)
	}
}

func TestFormatBlockQuote(t *testing.T) {
	bq := rich.NewEmbed(data.BlockQuote{Text: rich.NewText("hello", data.FontBold)})

	if x := html.Format(bq); x != "<blockquote><b>hello</b></blockquote>" {
		t.Error("Unexpected", x)
	}
}

func TestFormatFontStyle(t *testing.T) {
	styles := []data.FontStyle{
		data.FontStyleNormal,
		data.FontStyleItalic,
		data.FontStyleOblique,
	}

	for _, style := range styles {
		t.Run(string(style), func(t *testing.T) {
			s := rich.NewText("hello", style)
			expected := fmt.Sprintf(
				"<span style=\"font-style: %s\">hello</span>",
				style,
			)
			if style == data.FontStyleItalic {
				expected = "<i>hello</i>"
			}
			if x := html.Format(s); x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}

func TestFormatFontWeight(t *testing.T) {
	weights := []data.FontWeight{
		data.FontThin,
		data.FontExtraLight,
		data.FontLight,
		data.FontNormal,
		data.FontMedium,
		data.FontSemibold,
		data.FontBold,
		data.FontExtraBold,
		data.FontBlack,
	}
	for _, weight := range weights {
		str := fmt.Sprintf("%d", weight)
		t.Run(str, func(t *testing.T) {
			s := rich.NewText("hello", weight)
			expected := fmt.Sprintf(
				"<span style=\"font-weight: %d\">hello</span>",
				weight,
			)
			if weight == data.FontBold {
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
			h := rich.NewEmbed(data.Heading{
				Level: l,
				Text:  rich.NewText("x", data.FontBold),
			})
			expected := fmt.Sprintf("<%s><b>x</b></%s>", str, str)
			if x := html.Format(h); x != expected {
				t.Error("Unexpected", x, expected)
			}
		})
	}
}

func TestFormatImage(t *testing.T) {
	i := rich.NewEmbed(data.Image{Src: "quote\"d", AltText: "a < b"})

	if x := html.Format(i); x != "<img src=\"quote&#34;d\" alt=\"a &lt; b\"></img>" {
		t.Error("Unexpected", x)
	}
}

func TestFormatLink(t *testing.T) {
	s := rich.NewText("a < b")
	l := rich.NewEmbed(data.Link{URL: "quote\"d", Value: s})

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
		splits := strings.Split(pair[1], "\n")
		entries := types.A{}
		for _, v := range splits {
			entries = append(entries, types.S16(v))
		}
		l := rich.NewEmbed(data.List{Type: pair[0], Entries: entries})

		if x := html.Format(l); x != expected {
			t.Error("Unexpected", x, expected)
		}
	}

	x := html.Format(types.A{types.S16("hello"), types.S16("world")})
	if x != "<ul style=\"list-style-type: disc;\"><li>hello</li><li>world</li></ul>" {
		t.Error("Unexpected", x)
	}
}
