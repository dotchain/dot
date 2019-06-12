// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

func TestListApply(t *testing.T) {
	s1, s2 := rich.NewText("list1"), rich.NewText("list2")
	l1 := data.List{Type: "circle", Entries: types.A{s1}}
	l2 := data.List{Type: "square", Entries: types.A{s2}}

	if l1.Name() != "Embed" {
		t.Error("Unexpected name", l1.Name())
	}

	if x := l1.Apply(nil, nil); !reflect.DeepEqual(x, l1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: l2, After: l2}
	if x := l1.Apply(nil, replace); !reflect.DeepEqual(x, l2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path: []interface{}{"Type"},
		Change: changes.Replace{
			Before: types.S16("circle"),
			After:  types.S16("square"),
		},
	}
	if x := l1.Apply(nil, c).(data.List); x.Type != "square" {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"Entries"},
		Change: changes.Replace{
			Before: types.A{s1},
			After:  types.A{s2},
		},
	}
	if x := l1.Apply(nil, c).(data.List); !reflect.DeepEqual(x.Entries, types.A{s2}) {
		t.Error("Unexpected change", x.Entries)
	}
}
