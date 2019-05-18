// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package red

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

// Comma implements a comma operator
type Comma struct {
	Left, Right fred.Def
}

// Apply implements changes.Value
func (c *Comma) Apply(ctx changes.Context, cx changes.Change) changes.Value {
	return c
}

// Eval picks off the right side of the tree
func (c *Comma) Eval(e fred.Env) fred.Val {
	return c.Right.Eval(e)
}

// Args returns the flattened list of args
func (c *Comma) Args() []fred.Def {
	if x, ok := c.Left.(fnArgs); ok {
		return append(x.Args(), c.Right)
	}
	return []fred.Def{c.Left, c.Right}
}
