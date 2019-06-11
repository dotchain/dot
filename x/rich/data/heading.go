// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// Heading represents h1 to h6.
//
// Note that the contents of the heading tag can be any rich text.
type Heading struct {
	Level int // 1 => 6
	*rich.Text
}

// Name is the key to use with rich.Attrs
func (h Heading) Name() string {
	return "Embed"
}

// Apply implements changes.Value.
func (h Heading) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: h.set, Get: h.get}).Apply(ctx, c, h)
}

func (h Heading) get(key interface{}) changes.Value {
	if key == "Level" {
		return changes.Atomic{Value: h.Level}
	}
	return h.Text
}

func (h Heading) set(key interface{}, v changes.Value) changes.Value {
	if key == "Level" {
		h.Level = v.(changes.Atomic).Value.(int)
	} else {
		h.Text = v.(*rich.Text)
	}
	return h
}
