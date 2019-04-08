// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package meta_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"sync"

	"github.com/dotchain/dot"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"

	"github.com/dotchain/dot/x/meta"
)

func Example_clientServerUsingBoltDB() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	defer remove("file.bolt")()

	// start server
	url, close := startServer("file.bolt")
	defer close()

	// open stream
	session, updates, metas := dot.Reconnect(url, -1, nil)
	defer session.Close()

	// make a couple of changes
	c1 := changes.Replace{Before: changes.Nil, After: types.S8("hello")}
	c2 := changes.Replace{Before: types.S8("hello"), After: types.S8("hello2")}
	updates.Append(c1).Append(c2)

	// wait till version = 1 (i.e. version 0 => c1, version 1 => c2)
	metastream := &meta.DataStream{Stream: metas}
	wg.Add(1)
	metas.Nextf("key", func() {
		mu.Lock()
		defer mu.Unlock()

		metastream = metastream.Latest()
		if metastream.Value.Version == 1 && len(metastream.Value.Pending) == 0 {
			wg.Done()
		}
	})
	wg.Wait()

	// Validate
	latest := metastream.Latest().Value
	fmt.Println(latest.TransformedOp[0].Changes())
	fmt.Println(latest.MergeOps[0])
	fmt.Println(latest.TransformedOp[1].Changes())
	fmt.Println(latest.MergeOps[1])

	// Output:
	// {{} hello}
	// []
	// {hello hello2}
	// []
}

func remove(fname string) func() {
	if err := os.Remove(fname); err != nil {
		log.Println("Couldnt remove file", fname)
	}
	return func() {
		if err := os.Remove(fname); err != nil {
			log.Println("Couldnt remove file", fname)
		}
	}
}

func startServer(fname string) (url string, close func()) {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	srv := dot.WithLogger(dot.BoltServer("file.bolt"), logger)
	httpSrv := httptest.NewServer(srv)

	return httpSrv.URL, func() {
		dot.CloseServer(srv)
		httpSrv.Close()
	}
}
