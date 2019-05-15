// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/fred"
)

func TestScopeResolveMiss(t *testing.T) {
	s := fred.Scope{}
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.Scope{DefMap: &fred.DefMap{}}
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.Scope{DefMap: &fred.DefMap{"goop": fred.Nil{}}}
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}
}

func TestScopeResolveHit(t *testing.T) {
	s := fred.Scope{DefMap: &fred.DefMap{"goop": fred.Nil{}, "boop": nil}}
	if def, r := s.Resolve("goop"); def != (fred.Nil{}) || !reflect.DeepEqual(r, s) {
		t.Error("Unexpected", def, r)
	}
}

func TestScopeChainedResolveMiss(t *testing.T) {
	s := fred.ChainResolver(
		fred.Scope{},
		fred.ChainResolver(
			fred.Scope{DefMap: &fred.DefMap{}},
			fred.Scope{},
		),
	)
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.ChainResolver(
		fred.Scope{},
		fred.ChainResolver(
			fred.Scope{DefMap: &fred.DefMap{}},
			fred.Scope{DefMap: &fred.DefMap{}},
		),
	)
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.ChainResolver(
		fred.Scope{},
		fred.ChainResolver(
			fred.Scope{DefMap: &fred.DefMap{}},
			fred.Scope{DefMap: &fred.DefMap{"goop": fred.Nil{}}},
		),
	)

	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}
}

func TestScopeChainedResolveHit(t *testing.T) {
	expected := fred.Scope{DefMap: &fred.DefMap{"goop": fred.Nil{}}}
	s := fred.ChainResolver(
		fred.Scope{},
		fred.ChainResolver(
			expected,
			fred.Scope{&fred.DefMap{"goop": nil}},
		),
	)
	if def, r := s.Resolve("goop"); def != (fred.Nil{}) || !reflect.DeepEqual(r, expected) {
		t.Error("Unexpected", def, r)
	}
}
