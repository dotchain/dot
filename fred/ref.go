// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

type ref struct{}

func (r ref) Eval(e Env, args *Vals) Val {
	if def, env := e.Resolve((*args)[0]); def != nil {
		return e.CheckRecursion(nil, def, func(other Env) Val {
			return def.Eval(env.UseCheckerFrom(other))
		})
	}

	return Error("ref: no such ref")
}

// NewRef creates a new ref to whatever Def evaluates to
func NewRef(d Def) Def {
	return &Pure{Functor: &ref{}, Args: &Defs{d}}
}

// NewFixedRef creates a ref whose value is the lookup of the v in its environment
func NewFixedRef(v Val) Def {
	return &Pure{Functor: &ref{}, Args: &Defs{&Fixed{Val: v}}}
}
