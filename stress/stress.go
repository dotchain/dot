// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

// Package stress has the tools for building a stress test
package stress

import (
	"log"
	"net/http"
	"sync"

	"github.com/dotchain/dot"
)

// Starts the stress server at the provided address.
//
// It uses /stress/ as the path to serve.
// The returned function can be used to shutdown the server
func StartServer(addr string) func() {
	bolt := dot.BoltServer("stress.bolt")
	srv := &http.Server{Addr: ":8083", Handler: bolt}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	close := func() {
		//srv.Shutdown(context.Background())
		srv.Close()
		dot.CloseServer(bolt)
	}

	return close
}

// Run runs the required number of rounds of test
func Run(oldStates []SessionState, rounds, iterations, clients int) []SessionState {
	sessions := make([]*Session, clients)
	url := "http://localhost:8083/"

	var wg sync.WaitGroup
	for kk := range sessions {
		if oldStates != nil {
			sessions[kk] = oldStates[kk].Reconnect(url, clients, &wg)
		} else {
			sessions[kk] = NewSession(url, clients, &wg)
		}
	}

	for rr := 0; rr < rounds; rr++ {
		wg.Add(clients)
		for kk := range sessions {
			sessions[kk].MakeSomeRandomChanges(iterations)
		}
		log.Println("Waiting for round", rr+1)
		wg.Wait()
		log.Println("Finished round", rr+1)
	}

	states := make([]SessionState, len(sessions))
	for kk := range sessions {
		states[kk] = sessions[kk].Close()
	}
	return states
}
