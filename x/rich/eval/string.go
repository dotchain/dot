// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package eval implements evaluated objects
package eval

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

type str types.S16

func (s str) concat(scope Scope, args []changes.Value) changes.Value {
	return types.S16(string(s) + string(args[0].(types.S16)))
}

func (s str) getField(field changes.Value) changes.Value {
	switch field {
	case types.S16("count"):
		return changes.Atomic{Value: types.S16(s).Count()}
	case types.S16("concat"):
		return Callable(s.concat)
	}

	return changes.Atomic{Value: errUnknownField}
}
