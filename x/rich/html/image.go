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

// NewImage creates a rich text with embedded image
func NewImage(src, altText string) *rich.Text {
	return rich.NewText(" ", Image{src, altText})
}

// Image represents an image url
type Image struct {
	Src     string
	AltText string
}

// Name is the key to use with rich.Attrs
func (i Image) Name() string {
	return "Image"
}

// Apply implements changes.Value.
func (i Image) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: i.set, Get: i.get}).Apply(ctx, c, i)
}

func (i Image) get(key interface{}) changes.Value {
	switch key {
	case "Src":
		return types.S16(i.Src)
	case "AltText":
		return types.S16(i.AltText)
	}
	return changes.Nil
}

func (i Image) set(key interface{}, v changes.Value) changes.Value {
	switch key {
	case "Src":
		i.Src = string(v.(types.S16))
	case "AltText":
		i.AltText = string(v.(types.S16))
	}
	return i
}

// FormatHTML formats the image into HTML
func (i Image) FormatHTML(b *strings.Builder, f Formatter) {
	b.WriteString("<img src=\"")
	b.WriteString(html.EscapeString(i.Src))
	b.WriteString("\" alt=\"")
	b.WriteString(html.EscapeString(i.AltText))
	b.WriteString("\">")
	b.WriteString("</img>")
}
