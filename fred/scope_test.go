// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/fred"
)

func TestScopeResolveMiss(t *testing.T) {
	s := fred.Scope{}
	if def := s.Resolve("boo"); def != nil {
		t.Error("Unexpected", def)
	}

	s = fred.Scope{DefMap: &fred.DefMap{}}
	if def := s.Resolve("boo"); def != nil {
		t.Error("Unexpected", def)
	}

	s = fred.Scope{DefMap: &fred.DefMap{"goop": fred.Nil()}}
	if def := s.Resolve("boo"); def != nil {
		t.Error("Unexpected", def)
	}
}

func TestScopeResolveHit(t *testing.T) {
	s := fred.Scope{DefMap: &fred.DefMap{"goop": fred.Nil(), "boop": nil}}
	if def := s.Resolve("goop"); def != fred.Nil() {
		t.Error("Unexpected", def)
	}
}
