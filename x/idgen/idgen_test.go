// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package idgen_test

import (
	"github.com/dotchain/dot/x/idgen"
	"testing"
)

func TestCollision(t *testing.T) {
	count := 100000
	seen := map[interface{}]bool{}
	for kk := 0; kk < count; kk++ {
		x := idgen.New()
		if seen[x] {
			t.Fatal("Collided on attempt", kk)
		}
		seen[x] = true
	}
}
