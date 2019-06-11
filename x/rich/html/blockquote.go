// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"strings"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// NewBlockQuote creates a rich text with embedded block quote
func NewBlockQuote(r *rich.Text) *rich.Text {
	return rich.NewText(" ", BlockQuote{r})
}

// BlockQuote represents a block quote with any embedded content
type BlockQuote struct {
	Text *rich.Text
}

// Name is the key to use with rich.Attrs
func (bq BlockQuote) Name() string {
	return "BlockQuote"
}

// Apply implements changes.Value.
func (bq BlockQuote) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: bq.set, Get: bq.get}).Apply(ctx, c, bq)
}

func (bq BlockQuote) get(key interface{}) changes.Value {
	return bq.Text
}

func (bq BlockQuote) set(key interface{}, v changes.Value) changes.Value {
	bq.Text = v.(*rich.Text)
	return bq
}

// FormatHTML formats the blockQuote into HTML
func (bq BlockQuote) FormatHTML(b *strings.Builder, f Formatter) {
	b.WriteString("<blockquote>")
	FormatBuilder(b, bq.Text, f)
	b.WriteString("</blockquote>")
}
