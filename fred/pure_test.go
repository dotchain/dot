// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

// treat fred.Error as a vehicle for storing strings
type concatErrors struct{}

func (c concatErrors) Eval(e fred.Env, args *fred.Vals) fred.Val {
	result := ""
	for _, v := range *args {
		result += string(v.(fred.Error))
	}
	return fred.Error(result)
}

func TestPureEval(t *testing.T) {
	p := &fred.Pure{
		Functor: concatErrors{},
		Args: &fred.Defs{
			fred.Fixed(fred.Error("hello ")),
			fred.Fixed(fred.Error("world!")),
		},
	}

	expected := fred.Error("hello world!")
	if x := p.Eval(env); x != expected {
		t.Error("Unexpected eval", x)
	}
}

func TestPureUpdateArgs(t *testing.T) {
	p := &fred.Pure{
		Functor: concatErrors{},
		Args: &fred.Defs{
			fred.Fixed(fred.Error("hello ")),
			fred.Fixed(fred.Error("world!")),
		},
	}

	c := changes.PathChange{
		Path: []interface{}{"Args"},
		Change: changes.Splice{
			Offset: 0,
			Before: &fred.Defs{},
			After:  &fred.Defs{fred.Fixed(fred.Error("OK "))},
		},
	}
	expected := fred.Error("OK hello world!")
	if x := p.Apply(nil, c).(*fred.Pure).Eval(env); x != expected {
		t.Error("Unexpected eval", x)
	}
}
