// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// ValMap manages a map of values
type ValMap map[interface{}]Val

func (v *ValMap) get(key interface{}) changes.Value {
	if v != nil {
		if x, ok := (*v)[key]; ok {
			return x
		}
	}
	return changes.Nil
}

func (v *ValMap) set(key interface{}, value changes.Value) changes.Value {
	clone := ValMap{}
	if v != nil {
		for k, val := range *v {
			if k != key {
				clone[k] = val
			}
		}
	}
	if value != changes.Nil {
		clone[key] = value.(Val)
	}
	return &clone
}

// Apply implements changes.Value
func (v *ValMap) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: v.set, Get: v.get}).Apply(ctx, c, v)
}

// Text implements Val.Text
func (v *ValMap) Text() string {
	return "<map>"
}

// Visit implements Val.Visit
func (v *ValMap) Visit(visitor Visitor) {
	visitor.VisitChildrenBegin(v)
	if v != nil {
		for k, val := range *v {
			visitor.VisitChild(val, k)
		}
	}
	visitor.VisitChildrenEnd(v)
}
