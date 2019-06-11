// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"html"
	"strconv"
	"strings"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

// Formatter is a generic html formatter
type Formatter func(b *strings.Builder, v changes.Value)

// Format converts any value to HTML
func Format(v changes.Value) string {
	var b strings.Builder
	FormatBuilder(&b, v, nil)
	return b.String()
}

// FormatBuilder converts any value to HTML
//
// If a formatter is provided, it is used for embedded objects
func FormatBuilder(b *strings.Builder, v changes.Value, f Formatter) {
	if f == nil {
		f = func(b *strings.Builder, v changes.Value) {
			FormatBuilder(b, v, nil)
		}
	}
	switch v := v.(type) {
	case types.S16:
		b.WriteString(html.EscapeString(string(v)))
	case *rich.Text:
		formatRichText(b, v, f)
	case data.Link:
		formatLink(b, v, f)
	case data.BlockQuote:
		formatBlockQuote(b, v, f)
	case data.Heading:
		formatHeading(b, v, f)
	case data.Image:
		formatImage(b, v)
	case data.List:
		formatList(b, v, f)
	}
}

func formatLink(b *strings.Builder, l data.Link, f Formatter) {
	b.WriteString("<a href=\"")
	b.WriteString(html.EscapeString(l.Url))
	b.WriteString("\">")
	f(b, l.Value)
	b.WriteString("</a>")
}

func formatBlockQuote(b *strings.Builder, bq data.BlockQuote, f Formatter) {
	b.WriteString("<blockquote>")
	f(b, bq.Text)
	b.WriteString("</blockquote>")
}

func formatHeading(b *strings.Builder, h data.Heading, f Formatter) {
	l := h.Level
	if l < 1 || l > 6 {
		l = 1
	}
	b.WriteString("<h")
	b.WriteString(strconv.Itoa(l))
	b.WriteString(">")
	f(b, h.Text)
	b.WriteString("</h")
	b.WriteString(strconv.Itoa(l))
	b.WriteString(">")
}

func formatImage(b *strings.Builder, i data.Image) {
	b.WriteString("<img src=\"")
	b.WriteString(html.EscapeString(i.Src))
	b.WriteString("\" alt=\"")
	b.WriteString(html.EscapeString(i.AltText))
	b.WriteString("\">")
	b.WriteString("</img>")
}

func formatList(b *strings.Builder, l data.List, f Formatter) {
	tag := "ol"
	if l.Type == "disc" || l.Type == "circle" || l.Type == "square" || l.Type == "" {
		tag = "ul"
	}

	style := ""
	if l.Type != "" {
		style = " style=\"list-style-type: " + html.EscapeString(l.Type) + ";\""
	}
	b.WriteString("<" + tag + style + ">")
	writeListEntries(b, l.Text, f)
	b.WriteString("</" + tag + ">")
}

func writeListEntries(b *strings.Builder, t *rich.Text, f Formatter) {
mainloop:
	for len(*t) > 0 {
		seen := 0
		for _, x := range *t {
			if idx := strings.Index(x.Text, "\n"); idx >= 0 {
				b.WriteString("<li>")
				f(b, t.Slice(0, seen+idx).(*rich.Text))
				b.WriteString("</li>")
				t = t.Slice(seen+idx+1, t.Count()-seen-idx-1).(*rich.Text)
				continue mainloop
			}
			seen += x.Size
		}
		b.WriteString("<li>")
		f(b, t)
		b.WriteString("</li>")
		return
	}
}

var inlineStyles = []string{"FontStyle", "FontWeight"}

func formatRichText(b *strings.Builder, t *rich.Text, f Formatter) {
	last := rich.Attrs{}
	for _, x := range *t {
		for kk := range inlineStyles {
			name := inlineStyles[len(inlineStyles)-1-kk]
			if attr, ok := last[name]; ok {
				inlineClose(b, attr)
			}
		}
		for _, name := range inlineStyles {
			if attr, ok := x.Attrs[name]; ok {
				inlineOpen(b, attr)
			}
		}

		if attr, ok := x.Attrs["Embed"]; ok {
			f(b, attr)
		} else {
			f(b, types.S16(x.Text))
		}

		last = x.Attrs
	}
	if !last.Equal(rich.Attrs{}) {
		for _, name := range inlineStyles {
			if attr, ok := last[name]; ok {
				inlineClose(b, attr)
			}
		}
	}
}

func inlineOpen(b *strings.Builder, v changes.Value) {
	switch t := v.(type) {
	case data.FontStyle:
		if t == data.FontStyleItalic {
			b.WriteString("<i>")
		} else {
			b.WriteString("<span style=\"font-style: ")
			b.WriteString(html.EscapeString(string(t)))
			b.WriteString("\">")
		}
	case data.FontWeight:
		if t == data.FontBold {
			b.WriteString("<b>")
		} else {
			b.WriteString("<span style=\"font-weight: ")
			b.WriteString(strconv.Itoa(int(t)))
			b.WriteString("\">")
		}
	}
}

func inlineClose(b *strings.Builder, v changes.Value) {
	switch t := v.(type) {
	case data.FontStyle:
		if t == data.FontStyleItalic {
			b.WriteString("</i>")
		} else {
			b.WriteString("</span>")
		}
	case data.FontWeight:
		if t == data.FontBold {
			b.WriteString("</b>")
		} else {
			b.WriteString("</span>")
		}
	}
}
