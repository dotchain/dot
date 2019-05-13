// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/fred"
)

func TestTupleZeroElements(t *testing.T) {
	if x := fred.ToTuple([]fred.Object{}); x != nil {
		t.Error("Unexpected", x)
	}
	if x := fred.FromTuple(nil); len(x) != 0 {
		t.Error("Unexpected", x)
	}
}

func TestTupleTenElements(t *testing.T) {
	array := [...]fred.Object{fred.Error("0"), nil, nil, fred.Error("1"), nil, nil, nil, nil, nil, nil}
	slice := array[:]

	if x := fred.ToTuple(slice); !reflect.DeepEqual(x, array) {
		t.Error("Unexpected", x)
	}
	if x := fred.FromTuple(array); !reflect.DeepEqual(x, slice) {
		t.Error("Unexpected", x)
	}
}
