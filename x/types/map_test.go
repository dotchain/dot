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

func TestMApply(t *testing.T) {
	m := types.M{
		true: types.S8("bool"),
		5.3:  types.S8("float"),
	}

	x := m.Apply(nil)
	if !reflect.DeepEqual(x, m) {
		t.Error("Unexpected Apply.nil", x)
	}

	x = m.Apply(changes.Replace{m, changes.Nil})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = m.Apply(changes.Replace{m, types.S16("OK")})
	if x != types.S16("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	insert := changes.PathChange{[]interface{}{"new"}, changes.Replace{changes.Nil, types.S8("string")}}
	expected := types.M{
		true:  types.S8("bool"),
		5.3:   types.S8("float"),
		"new": types.S8("string"),
	}

	x = m.Apply(insert)
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected insert", x)
	}

	x = m.Apply(changes.ChangeSet{insert})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.ChangeSet", x)
	}

	x = m.Apply(changes.PathChange{nil, insert})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	modify := changes.PathChange{[]interface{}{true}, changes.Replace{types.S8("bool"), types.S8("BOOL")}}
	expected = types.M{
		true: types.S8("BOOL"),
		5.3:  types.S8("float"),
	}
	x = m.Apply(modify)
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	// validate that nil values can be replaced
	m = types.M{"nil": nil}
	rep := changes.Replace{changes.Nil, types.S8("OK")}
	x = m.Apply(changes.PathChange{[]interface{}{"nil"}, rep})
	if !reflect.DeepEqual(x, types.M{"nil": types.S8("OK")}) {
		t.Error("Unexpected apply with nil element", x)
	}

	remove := changes.Replace{types.S8("OK"), changes.Nil}
	x = x.Apply(changes.PathChange{[]interface{}{"nil"}, remove})
	if !reflect.DeepEqual(x, m) {
		t.Error("Unexpected apply with nil element", x)
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
		m := types.M{"x": types.S8("OK")}
		m.Slice(0, 0)
	})

	mustPanic(func() {
		m := types.M{"x": types.S8("OK")}
		m.Count()
	})

	mustPanic(func() {
		(types.M{}).Apply(poorlyDefinedChange{})
	})

	mustPanic(func() {
		(types.M{}).Apply(changes.Move{1, 1, 1})
	})

	mustPanic(func() {
		(types.M{}).Apply(changes.Splice{Before: types.S8(""), After: types.S8("OK")})
	})
}
