// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	_ "github.com/dotchain/dot/encoding/set"
	_ "github.com/dotchain/dot/encoding/sparse"
	"testing"
)

func TestUtilsAreSame(t *testing.T) {
	u := dot.Utils(dot.Transformer{})
	same := [][]interface{}{
		{nil, nil},
		{"", nil},
		{"", ""},
		{nil, []interface{}{}},
		{nil, map[string]interface{}{}},
		{nil, map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  []interface{}{},
		}},
		{"hello", "hello"},
		{[]interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
		{[]interface{}{1, 1}, map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  []interface{}{2, 1},
		}},
		{[]interface{}{nil}, []interface{}{""}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, pairs := range same {
		if !u.AreSame(pairs[0], pairs[1]) {
			t.Errorf("Expected '%#v' and '%#v' to be the same\n", pairs[0], pairs[1])
		}
		if !u.AreSame(pairs[1], pairs[0]) {
			t.Errorf("Expected '%#v' and '%#v' to be the same\n", pairs[1], pairs[0])
		}
	}
}

func TestUtilsAreNotSame(t *testing.T) {
	u := dot.Utils(dot.Transformer{})
	diff := [][]interface{}{
		{nil, 1},
		{nil, false},
		{nil, "hello"},
		{nil, []interface{}{1}},
		{nil, []int{}},
		{nil, map[string]interface{}{"a": 1}},
		{1.723, 1.721},
		{[]int{1}, []interface{}{1}},
		{[]interface{}{1, 2}, map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  []interface{}{2, 1},
		}},
		{"hello", "world"},
		{[]interface{}{1}, []interface{}{1, 2}},
		{[]interface{}{1}, map[string]interface{}{"a": 42}},
		{map[string]interface{}{"a": 42}, map[string]interface{}{"a": 43}},
		{map[string]interface{}{"a": 42}, map[string]interface{}{"b": 42}},
		{map[string]interface{}{"a": 42}, map[string]interface{}{"a": 42, "b": 40}},
	}

	for _, pairs := range diff {
		if u.AreSame(pairs[0], pairs[1]) {
			t.Errorf("Expected '%#v' and '%#v' to be different\n", pairs[0], pairs[1])
		}
		if u.AreSame(pairs[1], pairs[0]) {
			t.Errorf("Expected '%#v' and '%#v' to be different\n", pairs[1], pairs[0])
		}
	}
}
