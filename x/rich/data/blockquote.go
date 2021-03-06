// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// BlockQuote represents a block quote with any embedded content
type BlockQuote struct {
	Text *rich.Text
}

// Name is the key to use with rich.Attrs
func (bq BlockQuote) Name() string {
	return "Embed"
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
