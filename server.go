// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build !js

package dot

import (
	"net/http"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/bolt"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/ops/pg"
)

// BoltServer returns a http.Handler serving DOT requests backed by the db
func BoltServer(fileName string) http.Handler {
	store, err := bolt.New(fileName, "dot_root", nil)
	must(err)
	store = ops.Polled(store)
	return &nw.Handler{Store: store}
}

// PostgresServer returns a http.Handler serving DOT requests backed by the db
func PostgresServer(sourceName string) http.Handler {
	must(pg.Setup(sourceName))
	store, err := pg.New(sourceName, "dot_root", nil)
	must(err)
	return &nw.Handler{Store: store}
}

// WithLogger updates the logger for server
func WithLogger(h http.Handler, l log.Log) http.Handler {
	h.(*nw.Handler).Log = l
	return h
}

// CloseServer closes the http.Handler returned by this package
func CloseServer(h http.Handler) {
	h.(*nw.Handler).Store.Close()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
