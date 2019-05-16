// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/fred"
)

func TestRefMiss(t *testing.T) {
	scope := &fred.Scope{}
	env := &fred.Environ{Resolver: scope.Resolve, Cacher: &fred.Cache{}, Depth: 5}
	x := fred.Ref(fred.Fixed(fred.Error("boo"))).Eval(env)
	if x != fred.Error("ref: no such ref") {
		t.Error("Unexpected missing ref", x)
	}
}

func TestRefHit(t *testing.T) {
	scope := &fred.Scope{
		DefMap: &fred.DefMap{fred.Error("boo"): fred.Fixed(fred.Error("goo"))},
	}
	env := &fred.Environ{Resolver: scope.Resolve, Cacher: &fred.Cache{}, Depth: 5}
	x := fred.Ref(fred.Fixed(fred.Error("boo"))).Eval(env)
	if x != fred.Error("goo") {
		t.Error("Unexpected missing ref", x)
	}
}

func TestRefRecursion(t *testing.T) {
	scope := &fred.Scope{
		DefMap: &fred.DefMap{
			fred.Error("boo"): fred.Ref(fred.Fixed(fred.Error("goo"))),
			fred.Error("goo"): fred.Ref(fred.Fixed(fred.Error("boo"))),
		},
	}
	env := &fred.Environ{Resolver: scope.Resolve, Cacher: &fred.Cache{}, Depth: 5}

	var r interface{}
	func() {
		defer func() { r = recover() }()
		fred.Ref(fred.Fixed(fred.Error("boo"))).Eval(env)
	}()

	if r != fred.ErrRecursion {
		t.Errorf("Unexpected success %#v\n", r)
	}
}
