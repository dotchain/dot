// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// Nil represents an empty def and val
type Nil struct{}

// Apply implements changes.Value
func (n Nil) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return n
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, n)
}

// Eval returns itself
func (n Nil) Eval(e Env) Val {
	return n
}
