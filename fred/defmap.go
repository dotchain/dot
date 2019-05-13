// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// DefMap maanages a list of defs.  It returns a list as the value
type DefMap map[interface{}]Def

func (d *DefMap) get(key interface{}) changes.Value {
	if d != nil {
		if x, ok := (*d)[key]; ok {
			return x
		}
	}
	return changes.Nil
}

func (d *DefMap) set(key interface{}, value changes.Value) changes.Value {
	clone := DefMap{}
	if d != nil {
		for k, v := range *d {
			if k != key {
				clone[k] = v
			}
		}
	}
	if value != changes.Nil {
		clone[key] = value.(Def)
	}
	return &clone
}

// Apply implements changes.Value
func (d *DefMap) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: d.set, Get: d.get}).Apply(ctx, c, d)
}

// Eval evaluates each def
func (d *DefMap) Eval(e Env) Val {
	if d == nil {
		return &ValMap{}
	}

	// TODO: only cache if there is something varyable
	return e.ValueOf(d, func() Val {
		result := ValMap{}
		for kk, dd := range *d {
			result[kk] = dd.Eval(e)
		}
		return &result
	})
}
