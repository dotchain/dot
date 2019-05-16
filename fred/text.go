// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Text uses types.S16 to implement string-based values
type Text string

// Apply implements changes.Value
func (t Text) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if splice, ok := c.(changes.Splice); ok {
		splice.Before = types.S16(string(splice.Before.(Text)))
		splice.After = types.S16(string(splice.After.(Text)))
		c = splice
	}
	if custom, ok := c.(changes.Custom); ok {
		return custom.ApplyTo(ctx, t)
	}

	v := types.S16(string(t)).Apply(ctx, c)
	if x, ok := v.(types.S16); ok {
		return Text(x)
	}
	return v
}

// Count implements changes.Collection
func (t Text) Count() int {
	return types.S16(string(t)).Count()
}

// Slice implements changes.Collection
func (t Text) Slice(offset, count int) changes.Collection {
	v := types.S16(string(t)).Slice(offset, count)
	return Text(string(v.(types.S16)))
}

// Apply implements changes.Collection
func (t Text) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	if splice, ok := c.(changes.Splice); ok {
		splice.Before = types.S16(string(splice.Before.(Text)))
		splice.After = types.S16(string(splice.After.(Text)))
		c = splice
	}

	v := types.S16(string(t)).ApplyCollection(ctx, c)
	return Text(string(v.(types.S16)))
}

// Text implements Val.Text
func (t Text) Text() string {
	return string(t)
}

// Visit implements Val.Visit
func (t Text) Visit(v Visitor) {
	v.VisitLeaf(t)
}
