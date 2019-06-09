// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rich

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Attr is a single attribute value
type Attr interface {
	changes.Value
	Name() string
}

// Attrs is a collection of attribute name value pairs
type Attrs map[string]Attr

// Apply implements changes.Value
func (a Attrs) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: a.set, Get: a.get}).Apply(ctx, c, a)
}

// Equal does a deep comparison of all attributes
func (a Attrs) Equal(x interface{}) bool {
	o, ok := x.(Attrs)
	if !ok || len(a) != len(o) {
		return false
	}

	for k, v := range a {
		if v2, ok := o[k]; !ok || !equalAttr(v, v2) {
			return false
		}
	}
	return true
}

func (a Attrs) get(key interface{}) changes.Value {
	v, ok := a[key.(string)]
	if !ok {
		return changes.Nil
	}
	return v
}

func (a Attrs) set(key interface{}, v changes.Value) changes.Value {
	clone := Attrs{}
	for k, vx := range a {
		if k != key {
			clone[k] = vx
		}
	}

	if v != changes.Nil {
		clone[key.(string)] = v.(Attr)
	}
	return clone
}

func equalAttr(a, b changes.Value) bool {
	return a == b
}
