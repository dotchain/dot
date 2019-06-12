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

func TestCol(t *testing.T) {
	v1, v2 := rich.NewText("value1"), rich.NewText("value2")
	col := data.Col{ID: "col1", Value: v1}
	col = col.Apply(nil, changes.PathChange{
		Path: []interface{}{"Ord"},
		Change: changes.Replace{
			Before: types.S16(""),
			After:  types.S16("boo"),
		},
	}).(data.Col)

	if col.Ord != "boo" {
		t.Error("Unexpected ord change", col.Ord)
	}

	col = col.Apply(nil, changes.PathChange{
		Path:   []interface{}{"Value"},
		Change: changes.Replace{Before: v1, After: v2},
	}).(data.Col)

	if col.Value.PlainText() != v2.PlainText() {
		t.Error("Unexpected value", col.Value.PlainText())
	}
}
