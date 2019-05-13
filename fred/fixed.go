// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Fixed wraps a Val into a definition that does not change based on
// other definitions.
//
// Note that the definition itself can be changed as well as the
// underlying value (by simply passing a PathChange with key "Val")
// but the evaluated value for a given instance does not depend on the
// resolver.
//
// The provided value cannot be nil (though it can be fred.Nil{}).
type Fixed struct {
	Val
}

func (f *Fixed) get(key interface{}) changes.Value {
	if key != "Val" {
		panic("Unexpected key")
	}
	return f.Val
}

func (f *Fixed) set(key interface{}, val changes.Value) changes.Value {
	return &Fixed{Val: val.(Val)}
}

// Apply implements changes.Value
func (f *Fixed) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: f.set, Get: f.get}).Apply(ctx, c, f)
}

// Eval returns the wrapped val
func (f *Fixed) Eval(e Env) Val {
	return f.Val
}
