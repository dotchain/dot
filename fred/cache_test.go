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

func TestCacheValues(t *testing.T) {
	var c fred.Cache

	x := c.ValueOf("boo", func() fred.Val { return fred.Error("boo") })
	if x != fred.Error("boo") {
		t.Error("Unexpected value", x)
	}

	x = c.ValueOf("boo", func() fred.Val { return nil })
	if x != fred.Error("boo") {
		t.Error("Unexpected value", x)
	}
}

func TestCacheChanges(t *testing.T) {
	var c fred.Cache

	v := changes.Replace{Before: fred.Error("b"), After: fred.Error("a")}
	x := c.ChangeOf("boo", func() changes.Change { return v })
	if x != v {
		t.Error("Unexpected value", x)
	}

	x = c.ChangeOf("boo", func() changes.Change { return nil })
	if x != v {
		t.Error("Unexpected value", x)
	}
}

func TestCacheDefs(t *testing.T) {
	var c fred.Cache

	v := &fred.Defs{}
	x := c.DefOf("boo", func() fred.Def { return v })
	if x != v {
		t.Error("Unexpected value", x)
	}

	x = c.DefOf("boo", func() fred.Def { return nil })
	if x != v {
		t.Error("Unexpected value", x)
	}
}

func TestCacheResolvers(t *testing.T) {
	var c fred.Cache

	v := fred.Scope{}
	x := c.ResolverOf("boo", func() fred.Resolver { return v })
	if !reflect.DeepEqual(x, v) {
		t.Error("Unexpected value", x)
	}

	x = c.ResolverOf("boo", func() fred.Resolver { return nil })
	if !reflect.DeepEqual(x, v) {
		t.Error("Unexpected value", x)
	}
}

func TestCacheUntyped(t *testing.T) {
	var c fred.Cache

	v := struct{ z int }{5}
	x := c.UntypedOf("boo", func() interface{} { return v })
	if !reflect.DeepEqual(x, v) {
		t.Error("Unexpected value", x)
	}

	x = c.UntypedOf("boo", func() interface{} { return nil })
	if !reflect.DeepEqual(x, v) {
		t.Error("Unexpected value", x)
	}
}
