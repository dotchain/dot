// Code generated by github.com/tvastar/test/cmd/testmd/testmd.go. DO NOT EDIT.

package dot_test

import (
	"database/sql"
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
	defer os.Remove("file.bolt")

	srv := dot.BoltServer("file.bolt")
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
	stream1.Append(changes.Replace{changes.Nil, types.S8("hello")})

	wg.Wait()

	// Output:
}

func Example_clientServerUsingPostgresDB() {
	sourceName := "user=postgres dbname=dot_test sslmode=disable"
	maxPoll := pg.MaxPoll
	defer func() {
		pg.MaxPoll = maxPoll
		db, _ := sql.Open("postgres", sourceName)
		db.Exec("DROP TABLE operations")
		db.Close()
	}()

	pg.MaxPoll = time.Second
	srv := dot.PostgresServer(sourceName)
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
	stream1.Append(changes.Replace{changes.Nil, types.S8("hello")})

	wg.Wait()

	// Output:
}