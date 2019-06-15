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

type dict types.M

func (d dict) forEach(s Scope, args []changes.Value) changes.Value {
	if len(args) != 1 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	result := types.M{}
	for key, elt := range d {
		k, ok := key.(changes.Value)
		if !ok {
			k = changes.Atomic{Value: key}
		}
		result[key] = Eval(s, &data.Dir{
			Root: args[0],
			Objects: types.M{
				types.S16("key"):   k,
				types.S16("value"): elt,
			},
		})
	}
	return result
}

func (d dict) filter(s Scope, args []changes.Value) changes.Value {
	if len(args) != 1 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	result := types.M{}
	for key, elt := range d {
		// TODO: elt eval should be lazy?
		elt = Eval(s, elt)
		k, ok := key.(changes.Value)
		if !ok {
			k = changes.Atomic{Value: key}
		}
		v := Eval(s, &data.Dir{
			Root: args[0],
			Objects: types.M{
				types.S16("key"):   k,
				types.S16("value"): elt,
			},
		})
		a, _ := v.(changes.Atomic)
		b, _ := a.Value.(bool)
		if b {
			result[key] = elt
		}
	}
	return result
}

func (d dict) reduce(s Scope, args []changes.Value) changes.Value {
	if len(args) != 2 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	result := Eval(s, args[0])
	for key, elt := range d {
		k, ok := key.(changes.Value)
		if !ok {
			k = changes.Atomic{Value: key}
		}
		result = Eval(s, &data.Dir{
			Root: args[1],
			Objects: types.M{
				types.S16("key"):   k,
				types.S16("value"): elt,
				types.S16("last"):  result,
			},
		})
	}
	return result
}

func (d dict) getField(field changes.Value) changes.Value {
	if v, ok := d[field]; ok {
		return v
	}

	switch field {
	case types.S16("count"):
		return changes.Atomic{Value: len(d)}
	case types.S16("filter"):
		return Callable(d.filter)
	case types.S16("map"):
		return Callable(d.forEach)
	case types.S16("reduce"):
		return Callable(d.reduce)
	}

	return changes.Atomic{Value: errUnknownField}
}
