// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress_test

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/dotchain/dot/stress"
)

func TestStress(t *testing.T) {
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Couldn't delete stress.bolt file", err)
	}

	defer func() {
		if err := os.Remove("stress.bolt"); err != nil {
			log.Println("Couldn't delete stress.bolt file")
		}
	}()

	rand.Seed(time.Now().UTC().UnixNano())

	// rounds, iterations, clients
	stress.Run(nil, 12, 10, 4)
}

func TestStressAndReconnect(t *testing.T) {
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Couldn't delete stress.bolt file", err)
	}

	defer func() {
		if err := os.Remove("stress.bolt"); err != nil {
			log.Println("Couldn't delete stress.bolt file")
		}
	}()

	rand.Seed(time.Now().UTC().UnixNano())

	// rounds, iterations, clients
	states := stress.Run(nil, 12, 10, 4)
	stress.Run(states, 12, 10, 4)
}
