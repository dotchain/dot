// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/fred"
)

func TestIf(t *testing.T) {
	p := fred.If(
		fred.Fixed(fred.Bool(true)),
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("hoo")),
	)
	if x := p.Eval(env); x != fred.Error("boo") {
		t.Error("Unexpected if", x)
	}

	p = fred.If(
		fred.Fixed(fred.Bool(false)),
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("hoo")),
	)
	if x := p.Eval(env); x != fred.Error("hoo") {
		t.Error("Unexpected if", x)
	}

	p = fred.If(
		fred.Fixed(fred.Text("ok")),
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("hoo")),
	)
	if x := p.Eval(env); x != fred.ErrInvalidCondition {
		t.Error("Unexpected if", x)
	}

	p = fred.If(
		fred.Fixed(fred.Error("ok")),
		fred.Fixed(fred.Error("boo")),
		fred.Fixed(fred.Error("hoo")),
	)
	if x := p.Eval(env); x != fred.Error("ok") {
		t.Error("Unexpected if", x)
	}
}
