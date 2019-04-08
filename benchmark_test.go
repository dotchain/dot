// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func BenchmarkMergeNonConflicting(b *testing.B) {
	left, right := changes.ChangeSet{}, changes.ChangeSet{}
	count := 1500

	for kk := 0; kk < count; kk++ {
		c := changes.Splice{Offset: kk, Before: types.S8(" "), After: types.S8("")}
		left = append(left, c)
		c.Offset += count
		right = append(right, c)
	}

	changes.Merge(left, right)
}

func BenchmarkMergeConflicting(b *testing.B) {
	left, right := changes.ChangeSet{}, changes.ChangeSet{}
	count := 1500

	for kk := 0; kk < count; kk++ {
		c := changes.Splice{Offset: kk, Before: types.S8(" "), After: types.S8("")}
		left = append(left, c)

		c.Before, c.After = c.After, c.Before
		right = append(right, c)
	}

	changes.Merge(left, right)
}
