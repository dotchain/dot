// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/ops"
)

type S = types.S8

func TestTransformerBasic(t *testing.T) {
	store := MemStore(nil)
	xformed := ops.Transformed(store)
	results, err := xformed.GetSince(context.Background(), 0, 100)
	if err != nil || len(results) > 0 {
		t.Error("Unexpected results", err, results)
	}

	first := ops.Operation{"ID1", nil, 100, -1, nil}
	second := ops.Operation{"ID2", nil, 100, -1, nil}
	if err := xformed.Append(context.Background(), []ops.Op{first, second}); err != nil {
		t.Fatal("Unexpected append error", err)
	}

	results, err = xformed.GetSince(context.Background(), 0, 1)
	if err != nil || len(results) != 1 {
		t.Fatal("Unexpected getSince", err, results)
	}
	first.VerID = 0
	second.VerID = 1
	if results[0] != first {
		t.Fatal("Unexpected first value", results)
	}

	results, err = xformed.GetSince(context.Background(), 0, 100)
	if err != nil || len(results) != 2 {
		t.Fatal("Unexpected getSince", err, results)
	}

	if results[0] != first || results[1] != second {
		t.Fatal("Unexpected first value", results)
	}
}

func TestTransformerBranched(t *testing.T) {
	store := MemStore(nil)
	xformed := ops.Transformed(store)

	initial, final, items := getBranchedOps()
	if err := xformed.Append(context.Background(), items); err != nil {
		t.Error("Unexpected Append", err)
	}

	x, err := xformed.GetSince(context.Background(), 0, 100)
	if err != nil || len(x) != len(items) {
		t.Fatal("Unexpected results", x, err)
	}

	value := initial
	for kk, op := range x {
		item := items[kk].WithVersion(kk)
		if !reflect.DeepEqual(item.WithChanges(op.Changes()), op) {
			t.Fatal("Unexpected ID and other info", kk, op)
		}
		value = value.Apply(nil, op.Changes())
	}

	if !reflect.DeepEqual(value, final) {
		t.Fatal("Unexpected value", value)
	}
}

func TestTransformedBranchedErrors(t *testing.T) {
	n := 1

	_, _, items := getBranchedOps()
	store := MemStore(items)
	xformed := ops.Transformed(store)
	expectedItems, _ := xformed.GetSince(context.Background(), 0, 100)
	lastOp := expectedItems[len(expectedItems)-1]

	for {
		store = MemStore(items)
		f := &failNthGetSince{n, errors.New("My error"), store}
		xformed = ops.Transformed(f)
		n++

		x, err := xformed.GetSince(context.Background(), len(items)-1, 100)
		if err == nil {
			if (f.n < 0 || len(x) != 1) || !reflect.DeepEqual(x[0], lastOp) {
				t.Fatal("Unexpected success", n, f.n, x)
			}
			break
		}

		if err != f.err || len(x) != 0 {
			t.Fatal("Unexpected results", n, x, err)
		}
	}
}

func TestTransformedOpByOp(t *testing.T) {
	_, _, items := getBranchedOps()
	store := MemStore(items)
	xformed := ops.Transformed(store)
	expectedItems, _ := xformed.GetSince(context.Background(), 0, 100)
	for kk, op := range expectedItems {
		xformed = ops.Transformed(MemStore(items))
		x, err := xformed.GetSince(context.Background(), kk, 1)
		if err != nil || !reflect.DeepEqual(x, []ops.Op{op}) {
			t.Error("Unexpected behavior", kk, x, err, op)
		}
	}
}

func TestTransformedOpByOpWithCache(t *testing.T) {
	_, _, items := getBranchedOps()
	store := MemStore(items)
	xformed := ops.Transformed(store)
	expectedItems, _ := xformed.GetSince(context.Background(), 0, 100)
	cache := &sync.Map{}
	for kk, op := range expectedItems {
		xformed = ops.TransformedWithCache(MemStore(items), cache)
		x, err := xformed.GetSince(context.Background(), kk, 1)
		if err != nil || !reflect.DeepEqual(x, []ops.Op{op}) {
			t.Error("Unexpected behavior with cache", kk, x, err, op)
		}
	}
}

func getBranchedOps() (initial, final changes.Value, xformed []ops.Op) {
	// The following is the sequence of changes:
	// first = Replace empty => "Hello World"
	// client 1 => first + insert "A " and "B " => "A B Hello World"
	// client 2 => first + insert "X " and "Y " => "X Y Hello World"
	// client 1 => + merge (insert "X") + insert "C " => "A B X C Hello World"
	// client 2 => + merge (insert "A") + insert "Z " => "A X Y Z Hello World"

	first := changes.Replace{Before: changes.Nil, After: S("Hello World")}
	c1_1 := changes.Splice{Offset: 0, Before: S(""), After: S("A ")}
	c1_2 := changes.Splice{Offset: 2, Before: S(""), After: S("B ")}
	c2_1 := changes.Splice{Offset: 0, Before: S(""), After: S("X ")}
	c2_2 := changes.Splice{Offset: 2, Before: S(""), After: S("Y ")}
	c1_3 := changes.Splice{Offset: 6, Before: S(""), After: S("C ")}
	c2_3 := changes.Splice{Offset: 6, Before: S(""), After: S("Z ")}

	items := []ops.Op{
		ops.Operation{"first", nil, -1, -1, first},
		ops.Operation{"c1_1", nil, -1, 0, c1_1},
		ops.Operation{"c1_2", "c1_1", -1, 0, c1_2},
		ops.Operation{"c2_1", nil, -1, 0, c2_1},
		ops.Operation{"c2_2", "c2_1", -1, 0, c2_2},
		ops.Operation{"c1_3", nil, -1, 3, c1_3},
		ops.Operation{"c2_3", "c2_2", -1, 1, c2_3},
	}

	return changes.Nil, S("A B X Y C Z Hello World"), items
}

type failNthGetSince struct {
	n   int
	err error
	ops.Store
}

func (f *failNthGetSince) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	f.n--
	if f.n <= 0 {
		return nil, f.err
	}
	return f.Store.GetSince(ctx, version, limit)
}
