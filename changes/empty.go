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

// Apply can be called on empty but it only supports
// Replace{IsInsert:true} type changes.
func (e empty) Apply(ctx Context, c Change) Value {
	switch c := c.(type) {
	case nil:
		return e
	case Replace:
		if c.IsCreate() {
			return c.After
		}
	}
	return c.(Custom).ApplyTo(ctx, e)
}
