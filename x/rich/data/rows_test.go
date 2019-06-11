// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich/data"
)

func TestRows(t *testing.T) {
	rows := data.Rows{}
	rows = rows.Apply(nil, changes.PathChange{
		Path: []interface{}{"row1"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  data.Row{ID: "row1"},
		},
	}).(data.Rows)

	rows = rows.Apply(nil, changes.PathChange{
		Path: []interface{}{"row2"},
		Change: changes.Replace{
			Before: changes.Nil,
			After:  data.Row{ID: "row2"},
		},
	}).(data.Rows)

	if v, ok := rows["row1"]; !ok || v.ID != "row1" {
		t.Error("Unexpected apply", v, ok)
	}

	if v, ok := rows["row2"]; !ok || v.ID != "row2" {
		t.Error("Unexpected apply", v, ok)
	}

	rows = rows.Apply(nil, changes.PathChange{
		Path: []interface{}{"row1"},
		Change: changes.Replace{
			Before: data.Row{ID: "row1"},
			After:  changes.Nil,
		},
	}).(data.Rows)

	if v, ok := rows["row1"]; ok {
		t.Error("Unexpected apply", v, ok)
	}
}
