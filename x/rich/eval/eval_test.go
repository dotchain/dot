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
		types.S16("+"):      eval.Sum,
		types.S16("."):      eval.Dot,
		types.S16("<"):      eval.NumLess,
		types.S16("<="):     eval.NumLessThanEqual,
		types.S16("=="):     eval.Equal,
		types.S16(">"):      eval.NumMore,
		types.S16(">="):     eval.NumMoreThanEqual,
		types.S16("!="):     eval.NotEqual,
		types.S16("true"):   changes.Atomic{Value: true},
		types.S16("false"):  changes.Atomic{Value: false},
		types.S16("list"):   eval.Parse(scope, "(1, 2, 3)"),
		types.S16("dict"):   eval.Parse(scope, "obj(x = 1, y = 5)"),
		types.S16("intmap"): types.M{5: types.S16("five")},
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

		"'hello'.count":            "5",
		"'hello'.concat(' world')": "'hello world'",

		"intmap.filter(key != 5)": "obj()",
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

func TestEvalErrors(t *testing.T) {
	s := eval.Scope(func(v interface{}) changes.Value {
		if v == types.S16(".") {
			return eval.Dot
		}
		panic(string(v.(types.S16)))
	})
	fails := []string{
		"(1, 2, 3).map()",
		"(1, 2, 3).filter()",
		"(1, 2, 3).reduce(100)",
		"(1, 2, 3).boo",
		"(1, 2, 3)()",
		"(1, 2, 3).count.boo",
		"obj(x = 5).map()",
		"obj(x = 5).filter()",
		"obj(x = 5).reduce(100)",
		"obj(x = 5).boo",
		"'boo'.boo",
	}

	for _, fail := range fails {
		t.Run(fail, func(t *testing.T) {
			got := eval.Eval(s, eval.Parse(nil, fail))
			atomic, _ := got.(changes.Atomic)
			err, _ := atomic.Value.(error)
			if err == nil {
				t.Errorf("Got unexpected non-error %#v", got)
			}
		})
	}
}

func TestEvalPanics(t *testing.T) {
	catch := func(fn func()) (v interface{}) {
		defer func() {
			v = recover()
		}()
		fn()
		return nil
	}

	s := eval.Scope(func(v interface{}) changes.Value {
		if v == types.S16(".") {
			return eval.Dot
		}
		panic(string(v.(types.S16)))
	})

	v := catch(func() {
		eval.Eval(s, eval.Parse(nil, "do(x, x = y, y = x)"))
	})
	if v != "recursion detected" {
		t.Error("Unexpected recursion fail", v)
	}

	v = catch(func() {
		eval.Eval(s, eval.Parse(nil, "obj(x = y, y = x).x"))
	})

	if v != "recursion detected" {
		t.Error("Unexpected recursion fail", v)
	}
}
