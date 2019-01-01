// +build integration
// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package pg_test

import (
	"context"
	"database/sql"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/ops/pg"
	"reflect"
	"testing"
	"time"
)

var sourceName = "user=postgres dbname=dot_test sslmode=disable"

func dropTable() {
	db, _ := sql.Open("postgres", sourceName)
	db.Exec("DROP TABLE operations")
	db.Close()
}

func TestSetup(t *testing.T) {
	defer dropTable()
	if err := pg.Setup(sourceName); err != nil {
		t.Error("setup failure", err)
	}

	name := "user=postgres dbname=none sslmode=disable"
	if err := pg.Setup(name); err == nil {
		t.Error("setup failure", err)
	}
}

func TestSimple(t *testing.T) {
	defer dropTable()
	pg.Setup(sourceName)
	s, err := pg.New(sourceName, "hello", nil)
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
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
	pg.Setup(sourceName)
	s, err := pg.New(sourceName, "hello", nil)
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
	defer s.Close()

	go func() {
		time.Sleep(time.Millisecond * 100)
		s.Append(context.Background(), []ops.Op{ops.Operation{OpID: "ok"}})
	}()

	wait, cancel := context.WithTimeout(context.Background(), time.Second*5)
	before := time.Now()
	if err := s.Poll(wait, 0); err != nil {
		t.Error("unexpected error", err)
	}
	if time.Since(before) > time.Second*2/3 {
		t.Error("Waited too long")
	}
	cancel()

	wait, cancel = context.WithTimeout(context.Background(), time.Second)
	before = time.Now()
	if err := s.Poll(wait, 1); err != wait.Err() {
		t.Error("unexpected error", err)
	}
	if time.Since(before) < time.Second*2/3 {
		t.Error("Waited too little")
	}
	cancel()
}

func TestErrors(t *testing.T) {
	s, err := pg.New("\u2312", "hello", nil)
	if err == nil {
		t.Error("succeeded with invalid db", err)
	}

	defer dropTable()
	pg.Setup(sourceName)
	s, err = pg.New(sourceName, "hello", nw.DefaultCodecs["application/x-gob"])
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
	defer s.Close()

	// empty appends should always succeed
	if err := s.Append(context.Background(), nil); err != nil {
		t.Fatal("Append fail", err)
	}

	// encode error due to myChange
	opx := []ops.Op{ops.Operation{OpID: "ok", Change: myChange{nil}}}
	if err := s.Append(context.Background(), opx); err == nil {
		t.Fatal("Append fail", err)
	}

	if _, err := s.GetSince(context.Background(), 0, -5); err == nil {
		t.Error("succeeded with invalid limit", err)
	}

	// write a fake op with bad data
	db, _ := sql.Open("postgres", sourceName)
	db.Exec("INSERT into operations (id, op_id, data) VALUES ($1, $2, $3);",
		"hello", []byte("opid"), []byte("data"))
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
