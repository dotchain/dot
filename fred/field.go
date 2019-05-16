// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// ErrNoSuchField is returned when the field does not exist
var ErrNoSuchField = Error("no such field")

// ErrNoFields is returned if there are no fields at all
var ErrNoFields = Error("no fields")

type fieldable struct {
	base Def
	args *Defs
}

func (f *fieldable) get(key interface{}) changes.Value {
	if key == "Base" {
		return f.base
	}
	return f.args
}

func (f *fieldable) set(key interface{}, val changes.Value) changes.Value {
	if key == "Base" {
		return &fieldable{val.(Def), f.args}
	}
	return &fieldable{f.base, val.(*Defs)}
}

func (f *fieldable) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: f.set, Get: f.get}).Apply(ctx, c, f)
}

func (f *fieldable) Eval(e Env) Val {
	return e.ValueOf(f, func() Val {
		v := f.base.Eval(e)

		if f.args != nil {
			for _, arg := range *f.args {
				key := arg.Eval(e)
				if fieldable, ok := v.(Fieldable); ok {
					v = fieldable.Field(e, key)
				} else {
					return ErrNoFields
				}
			}
		}
		return v
	})
}

// Field evaluates base.arg1.arg2 etc.
func Field(base Def, args ...Def) Def {
	a := Defs(args)
	return &fieldable{base, &a}
}
