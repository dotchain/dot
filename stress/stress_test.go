// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress_test

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/dotchain/dot/stress"
)

var rounds = flag.Int("rounds", 10, "number of rounds")
var iterations = flag.Int("iterations", 20, "number of iterations per round")
var clients = flag.Int("clients", 2, "number of clients per round")

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

	stress.Run(nil, *rounds, *iterations, *clients)
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

	states := stress.Run(nil, *rounds, *iterations, *clients)
	stress.Run(states, *rounds, *iterations, *clients)
}
