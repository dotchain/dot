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

// NewLink creates a rich text that represents a link element
func NewLink(url string, contents rich.Text) rich.Text {
	return rich.NewText(" ", Link{url, &contents})
}

// Link represents a url link
//
// Note that the contents of the link can be any rich text.
type Link struct {
	Url string
	*rich.Text
}

// Name is the key to use with rich.Attrs
func (l Link) Name() string {
	return "Link"
}

// Apply implements changes.Value.
func (l Link) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: l.set, Get: l.get}).Apply(ctx, c, l)
}

func (l Link) get(key interface{}) changes.Value {
	switch key {
	case "Url":
		return types.S16(l.Url)
	case "Text":
		return *l.Text
	}
	return changes.Nil
}

func (l Link) set(key interface{}, v changes.Value) changes.Value {
	switch key {
	case "Url":
		l.Url = string(v.(types.S16))
	case "Text":
		x := v.(rich.Text)
		l.Text = &x
	}
	return l
}

// FormatHTML formats the link into HTML
func (l Link) FormatHTML(b *strings.Builder, f Formatter) {
	b.WriteString("<a href=\"")
	b.WriteString(html.EscapeString(l.Url))
	b.WriteString("\">")
	FormatBuilder(b, *l.Text, f)
	b.WriteString("</a>")
}
