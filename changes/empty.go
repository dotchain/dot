// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Nil represents an empty value. It can be used with Replace or
// Splice to indicate that Before or After is empty. The only
// operation that can be applied is a Replace.  Count and Slice cannot
// be called on it.
var Nil = empty{}

// empty represents an empty value. This should be used instead of nil
// with Replace. The only operation it supports is replacing with
// another value.
type empty struct{}

// Slice should not be called on empty
func (e empty) Slice(offset, count int) Value {
	panic("Unexpected Slice call on empty value")
}

// Count should not be called on empty
func (e empty) Count() int {
	panic("Unexpected Count call on empty value")
}

// Apply can be called on empty but it only supports
// Replace{IsInsert:true} type changes.
func (e empty) Apply(c Change) Value {
	switch c := c.(type) {
	case nil:
		return e
	case Replace:
		if c.IsInsert && c.After != Nil {
			return c.After
		}
	case PathChange:
		if len(c.Path) == 0 {
			return e.Apply(c.Change)
		}
	case ChangeSet:
		v := Value(e)
		for _, cx := range c {
			if cx != nil {
				v = v.Apply(cx)
			}
		}
		return v
	}
	panic("Unexpected change applied on empty")
}
