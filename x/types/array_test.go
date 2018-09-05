// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestASlice(t *testing.T) {
	a := types.A{types.S8("a"), types.S8("b"), types.S8("c"), types.S8("d"), types.S8("e")}
	if x := a.Slice(3, 0); x.Count() != 0 {
		t.Error("Unexpected Slice(3,0)", x)
	}

	if x := a.Slice(3, 1); !reflect.DeepEqual(x, types.A{types.S8("d")}) {
		t.Error("Unexpected Slice(3,1)", x)
	}
}

func TestACount(t *testing.T) {
	a := types.A(nil)
	if x := a.Count(); x != 0 {
		t.Error("Unexpected Count", x)
	}

	a = types.A{types.S8("a")}
	if x := a.Count(); x != 1 {
		t.Error("Unexpected Count", x)
	}
}

func TestAApply(t *testing.T) {
	a := types.A{types.S8("a"), types.S8("b"), types.S8("c"), types.S8("d"), types.S8("e")}

	x := a.Apply(nil)
	if !reflect.DeepEqual(x, a) {
		t.Error("Unexpected Apply.nil", x)
	}

	x = a.Apply(changes.Replace{Before: a, IsDelete: true})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = a.Apply(changes.Replace{Before: a, After: types.S16("OK")})
	if x != types.S16("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	x = a.Apply(changes.Splice{1, a.Slice(1, 3), types.A{types.S8("-")}})
	expected := types.A{types.S8("a"), types.S8("-"), types.S8("e")}
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = a.Apply(changes.Move{2, 2, -1})
	expected = types.A{types.S8("a"), types.S8("c"), types.S8("d"), types.S8("b"), types.S8("e")}
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.Move", x)
	}

	x = a.Apply(changes.ChangeSet{changes.Move{2, 2, -1}})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.ChangeSet", x)
	}

	x = a.Apply(changes.PathChange{nil, changes.Move{2, 2, -1}})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	insert := changes.Splice{0, types.S8(""), types.S8("**")}
	x = a.Apply(changes.PathChange{[]interface{}{0}, insert})
	a[0] = types.S8("**a")
	if !reflect.DeepEqual(x, a) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	// validate that nil values can be replaced
	a = types.A{nil}
	rep := changes.Replace{IsInsert: true, Before: changes.Nil, After: types.S8("OK")}
	x = a.Apply(changes.PathChange{[]interface{}{0}, rep})
	if !reflect.DeepEqual(x, types.A{types.S8("OK")}) {
		t.Error("Unexpected apply with nil element", x)
	}
}

func TestAPanics(t *testing.T) {
	mustPanic := func(fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Failed to panic")
			}
		}()
		fn()
	}

	mustPanic(func() {
		(types.A{}).Apply(poorlyDefinedChange{})
	})

	mustPanic(func() {
		(types.A{}).Apply(changes.ChangeSet{changes.Move{7, 3, -1}})
	})
}
