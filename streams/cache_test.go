// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"github.com/dotchain/dot/streams"
	"testing"
)

func TestCache(t *testing.T) {
	c := streams.Cache{}

	n := &streams.Notifier{}
	c.Begin()
	if v, _, ok := c.GetSubstream(n, "key"); ok {
		t.Fatal("Get should not return a value", v)
	}
	closed := 0
	c.SetSubstream(n, "key1", 1, &streams.Handler{func() {}}, func() { closed++ })
	c.SetSubstream(n, "key2", 2, &streams.Handler{func() {}}, func() { closed++ })
	c.End()

	c.Begin()
	if v, _, ok := c.GetSubstream(n, "key"); ok {
		t.Fatal("Get should not return a value", v)
	}
	if v, _, ok := c.GetSubstream(n, "key1"); !ok || v != 1 {
		t.Fatal("Unexpected key1 value", v, ok)
	}
	if v, _, ok := c.GetSubstream(n, "key2"); !ok || v != 2 {
		t.Fatal("Unexpected key2 value", v, ok)
	}

	// update only key1
	c.SetSubstream(n, "key1", 1, &streams.Handler{func() {}}, func() { closed++ })
	c.End()

	if closed != 1 {
		t.Error("Unexpected", closed)
	}
}
