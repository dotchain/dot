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

// Cache implements Cacher
type Cache struct {
	Values    map[interface{}]Val
	Changes   map[interface{}]changes.Change
	Defs      map[interface{}]Def
	Resolvers map[interface{}]Resolver
	Untypeds  map[interface{}]interface{}
}

// ValueOf implements a value fetcher
func (cx *Cache) ValueOf(key interface{}, val func() Val) Val {
	cx.init()
	if _, ok := cx.Values[key]; !ok {
		cx.Values[key] = val()
	}
	return cx.Values[key]
}

// ChangeOf implements a change fetcher
func (cx *Cache) ChangeOf(key interface{}, c func() changes.Change) changes.Change {
	cx.init()
	if _, ok := cx.Changes[key]; !ok {
		cx.Changes[key] = c()
	}
	return cx.Changes[key]
}

// DefOf implements a def fetch
func (cx *Cache) DefOf(key interface{}, def func() Def) Def {
	cx.init()
	if _, ok := cx.Defs[key]; !ok {
		cx.Defs[key] = def()
	}
	return cx.Defs[key]
}

// ResolverOf implements a resolver fetcher
func (cx *Cache) ResolverOf(key interface{}, resolver func() Resolver) Resolver {
	cx.init()
	if _, ok := cx.Resolvers[key]; !ok {
		cx.Resolvers[key] = resolver()
	}
	return cx.Resolvers[key]
}

// UntypedOf implements a generic fetcher
func (cx *Cache) UntypedOf(key interface{}, untyped func() interface{}) interface{} {
	cx.init()
	if _, ok := cx.Untypeds[key]; !ok {
		cx.Untypeds[key] = untyped()
	}
	return cx.Untypeds[key]
}

func (cx *Cache) init() {
	if cx.Values == nil {
		cx.Values = map[interface{}]Val{}
		cx.Changes = map[interface{}]changes.Change{}
		cx.Defs = map[interface{}]Def{}
		cx.Resolvers = map[interface{}]Resolver{}
		cx.Untypeds = map[interface{}]interface{}{}
	}

}
