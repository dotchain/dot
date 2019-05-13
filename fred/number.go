// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// Number uses a rational string format (i.e. p/q)
type Number string

// Apply implements changes.Value
func (n Number) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return n
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, n)
}

// Eval evaluates to itself
func (n Number) Eval(dir *DirStream) Object {
	return n
}

// Diff returns the change from old to new
func (n Number) Diff(old, next *DirStream, c changes.Change) changes.Change {
	return nil
}
