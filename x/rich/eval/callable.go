// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval

import (
	"errors"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Callable wraps a function into a value type that can be called
type Callable func(s Scope, args []changes.Value) changes.Value

// Apply implements changes.Value
func (cx Callable) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return cx
	case changes.Replace:
		return c.After
	}
	return c.(changes.Custom).ApplyTo(ctx, cx)
}

// Sum is the "callable" which evaluates the sum of all args
//
// Non-numeric values are treated as zeros
var Sum Callable = func(s Scope, args []changes.Value) changes.Value {
	var sum int
	for _, arg := range args {
		atomic, _ := Eval(s, arg).(changes.Atomic)
		val, _ := atomic.Value.(int)
		sum += val
	}
	return changes.Atomic{Value: sum}
}

// Dot is the "callable" which evaluates args[0].args[1].args[2]...
var Dot Callable = func(s Scope, args []changes.Value) changes.Value {
	var receiver changes.Value = changes.Nil
	if len(args) > 0 {
		receiver = args[0]
		for _, arg := range args[1:] {
			receiver = dot(s, receiver, arg)
		}
	}
	return receiver
}

func dot(s Scope, receiver, field changes.Value) changes.Value {
	receiver = Eval(s, receiver)
	field = Eval(s, field)
	if r, ok := receiver.(types.A); ok {
		switch field {
		case types.S16("count"):
			return changes.Atomic{Value: r.Count()}
		case types.S16("map"):
			return Callable(array(r).forEach)
		case types.S16("reduce"):
			return Callable(array(r).reduce)
		default:
			return changes.Atomic{Value: errors.New("unknown field")}
		}
	}
	return changes.Atomic{Value: errors.New("unknown receiver")}
}
