// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress

//go:generate go run codegen.go

import (
	"math/rand"

	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/ops/nw"
)

type State struct {
	Text  string
	Count types.Counter
}

func init() {
	nw.Register(State{})
	rand.Seed(42)
}
