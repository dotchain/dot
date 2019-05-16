// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// ErrNotCallable is returned when calling a non-callable function
var ErrNotCallable = Error("not callable")

type callable struct {
	fn   Def
	args *Defs
}

func (cx *callable) get(key interface{}) changes.Value {
	if key == "Func" {
		return cx.fn
	}
	return cx.args
}

func (cx *callable) set(key interface{}, val changes.Value) changes.Value {
	if key == "Func" {
		return &callable{val.(Def), cx.args}
	}
	return &callable{cx.fn, val.(*Defs)}
}

func (cx *callable) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: cx.set, Get: cx.get}).Apply(ctx, c, cx)
}

func (cx *callable) Eval(e Env) Val {
	return e.ValueOf(cx, func() Val {
		v := cx.fn.Eval(e)
		if callable, ok := v.(Callable); ok {
			return callable.Call(e, cx.args)
		}
		return ErrNotCallable
	})
}

// Call evaluates the first arg and passes the rest of the defs as args to it.
//
// The provided fn should evaluate to a Callable or this results in an error.
func Call(fn Def, args ...Def) Def {
	a := Defs(args)
	return &callable{fn: fn, args: &a}
}

// method implements a Val type that when called, passes through to the provided fn
type method func(Env, *Defs) Val

func (m method) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return changes.Nil
}

func (m method) Call(e Env, args *Defs) Val {
	return m(e, args)
}

func (m method) Text() string {
	return "<method>"
}

func (m method) Visit(v Visitor) {
	v.VisitLeaf(m)
}
