// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Vals manages a list of values
type Vals []Val

func (v *Vals) get(key interface{}) changes.Value {
	return (*v)[key.(int)]
}

func (v *Vals) set(key interface{}, value changes.Value) changes.Value {
	clone := append(Vals(nil), *v...)
	clone[key.(int)] = value.(Val)
	return &clone
}

func (v *Vals) splice(offset, count int, insert changes.Collection) changes.Collection {
	if v == nil {
		v = &Vals{}
	}
	vx := insert.(*Vals)
	if vx == nil {
		vx = &Vals{}
	}
	clone := append((*v)[0:offset:offset], *vx...)
	clone = append(clone, (*v)[offset+count:]...)
	return &clone
}

// Apply implements changes.Value
func (v *Vals) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: v.set, Get: v.get, Splice: v.splice}).Apply(ctx, c, v)
}

// ApplyCollection implements changes.Collection
func (v *Vals) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Set: v.set, Get: v.get, Splice: v.splice}).
		ApplyCollection(ctx, c, v)
}

// Count return number of vals
func (v *Vals) Count() int {
	if v == nil {
		return 0
	}
	return len(*v)
}

// Slice implements changes.Collection
func (v *Vals) Slice(offset, count int) changes.Collection {
	if count == 0 {
		return &Vals{}
	}
	result := (*v)[offset : offset+count : offset+count]
	return &result
}

// Text implements Val.Text
func (v *Vals) Text() string {
	return "<list>"
}

// Visit implements Val.Visite
func (v *Vals) Visit(visitor Visitor) {
	visitor.VisitChildrenBegin(v)
	if v != nil {
		for kk, val := range *v {
			visitor.VisitChild(val, kk)
		}
	}
	visitor.VisitChildrenEnd(v)
}
