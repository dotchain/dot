// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt_test

import (
	"testing"

	"github.com/dotchain/dot/changes/crdt"
)

func TestNewRank(t *testing.T) {
	last := crdt.NewRank()
	for kk := 0; kk < 1000; kk++ {
		next := crdt.NewRank()
		if last.Less(next) && next.Less(last) {
			t.Fatal("Invalid rank comparisons", last, next)
		}
		last = next
	}
}
