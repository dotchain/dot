// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

func TestCells(t *testing.T) {
	cells := data.Cells{}
	cells = cells.Apply(nil, changes.PathChange{
		Path: []interface{}{"col1"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  rich.NewText("cell1"),
		},
	}).(data.Cells)

	if v, ok := cells["col1"]; !ok || v.PlainText() != "cell1" {
		t.Error("Unexpected apply", v, ok)
	}

	cells = cells.Apply(nil, changes.PathChange{
		Path: []interface{}{"col1"},
		Change: changes.Replace{
			Before: rich.NewText("cell1"),
			After:  changes.Nil,
		},
	}).(data.Cells)

	if v, ok := cells["col1"]; ok {
		t.Error("Unexpected apply", v, ok)
	}
}
