// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package pg implements the dot storage for postgres 9.5+
//
// A http server can be implemented like so:
//      import "github.com/dotchain/dot/ops/pg"
//      import "github.com/dotchain/dot/ops/nw"
//      dataSource := "dbname=mydb user=xyz"
//      store, _ := sql.New(dataSource, "instance", nil)
//      defer  store.Close()
//      handler := &nw.Handler{Store: store}
//      h := func(w http.ResponseWriter, req  *http.Request) {
//              // Enable CORS
//              w.Header().Set("Access-Control-Allow-Origin", "*")
//              w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//              if req.Method == "OPTIONS" {
//                    return
//              }
//              handler.ServeHTTP(w, req)
//      }
//      http.HandleFunc("/api/", h)
//      http.ListenAndServe()
//
package pg

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/lib/pq"
)

// New returns a store connected to the provided data stource
func New(dataSourceName, id string, codec nw.Codec) (ops.Store, error) {
	s := &store{id: id, Codec: codec}
	if err := s.init(dataSourceName); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *store) init(name string) error {
	var err error
	delay := time.Second * 30

	s.db, err = sql.Open("postgres", name)
	if err == nil {
		err = s.db.Ping()
	}
	if err == nil {
		s.l = pq.NewListener(name, delay, delay, nil)
		err = s.l.Listen(s.id)
	}

	if err != nil {
		s.Close()
	} else {
		ch := s.l.NotificationChannel()
		go func() {
			for range ch {
				s.broadcast()
			}
		}()
	}
	return err
}

type store struct {
	id      string
	db      *sql.DB
	l       *pq.Listener
	waiters []chan struct{}
	lock    sync.Mutex
	nw.Codec
}

var createTableCommand = `
CREATE TABLE IF NOT EXISTS operations (
	id BYTEA NOT NULL,
        seq BIGSERIAL,
	op_id BYTEA NOT NULL,
	data BYTEA,
        PRIMARY KEY (id, seq),
        CONSTRAINT unique_op_id UNIQUE (id, op_id)
);
`

var fetchCommand = `
SELECT data from operations
WHERE id = $1
ORDER BY seq
`

// Setup creates the tables and indices
func Setup(dataSourceName string) error {
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		_, err = db.Exec(createTableCommand)
	}

	if err != nil && db != nil {
		must(db.Close())
		db = nil
	}

	return err
}

// Close releases any allocated DB resources.  It is not safe to call
// Close when other calls may be in progress.
func (s *store) Close() {
	if s.db != nil {
		must(s.db.Close())
		s.db = nil
	}
	if s.l != nil {
		must(s.l.Close())
		s.l = nil
	}
}

// Append implements store.Append.  It uses gob/encoding to serialize
// the provided operation.
func (s *store) Append(ctx context.Context, ops []ops.Op) error {
	if len(ops) == 0 {
		return nil
	}

	cmd := "INSERT into operations (id, op_id, data) VALUES"
	args := []interface{}{}

	for kk, op := range ops {
		id, data, err := s.encode(op)
		if err != nil {
			return err
		}

		l := len(args)
		if kk > 0 {
			cmd = cmd + ", "
		}
		cmd = cmd + fmt.Sprintf("($%d, $%d, $%d)", l+1, l+2, l+3)
		args = append(args, []byte(s.id), id, data)
	}

	cmd += " ON CONFLICT ON CONSTRAINT unique_op_id DO NOTHING;"

	_, err := s.db.ExecContext(ctx, cmd+";", args...)
	if err == nil {
		log.Println("Notifying", s.id)
		_, err = s.db.ExecContext(ctx, "NOTIFY "+s.id+";")
	}
	return err
}

// GetSince implements store.GetSince
func (s *store) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	cmd := fetchCommand + fmt.Sprintf("LIMIT %d OFFSET %d;", limit, version)
	rows, err := s.db.QueryContext(ctx, cmd, []byte(s.id))
	if err != nil {
		return nil, err
	}
	result := []ops.Op{}
	for rows.Next() {
		var data []byte
		var op ops.Op
		err := rows.Scan(&data)
		if err == nil {
			op, err = s.decode(data)
		}
		if err != nil {
			return nil, err
		}
		op = op.WithVersion(version)
		version++
		result = append(result, op)
	}

	return result, nil
}

// Poll uses postgres LISTEN to wait get notified on changes.
func (s *store) Poll(ctx context.Context, version int) error {
	ch := make(chan struct{}, 1)
	s.lock.Lock()
	s.waiters = append(s.waiters, ch)
	s.lock.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
	case <-time.After(time.Minute): // max wait is 1min
	}
	return nil
}

type opdata struct {
	ops.Op
}

func (s *store) encode(op ops.Op) ([]byte, []byte, error) {
	codec := nw.DefaultCodecs["application/x-gob"]
	if s.Codec != nil {
		codec = s.Codec
	}
	var data, id bytes.Buffer
	err := codec.Encode(opdata{op}, &data)
	if err == nil {
		err = codec.Encode(op.ID(), &id)
	}
	if err != nil {
		return nil, nil, err
	}

	return id.Bytes(), data.Bytes(), nil
}

func (s *store) decode(data []byte) (ops.Op, error) {
	codec := nw.DefaultCodecs["application/x-gob"]
	if s.Codec != nil {
		codec = s.Codec
	}
	var opd opdata
	if err := codec.Decode(&opd, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return opd.Op, nil
}

func (s *store) broadcast() {
	s.lock.Lock()
	waiters := s.waiters
	s.waiters = nil
	s.lock.Unlock()
	for _, ch := range waiters {
		ch <- struct{}{}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
