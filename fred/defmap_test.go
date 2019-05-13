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

func TestDefMap(t *testing.T) {
	var d *fred.DefMap

	if x := d.Eval(env); !reflect.DeepEqual(x, &fred.ValMap{}) {
		t.Error("Unexpected nil eval", x)
	}

	var c changes.Change = changes.ChangeSet{
		changes.PathChange{
			Path: []interface{}{"boo"},
			Change: changes.Replace{
				Before: changes.Nil,
				After:  &fred.Fixed{Val: fred.Error("hoo")},
			},
		},
		changes.PathChange{
			Path: []interface{}{"goo"},
			Change: changes.Replace{
				Before: changes.Nil,
				After:  &fred.Fixed{Val: fred.Error("goo")},
			},
		},
	}
	expected := &fred.DefMap{
		"boo": &fred.Fixed{Val: fred.Error("hoo")},
		"goo": &fred.Fixed{Val: fred.Error("goo")},
	}
	if x := d.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Defs..Apply", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"boo"},
		Change: changes.Replace{
			Before: &fred.Fixed{Val: fred.Error("hoo")},
			After:  &fred.Fixed{Val: fred.Error("hoo2")},
		},
	}
	expected2 := &fred.DefMap{
		"boo": &fred.Fixed{Val: fred.Error("hoo2")},
		"goo": &fred.Fixed{Val: fred.Error("goo")},
	}
	if x := expected.Apply(nil, c); !reflect.DeepEqual(x, expected2) {
		t.Error("Unexpected Defs..Apply", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"boo"},
		Change: changes.Replace{
			Before: &fred.Fixed{Val: fred.Error("hoo2")},
			After:  changes.Nil,
		},
	}

	expected3 := &fred.DefMap{
		"goo": &fred.Fixed{Val: fred.Error("goo")},
	}
	if x := expected2.Apply(nil, c); !reflect.DeepEqual(x, expected3) {
		t.Error("Unexpected Defs..Apply", x)
	}

	if x := expected3.Eval(env); !reflect.DeepEqual(x, &fred.ValMap{"goo": fred.Error("goo")}) {
		t.Error("Unexpected eval result", x)
	}
}
