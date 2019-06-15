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

// NumLess compares if numbers are in ascending order
var NumLess Callable = compare(func(x, y int) bool {
	return x < y
})

// NumLessThanEqual is like NumEqual but with equality allowed
var NumLessThanEqual Callable = compare(func(x, y int) bool {
	return x <= y
})

// NumMore compares if numbers are in descending order
var NumMore Callable = compare(func(x, y int) bool {
	return x > y
})

// NumMoreThanEqual compares if numbers are in descending order
var NumMoreThanEqual Callable = compare(func(x, y int) bool {
	return x >= y
})

// Equal compares if value are all the same
var Equal Callable = equality(func(x, y changes.Value) bool {
	// TODO: do deep comparison
	return x == y
})

// NotEqual compares if values are all different
var NotEqual Callable = equality(func(x, y changes.Value) bool {
	// TODO: do deep difference
	return x != y
})

func compare(fn func(x, y int) bool) func(s Scope, args []changes.Value) changes.Value {
	return func(s Scope, args []changes.Value) changes.Value {
		args = evalArray(s, args).(types.A)
		for kk, current := range args[1:] {
			before, ok1 := args[kk].(changes.Atomic)
			b, ok2 := before.Value.(int)
			after, ok3 := current.(changes.Atomic)
			a, ok4 := after.Value.(int)
			if !ok1 || !ok2 || !ok3 || !ok4 || !fn(b, a) {
				return changes.Atomic{Value: false}
			}
		}
		return changes.Atomic{Value: true}
	}
}

func equality(fn func(x, y changes.Value) bool) func(s Scope, args []changes.Value) changes.Value {
	return func(s Scope, args []changes.Value) changes.Value {
		args = evalArray(s, args).(types.A)
		for kk, current := range args[1:] {
			if !fn(args[kk], current) {
				return changes.Atomic{Value: false}
			}
		}
		return changes.Atomic{Value: true}
	}
}

var errUnknownField = errors.New("unknown field")

func dot(s Scope, receiver, field changes.Value) changes.Value {
	receiver = Eval(s, receiver)
	field = Eval(s, field)
	if r, ok := receiver.(types.A); ok {
		return array(r).getField(field)
	}

	if d, ok := receiver.(types.M); ok {
		return dict(d).getField(field)
	}

	return changes.Atomic{Value: errors.New("unknown receiver")}
}
