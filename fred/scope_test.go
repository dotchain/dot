// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func TestScopeApply(t *testing.T) {
	s := fred.Scope(nil)
	s2 := s.Apply(nil, changes.PathChange{
		Path: []interface{}{"goop"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  fred.Nil{},
		},
	})
	expected := fred.Scope{"goop": fred.Nil{}}
	if !reflect.DeepEqual(expected, s2) {
		t.Error("Unexpected", s2)
	}
	s3 := s2.Apply(nil, changes.PathChange{
		Path: []interface{}{"goop"},
		Change: changes.Replace{
			Before: fred.Nil{},
			After:  changes.Nil,
		},
	})
	if !reflect.DeepEqual(s3, fred.Scope{}) {
		t.Error("Unexpected", s3)
	}

	s4 := s2.Apply(nil, changes.PathChange{
		Path: []interface{}{"boop"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  fred.Nil{},
		},
	})
	if !reflect.DeepEqual(s4, fred.Scope{"goop": fred.Nil{}, "boop": fred.Nil{}}) {
		t.Error("Unexpected", s4)
	}
}

func TestScopeResolveMiss(t *testing.T) {
	s := fred.Scope(nil)
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.Scope{}
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.Scope{"goop": fred.Nil{}}
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}
}

func TestScopeResolveHit(t *testing.T) {
	s := fred.Scope{"goop": fred.Nil{}, "boop": nil}
	if def, r := s.Resolve("goop"); def != (fred.Nil{}) || !reflect.DeepEqual(r, s) {
		t.Error("Unexpected", def, r)
	}
}

func TestScopeChainedResolveMiss(t *testing.T) {
	s := fred.Scope(nil).WithParent((fred.Scope{}).WithParent(fred.Scope(nil)))
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.Scope(nil).WithParent((fred.Scope{}).WithParent(fred.Scope{}))
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}

	s = fred.Scope(nil).WithParent((fred.Scope{}).WithParent(fred.Scope{"goop": fred.Nil{}}))
	if def, r := s.Resolve("boo"); def != nil || r != nil {
		t.Error("Unexpected", def, r)
	}
}

func TestScopeChainedResolveHit(t *testing.T) {
	expected := fred.Scope{"goop": fred.Nil{}}
	s := fred.Scope(nil).WithParent(expected.WithParent(fred.Scope{"goop": nil}))
	if def, r := s.Resolve("goop"); def != (fred.Nil{}) || !reflect.DeepEqual(r, expected) {
		t.Error("Unexpected", def, r)
	}
}
