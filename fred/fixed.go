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
func Fixed(v Val) Def {
	return &fixed{val: v}
}

type fixed struct {
	val Val
}

func (f *fixed) get(key interface{}) changes.Value {
	if key != "Val" {
		panic("Unexpected key")
	}
	return f.val
}

func (f *fixed) set(key interface{}, val changes.Value) changes.Value {
	return &fixed{val.(Val)}
}

func (f *fixed) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: f.set, Get: f.get}).Apply(ctx, c, f)
}

func (f *fixed) Eval(e Env) Val {
	return f.val
}
