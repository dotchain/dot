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

func TestCols(t *testing.T) {
	col1, col2 := rich.NewText("col1"), rich.NewText("col2")
	cols := data.Cols{}
	cols = cols.Apply(nil, changes.PathChange{
		Path: []interface{}{"col1"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  data.Col{ID: "col1", Value: &col1},
		},
	}).(data.Cols)

	cols = cols.Apply(nil, changes.PathChange{
		Path: []interface{}{"col2"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  data.Col{ID: "col2", Value: &col2},
		},
	}).(data.Cols)

	if v, ok := cols["col1"]; !ok || v.ID != "col1" {
		t.Error("Unexpected apply", v, ok)
	}

	if v, ok := cols["col2"]; !ok || v.ID != "col2" {
		t.Error("Unexpected apply", v, ok)
	}

	cols = cols.Apply(nil, changes.PathChange{
		Path: []interface{}{"col1"},
		Change: changes.Replace{
			Before: data.Col{ID: "col1", Value: &col1},
			After:  data.Col{ID: "col1", Value: &col2},
		},
	}).(data.Cols)

	if v, ok := cols["col1"]; v.Value.PlainText() != "col2" {
		t.Error("Unexpected apply", v, ok)
	}
}
