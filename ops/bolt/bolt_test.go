// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package bolt_test

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/bolt"
	"github.com/etcd-io/bbolt"
)

var fname = "bolt.data"

func TestInvalidFile(t *testing.T) {
	_, err := bolt.New(".", "hello", nil)
	if err == nil {
		t.Fatal("Unexpected invalid file success")
	}
}

func TestEmpty(t *testing.T) {
	defer os.Remove(fname)
	s, err := bolt.New(fname, "hello", nil)
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
	defer s.Close()

	if err := s.Append(context.Background(), nil); err != nil {
		t.Fatal("EmptyAppend", err)
	}

	ops, err := s.GetSince(context.Background(), 0, 100)
	if err != nil || len(ops) > 0 {
		t.Error("Unexpected GetSince response", ops, err)
	}

	err = s.Poll(context.Background(), 0)
	if err != nil {
		t.Error("Unexpected GetSince response", err)
	}
}

func TestSimple(t *testing.T) {
	defer os.Remove(fname)
	s, err := bolt.New(fname, "hello", nil)
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

	if err := s.Append(context.Background(), opx); err != nil {
		t.Fatal("Append fail", err)
	}

	result, err := s.GetSince(context.Background(), 0, 100)
	if err != nil {
		t.Fatal("GetSince fail", err)
	}

	expected := []ops.Op{op1.WithVersion(0), op2.WithVersion(1)}
	if !reflect.DeepEqual(result, expected) {
		t.Error("result did not match", result, expected)
	}

	result, err = s.GetSince(context.Background(), 0, 1)
	if err != nil {
		t.Fatal("GetSince fail", err)
	}

	if !reflect.DeepEqual(result, expected[:1]) {
		t.Error("result did not match", result, expected)
	}

	result, err = s.GetSince(context.Background(), 1, 1)
	if err != nil {
		t.Fatal("GetSince fail", err)
	}

	if !reflect.DeepEqual(result, expected[1:]) {
		t.Error("result did not match", result, expected)
	}
}

func TestEncodeError(t *testing.T) {
	defer os.Remove(fname)
	s, err := bolt.New(fname, "hello", nil)
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
	defer s.Close()

	op1 := ops.Operation{OpID: "one"}
	op2 := ops.Operation{OpID: "two", Change: myChange{}}
	opx := []ops.Op{op1, op2}
	if err := s.Append(context.Background(), opx); err == nil {
		t.Fatal("Append fail", err)
	}

	result, err := s.GetSince(context.Background(), 0, 100)
	if err != nil || len(result) > 0 {
		t.Fatal("GetSince fail", err, result)
	}

}

func TestDecodeError(t *testing.T) {
	defer os.Remove(fname)
	db, err := bbolt.Open(fname, 0666, nil)
	if err == nil {
		err = db.Update(func(tx *bbolt.Tx) error {
			root, err := tx.CreateBucketIfNotExists([]byte("hello"))
			if err == nil {
				root.NextSequence()
				root.Put([]byte("one"), []byte{0})
				root.Put([]byte(strconv.FormatUint(0, 16)), []byte{1, 2, 3})
			}
			return err
		})
	}
	if err != nil {
		t.Fatal("Unable to open bolt file", err)
	}
	db.Close()

	s, err := bolt.New(fname, "hello", nil)
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
	defer s.Close()

	result, err := s.GetSince(context.Background(), 0, 100)
	if err == nil || len(result) > 0 {
		t.Fatal("GetSince fail", err, result)
	}

}

type myChange struct{}

func (myChange) Merge(c changes.Change) (cx, ox changes.Change) {
	return nil, nil
}

func (myChange) Revert() changes.Change {
	return nil
}
