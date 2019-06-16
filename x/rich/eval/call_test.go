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

func TestCall(t *testing.T) {
	cx := &eval.Call{types.A{types.S16("boo")}}
	if x := cx.Apply(nil, nil); !reflect.DeepEqual(cx, x) {
		t.Error("Unexpected nil apply", x)
	}

	c := changes.PathChange{
		Path:   []interface{}{"A"},
		Change: changes.Replace{Before: cx.A, After: types.A{types.S16("hoo")}},
	}
	expected := &eval.Call{types.A{types.S16("hoo")}}
	if x := cx.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected apply", x)
	}
}
