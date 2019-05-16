// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// Bool implements Val
type Bool bool

// Apply implements changes.Value
func (b Bool) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return b
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, b)
}

// Text implements Val.Text
func (b Bool) Text() string {
	if b {
		return "true"
	}
	return "false"
}

// Visit implements Val.Visit
func (b Bool) Visit(v Visitor) {
	v.VisitLeaf(b)
}
