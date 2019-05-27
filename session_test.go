// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"database/sql"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
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

	stream1, store1 := dot.NewSession().NonBlockingStream(httpSrv.URL, nil)
	stream2, store2 := dot.NewSession().Stream(httpSrv.URL, nil)

	defer store1.Close()
	defer store2.Close()

	stream1.Append(changes.Replace{Before: changes.Nil, After: types.S8("hello")})
	fmt.Println("push", stream1.Push())
	fmt.Println("pull", stream2.Pull())

	// Output:
	// push <nil>
	// pull <nil>
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

	stream1, store1 := dot.NewSession().Stream(httpSrv.URL, logger)
	stream2, store2 := dot.NewSession().Stream(httpSrv.URL, logger)

	defer store1.Close()
	defer store2.Close()

	stream1.Append(changes.Replace{Before: changes.Nil, After: types.S8("hello")})
	fmt.Println("push", stream1.Push())
	fmt.Println("pull", stream2.Pull())

	// Output:
	// push <nil>
	// pull <nil>
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
