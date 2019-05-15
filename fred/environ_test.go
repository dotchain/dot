// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/fred"
)

var env = &fred.Environ{
	Cacher: &fred.Cache{},
	Depth:  5,
}

func TestEnvironRecursion(t *testing.T) {

	keys := []string{"goop", "world", "hello", "boo", "booya", "goop", "gomer"}
	idx := 0

	var fn func(e fred.Env) fred.Val
	fn = func(e fred.Env) fred.Val {
		key := keys[idx]
		idx++
		return e.CheckRecursion(nil, key, fn)
	}

	var r interface{}
	func() {
		defer func() { r = recover() }()
		env.CheckRecursion(nil, "hello", fn)
	}()

	// idx should not go past the second hello
	if idx != 3 {
		t.Error("Unexpected idx", idx, r)
	}
}

func TestEnvironDepth(t *testing.T) {
	keys := []string{"goop", "world", "hello", "boo", "booya", "goop", "gomer"}
	var idx uint = 0

	var fn func(e fred.Env) fred.Val
	fn = func(e fred.Env) fred.Val {
		key := keys[idx]
		idx++
		return e.CheckRecursion(nil, key, fn)
	}

	var r interface{}
	func() {
		defer func() { r = recover() }()
		env.CheckRecursion(nil, "zing", fn)
	}()

	// idx should not go past 5
	if idx != env.Depth+1 {
		t.Error("Unexpected idx", idx, r)
	}
}

func TestEnvironChainedMiss(t *testing.T) {
	goo := &fred.Fixed{Val: fred.Error("goo")}
	boo := &fred.Fixed{Val: fred.Error("boo")}
	parent := fred.Scope{DefMap: &fred.DefMap{"goo": goo}}
	child := fred.Scope{DefMap: &fred.DefMap{"boo": boo}}
	e := fred.Environ{
		Resolver: child.Resolve,
		Parent:   &fred.Environ{Resolver: parent.Resolve},
	}
	if def, env := e.Resolve("noo"); def != nil || env != nil {
		t.Error("Unexpected", def, env)
	}
}

func TestScopeChainedResolveHit(t *testing.T) {
	goo := &fred.Fixed{Val: fred.Error("goo")}
	boo := &fred.Fixed{Val: fred.Error("boo")}
	parent := fred.Scope{DefMap: &fred.DefMap{"goo": goo}}
	child := fred.Scope{DefMap: &fred.DefMap{"boo": boo}}
	e := &fred.Environ{
		Resolver: child.Resolve,
		Parent:   &fred.Environ{Resolver: parent.Resolve},
	}
	def, env := e.Resolve("boo")
	if env != e || !reflect.DeepEqual(def, boo) {
		t.Error("Unexpected", def, env)
	}

	def, env = e.Resolve("goo")
	if env != e.Parent || !reflect.DeepEqual(def, goo) {
		t.Error("Unexpected", def)
	}
}
