// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func TestValMap(t *testing.T) {
	var v *fred.ValMap

	var c changes.Change = changes.ChangeSet{
		changes.PathChange{
			Path: []interface{}{"boo"},
			Change: changes.Replace{
				Before: changes.Nil,
				After:  fred.Error("hoo"),
			},
		},
		changes.PathChange{
			Path: []interface{}{"goo"},
			Change: changes.Replace{
				Before: changes.Nil,
				After:  fred.Error("goo"),
			},
		},
	}
	expected := &fred.ValMap{
		"boo": fred.Error("hoo"),
		"goo": fred.Error("goo"),
	}
	if x := v.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Vals..Apply", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"boo"},
		Change: changes.Replace{
			Before: fred.Error("hoo"),
			After:  fred.Error("hoo2"),
		},
	}
	expected2 := &fred.ValMap{
		"boo": fred.Error("hoo2"),
		"goo": fred.Error("goo"),
	}
	if x := expected.Apply(nil, c); !reflect.DeepEqual(x, expected2) {
		t.Error("Unexpected Vals..Apply", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"boo"},
		Change: changes.Replace{
			Before: fred.Error("hoo2"),
			After:  changes.Nil,
		},
	}

	expected3 := &fred.ValMap{
		"goo": fred.Error("goo"),
	}
	if x := expected2.Apply(nil, c); !reflect.DeepEqual(x, expected3) {
		t.Error("Unexpected Vals..Apply", x)
	}
}
