// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Atomic is an atomic Value. It can wrap any particular value and can
// be used in the Before, After fields of Replace or Splice.
type Atomic struct {
	Value interface{}
}

// Apply only accepts one type of change: one that Replace's the
// value.
func (a Atomic) Apply(ctx Context, c Change) Value {
	switch c := c.(type) {
	case nil:
		return a
	case Replace:
		if !c.IsCreate() {
			return c.After
		}
	}
	return c.(Custom).ApplyTo(ctx, a)
}
