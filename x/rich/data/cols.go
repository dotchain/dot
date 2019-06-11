// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Cols represents all the columns in the table
type Cols map[interface{}]*Col

// Apply implements changes.Value
func (cols Cols) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: cols.set, Get: cols.get}).Apply(ctx, c, cols)
}

func (cols Cols) get(key interface{}) changes.Value {
	if v, ok := cols[key]; ok {
		return *v
	}
	return changes.Nil
}

func (cols Cols) set(key interface{}, v changes.Value) changes.Value {
	clone := Cols{}
	for k, val := range cols {
		if k != key {
			clone[k] = val
		}
	}
	if v != changes.Nil {
		x := v.(Col)
		clone[key] = &x
	}
	return clone
}
