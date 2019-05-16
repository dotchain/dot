// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func TestNum(t *testing.T) {
	t1 := fred.Num("1")
	t2 := fred.Num("1")
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

func TestNumArithmetic(t *testing.T) {
	cases := map[string][]string{
		"5":    {"+", "1", "3", "2", "-1"},
		"0":    {"-", "0"},
		"-1/3": {"-", "2/3", "1"},
		"1":    {"*", "1/2", "2/1"},
		"-1/6": {"/", "2", "3", "4", "-1"},
	}
	for name, list := range cases {
		t.Run(name, func(t *testing.T) {
			expected := fred.Num(name)
			op := fred.Fixed(fred.Text(list[0]))
			first := fred.Fixed(fred.Num(list[1]))
			args := fred.Defs{}
			for _, arg := range list[2:] {
				args = append(args, fred.Fixed(fred.Num(arg)))
			}
			expr := fred.Call(fred.Field(first, op), args...)
			got := expr.Eval(env)
			if got != expected {
				t.Error("unexpected", got)
			}
		})
	}
}

func TestNumErrors(t *testing.T) {
	n := fred.Fixed(fred.Num("5"))
	z := fred.Field(n, fred.Fixed(fred.Text("?")))
	if x := z.Eval(env); x != fred.ErrNoSuchField {
		t.Error("Unexpected", x)
	}

	sum := fred.Field(n, fred.Fixed(fred.Text("+")))
	x := fred.Call(sum, fred.Fixed(fred.Text("boo"))).Eval(env)
	if x != fred.ErrNotNumber {
		t.Error("Unexpected", x)
	}

	x = fred.Call(sum, fred.Fixed(fred.Num("boo"))).Eval(env)
	if x != fred.Error("math/big: cannot unmarshal \"boo\" into a *big.Rat") {
		t.Error("Unexpected", x)
	}

	x = fred.Call(sum, fred.Fixed(fred.Error("boo"))).Eval(env)
	if x != fred.Error("boo") {
		t.Error("Unexpected", x)
	}

	n = fred.Fixed(fred.Num("z"))
	z = fred.Field(n, fred.Fixed(fred.Text("+")))
	x = fred.Call(z, fred.Fixed(fred.Num("0"))).Eval(env)
	if x != fred.Error("math/big: cannot unmarshal \"z\" into a *big.Rat") {
		t.Error("Unexpected", x)
	}

	x = fred.Call(fred.Field(fred.Fixed(fred.Num("9")), fred.Fixed(fred.Text("+")))).Eval(env)
	if x != fred.Num("9") {
		t.Error("Unexpected", x)
	}
}
