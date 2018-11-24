// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import "github.com/dotchain/dot/changes"

// The M type represents a map of arbitrary key/values. It implements
// the changes.Value interface.  As with all values, a nil value
// should be expressed via changes.Atomic{nil} to distinguish it from
// an attempt to remove a key.
type M map[interface{}]changes.Value

// Apply applies the change and returns the updated value
func (m M) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return m
	case changes.Replace:
		if c.IsDelete() {
			return changes.Nil
		}
		return c.After
	case changes.PathChange:
		if len(c.Path) > 0 {
			key := c.Path[0]
			c.Path = c.Path[1:]
			return m.applyKey(ctx, key, c)
		}
	}
	return c.(changes.Custom).ApplyTo(ctx, m)
}

func (m M) applyKey(ctx changes.Context, key interface{}, c changes.Change) changes.Value {
	v, ok := m[key]
	if !ok {
		v = changes.Nil
	}
	v = v.Apply(ctx, c)
	result := make(M, len(m))
	for k, v := range m {
		if k != key {
			result[k] = v
		}
	}
	if v != changes.Nil {
		result[key] = v
	}
	return result
}
