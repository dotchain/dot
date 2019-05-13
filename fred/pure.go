// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// EvalWithArgs are functions that support evaluating based on
// provided args
type EvalWithArgs interface {
	Eval(e Env, args *Vals) Val
}

// Pure functions take a functor that implements EvalWithArgs
type Pure struct {
	Functor EvalWithArgs
	Args    *Defs
}

func (p *Pure) get(key interface{}) changes.Value {
	return p.Args
}

func (p *Pure) set(key interface{}, val changes.Value) changes.Value {
	return &Pure{p.Functor, val.(*Defs)}
}

// Apply implements changes.Value
func (p *Pure) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: p.set, Get: p.get}).Apply(ctx, c, p)
}

// Eval calls EvalWithArgs after resolving the arg defs
func (p *Pure) Eval(e Env) Val {
	return e.ValueOf(p, func() Val {
		return p.Functor.Eval(e, p.Args.Eval(e).(*Vals))
	})
}
