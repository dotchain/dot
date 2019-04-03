// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/dotchain/dot/stress"
)

func TestStress(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	// rounds, iterations, clients
	stress.Run(12, 10, 4)
}
