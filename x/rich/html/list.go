// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// NewList creates a rich text that represents a list element
func NewList(listType string, contents *rich.Text) *rich.Text {
	return rich.NewText(" ", List{listType, contents})
}

// List represents an ordered or unordered list
//
// The type can be one of the string values defined here:
// https://developer.mozilla.org/en-US/docs/Web/CSS/list-style-type
//
// (such as disc, circle etc)
type List struct {
	Type string
	*rich.Text
}

// Name is the key to use with rich.Attrs
func (l List) Name() string {
	return "List"
}

// Apply implements changes.Value.
func (l List) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: l.set, Get: l.get}).Apply(ctx, c, l)
}

func (l List) get(key interface{}) changes.Value {
	if key == "Type" {
		return types.S16(l.Type)
	}
	return l.Text
}

func (l List) set(key interface{}, v changes.Value) changes.Value {
	if key == "Type" {
		l.Type = string(v.(types.S16))
	} else {
		l.Text = v.(*rich.Text)
	}
	return l
}

// FormatHTML formats the list into HTML
func (l List) FormatHTML(b *strings.Builder, f Formatter) {
	tag := "ol"
	if l.Type == "disc" || l.Type == "circle" || l.Type == "square" || l.Type == "" {
		tag = "ul"
	}

	style := ""
	if l.Type != "" {
		style = " style=\"list-style-type: " + html.EscapeString(l.Type) + ";\""
	}
	b.WriteString("<" + tag + style + ">")
	l.writeListEntries(b, f, l.Text)
	b.WriteString("</" + tag + ">")
}

func (l List) writeListEntries(b *strings.Builder, f Formatter, t *rich.Text) {
mainloop:
	for len(*t) > 0 {
		seen := 0
		for _, x := range *t {
			if idx := strings.Index(x.Text, "\n"); idx >= 0 {
				b.WriteString("<li>")
				FormatBuilder(b, t.Slice(0, seen+idx).(*rich.Text), f)
				b.WriteString("</li>")
				t = t.Slice(seen+idx+1, t.Count()-seen-idx-1).(*rich.Text)
				continue mainloop
			}
			seen += x.Size
		}
		b.WriteString("<li>")
		FormatBuilder(b, t, f)
		b.WriteString("</li>")
		return
	}
}
