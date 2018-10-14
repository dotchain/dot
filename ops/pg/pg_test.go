// +build integration
// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package pg_test

import (
	"context"
	"database/sql"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/pg"
	"github.com/dotchain/dot/x/nw"
	"reflect"
	"testing"
)

var sourceName = "user=postgres dbname=dot_test sslmode=disable"

func dropTable() {
	db, _ := sql.Open("postgres", sourceName)
	db.Exec("DROP TABLE operations")
	db.Close()
}

func TestSetup(t *testing.T) {
	defer dropTable()
	s := &pg.Store{DataSourceName: sourceName, ID: []byte("hello")}
	if err := s.Setup(); err != nil {
		t.Error("setup failure", err)
	}

	s = &pg.Store{DataSourceName: "user=postgres dbname=none sslmode=disable", ID: []byte("hello")}
	if err := s.Setup(); err == nil {
		t.Error("setup failure", err)
	}
}

func TestSimple(t *testing.T) {
	defer dropTable()
	s := &pg.Store{DataSourceName: sourceName, ID: []byte("hello")}
	s.Setup()
	defer s.Close()

	c := changes.PathChange{[]interface{}{5}, changes.Move{2, 2, 2}}
	op1 := ops.Operation{OpID: "one", Change: c}
	op2 := ops.Operation{OpID: "two", Change: c}
	opx := []ops.Op{op1, op2, op1, op2}
	if err := s.Append(context.Background(), opx); err != nil {
		t.Fatal("Append fail", err)
	}

	result, err := s.GetSince(context.Background(), 0, 100)
	if err != nil {
		t.Fatal("GetSince fail", err)
	}

	expected := []ops.Op{op1.WithVersion(0), op2.WithVersion(1)}
	if !reflect.DeepEqual(result, expected) {
		t.Error("result did not match", result, opx)
	}
}

func TestPoll(t *testing.T) {
	defer dropTable()
	s := &pg.Store{DataSourceName: sourceName, ID: []byte("hello")}
	s.Setup()
	defer s.Close()

	if err := s.Poll(context.Background(), 0); err != nil {
		t.Error("unexpected error", err)
	}
}

func TestErrors(t *testing.T) {
	s := &pg.Store{DataSourceName: "\u2312", ID: []byte("hello")}
	if err := s.Append(context.Background(), nil); err == nil {
		t.Error("succeeded with invalid db", err)
	}

	if _, err := s.GetSince(context.Background(), 0, 5); err == nil {
		t.Error("succeeded with invalid db", err)
	}

	defer dropTable()
	s = &pg.Store{DataSourceName: sourceName, ID: []byte("hello")}
	s.Codec = nw.DefaultCodecs["application/x-gob"]
	s.Setup()
	defer s.Close()

	// encode error due to myChange
	opx := []ops.Op{ops.Operation{OpID: "ok", Change: myChange{nil}}}
	if err := s.Append(context.Background(), opx); err == nil {
		t.Fatal("Append fail", err)
	}

	if _, err := s.GetSince(context.Background(), 0, -5); err == nil {
		t.Error("succeeded with invalid limit", err)
	}

	// write a fake op with bad data
	db, _ := sql.Open("postgres", s.DataSourceName)
	db.Exec("INSERT into operations (id, op_id, data) VALUES ($1, $2, $3);",
		s.ID, []byte("opid"), []byte("data"))
	db.Close()

	if _, err := s.GetSince(context.Background(), 0, 100); err == nil {
		t.Error("unexpected success", err)
	}

}

type myChange struct{ fn func() }

func (myChange) Merge(c changes.Change) (cx, ox changes.Change) {
	return nil, nil
}

func (myChange) Revert() changes.Change {
	return nil
}
