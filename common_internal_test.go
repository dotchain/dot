// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"github.com/dotchain/dot/encoding"
	"strings"
	"testing"
)

var x = Transformer{}

func detailsFromChanges(info []Change) []interface{} {
	result := []interface{}{}
	for _, ii := range info {
		result = append(result, strings.Join(ii.Path, "."))
		if ii.Splice != nil {
			result = append(result, *ii.Splice)
		} else if ii.Set != nil {
			result = append(result, *ii.Set)
		} else if ii.Move != nil {
			result = append(result, *ii.Move)
		}
	}
	return result
}

func testMerge(t *testing.T, input interface{}, leftOps []Change, rightOps []Change) {
	testMergeFiltered(t, input, leftOps, rightOps, nil)
}

func testMergeFiltered(t *testing.T, input interface{}, leftOps []Change, rightOps []Change, shouldIgnore func(left, right Change) bool) {
	// now for each pair of operations, test if they converge
	for _, left := range leftOps {
		for _, right := range rightOps {
			if shouldIgnore != nil && shouldIgnore(left, right) {
				continue
			}
			left1, right1 := x.mergeChange(left, right)
			leftAll := append([]Change{left}, left1...)
			rightAll := append([]Change{right}, right1...)

			leftOutput := Utils(x).Apply(input, leftAll)
			rightOutput := Utils(x).Apply(input, rightAll)

			if !Utils(x).AreSame(leftOutput, rightOutput) {
				t.Fatalf("Failed %#v %#v %s != %s\n", detailsFromChanges(leftAll), detailsFromChanges(rightAll), leftOutput, rightOutput)
			}
		}
	}
}

func array(a ...interface{}) []interface{} {
	return a
}

func arraySize(i interface{}) int {
	if i == nil {
		return encoding.Get(i).Count()
	}
	return 0
}
