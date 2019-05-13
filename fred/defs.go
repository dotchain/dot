// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Defs maanages a list of defs.  It returns a list as the value
type Defs []Def

func (d *Defs) get(key interface{}) changes.Value {
	return (*d)[key.(int)]
}

func (d *Defs) set(key interface{}, value changes.Value) changes.Value {
	clone := append(Defs(nil), *d...)
	clone[key.(int)] = value.(Def)
	return &clone
}

func (d *Defs) splice(offset, count int, insert changes.Collection) changes.Collection {
	if d == nil {
		d = &Defs{}
	}
	dx := insert.(*Defs)
	if dx == nil {
		dx = &Defs{}
	}
	clone := append((*d)[0:offset:offset], *dx...)
	clone = append(clone, (*d)[offset+count:]...)
	return &clone
}

// Apply implements changes.Value
func (d *Defs) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: d.set, Get: d.get, Splice: d.splice}).Apply(ctx, c, d)
}

// ApplyCollection implements changes.Collection
func (d *Defs) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Set: d.set, Get: d.get, Splice: d.splice}).
		ApplyCollection(ctx, c, d)
}

// Count returns the number of defs
func (d *Defs) Count() int {
	if d == nil {
		return 0
	}
	return len(*d)
}

// Slice implements changes.Collection
func (d *Defs) Slice(offset, count int) changes.Collection {
	if count == 0 {
		return &Defs{}
	}
	result := (*d)[offset : offset+count : offset+count]
	return &result
}

var empty = &Vals{}

// Eval evaluates each def
func (d *Defs) Eval(e Env) Val {
	if d == nil {
		return empty
	}
	return e.ValueOf(d, func() Val {
		result := Vals(make([]Val, len(*d)))
		for kk, dd := range *d {
			result[kk] = dd.Eval(e)
		}
		return &result
	})
}
