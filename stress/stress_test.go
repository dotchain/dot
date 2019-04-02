// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress_test

//go:generate go run codegen.go

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dotchain/dot"
	"github.com/dotchain/dot/changes/types"
)

func Server() {
	// uses a local-file backed bolt DB backend
	http.Handle("/stress/", dot.BoltServer("stress.bolt"))
	http.ListenAndServe(":8083", nil)
}

func Client() {
	streams := make([]*StateStream, 4)
	sessions := make([]*dot.Session, 4)
	initial := State{Text: "hello world", Count: types.Counter(0)}
	for kk := range streams {
		session, str := dot.Connect("http://localhost:8083/stress/")
		sessions[kk] = session
		streams[kk] = &StateStream{Stream: str, Value: initial}
	}
	rounds(streams, 5, 100)
}

func TestStress(t *testing.T) {
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Didnt remove stress.bolt file", err)
	}

	go Server()
	time.Sleep(time.Second * 2)
	Client()
	if err := os.Remove("stress.bolt"); err != nil {
		log.Println("Didnt remove stress.bolt file", err)
	}
}
