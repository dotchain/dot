// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/eval"
)

func TestEval(t *testing.T) {
	var globals types.M
	var scope eval.Scope

	scope = eval.Scope(func(v interface{}) changes.Value {
		return eval.Eval(scope, globals[v])
	})
	globals = types.M{
		types.S16("+"):     eval.Sum,
		types.S16("."):     eval.Dot,
		types.S16("<"):     eval.NumLess,
		types.S16("<="):    eval.NumLessThanEqual,
		types.S16("=="):    eval.Equal,
		types.S16(">"):     eval.NumMore,
		types.S16(">="):    eval.NumMoreThanEqual,
		types.S16("!="):    eval.NotEqual,
		types.S16("true"):  changes.Atomic{Value: true},
		types.S16("false"): changes.Atomic{Value: false},
		types.S16("list"):  eval.Parse(scope, "(1, 2, 3)"),
		types.S16("dict"):  eval.Parse(scope, "obj(x = 1, y = 5)"),
	}

	// map of code => expected result
	tests := map[string]string{
		"list.map(value+10)":           "(11, 12, 13)",
		"list.reduce(100, value+last)": "106",
		"list.filter(value >= 2)":      "(2, 3)",
		"list.count":                   "3",
		"dict.x + dict.y":              "6",
		"dict.count":                   "2",
		"dict.map(value+10).x":         "11",
		"dict.reduce(100,value+last)":  "106",
		"dict.filter(value >= 2)":      "obj(y=5)",
		"1 < 2":                        "true",
		"1 <= 2":                       "true",
		"1 < 0":                        "false",
		"1 <= 0":                       "false",
		"2 == 2":                       "true",
		"2 != 2":                       "false",
		"2 > 1":                        "true",
		"2 >= 2":                       "true",
		"0 > 1":                        "false",
		"0 >= 1":                       "false",
		"do(z + 2, z = list.count)":    "5",
	}

	for code, result := range tests {
		t.Run(code, func(t *testing.T) {
			got := eval.Eval(scope, eval.Parse(scope, code))
			expected := eval.Eval(scope, eval.Parse(scope, result))
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("got %#v", got)
			}
		})
	}
}
