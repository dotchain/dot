// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Link represents a url link
type Link struct {
	Url string
	changes.Value
}

// Name is the key to use with rich.Attrs
func (l Link) Name() string {
	return "Embed"
}

// Apply implements changes.Value.
func (l Link) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: l.set, Get: l.get}).Apply(ctx, c, l)
}

func (l Link) get(key interface{}) changes.Value {
	if key == "Url" {
		return types.S16(l.Url)
	}
	return l.Value
}

func (l Link) set(key interface{}, v changes.Value) changes.Value {
	if key == "Url" {
		l.Url = string(v.(types.S16))
	} else {
		l.Value = v
	}
	return l
}
