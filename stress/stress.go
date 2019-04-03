// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

// Package stress has the tools for building a stress test
package stress

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/dotchain/dot"
)

// Starts the stress server at the provided address.
//
// It uses /stress/ as the path to serve.
// The returned function can be used to shutdown the server
func StartServer(addr string) func() {
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Couldn't delete stress.bolt file", err)
	}

	bolt := dot.BoltServer("stress.bolt")
	http.Handle("/stress/", bolt)

	srv := &http.Server{Addr: addr}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	close := func() {
		//srv.Shutdown(context.Background())
		srv.Close()
		dot.CloseServer(bolt)
		if err := os.Remove("stress.bolt"); err != nil {
			log.Println("Couldn't delete stress.bolt file")
		}
	}

	return close
}

// Run runs the required number of rounds of test
func Run(rounds, iterations, clients int) {
	defer StartServer(":8083")()
	log.Println("Started server...")
	sessions := make([]*Session, clients)
	url := "http://localhost:8083/stress/"

	var wg sync.WaitGroup
	for kk := range sessions {
		sessions[kk] = NewSession(url, clients, &wg)
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

	for kk := range sessions {
		sessions[kk].Close()
	}
}
