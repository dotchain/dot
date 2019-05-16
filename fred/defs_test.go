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

func TestDefsNil(t *testing.T) {
	var d *fred.Defs

	if d.Count() != 0 {
		t.Error("Unexpected count", d.Count())
	}

	if x := d.Slice(0, 0); !reflect.DeepEqual(x, &fred.Defs{}) {
		t.Error("Unexpected slice", x)
	}

	if x := d.Eval(env); !reflect.DeepEqual(x, &fred.Vals{}) {
		t.Error("Unexpected eval", x)
	}

	expected := &fred.Defs{fred.Fixed(fred.Error("one"))}
	c := changes.Splice{Before: d.Slice(0, 0), After: expected}
	if x := d.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Defs.Apply", x)
	}
}

func TestDefsSplice(t *testing.T) {
	d := &fred.Defs{}

	if d.Count() != 0 {
		t.Error("Unexpected count", d.Count())
	}

	if x := d.Slice(0, 0); !reflect.DeepEqual(x, &fred.Defs{}) {
		t.Error("Unexpected slice", x)
	}

	if x := d.Eval(env); !reflect.DeepEqual(x, &fred.Vals{}) {
		t.Error("Unexpected eval", x)
	}

	expected := &fred.Defs{
		fred.Fixed(fred.Error("one")),
		fred.Fixed(fred.Error("two")),
	}
	c := changes.Splice{Before: d.Slice(0, 0), After: expected}
	if x := d.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Defs.Apply", x)
	}

	expected2 := &fred.Defs{
		fred.Fixed(fred.Error("one")),
	}

	c = changes.Splice{Offset: 1, Before: expected.Slice(1, 1), After: (*fred.Defs)(nil)}
	if x := expected.ApplyCollection(nil, c); !reflect.DeepEqual(x, expected2) {
		t.Error("Unexpected splice", x)
	}

	vals := &fred.Vals{fred.Error("one"), fred.Error("two")}
	if x := expected.Eval(env); !reflect.DeepEqual(vals, x) {
		t.Error("Unexpected eval", x)
	}
}

func TestDefsReplaceItem(t *testing.T) {
	d := &fred.Defs{fred.Fixed(fred.Error("before"))}
	c := changes.PathChange{
		Path: []interface{}{0, "Val"},
		Change: changes.Replace{
			Before: fred.Error("before"),
			After:  fred.Error("after"),
		},
	}
	expected := &fred.Defs{fred.Fixed(fred.Error("after"))}
	if x := d.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected apply", x)
	}
}
