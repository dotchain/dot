// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Atomic is an atomic Value. It can wrap any particular value and can
// be used in the Before, After fields of Replace or Splice.
type Atomic struct {
	Value interface{}
}

// Slice should not be called on Atomic values
func (a Atomic) Slice(offset, count int) Value {
	panic("Slice should  not be called on atomic values")
}

// Count should not be called on Atomic values
func (a Atomic) Count() int {
	panic("Count should not be called on atomic values")
}

// Apply only accepts one type of change: one that Replace's the
// value.
func (a Atomic) Apply(c Change) Value {
	switch c := c.(type) {
	case nil:
		return a
	case Replace:
		if !c.IsInsert {
			return c.After
		}
	case PathChange:
		if len(c.Path) == 0 {
			return a.Apply(c.Change)
		}
	case ChangeSet:
		v := Value(a)
		for _, cx := range c {
			if cx != nil {
				v = v.Apply(cx)
			}
		}
		return v
	}
	panic("Unexpected change applied on atomic")
}
