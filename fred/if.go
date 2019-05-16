// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

// ErrInvalidCondition is returned if the condition is not valid
var ErrInvalidCondition = Error("invalid condition")

// If evaluates success if condition is true, otherwise evaluates fail.
//
// If condition is not boolean, it returns an error
func If(cond, success, fail Def) Def {
	check := Fixed(method(func(e Env, defs *Defs) Val {
		return e.ValueOf(defs, func() Val {
			d := *defs
			cond := d[0].Eval(e)
			b, ok := cond.(Bool)
			if !ok {
				if err, ok := cond.(Error); ok {
					return err
				}
				return ErrInvalidCondition
			}
			if b {
				return d[1].Eval(e)
			}
			return d[2].Eval(e)
		})
	}))
	return Call(check, cond, success, fail)
}
