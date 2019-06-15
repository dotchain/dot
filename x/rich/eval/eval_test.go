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

func TestArrayMap(t *testing.T) {
	var globals types.M
	var scope eval.Scope

	scope = eval.Scope(func(v interface{}) changes.Value {
		return eval.Eval(scope, globals[v])
	})
	globals = types.M{
		types.S16("+"):    eval.Sum,
		types.S16("."):    eval.Dot,
		types.S16("list"): eval.Parse(scope, "(1, 2, 3)"),
	}

	tests := map[string]string{
		"list.map(value+10)":           "(11, 12, 13)",
		"list.reduce(100, value+last)": "106",
	}

	t.Run("basic", func(t *testing.T) {
		got := eval.Eval(scope, eval.Parse(scope, "list"))
		expected := changes.Value(types.A{
			changes.Atomic{Value: 1},
			changes.Atomic{Value: 2},
			changes.Atomic{Value: 3},
		})

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Got %#v\n", got)
		}
	})

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
