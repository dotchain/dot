// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package meta describes the DOT session meta data
package meta

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/ops"
)

// Data represents the session meta data
type Data struct {
	Version       int
	Pending       []ops.Op
	TransformedOp CachedOp
	MergeOps      CachedOps
}

// Apply implements changes.Value
func (d Data) Apply(ctx changes.Context, c changes.Change) changes.Value {
	g := types.Generic{Get: d.get, Set: d.set}
	return g.Apply(ctx, c, d)
}

func (d Data) get(key interface{}) changes.Value {
	switch key {
	case "Pending":
		return changes.Atomic{Value: d.Pending}
	case "TransformedOp":
		return d.TransformedOp
	case "MergeOps":
		return d.MergeOps
	}

	// default = "Version"
	return changes.Atomic{Value: d.Version}
}

func (d Data) set(key interface{}, v changes.Value) changes.Value {
	switch key {
	case "Version":
		d.Version = v.(changes.Atomic).Value.(int)
	case "Pending":
		d.Pending = v.(changes.Atomic).Value.([]ops.Op)
	case "TransformedOp":
		d.TransformedOp = v.(CachedOp)
	case "MergeOps":
		d.MergeOps = v.(CachedOps)
	}
	return d
}

// CachedOp represents a map of version => operation
type CachedOp map[int]ops.Op

// Apply implements changes.Value
func (o CachedOp) Apply(ctx changes.Context, c changes.Change) changes.Value {
	g := types.Generic{Get: o.get, Set: o.set}
	return g.Apply(ctx, c, o)
}

func (o CachedOp) get(key interface{}) changes.Value {
	if x, ok := o[key.(int)]; ok {
		return changes.Atomic{Value: x}
	}
	return changes.Nil
}

func (o CachedOp) set(key interface{}, v changes.Value) changes.Value {
	result := CachedOp{}
	idx := key.(int)
	for k, v := range o {
		if k != idx {
			result[k] = v
		}
	}

	if v != changes.Nil {
		result[idx] = v.(changes.Atomic).Value.(ops.Op)
	}
	return result
}

// CachedOps represents a map of version => operations
type CachedOps map[int][]ops.Op

// Apply implements changes.Value
func (o CachedOps) Apply(ctx changes.Context, c changes.Change) changes.Value {
	g := types.Generic{Get: o.get, Set: o.set}
	return g.Apply(ctx, c, o)
}

func (o CachedOps) get(key interface{}) changes.Value {
	if x, ok := o[key.(int)]; ok {
		return changes.Atomic{Value: x}
	}
	return changes.Nil
}

func (o CachedOps) set(key interface{}, v changes.Value) changes.Value {
	result := CachedOps{}
	idx := key.(int)
	for k, v := range o {
		if k != idx {
			result[k] = v
		}
	}

	if v != changes.Nil {
		result[idx] = v.(changes.Atomic).Value.([]ops.Op)
	}
	return result
}
