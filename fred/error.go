// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// Error implements Val (as well as Error())
type Error string

// Error satisfies the error interface
func (e Error) Error() string {
	return string(e)
}

// Apply implements changes.Value
func (e Error) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return e
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, e)
}

// Text implements Val.Text
func (e Error) Text() string {
	return "err: " + string(e)
}

// Visit implements Val.Visit
func (e Error) Visit(v Visitor) {
	v.VisitLeaf(e)
}

// Field implements Fieldable but only returns itself no matter what
func (e Error) Field(ex Env, key Val) Val {
	return e
}

// Call implements Callable but only returns itself no matter what
func (e Error) Call(ex Env, args *Defs) Val {
	return e
}
