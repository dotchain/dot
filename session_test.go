// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"sync"
	"time"

	"github.com/dotchain/dot"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/ops/pg"
)

func Example_clientServerUsingBoltDB() {
	defer remove("file.bolt")()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	srv := dot.WithLogger(dot.BoltServer("file.bolt"), logger)
	defer dot.CloseServer(srv)
	httpSrv := httptest.NewServer(srv)
	defer httpSrv.Close()

	session1, stream1 := dot.Connect(httpSrv.URL)
	session2, stream2 := dot.Connect(httpSrv.URL)

	defer session1.Close()
	defer session2.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	stream2.Nextf("key", wg.Done)
	stream1.Append(changes.Replace{Before: changes.Nil, After: types.S8("hello")})

	wg.Wait()

	// Output:
}

func Example_clientServerUsingPostgresDB() {
	sourceName := "user=postgres dbname=dot_test sslmode=disable"
	maxPoll := pg.MaxPoll
	defer func() {
		pg.MaxPoll = maxPoll
		db, err := sql.Open("postgres", sourceName)
		must(err)
		_, err = db.Exec("DROP TABLE operations")
		must(err)
		must(db.Close())
	}()

	pg.MaxPoll = time.Second
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	srv := dot.WithLogger(dot.PostgresServer(sourceName), logger)
	defer dot.CloseServer(srv)
	httpSrv := httptest.NewServer(srv)
	defer httpSrv.Close()

	session1, stream1 := dot.Connect(httpSrv.URL)
	session2, stream2 := dot.Connect(httpSrv.URL)

	defer session1.Close()
	defer session2.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	stream2.Nextf("key", wg.Done)
	stream1.Append(changes.Replace{Before: changes.Nil, After: types.S8("hello")})

	wg.Wait()

	// Output:
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
