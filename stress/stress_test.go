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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dotchain/dot/stress"
)

var rounds = flag.Int("rounds", 10, "number of rounds")
var iterations = flag.Int("iterations", 20, "number of iterations per round")
var clients = flag.Int("clients", 2, "number of clients per round")

func TestSimple(t *testing.T) {
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Couldn't delete stress.bolt file", err)
	}

	defer func() {
		if err := os.Remove("stress.bolt"); err != nil {
			log.Println("Couldn't delete stress.bolt file")
		}
	}()

	defer stress.StartServer(":8083")()
	rand.Seed(time.Now().UTC().UnixNano())

	stress.Run(nil, *rounds, *iterations, *clients)
}

func TestServerRestarts(t *testing.T) {
	var done int32
	var wg sync.WaitGroup

	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Couldn't delete stress.bolt file", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		first := true
		for atomic.LoadInt32(&done) != 1 {
			close := stress.StartServer(":8083")
			log.Println("Started server")
			if first {
				time.Sleep(time.Second * 10)
				first = false
			} else {
				// must give atleeast 30s of server up time to let
				// the last call succeed
				time.Sleep(time.Second * 60)
			}
			close()
			log.Println("Stopped server")
		}

		if err := os.Remove("stress.bolt"); err != nil {
			log.Println("Couldn't delete stress.bolt file")
		}
	}()

	rand.Seed(time.Now().UTC().UnixNano())
	stress.Run(nil, *rounds, *iterations, *clients)
	atomic.StoreInt32(&done, 1)
	wg.Wait()
}

func TestReconnect(t *testing.T) {
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Couldn't delete stress.bolt file", err)
	}

	defer func() {
		fi, err := os.Stat("stress.bolt")
		if err == nil {
			log.Println("DB Size ", fi.Size(), "bytes")
		}
		if err := os.Remove("stress.bolt"); err != nil {
			log.Println("Couldn't delete stress.bolt file")
		}
	}()

	defer stress.StartServer(":8083")()
	rand.Seed(time.Now().UTC().UnixNano())

	versions := [5]int{0, 0, 0, 0, 0}
	states := stress.Run(nil, 1, 4, 5)
	for kk := 0; kk < *rounds; kk++ {
		for nn := 0; nn < 5; nn++ {
			versions[nn] = states[nn].Version
		}
		log.Println("Version", versions)
		states = stress.Run(states, 1, 4, 5)
	}
}
