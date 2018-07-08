// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import (
	"github.com/dotchain/dot/encoding"
	"testing"
)

func TestCatalogFallback(t *testing.T) {
	cat := encoding.NewCatalog()
	external := SparseTest().initial
	direct := encoding.Get(external)
	indirect := cat.Get(external)
	ensureEqual(t, direct, indirect)
}

func TestCatalogInvalidType(t *testing.T) {
	shouldPanic(t, "unknown encoding", func() { encoding.Get([]int{42}) })
}

func TestCatalogRegistration(t *testing.T) {
	register := encoding.NewCatalog().RegisterConstructor
	shouldPanic(t, "non func", func() { register("ok", "hello") })
	shouldPanic(t, "one arg", func() { register("ok", func(encoding.Catalog) {}) })
	shouldPanic(t, "wrong type", func() { register("ok", func(encoding.Catalog, int) {}) })
	shouldPanic(t, "wrong type 2", func() { register("ok", func(int, map[string]interface{}) {}) })
	shouldPanic(t, "no return type", func() { register("ok", func(encoding.Catalog, map[string]interface{}) {}) })
	shouldPanic(t, "two returns", func() {
		register("ok", func(encoding.Catalog, map[string]interface{}) (int, int) { return 0, 0 })
	})
	shouldPanic(t, "invalid return type", func() {
		register("ok", func(encoding.Catalog, map[string]interface{}) int { return 0 })
	})
}
