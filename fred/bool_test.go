// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func TestBool(t *testing.T) {
	t1 := fred.Bool(false)
	t2 := fred.Bool(false)
	if t1 != t2 {
		t.Error("Unexpected inequality")
	}

	t3 := t1.Apply(nil, changes.ChangeSet{
		changes.Replace{
			Before: t1,
			After:  changes.Nil,
		},
	})
	if t3 != changes.Nil {
		t.Error("Unexpected apply", t3)
	}

	if t1.Apply(nil, nil) != t1 {
		t.Error("Unexpected apply", t3)
	}
}

func TestBoolOperator(t *testing.T) {
	bTrue := fred.Bool(true)
	bFalse := fred.Bool(false)
	wat := fred.Error("wat")

	cases := map[string][]fred.Val{
		"and1": {bTrue, fred.Text("&"), bTrue},
		"and2": {bTrue, fred.Text("&"), bTrue, bTrue},
		"and3": {bFalse, fred.Text("&"), bTrue, bFalse, wat},
		"and4": {bFalse, fred.Text("&"), bFalse, wat},
		"and5": {wat, fred.Text("&"), bTrue, wat},

		"or1": {bTrue, fred.Text("|"), bTrue},
		"or2": {bTrue, fred.Text("|"), bTrue, bFalse},
		"or3": {bTrue, fred.Text("|"), bTrue, wat},
		"or4": {bTrue, fred.Text("|"), bFalse, bTrue, wat},
		"or5": {wat, fred.Text("|"), bFalse, wat},

		"err1": {fred.ErrNoSuchField, fred.Text("?"), bTrue, bTrue},
		"err2": {fred.ErrNotBool, fred.Text("|"), bFalse, fred.Text("boo")},
	}
	for name, list := range cases {
		t.Run(name, func(t *testing.T) {
			expected := list[0]
			op := fred.Fixed(list[1])
			first := fred.Fixed(list[2])
			args := fred.Defs{}
			for _, arg := range list[3:] {
				args = append(args, fred.Fixed(arg))
			}
			expr := fred.Call(fred.Field(first, op), args...)
			got := expr.Eval(env)
			if got != expected {
				t.Error("unexpected", got, list)
			}
		})
	}
}
