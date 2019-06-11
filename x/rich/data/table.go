// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package data impleements data structures for use with rich text
package data

import (
	"sort"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"

	// crdt is only used for ord management
	"github.com/dotchain/dot/changes/crdt"
)

// Table represents a table with named columns and rows
//
// This is an immutable type and all mutation methods (like AppendCol)
// return changes.Change (which can be applied to get the new Table)
//
// The rows and columns are stored in a map structure but ordered IDs
// can be obtained via RowIDs and ColIDs
type Table struct {
	Cols
	Rows
}

// ColIDs returns all the column IDs sorted in order
func (t Table) ColIDs() []interface{} {
	ids := make([]interface{}, 0, len(t.Cols))
	for id := range t.Cols {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		c1, c2 := *t.Cols[ids[i]], *t.Cols[ids[j]]
		return crdt.LessOrd(c1.Ord, c2.Ord)
	})
	return ids
}

// RowIDs returns all the row IDs sorted in order
func (t Table) RowIDs() []interface{} {
	ids := make([]interface{}, 0, len(t.Rows))
	for id := range t.Rows {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		r1, r2 := *t.Rows[ids[i]], *t.Rows[ids[j]]
		return crdt.LessOrd(r1.Ord, r2.Ord)
	})
	return ids
}

// UpdateCol takes a change meant for a columm value and wraps it
// with the path of the column -- the resulting change can be applied
// directly on the table.
func (t Table) UpdateCol(colID interface{}, c changes.Change) changes.Change {
	path := []interface{}{"Cols", colID, "Value"}
	return changes.PathChange{Path: path, Change: c}
}

// UpdateCellValue takes a change meant for a cell value and wraps it
// with the path of the column -- the resulting change can be applied
// directly on the table.
//
// Note that the cell must exist (i.e the row must have a entry for
// the colID)
func (t Table) UpdateCellValue(rowID, colID interface{}, c changes.Change) changes.Change {
	path := []interface{}{"Rows", rowID, "Cells", colID}
	return changes.PathChange{Path: path, Change: c}
}

// SetCellValue sets the rich text for a cell. It works for rows which
// don't have values for specific columns as well as rows where the
// columns have a value already.
func (t Table) SetCellValue(rowID, colID interface{}, r rich.Text) changes.Change {
	c := changes.Replace{Before: changes.Nil, After: r}
	row := *t.Rows[rowID]
	if v, ok := row.Cells[colID]; ok {
		c.Before = *v
	}
	path := []interface{}{"Rows", rowID, "Cells", colID}
	return changes.PathChange{Path: path, Change: c}
}

// AppendCol adds a new column to the end of the sorted list.
//
// The column can already exist in which case it is simply moved.
func (t Table) AppendCol(col Col) changes.Change {
	if len(t.Cols) == 0 {
		col.Ord = ""
	} else {
		ids := t.ColIDs()
		last := *t.Cols[ids[len(ids)-1]]
		col.Ord = crdt.NextOrd(last.Ord)
	}
	return t.insertOrReorderCol(col)
}

// InsertColBefore inserts a new column before another
//
// The column may already exist in which case it is simply moved
func (t Table) InsertColBefore(col Col, beforeCol Col) changes.Change {
	ids := t.ColIDs()
	for kk, id := range ids {
		if id == beforeCol.ID {
			if kk == 0 {
				col.Ord = crdt.PrevOrd(t.Cols[ids[0]].Ord)
			} else {
				l := t.Cols[ids[kk-1]].Ord
				r := t.Cols[ids[kk]].Ord
				col.Ord = crdt.BetweenOrd(l, r, 1)[0]
			}
			return t.insertOrReorderCol(col)
		}
	}
	return nil
}

func (t Table) insertOrReorderCol(col Col) changes.Change {
	if before := t.Cols[col.ID]; before != nil {
		b, a := types.S16(before.Ord), types.S16(col.Ord)
		c := changes.Replace{Before: b, After: a}
		path := []interface{}{"Cols", col.ID, "Ord"}
		return changes.PathChange{Path: path, Change: c}
	}

	c := changes.Replace{Before: changes.Nil, After: col}
	path := []interface{}{"Cols", col.ID}
	return changes.PathChange{Path: path, Change: c}
}

// AppendRow adds a new row to the end of the sorted list.
//
// The row can already exist in which case it is simply moved.
func (t Table) AppendRow(row Row) changes.Change {
	if len(t.Rows) == 0 {
		row.Ord = ""
	} else {
		ids := t.RowIDs()
		last := *t.Rows[ids[len(ids)-1]]
		row.Ord = crdt.NextOrd(last.Ord)
	}
	return t.insertOrReorderRow(row)
}

// InsertRowBefore inserts a new row before another
//
// The row may already exist in which case it is simply moved
func (t Table) InsertRowBefore(row Row, beforeRow Row) changes.Change {
	ids := t.RowIDs()
	for kk, id := range ids {
		if id == beforeRow.ID {
			if kk == 0 {
				row.Ord = crdt.PrevOrd(t.Rows[ids[0]].Ord)
			} else {
				l := t.Rows[ids[kk-1]].Ord
				r := t.Rows[ids[kk]].Ord
				row.Ord = crdt.BetweenOrd(l, r, 1)[0]
			}
			return t.insertOrReorderRow(row)
		}
	}
	return nil
}

func (t Table) insertOrReorderRow(row Row) changes.Change {
	if before := t.Rows[row.ID]; before != nil {
		b, a := types.S16(before.Ord), types.S16(row.Ord)
		c := changes.Replace{Before: b, After: a}
		path := []interface{}{"Rows", row.ID, "Ord"}
		return changes.PathChange{Path: path, Change: c}
	}

	c := changes.Replace{Before: changes.Nil, After: row}
	path := []interface{}{"Rows", row.ID}
	return changes.PathChange{Path: path, Change: c}
}

// DeleteCol deletes a column by ID
func (t Table) DeleteCol(colID interface{}) changes.Change {
	var result changes.Change
	if col, ok := t.Cols[colID]; ok {
		c := changes.Replace{Before: *col, After: changes.Nil}
		path := []interface{}{"Cols", colID}
		result = changes.PathChange{Path: path, Change: c}
	}
	return result
}

// DeleteRow deletes a row by ID
func (t Table) DeleteRow(rowID interface{}) changes.Change {
	var result changes.Change
	if r, ok := t.Rows[rowID]; ok {
		c := changes.Replace{Before: *r, After: changes.Nil}
		path := []interface{}{"Rows", rowID}
		result = changes.PathChange{Path: path, Change: c}
	}
	return result
}

// Apply implements changes.Value
func (t Table) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: t.set, Get: t.get}).Apply(ctx, c, t)
}

func (t Table) get(key interface{}) changes.Value {
	if key == "Cols" {
		return t.Cols
	}
	return t.Rows
}

func (t Table) set(key interface{}, v changes.Value) changes.Value {
	if key == "Cols" {
		t.Cols = v.(Cols)
	} else {
		t.Rows = v.(Rows)
	}
	return t
}
