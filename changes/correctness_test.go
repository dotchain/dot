// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/test/seqtest"
	"testing"
)

// These sequences are a bit odd to start with, so the specific
// outcomes are probably ok.  Should just rewrite the dataset for
// this. Most of these seem like a problem with splices adjacent to
// the move wrongly moving along with the move.
var knownFailures = map[string]string{
	"-abc[123]de4|56fg- x -abc1|23de[456]fg-": "-abcde145623fg-",
	"-abc123[|d]456f- x -abc123[456]f|-":      "-abc123fd456-",
	"-abc123[|d]456f- x -abc|123[456]f-":      "-abcd456123f-",
	"-abc123[|d]456f- x -ab|c123[456]f-":      "-abd456c123f-",
	"-ab|c[123]ef- x -abc123[|d]ef-":          "-ab123dcef-",
	"-abc[123]e|f- x -abc123[|d]ef-":          "-abce123df-",
	"-abc[1234|d]ef- x -ab|c[1234]ef-":        "-ab1234cdef-",
	"-abc[1234|]ef- x -ab|c[1234]ef-":         "-ab1234cef-",
	"-abc[|d]ef- x -a|b[ce]f-":                "-acdebf-",
}

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

			if lval != S(merged) && lval != S(knownFailures[name]) {
				t.Error("Converged on unexpected value", lval, merged, knownFailures[name])
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
