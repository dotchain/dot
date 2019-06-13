// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Dir wraps a value along with a map of objects.  The contents of the
// value can refer to entries in the map using Ref
type Dir struct {
	Root    changes.Value
	Objects types.M
}

// Name returns the key name for use with Attrs
func (d *Dir) Name() string {
	return "Embed"
}

// Apply implements channges.Value
func (d *Dir) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: d.set, Get: d.get}).Apply(ctx, c, d)
}

func (d *Dir) get(key interface{}) changes.Value {
	if key == "Root" {
		return d.Root
	}
	return d.Objects
}

func (d *Dir) set(key interface{}, v changes.Value) changes.Value {
	clone := *d
	if key == "Root" {
		clone.Root = v
	} else {
		clone.Objects = v.(types.M)
	}
	return &clone
}
