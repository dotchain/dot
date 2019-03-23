// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func TestMApply(t *testing.T) {
	m := types.M{
		true: types.S8("bool"),
		5.3:  types.S8("float"),
	}

	x := m.Apply(nil, nil)
	if !reflect.DeepEqual(x, m) {
		t.Error("Unexpected Apply.nil", x)
	}

	x = m.Apply(nil, changes.Replace{m, changes.Nil})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = m.Apply(nil, changes.Replace{m, types.S16("OK")})
	if x != types.S16("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	insert := changes.PathChange{[]interface{}{"new"}, changes.Replace{changes.Nil, types.S8("string")}}
	expected := types.M{
		true:  types.S8("bool"),
		5.3:   types.S8("float"),
		"new": types.S8("string"),
	}

	x = m.Apply(nil, insert)
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected insert", x)
	}

	x = m.Apply(nil, changes.ChangeSet{insert})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.ChangeSet", x)
	}

	x = m.Apply(nil, changes.PathChange{nil, insert})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	modify := changes.PathChange{[]interface{}{true}, changes.Replace{types.S8("bool"), types.S8("BOOL")}}
	expected = types.M{
		true: types.S8("BOOL"),
		5.3:  types.S8("float"),
	}
	x = m.Apply(nil, modify)
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	remove := changes.PathChange{[]interface{}{5.3}, changes.Replace{types.S8("float"), changes.Nil}}
	expected = types.M{true: types.S8("BOOL")}
	x = x.Apply(nil, remove)
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.PathChange", x)
	}
}

func TestMPanics(t *testing.T) {
	mustPanic := func(fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Failed to panic")
			}
		}()
		fn()
	}

	mustPanic(func() {
		(types.M{}).Apply(nil, poorlyDefinedChange{})
	})

	mustPanic(func() {
		(types.M{}).Apply(nil, changes.Move{1, 1, 1})
	})

	mustPanic(func() {
		(types.M{}).Apply(nil, changes.Splice{Before: types.S8(""), After: types.S8("OK")})
	})
}
