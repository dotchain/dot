// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import "github.com/dotchain/dot/changes"

// The M type represents a map of arbitrary key/values. It implements
// the changes.Value interface.  As with all values, a nil value
// should be expressed via changes.Atomic{nil} to distinguish it from
// an attempt to remove a key.
type M map[interface{}]changes.Value

func (m M) get(key interface{}) changes.Value {
	if v, ok := m[key]; ok {
		return v
	}
	return changes.Nil
}

func (m M) set(key interface{}, v changes.Value) changes.Value {
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

// Apply applies the change and returns the updated value
func (m M) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (Generic{Get: m.get, Set: m.set}).Apply(ctx, c, m)
}
