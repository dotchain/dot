// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// Nil returns an empty def
func Nil() Def {
	return none{}
}

type none struct{}

func (n none) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return n
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, n)
}

func (n none) Eval(e Env) Val {
	return n
}

func (n none) Text() string {
	return "<nil>"
}

func (n none) Visit(v Visitor) {
	v.VisitLeaf(n)
}
