// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package red

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

// Group represents a bracketed item
type Group struct {
	inner fred.Def
}

// Apply is just a stub
func (g *Group) Apply(ctx changes.Context, cx changes.Change) changes.Value {
	return g
}

// Eval implements fred.Def
func (g *Group) Eval(e fred.Env) fred.Val {
	return g.inner.Eval(e)
}

// Args returns the items
func (g *Group) Args() []fred.Def {
	if x, ok := g.inner.(fnArgs); ok {
		return x.Args()
	}
	return []fred.Def{g.inner}
}

type fnArgs interface {
	Args() []fred.Def
}
