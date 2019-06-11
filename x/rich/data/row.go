// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Row represents the information about a row
//
// The ID field is immutable.  The Ord field is meant to help sort the
// rows and is not meant to be directly modified by the
// application. Instead, Table methods can be used to effectively
// manipulate these.
//
// The Cells field contains a map of Column.ID to actual cell value
type Row struct {
	ID  interface{}
	Ord string
	Cells
}

// Apply implements changes.Value
func (r Row) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: r.set, Get: r.get}).Apply(ctx, c, r)
}

func (r Row) get(key interface{}) changes.Value {
	if key == "Ord" {
		return types.S16(r.Ord)
	}
	return r.Cells
}

func (r Row) set(key interface{}, v changes.Value) changes.Value {
	if key == "Ord" {
		r.Ord = string(v.(types.S16))
	} else {
		r.Cells = v.(Cells)
	}
	return r
}
