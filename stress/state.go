// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress

//go:generate go run codegen.go

import (
	"encoding/gob"
	"math/rand"

	"github.com/dotchain/dot/changes/types"
)

type State struct {
	Text  string
	Count types.Counter
}

func init() {
	gob.Register(State{})
	rand.Seed(42)
}
