// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// List represents an ordered or unordered list
//
// The type can be one of the string values defined here:
// https://developer.mozilla.org/en-US/docs/Web/CSS/list-style-type
//
// (such as disc, circle etc)
type List struct {
	Type    string
	Entries types.A
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
	return l.Entries
}

func (l List) set(key interface{}, v changes.Value) changes.Value {
	if key == "Type" {
		l.Type = string(v.(types.S16))
	} else {
		l.Entries = v.(types.A)
	}
	return l
}
