// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package eval implements evaluated objects
package eval

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/data"
)

type array types.A

func (a array) forEach(s Scope, args []changes.Value) changes.Value {
	if len(args) != 1 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	result := make(types.A, len(a))
	for kk, elt := range a {
		result[kk] = Eval(s, &data.Dir{
			Root: args[0],
			Objects: types.M{
				types.S16("index"): changes.Atomic{Value: kk},
				types.S16("value"): elt,
			},
		})
	}
	return result
}

func (a array) filter(s Scope, args []changes.Value) changes.Value {
	if len(args) != 1 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	result := make(types.A, 0, len(a))
	for kk, elt := range a {
		// TODO: elt eval should be lazy?
		elt = Eval(s, elt)
		v := Eval(s, &data.Dir{
			Root: args[0],
			Objects: types.M{
				types.S16("index"): changes.Atomic{Value: kk},
				types.S16("value"): elt,
			},
		})
		a, _ := v.(changes.Atomic)
		b, _ := a.Value.(bool)
		if b {
			result = append(result, elt)
		}
	}
	return result
}

func (a array) reduce(s Scope, args []changes.Value) changes.Value {
	if len(args) != 2 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	result := Eval(s, args[0])
	for kk, elt := range a {
		result = Eval(s, &data.Dir{
			Root: args[1],
			Objects: types.M{
				types.S16("index"): changes.Atomic{Value: kk},
				types.S16("value"): elt,
				types.S16("last"):  result,
			},
		})
	}
	return result
}

func (a array) getField(field changes.Value) changes.Value {
	switch field {
	case types.S16("count"):
		return changes.Atomic{Value: len(a)}
	case types.S16("map"):
		return Callable(a.forEach)
	case types.S16("reduce"):
		return Callable(a.reduce)
	case types.S16("filter"):
		return Callable(a.filter)
	}

	return changes.Atomic{Value: errUnknownField}
}
