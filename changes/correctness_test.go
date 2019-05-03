// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/test/seqtest"
)

func TestCorrectnessOfStandardSequences(t *testing.T) {
	s := seqTester{}
	seqtest.ForEachTest(s, s, func(name, initial string, l, r interface{}, merged string) {
		if l == nil || r == nil {
			return
		}
		t.Run(name, func(t *testing.T) {
			left := l.(changes.Change)
			right := r.(changes.Change)
			leftx, rightx := left.Merge(right)
			lval := S(initial).Apply(nil, changes.ChangeSet{left, leftx})
			rval := S(initial).Apply(nil, changes.ChangeSet{right, rightx})
			if lval != rval {
				t.Error("Diverged", lval, rval, left, right)
			}

			if lval != S(merged) {
				t.Error("Converged on unexpected value", lval, merged)
			}
		})
	})
}

type seqTester struct{}

func (s seqTester) Splice(initial string, offset, count int, replacement string) interface{} {
	before := S(initial[offset : offset+count])
	after := S(replacement)
	return changes.Splice{offset, before, after}
}

func (s seqTester) Move(initial string, offset, count, distance int) interface{} {
	return changes.Move{offset, count, distance}
}

func (s seqTester) Range(initial string, offset, count int, attribute string) interface{} {
	return nil
}
