// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Call represents a "call" expression. The first element of the array
// is expected to be "Callable" which can be called with the rest as
// args.
//
// All errors result in a nil result.
type Call struct {
	types.A
}

// Eval evaluates a call expression in the provided scope
func (cx *Call) Eval(s Scope) changes.Value {
	if len(cx.A) == 0 {
		return changes.Nil
	}

	v := Eval(s, cx.A[0])
	if fn, ok := v.(Callable); ok {
		return fn(s, cx.A[1:])
	}

	return changes.Nil
}

// Apply implements changes.Value.
func (cx *Call) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: cx.set, Get: cx.get}).Apply(ctx, c, cx)
}

func (cx *Call) get(key interface{}) changes.Value {
	return cx.A
}

func (cx *Call) set(key interface{}, v changes.Value) changes.Value {
	return &Call{v.(types.A)}
}
