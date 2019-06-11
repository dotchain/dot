// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Rows represents all the rows in the table
type Rows map[interface{}]*Row

// Apply implements changes.Value
func (rows Rows) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: rows.set, Get: rows.get}).Apply(ctx, c, rows)
}

func (rows Rows) get(key interface{}) changes.Value {
	if v, ok := rows[key]; ok {
		return *v
	}
	return changes.Nil
}

func (rows Rows) set(key interface{}, v changes.Value) changes.Value {
	clone := Rows{}
	for k, val := range rows {
		if k != key {
			clone[k] = val
		}
	}
	if v != changes.Nil {
		x := v.(Row)
		clone[key] = &x
	}
	return clone
}
