// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package fred implements a convergent functional reactive engine.
//
// It uses a pull-based functional reactive system (instead of push
// based setup) to make it easy to only rebuild the parts of the graph
// needed.
//
// The main entrypoint is Dir which maintains a directory of object
// defintions.  Objects are indexed by an ID (string) and values can
// be fixed or derived.  Fixed values evaluate to themselves while
// derived values evaluate to a different result.
package fred

import "github.com/dotchain/dot/changes"

// Env is the runtime environment.
type Env interface {
	Resolver
	Cacher
	RecursionChecker
}

// Cacher is a generic cache
type Cacher interface {
	ValueOf(key interface{}, val func() Val) Val
	ChangeOf(key interface{}, c func() changes.Change) changes.Change
	DefOf(key interface{}, def func() Def) Def
	ResolverOf(key interface{}, resolver func() Resolver) Resolver
	UntypedOf(key interface{}, fn func() interface{}) interface{}
}

// RecursionChecker is used by Eval to track recursion.
type RecursionChecker interface {
	CheckRecursion(scope interface{}, key interface{}, fn func(inner Env) Val) Val
}

// Resolver resolves a "name"
type Resolver interface {
	// Lookup returns the definition and the associated scope
	// it was found in.  Both can be nil if the key is not found.
	Resolve(key interface{}) (Def, Resolver)
}

// Def holds the definition of an object value to be evaluated.
type Def interface {
	changes.Value

	// Eval evaluates a definition in the provided scope.
	Eval(Env) Val
}

// Val implements an actual immutable value.
type Val interface {
	changes.Value
}
