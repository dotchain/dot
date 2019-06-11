// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

func TestRow(t *testing.T) {
	row := data.Row{ID: "row1"}
	row = row.Apply(nil, changes.PathChange{
		Path: []interface{}{"Ord"},
		Change: changes.Replace{
			Before: types.S16(""),
			After:  types.S16("boo"),
		},
	}).(data.Row)

	if row.Ord != "boo" {
		t.Error("Unexpected ord change", row.Ord)
	}

	row = row.Apply(nil, changes.PathChange{
		Path: []interface{}{"Cells", "col1"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  rich.NewText("cell1"),
		},
	}).(data.Row)

	if x := row.Cells["col1"].(*rich.Text).PlainText(); x != "cell1" {
		t.Error("Unexpected value", x)
	}
}
