// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/fred"
)

var env = fred.Environ{
	Resolver: fred.Scope{},
	Cacher:   &fred.Cache{},
	Depth:    5,
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
