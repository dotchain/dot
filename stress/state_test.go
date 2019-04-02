// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress_test

//go:generate go run codegen.go

import (
	"encoding/gob"
	"log"
	"math/rand"
	"sync"

	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

type State struct {
	Text  string
	Count types.Counter
}

func rounds(clients []*StateStream, iterations, rounds int) {
	var wg sync.WaitGroup
	var waitCount int

	for jj := 0; jj < rounds; jj++ {
		wg.Add(len(clients))
		waitCount += len(clients)
		for kk, client := range clients {
			round(kk, client, iterations, waitCount, &wg)
			clients[kk] = client.Latest()
		}
		wg.Wait()
		log.Println("finished round", jj)
	}
}

func round(idx int, s *StateStream, iterations, waitCount int, wait *sync.WaitGroup) {
	go func() {
		modifyTextRandomly(idx, s.Text(), iterations)
		s = s.Latest()
		counter := s.Count()
		counter.Increment(1)

		finished := false
		counter.Stream.Nextf(idx, func() {
			s = s.Latest()
			if !finished && int32(s.Value.Count) >= int32(waitCount) {
				finished = true
				counter.Stream.Nextf(idx, nil)
				wait.Done()
			}
		})
	}()
}

func modifyTextRandomly(idx int, s *streams.S8, iterations int) {
	for kk := 0; kk < iterations; kk++ {
		s = s.Latest()
		l := len(s.Value)
		insert := randString(3)

		if l == 0 {
			s = s.Splice(0, 0, insert)
		} else {
			var offset, count int
			if l > 0 {
				offset = rand.Intn(l)
			}
			if l-offset > 0 {
				count = rand.Intn(l - offset)
			}
			s = s.Splice(offset, count, insert)
		}
	}
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	gob.Register(State{})
	rand.Seed(42)
}
