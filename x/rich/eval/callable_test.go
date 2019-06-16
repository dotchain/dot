// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/eval"
)

func TestCallable(t *testing.T) {
	before := eval.Callable(func(s eval.Scope, args []changes.Value) changes.Value {
		return types.S16("before")
	})
	after := eval.Callable(func(s eval.Scope, args []changes.Value) changes.Value {
		return types.S16("after")
	})

	c := changes.ChangeSet{
		changes.Replace{Before: before, After: after},
		nil,
	}

	x := before.Apply(nil, c)
	if fn, ok := x.(eval.Callable); !ok || fn(nil, nil) != types.S16("after") {
		t.Error("Unexpected apply result", x)
	}
}
