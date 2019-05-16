// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func TestCallUpdateFunc(t *testing.T) {
	p := fred.Call(
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("hoo")),
	)
	c := changes.PathChange{
		Path: []interface{}{"Func"},
		Change: changes.Replace{
			Before: fred.Fixed(fred.Error("boo")),
			After:  fred.Fixed(fred.Error("goo")),
		},
	}
	expected := fred.Call(
		fred.Fixed(fred.Error("goo")),
		fred.Fixed(fred.Error("hoo")),
	)

	if x := p.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Errorf("Unexpected eval %#v %#v\n", x, expected)
	}
}

func TestCallUpdateArgs(t *testing.T) {
	p := fred.Call(
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("hoo")),
	)
	c := changes.PathChange{
		Path: []interface{}{"Args"},
		Change: changes.Splice{
			Offset: 0,
			Before: &fred.Defs{},
			After:  &fred.Defs{fred.Fixed(fred.Error("OK "))},
		},
	}

	expected := fred.Call(
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("OK ")),
		fred.Fixed(fred.Error("hoo")),
	)

	if x := p.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected eval", x)
	}
}

func TestCallUncallable(t *testing.T) {
	x := fred.Call(fred.Fixed(fred.Text("boo"))).Eval(env)
	if x != fred.ErrNotCallable {
		t.Error("Unexpected calling strings", x)
	}
}
