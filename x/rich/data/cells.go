// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// Cells represents a row of data.
//
// The key in the map is the column ID. If a column ID exists in the
// map, the value cannot be nil.
type Cells map[interface{}]*rich.Text

// Apply implements changes.Value
func (cells Cells) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: cells.set, Get: cells.get}).Apply(ctx, c, cells)
}

func (cells Cells) get(key interface{}) changes.Value {
	if v, ok := cells[key]; ok {
		return *v
	}
	return changes.Nil
}

func (cells Cells) set(key interface{}, v changes.Value) changes.Value {
	clone := Cells{}
	for k, val := range cells {
		if k != key {
			clone[k] = val
		}
	}
	if v != changes.Nil {
		x := v.(rich.Text)
		clone[key] = &x
	}
	return clone
}
