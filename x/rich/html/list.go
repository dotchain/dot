// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
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
	return "Embed"
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
