// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import "github.com/dotchain/dot/changes"

// Generic is a helper to build value and collection types
type Generic struct {
	Get    func(key interface{}) changes.Value
	Set    func(key interface{}, v changes.Value) changes.Value
	Splice func(offset, count int, insert changes.Collection) changes.Collection
}

// Apply applies a change to produce a new value. Note that this
// requires the current value as input unlike changes.Value.
func (g Generic) Apply(ctx changes.Context, c changes.Change, current changes.Value) changes.Value {
	switch c := c.(type) {
	case nil:
		return current
	case changes.Replace:
		if c.IsDelete() {
			return changes.Nil
		}
		return c.After
	case changes.PathChange:
		if len(c.Path) == 0 {
			return g.Apply(ctx, c.Change, current)
		}
		inner := changes.PathChange{c.Path[1:], c.Change}
		return g.Set(c.Path[0], g.Get(c.Path[0]).Apply(ctx, inner))
	case changes.Splice, changes.Move:
		return g.ApplyCollection(ctx, c, current.(changes.Collection))
	}

	return c.(changes.Custom).ApplyTo(ctx, current)
}

// ApplyCollection implements changes.ApplyCollection but with
// using the provided current value and the helper functions
func (g Generic) ApplyCollection(ctx changes.Context, c changes.Change, current changes.Collection) changes.Collection {
	switch c := c.(type) {
	case changes.Splice:
		return g.Splice(c.Offset, c.Before.Count(), c.After)
	case changes.Move:
		c = c.Normalize()
		slice := current.Slice(c.Offset, c.Count)
		empty := current.Slice(0, 0)
		splice := changes.Splice{c.Offset + c.Distance, empty, slice}
		return g.Splice(c.Offset, c.Count, empty).ApplyCollection(ctx, splice)
	}
	return g.Apply(ctx, c, current).(changes.Collection)
}
