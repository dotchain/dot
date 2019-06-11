// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

func ExampleTable() {
	t := &data.Table{}
	col1, col2 := rich.NewText("col1"), rich.NewText("col2")
	t = t.Apply(nil, t.AppendCol(data.Col{ID: "col2", Value: col2})).(*data.Table)
	t = t.Apply(nil, t.InsertColBefore(data.Col{ID: "col1", Value: col1}, *t.Cols["col2"])).(*data.Table)

	t = t.Apply(nil, t.AppendRow(data.Row{ID: "row1"})).(*data.Table)
	t = t.Apply(nil, t.SetCellValue("row1", "col1", rich.NewText("1-1"))).(*data.Table)
	t = t.Apply(nil, t.SetCellValue("row1", "col2", rich.NewText("1-2"))).(*data.Table)

	fmt.Println(tableToText(t))

	// Output: 1-1 1-2
}

func tableToText(t *data.Table) string {
	lines := []string{}
	colIDs := t.ColIDs()
	for _, rowID := range t.RowIDs() {
		cells := []string{}
		r := *t.Rows[rowID]
		for _, colID := range colIDs {
			cell := ""
			if v, ok := r.Cells[colID]; ok {
				cell = v.(*rich.Text).PlainText()
			}
			cells = append(cells, cell)
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	return strings.Join(lines, "\n")
}

func TestTableUpdateCol(t *testing.T) {
	col1 := rich.NewText("col1")
	tbl := &data.Table{
		Cols: data.Cols{
			"col1": &data.Col{
				ID:    "col1",
				Value: col1,
			},
		},
	}

	c := changes.Move{Offset: 3, Count: 1, Distance: -3}
	tbl = tbl.Apply(nil, tbl.UpdateCol("col1", c)).(*data.Table)

	if x := tbl.Cols["col1"].Value.PlainText(); x != "1col" {
		t.Error("Unexpected", x)
	}
}

func TestTableDeleteCol(t *testing.T) {
	col1 := rich.NewText("col1")
	tbl := &data.Table{
		Cols: data.Cols{
			"col1": &data.Col{
				ID:    "col1",
				Value: col1,
			},
		},
	}

	tbl = tbl.Apply(nil, tbl.DeleteCol("col1")).(*data.Table)

	if x, ok := tbl.Cols["col1"]; ok {
		t.Error("Unexpected", x)
	}
}

func TestTableUpdateCellValue(t *testing.T) {
	val := rich.NewText("val1")
	tbl := &data.Table{
		Rows: data.Rows{
			"row1": &data.Row{
				ID:    "row1",
				Cells: types.M{"col1": val},
			},
		},
	}

	c := changes.Move{Offset: 3, Count: 1, Distance: -3}
	tbl = tbl.Apply(nil, tbl.UpdateCellValue("row1", "col1", c)).(*data.Table)

	if x := tbl.Rows["row1"].Cells["col1"].(*rich.Text).PlainText(); x != "1val" {
		t.Error("Unexpected", x)
	}
}

func TestTableSetCellValue(t *testing.T) {
	val := rich.NewText("val1")
	tbl := &data.Table{Rows: data.Rows{"row1": &data.Row{ID: "row1"}}}

	tbl = tbl.Apply(nil, tbl.SetCellValue("row1", "col1", val)).(*data.Table)

	if x := tbl.Rows["row1"].Cells["col1"].(*rich.Text).PlainText(); x != "val1" {
		t.Error("Unexpected", x)
	}

	val = rich.NewText("val2")

	tbl = tbl.Apply(nil, tbl.SetCellValue("row1", "col1", val)).(*data.Table)

	if x := tbl.Rows["row1"].Cells["col1"].(*rich.Text).PlainText(); x != "val2" {
		t.Error("Unexpected", x)
	}
}

func TestTableDeleteRow(t *testing.T) {
	tbl := &data.Table{Rows: data.Rows{"row1": &data.Row{ID: "row1"}}}
	tbl = tbl.Apply(nil, tbl.DeleteRow("row1")).(*data.Table)

	if x, ok := tbl.Rows["row1"]; ok {
		t.Error("Unexpected", x)
	}
}

func TestTableAppendCol(t *testing.T) {
	tbl := &data.Table{}
	tbl = tbl.Apply(nil, tbl.AppendCol(data.Col{ID: "col1"})).(*data.Table)
	tbl = tbl.Apply(nil, tbl.AppendCol(data.Col{ID: "col2"})).(*data.Table)
	if x := tbl.ColIDs(); !reflect.DeepEqual(x, []interface{}{"col1", "col2"}) {
		t.Error("Unexpected", x)
	}

	tbl = tbl.Apply(nil, tbl.AppendCol(*tbl.Cols["col1"])).(*data.Table)
	if x := tbl.ColIDs(); !reflect.DeepEqual(x, []interface{}{"col2", "col1"}) {
		t.Error("Unexpected", x)
	}
}

func TestTableInsertColBefore(t *testing.T) {
	tbl := &data.Table{Cols: data.Cols{"col3": &data.Col{ID: "col3"}}}

	col3 := *tbl.Cols["col3"]
	tbl = tbl.Apply(nil, tbl.InsertColBefore(data.Col{ID: "col1"}, col3)).(*data.Table)
	tbl = tbl.Apply(nil, tbl.InsertColBefore(data.Col{ID: "col2"}, col3)).(*data.Table)

	if x := tbl.ColIDs(); !reflect.DeepEqual(x, []interface{}{"col1", "col2", "col3"}) {
		t.Error("Unexpected", x)
	}

	col1 := *tbl.Cols["col1"]
	tbl = tbl.Apply(nil, tbl.InsertColBefore(col1, col3)).(*data.Table)
	if x := tbl.ColIDs(); !reflect.DeepEqual(x, []interface{}{"col2", "col1", "col3"}) {
		t.Error("Unexpected", x)
	}
}

func TestTableAppendRow(t *testing.T) {
	tbl := &data.Table{}
	tbl = tbl.Apply(nil, tbl.AppendRow(data.Row{ID: "row1"})).(*data.Table)
	tbl = tbl.Apply(nil, tbl.AppendRow(data.Row{ID: "row2"})).(*data.Table)
	if x := tbl.RowIDs(); !reflect.DeepEqual(x, []interface{}{"row1", "row2"}) {
		t.Error("Unexpected", x)
	}

	tbl = tbl.Apply(nil, tbl.AppendRow(*tbl.Rows["row1"])).(*data.Table)
	if x := tbl.RowIDs(); !reflect.DeepEqual(x, []interface{}{"row2", "row1"}) {
		t.Error("Unexpected", x)
	}
}

func TestTableInsertRowBefore(t *testing.T) {
	tbl := &data.Table{Rows: data.Rows{"row3": &data.Row{ID: "row3"}}}

	row3 := *tbl.Rows["row3"]
	tbl = tbl.Apply(nil, tbl.InsertRowBefore(data.Row{ID: "row1"}, row3)).(*data.Table)
	tbl = tbl.Apply(nil, tbl.InsertRowBefore(data.Row{ID: "row2"}, row3)).(*data.Table)

	if x := tbl.RowIDs(); !reflect.DeepEqual(x, []interface{}{"row1", "row2", "row3"}) {
		t.Error("Unexpected", x)
	}

	row1 := *tbl.Rows["row1"]
	tbl = tbl.Apply(nil, tbl.InsertRowBefore(row1, row3)).(*data.Table)
	if x := tbl.RowIDs(); !reflect.DeepEqual(x, []interface{}{"row2", "row1", "row3"}) {
		t.Error("Unexpected", x)
	}
}
