// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

// ErrMaxDepthReached is fired when stack depth exceeds configured value
var ErrMaxDepthReached = Error("max stack depth reached")

// ErrRecursion is fired when recursion is detected
var ErrRecursion = Error("recursion not allowed")

// Environ implements Env.
//
// Use Scope for an implementation of Resolver and Cache for Cacher.
//
// If MaxDepth is specified, that is used to check max depth
type Environ struct {
	Parent   *Environ
	Resolver func(key interface{}) Def
	Cacher
	Depth uint
}

// Resolve maps the key to the def and returns the environment where
// it was found in
func (e *Environ) Resolve(key interface{}) (Def, Env) {
	if def := e.Resolver(key); def != nil {
		return def, e
	}
	if e.Parent != nil {
		return e.Parent.Resolve(key)
	}
	return nil, nil
}

// CheckRecursion checks if the provided scope/key pair were already
// used in the current invocation stack
func (e *Environ) CheckRecursion(scope interface{}, key interface{}, fn func(inner Env) Val) Val {
	return fn(&environ{e, e.Cacher, e.Depth, [2]interface{}{scope, key}, nil})
}

// UseCheckerFrom switches just the recursion/env from the other environment
func (e *Environ) UseCheckerFrom(other Env) Env {
	result := &environ{e, e.Cacher, e.Depth, nil, nil}
	return result.UseCheckerFrom(other)
}

type environ struct {
	Resolver
	Cacher
	depth  uint
	key    interface{}
	parent *environ
}

func (e *environ) CheckRecursion(scope interface{}, key interface{}, fn func(inner Env) Val) Val {
	if e.depth == 0 {
		panic(ErrMaxDepthReached)
	}

	key = [2]interface{}{scope, key}
	for x := e; x != nil; x = x.parent {
		if x.key == key {
			panic(ErrRecursion)
		}
	}
	return fn(&environ{e.Resolver, e.Cacher, e.depth - 1, key, e})
}

func (e *environ) UseCheckerFrom(other Env) Env {
	result := *e
	result.key = other.(*environ).key
	result.parent = other.(*environ).parent
	result.depth = other.(*environ).depth
	return &result
}
