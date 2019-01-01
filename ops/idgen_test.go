// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"github.com/dotchain/dot/ops"
	"testing"
)

func TestIDCollision(t *testing.T) {
	count := 100000
	seen := map[interface{}]bool{}
	for kk := 0; kk < count; kk++ {
		x := ops.NewID()
		if seen[x] {
			t.Fatal("Collided on attempt", kk)
		}
		seen[x] = true
	}
}
