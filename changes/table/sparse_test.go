// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package table_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/table"
	"github.com/dotchain/dot/changes/types"
	"reflect"
	"testing"
)

func TestSparse(t *testing.T) {
	s := Sparse{}
	t.Run("SpliceRows", s.TestSpliceRows)
	t.Run("SpliceCols", s.TestSpliceCols)
	t.Run("Cell", s.TestCell)
	t.Run("GC", s.TestGC)
	t.Run("EdgeCases", s.TestEdgeCases)
}

type Sparse struct{}

func (x Sparse) TestSpliceRows(t *testing.T) {
	s := table.Sparse{}
	abc, c1 := s.SpliceRows(0, 0, []interface{}{"a", "b", "c"})
	ac, c2 := abc.SpliceRows(1, 1, nil)
	expected := table.Sparse{RowIDs: types.A{changes.Atomic{"a"}, changes.Atomic{"c"}}}
	if !reflect.DeepEqual(ac, expected) {
		t.Error("Unexpected splice result", ac)
	}

	s2 := s.Apply(c1).Apply(c2)
	if !reflect.DeepEqual(s2, expected) {
		t.Error("Unexpected changes", s2)
	}
}

func (x Sparse) TestSpliceCols(t *testing.T) {
	s := table.Sparse{}
	abc, c1 := s.SpliceCols(0, 0, []interface{}{"a", "b", "c"})
	ac, c2 := abc.SpliceCols(1, 1, nil)
	expected := table.Sparse{ColIDs: types.A{changes.Atomic{"a"}, changes.Atomic{"c"}}}
	if !reflect.DeepEqual(ac, expected) {
		t.Error("Unexpected splice result", ac)
	}

	s2 := s.Apply(c1).Apply(c2)
	if !reflect.DeepEqual(s2, expected) {
		t.Error("Unexpected changes", s2)
	}
}

func (x Sparse) TestCell(t *testing.T) {
	s := table.Sparse{}
	if v, ok := s.Cell("row", "col"); ok {
		t.Fatal("Unexpected cell value", v)
	}

	if x, c := s.RemoveCell("row", "col"); c != nil || !reflect.DeepEqual(x, s) {
		t.Error("Non-existent remove did something odd", x, c)
	}

	s1, c1 := s.UpdateCell("row", "col", "hello")
	if v, ok := s1.Cell("row", "col"); !ok || v != "hello" {
		t.Error("Update resulted in wrong value", v, ok)
	}

	if z := s.Apply(c1); !reflect.DeepEqual(z, s1) {
		t.Error("Update change does not match result", z)
	}

	s2, c2 := s1.UpdateCell("row", "col", types.S8("boo"))
	if v, ok := s2.Cell("row", "col"); !ok || v != types.S8("boo") {
		t.Error("Update resulted in wrong value", v, ok)
	}

	if z := s1.Apply(c2); !reflect.DeepEqual(z, s2) {
		t.Error("Update change does not match result", z)
	}

	s3, c3 := s2.RemoveCell("row", "col")
	if v, ok := s3.Cell("row", "col"); ok {
		t.Fatal("Unexpected cell value", v)
	}

	if z := s2.Apply(c3); !reflect.DeepEqual(z, s3) {
		t.Error("Update change does not match result", z)
	}
}

func (x Sparse) TestGC(t *testing.T) {
	s := table.Sparse{Data: types.M{}}
	s1, c1 := s.SpliceRows(0, 0, []interface{}{"row1"})
	s2, c2 := s1.SpliceCols(0, 0, []interface{}{"col1"})
	s3, c3 := s2.UpdateCell("row1", "col1", "hello")

	if x, c := s3.GC(nil, nil); c != nil || !reflect.DeepEqual(x, s3) {
		t.Error("GC threw out valid entries", x, c)
	}

	s4, c4 := s3.UpdateCell("row1", "col2", "a")
	s5, c5 := s4.UpdateCell("row2", "col1", "b")

	s6, c6 := s5.GC(nil, nil)

	if !reflect.DeepEqual(s6, s3) {
		t.Error("GC didnt work", s6)
	}

	z := s.Apply(changes.ChangeSet{c1, c2, c3, c4, c5, c6})
	if !reflect.DeepEqual(z, s3) {
		t.Error("GC didnt work", x, c6)
	}
}

func (x Sparse) TestEdgeCases(t *testing.T) {
	s := table.Sparse{Data: types.M{}}
	s1, c1 := s.SpliceRows(0, 0, []interface{}{"row1"})
	s2 := s.Apply(changes.PathChange{nil, c1})
	if !reflect.DeepEqual(s2, s1) {
		t.Error("Nonconsequential path change had wrong effect", s2)
	}

	s2 = s.Apply(nil)
	if !reflect.DeepEqual(s2, s) {
		t.Error("Nil change had wrong effect", s2)
	}
}
