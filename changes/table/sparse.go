// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package table implements a loose 2d collection of values
package table

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Sparse represents a table where the data is stored in a map
// accessed by the row and column ids.  Not all row-col pairs may have
// a corresponding value (hence the sparse).  In addition, deletion of
// rowIDs and column IDs do not automatically delete the corresponding
// data --  that is expected to be handled in a late
// garbage-collection phase.
//
// The row ID and column ID collections can also be used to store
// other row-ID specific or column-ID specific data but in that case,
// the caller must implement a rowKey and/or colKey function that maps
// these to the corresponding key in the data map. This is only used
// by GC().
//
// Under some very contrived conditions, RowIDs and ColIDs may contain
// duplicates. For this reason, the callers should take care when
// iterating these arrays. To maintain consistency, only the first
// occurrence of an ID should be considered and all further
// occurrences ignored.
type Sparse struct {
	RowIDs, ColIDs types.A
	Data           types.M
}

// SpliceRows splices a set of rowIDs at the provided offset, first
// removing the specified count of IDs.
func (s Sparse) SpliceRows(offset, remove int, ids []interface{}) (Sparse, changes.Change) {
	before := s.RowIDs.Slice(offset, remove)
	c := changes.Splice{offset, before, s.toValues(ids)}
	s.RowIDs = s.RowIDs.Apply(nil, c).(types.A)
	return s, changes.PathChange{[]interface{}{"RowIDs"}, c}
}

// SpliceCols splices a set of colIDs at the provided offset, first
// removing the specified count of IDs
func (s Sparse) SpliceCols(offset, remove int, ids []interface{}) (Sparse, changes.Change) {
	before := s.ColIDs.Slice(offset, remove)
	c := changes.Splice{offset, before, s.toValues(ids)}
	s.ColIDs = s.ColIDs.Apply(nil, c).(types.A)
	return s, changes.PathChange{[]interface{}{"ColIDs"}, c}
}

// Cell returns the value at a given cell position. The bool return
// value indicates if the element exists or not
func (s Sparse) Cell(row, col interface{}) (interface{}, bool) {
	return s.fromValue(s.Data[[2]interface{}{row, col}])
}

// UpdateCell updates the value of a cell (which need not exist)
func (s Sparse) UpdateCell(row, col, value interface{}) (Sparse, changes.Change) {
	key := [2]interface{}{row, col}
	c := changes.Replace{changes.Nil, s.toValue(value)}
	if before, ok := s.Data[key]; ok {
		c.Before = before
	}
	pc := changes.PathChange{[]interface{}{key}, c}
	s.Data = s.Data.Apply(nil, pc).(types.M)
	return s, changes.PathChange{[]interface{}{"Data", key}, c}
}

// RemoveCell removes the value of a cell (which need not exist)
func (s Sparse) RemoveCell(row, col interface{}) (Sparse, changes.Change) {
	key := [2]interface{}{row, col}
	c := changes.Replace{changes.Nil, changes.Nil}
	if before, ok := s.Data[key]; ok {
		c.Before = before
	} else {
		return s, nil
	}
	pc := changes.PathChange{[]interface{}{key}, c}
	s.Data = s.Data.Apply(nil, pc).(types.M)
	return s, changes.PathChange{[]interface{}{"Data", key}, c}
}

// GC finds all cells that are orphaned by prior row/column deletes
// and clears them out.  It returns the changes matching the delete
// as well. The rowKey and colKey functions provide the mapping from
// the corresponding IDs to the key. These can be nil to indicate
// the IDs are themselves the keys into the Data map.
func (s Sparse) GC(rowKey, colKey func(interface{}) interface{}) (Sparse, changes.Change) {
	if rowKey == nil {
		rowKey = func(x interface{}) interface{} { return x }
	}
	if colKey == nil {
		colKey = func(x interface{}) interface{} { return x }
	}

	c := changes.ChangeSet(nil)
	data := types.M{}
	s.forEachRow(rowKey, func(rowID interface{}) {
		s.forEachCol(colKey, func(colID interface{}) {
			key := [2]interface{}{rowKey(rowID), colKey(colID)}
			if v, ok := s.Data[key]; ok {
				data[key] = v
			}
		})
	})
	for key, val := range s.Data {
		if _, ok := data[key]; !ok {
			replace := changes.Replace{val, changes.Nil}
			pc := changes.PathChange{[]interface{}{key}, replace}
			c = append(c, pc)
		}
	}
	if len(c) == 0 {
		return s, nil
	}
	s.Data = data
	return s, changes.PathChange{[]interface{}{"Data"}, c}
}

// Apply implements changes.Value:Apply
func (s Sparse) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return s
	case changes.PathChange:
		if len(c.Path) == 0 {
			return s.Apply(ctx, c.Change)
		}
		field := c.Path[0].(string)
		c.Path = c.Path[1:]
		switch field {
		case "RowIDs":
			s.RowIDs = s.RowIDs.Apply(ctx, c).(types.A)
		case "ColIDs":
			s.ColIDs = s.ColIDs.Apply(ctx, c).(types.A)
		case "Data":
			s.Data = s.Data.Apply(ctx, c).(types.M)
		}
		return s
	}
	return c.(changes.Custom).ApplyTo(ctx, s)
}

func (s Sparse) toValues(v []interface{}) types.A {
	result := make(types.A, len(v))
	for kk, vv := range v {
		result[kk] = s.toValue(vv)
	}
	return result
}

func (s Sparse) toValue(v interface{}) changes.Value {
	if val, ok := v.(changes.Value); ok {
		return val
	}
	return changes.Atomic{v}
}

func (s Sparse) fromValue(v changes.Value) (interface{}, bool) {
	switch v := v.(type) {
	case nil:
		return nil, false
	case changes.Atomic:
		return v.Value, true
	}

	return v, true
}

func (s Sparse) forEachRow(key func(interface{}) interface{}, fn func(rowID interface{})) {
	seen := map[interface{}]bool{}
	for _, row := range s.RowIDs {
		r, _ := s.fromValue(row)
		k := key(r)
		if !seen[k] {
			seen[k] = true
			fn(r)
		}
	}
}

func (s Sparse) forEachCol(key func(interface{}) interface{}, fn func(colID interface{})) {
	seen := map[interface{}]bool{}
	for _, col := range s.ColIDs {
		c, _ := s.fromValue(col)
		k := key(c)
		if !seen[k] {
			seen[k] = true
			fn(c)
		}
	}
}
