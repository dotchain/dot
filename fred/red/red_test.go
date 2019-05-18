// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package red implements a parser for fred
package red_test

import (
	"testing"

	"github.com/dotchain/dot/fred"
	"github.com/dotchain/dot/fred/red"
)

func TestRed(t *testing.T) {
	cases := map[string]string{
		"":                                      "<nil>",
		"1":                                     "1",
		"1 + 2":                                 "3",
		" + 2":                                  "2",
		" - 3  ":                                "-3",
		"2*2 + 3/-1":                            "1",
		"2 * (3 + 4)  ":                         "14",
		` "Hello" `:                             `Hello`,
		`1 +"Hello".length + 2`:                 "8",
		`"Hello".concat(" World")`:              "Hello World",
		`"Hello".splice(0, 2, "Je")`:            "Jello",
		`"Hello".splice(2-2, 1+1, "Je").length`: "5",
		"2 > 3 > 4":                             "false",
		"3 < 5 | booya":                         "true",

		"x+100": "101",

		"  +":           "err: 3: incomplete",
		"  * 5":         "err: 2: unexpected op",
		"3 < 5 & booya": "err: ref: no such ref",
		"3 > 5 | booya": "err: ref: no such ref",
	}

	r := func(key interface{}) fred.Def {
		if key == fred.Text("x") {
			return fred.Fixed(fred.Num("1"))
		}
		return nil
	}
	env := &fred.Environ{Resolver: r, Cacher: &fred.Cache{}, Depth: 5}
	for name, expected := range cases {
		t.Run(name, func(t *testing.T) {
			got := red.Parse(name).Eval(env).Text()
			if got != expected {
				t.Error("Unexpected", got)
			}
		})
	}
}
