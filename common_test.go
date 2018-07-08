// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"testing"
)

var x = dot.Transformer{}

//
// This file has some common test routines
//

func concat(a ...[]interface{}) []interface{} {
	result := []interface{}{}
	for _, ii := range a {
		result = append(result, ii...)
	}
	return result
}

func testOps(t *testing.T, input, output interface{}, left []dot.Change, right []dot.Change) {
	left1, right1 := x.MergeChanges(left, right)
	leftOutput := applyMany(applyMany(input, left), left1)
	rightOutput := applyMany(applyMany(input, right), right1)
	if !dot.Utils(x).AreSame(leftOutput, rightOutput) {
		t.Errorf("Failed to converge! %#v %#v\n", leftOutput, rightOutput)
	} else if !dot.Utils(x).AreSame(leftOutput, output) {
		t.Errorf("Unexpected convergence. Expected %#v but got %#v\n", output, leftOutput)
	}
}

func applyMany(input interface{}, changes []dot.Change) interface{} {
	u := dot.Utils(dot.Transformer{})
	result := u.Apply(input, changes)
	undoChanges := dot.Operation{Changes: changes}.Undo().Changes
	undone := u.Apply(result, undoChanges)
	if !dot.Utils(x).AreSame(input, undone) {
		panic("Undo failed")
	}

	return result
}

func makeArray(input string, asArray bool) interface{} {
	if !asArray {
		return input
	}
	result := []interface{}{}
	for _, b := range []byte(input) {
		result = append(result, b)
	}
	return result
}
