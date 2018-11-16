// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import "github.com/dotchain/dot/changes"

// The M type represents a map of arbitrary key/values. It implements
// the changes.Value interface.  The actual values can be non-nil
// (unlike the regular requirement for values to be non-nil. In this
// case, the value is treated as if it were changes.Nil
type M map[interface{}]changes.Value

// Apply applies the change and returns the updated value
func (m M) Apply(c changes.Change) changes.Value {
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
			clone := map[interface{}]changes.Value{}
			for k, v := range m {
				clone[k] = v
			}
			if clone[c.Path[0]] == nil {
				clone[c.Path[0]] = changes.Nil
			}
			clone[c.Path[0]] = clone[c.Path[0]].Apply(changes.PathChange{c.Path[1:], c.Change})
			if clone[c.Path[0]] == changes.Nil {
				clone[c.Path[0]] = nil
			}
			return M(clone)
		}
		return c.ApplyTo(m)
	case changes.Custom:
		return c.ApplyTo(m)
	}
	panic("Unexpected change on Apply")
}
