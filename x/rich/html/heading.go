// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"fmt"
	"strings"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// NewHeading creates a rich text that represents a heading element
func NewHeading(level int, r rich.Text) rich.Text {
	return rich.NewText(" ", Heading{level, &r})
}

// Heading represents h1 to h6.
//
// Note that the contents of the heading tag can be any rich text.
type Heading struct {
	Level int // 1 => 6
	*rich.Text
}

// Name is the key to use with rich.Attrs
func (h Heading) Name() string {
	return "Heading"
}

// Apply implements changes.Value.
func (h Heading) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: h.set, Get: h.get}).Apply(ctx, c, h)
}

func (h Heading) get(key interface{}) changes.Value {
	switch key {
	case "Level":
		return changes.Atomic{Value: h.Level}
	case "Text":
		return *h.Text
	}
	return changes.Nil
}

func (h Heading) set(key interface{}, v changes.Value) changes.Value {
	switch key {
	case "Level":
		h.Level = v.(changes.Atomic).Value.(int)
	case "Text":
		x := v.(rich.Text)
		h.Text = &x
	}
	return h
}

// FormatHTML formats the heading into HTML
func (h Heading) FormatHTML(b *strings.Builder, f Formatter) {
	l := h.Level
	if l < 1 || l > 6 {
		l = 1
	}
	b.WriteString(fmt.Sprintf("<h%d>", l))
	FormatBuilder(b, *h.Text, f)
	b.WriteString(fmt.Sprintf("</h%d>", l))
}
