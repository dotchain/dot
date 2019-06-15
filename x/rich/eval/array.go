// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package eval implements evaluated objects
package eval

import (
	"errors"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/data"
)

var errInvalidArgs = errors.New("invalid arguments")

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

func (a array) reduce(s Scope, args []changes.Value) changes.Value {
	if len(args) != 2 {
		return changes.Atomic{Value: errInvalidArgs}
	}

	var result = Eval(s, args[0])
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
