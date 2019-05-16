// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// ErrNotBool is returned with & and | expressions where args are not bool
var ErrNotBool = Error("not boolean")

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

// Field implements "&" and "|"
func (b Bool) Field(e Env, key Val) Val {
	switch key {
	case Text("&"):
		return b.op(true, func(b bool) *bool {
			if !b {
				return &b
			}
			return nil
		})
	case Text("|"):
		return b.op(false, func(b bool) *bool {
			if b {
				return &b
			}
			return nil
		})
	}
	return ErrNoSuchField
}

func (b Bool) op(zero bool, fn func(b bool) *bool) Val {
	return method(func(e Env, defs *Defs) Val {
		if defs != nil {
			if x := fn(bool(b)); x != nil {
				return Bool(*x)
			}

			for _, arg := range *defs {
				r := arg.Eval(e)
				switch r := r.(type) {
				case Bool:
					if x := fn(bool(r)); x != nil {
						return Bool(*x)
					}
				case Error:
					return r
				default:
					return ErrNotBool
				}
			}
		}
		return Bool(zero)
	})
}
