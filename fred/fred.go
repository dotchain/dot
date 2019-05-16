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
	Cacher
	RecursionChecker
	Resolver
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
	UseCheckerFrom(Env) Env
}

// Resolver resolves a "name"
type Resolver interface {
	Resolve(name interface{}) (Def, Env)
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

	// Text converts the val into some string form
	Text() string

	// Visit is used to visit the val.  Val then calls
	// Vistor.Leaf (for leaf nodes) or Visitor.Child (for each child)
	Visit(visitor Visitor)
}

// Fieldable values implement Field
type Fieldable interface {
	// Field returns a field (say, via .key).  Also used for index
	Field(e Env, key Val) Val
}

// Callable values implement Call
type Callable interface {
	// Call calls the value with the args.  Note that the args are
	// not evaluated (this allows the callable to lazy evaluate them)
	Call(e Env, args *Defs) Val
}

// Visitor is the interface that all value visitors should implement
type Visitor interface {
	VisitLeaf(v Val)
	VisitChildrenBegin(v Val)
	VisitChild(v Val, key interface{})
	VisitChildrenEnd(v Val)
}
