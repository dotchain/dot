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

func TestValsNil(t *testing.T) {
	var v *fred.Vals

	if v.Count() != 0 {
		t.Error("Unexpected count", v.Count())
	}

	if x := v.Slice(0, 0); !reflect.DeepEqual(x, &fred.Vals{}) {
		t.Error("Unexpected slice", x)
	}

	expected := &fred.Vals{fred.Error("one")}
	c := changes.Splice{Before: v.Slice(0, 0), After: expected}
	if x := v.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Vals.Apply", x)
	}
}

func TestValsSplice(t *testing.T) {
	v := &fred.Vals{}

	if v.Count() != 0 {
		t.Error("Unexpected count", v.Count())
	}

	if x := v.Slice(0, 0); !reflect.DeepEqual(x, &fred.Vals{}) {
		t.Error("Unexpected slice", x)
	}

	expected := &fred.Vals{fred.Error("one"), fred.Error("two")}
	c := changes.Splice{Before: v.Slice(0, 0), After: expected}
	if x := v.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Vals.Apply", x)
	}

	expected2 := &fred.Vals{fred.Error("one")}
	c = changes.Splice{Offset: 1, Before: expected.Slice(1, 1), After: (*fred.Vals)(nil)}
	if x := expected.ApplyCollection(nil, c); !reflect.DeepEqual(x, expected2) {
		t.Error("Unexpected splice", x)
	}
}

func TestValsReplaceItem(t *testing.T) {
	v := &fred.Vals{fred.Error("before")}
	c := changes.PathChange{
		Path: []interface{}{0},
		Change: changes.Replace{
			Before: fred.Error("before"),
			After:  fred.Error("after"),
		},
	}
	expected := &fred.Vals{fred.Error("after")}
	if x := v.Apply(nil, c); !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected apply", x)
	}
}
