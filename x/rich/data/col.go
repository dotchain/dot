// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
)

// Col represents the information about the column.
//
// The ID field is immutable.  The Ord field is meant to help sort the
// columns and is not meant to be directly modified by the
// application. Instead, Table methods can be used to effectively
// manipulate these.
//
// The value field tracks the actual heading/name for the column and
// should not be nil.
type Col struct {
	ID    interface{}
	Ord   string
	Value *rich.Text
}

// Apply implements changes.Value
func (col Col) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: col.set, Get: col.get}).Apply(ctx, c, col)
}

func (col Col) get(key interface{}) changes.Value {
	if key == "Ord" {
		return types.S16(col.Ord)
	}
	return *col.Value
}

func (col Col) set(key interface{}, v changes.Value) changes.Value {
	if key == "Ord" {
		col.Ord = string(v.(types.S16))
	} else {
		x := v.(rich.Text)
		col.Value = &x
	}
	return col
}
