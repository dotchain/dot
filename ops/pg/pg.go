// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package pg implements the database storage layer for DOT
//
// TODO: optimize the inefficient implementation of Append
package pg

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/x/nw"
	_ "github.com/lib/pq" // pq registers with sql
	"sync"
)

// Store implements ops.Store
type Store struct {
	DataSourceName string
	ID             []byte
	db             *sql.DB
	lock           sync.Mutex
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
func (s *Store) Setup() error {
	db, err := sql.Open("postgres", s.DataSourceName)
	if err == nil {
		_, err = db.Exec(createTableCommand)
	}

	if err != nil && db != nil {
		must(db.Close())
		db = nil
	}

	return err
}

func (s *Store) getDB() (*sql.DB, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var err error
	if s.db == nil {
		s.db, err = sql.Open("postgres", s.DataSourceName)
		if s.db != nil {
			err = s.db.Ping()
		}
		if err != nil {
			s.Close()
		}
	}
	return s.db, err
}

// Close releases any allocated DB resources.  It is not safe to call
// Close when other calls may be in progress.
func (s *Store) Close() {
	if s.db != nil {
		must(s.db.Close())
		s.db = nil
	}
}

// Append implements Store.Append.  It uses gob/encoding to serialize
// the provided operation.
func (s *Store) Append(ctx context.Context, ops []ops.Op) error {
	db, err := s.getDB()
	if err != nil {
		return err
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
		args = append(args, s.ID, id, data)
	}

	cmd += " ON CONFLICT ON CONSTRAINT unique_op_id DO NOTHING;"

	_, err = db.Exec(cmd+";", args...)
	return err
}

// GetSince implements Store.GetSince
func (s *Store) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	db, err := s.getDB()
	if err != nil {
		return nil, err
	}

	cmd := fetchCommand + fmt.Sprintf("LIMIT %d OFFSET %d;", limit, version)
	rows, err := db.Query(cmd, s.ID)
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

type opdata struct {
	ops.Op
}

func (s *Store) encode(op ops.Op) ([]byte, []byte, error) {
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

func (s *Store) decode(data []byte) (ops.Op, error) {
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

// Poll is not implemented
func (s *Store) Poll(ctx context.Context, version int) error {
	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
