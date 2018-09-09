// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rle_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rle"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestEncoding(t *testing.T) {
	a := rle.Encode(types.A{types.S8("a"), types.S8("a"), types.S8("b"), types.S8("b"), types.S8("c")})
	expected := rle.A{{types.S8("a"), 2}, {types.S8("b"), 2}, {types.S8("c"), 1}}
	if !reflect.DeepEqual(a, expected) {
		t.Fatal("Failed to encode", a)
	}

	sliced := a.Slice(1, 4)
	expected = rle.A{{types.S8("a"), 1}, {types.S8("b"), 2}, {types.S8("c"), 1}}

	if !reflect.DeepEqual(sliced, expected) {
		t.Fatal("Failed to encode", a)
	}

	sliced = a.Slice(1, 2)
	expected = rle.A{{types.S8("a"), 1}, {types.S8("b"), 1}}

	if !reflect.DeepEqual(sliced, expected) {
		t.Fatal("Failed to encode", a)
	}

	replace := rle.A{{types.S8("a"), 3}, {types.S8("z"), 2}, {types.S8("b"), 9}}
	spliced := a.Apply(changes.Splice{1, sliced, replace})
	expected = rle.A{{types.S8("a"), 4}, {types.S8("z"), 2}, {types.S8("b"), 10}, {types.S8("c"), 1}}

	if !reflect.DeepEqual(spliced, expected) {
		t.Fatal("Failed to encode", a)
	}

	a2b := changes.Replace{types.S8("a"), types.S8("b")}
	v := a.Apply(changes.PathChange{[]interface{}{1}, a2b})
	expected = rle.A{{types.S8("a"), 1}, {types.S8("b"), 3}, {types.S8("c"), 1}}

	if !reflect.DeepEqual(v, expected) {
		t.Fatal("Failed to encode", a)
	}

	b2a := changes.Replace{types.S8("b"), types.S8("a")}
	v = a.Apply(changes.PathChange{[]interface{}{2}, b2a})
	expected = rle.A{{types.S8("a"), 3}, {types.S8("b"), 1}, {types.S8("c"), 1}}

	if !reflect.DeepEqual(v, expected) {
		t.Fatal("Failed to encode", a)
	}
}

func TestIsEqual(t *testing.T) {
	base := types.A{ratio{5, 2}, ratio{10, 4}, ratio{10, 3}}
	v1 := rle.Encode(base)
	if v1.IsEqual(base) {
		t.Error("Wrongly equates encoded with unencoded")
	}

	v2 := rle.Encode(types.A{ratio{5, 2}, ratio{5, 2}, ratio{10, 3}})
	if !v1.IsEqual(v2) {
		t.Error("Failed to coalesce", v1, v2)
	}

	if v1.IsEqual(rle.A{}) {
		t.Error("Wrongly compares to nil", v1)
	}

	v3 := rle.A{{ratio{5, 2}, 4}, {ratio{10, 3}, 1}}
	if v1.IsEqual(v3) {
		t.Error("Wrongly compares counts", v1)
	}

	v4 := rle.A{{ratio{5, 3}, 2}, {ratio{10, 3}, 1}}
	if v1.IsEqual(v4) {
		t.Error("Wrongly compares counts", v1)
	}

}

func TestSlice(t *testing.T) {
	a := rle.Encode(types.A{types.S8("a"), types.S8("b"), types.S8("c"), types.S8("d"), types.S8("e")})
	if x := a.Slice(3, 0); x.Count() != 0 {
		t.Error("Unexpected Slice(3,0)", x)
	}

	if x := a.Slice(3, 1); !reflect.DeepEqual(x, rle.Encode(types.A{types.S8("d")})) {
		t.Error("Unexpected Slice(3,1)", x)
	}
}

func TestACount(t *testing.T) {
	a := rle.A(nil)
	if x := a.Count(); x != 0 {
		t.Error("Unexpected Count", x)
	}

	a = rle.Encode(types.A{types.S8("a")})
	if x := a.Count(); x != 1 {
		t.Error("Unexpected Count", x)
	}
}

func TestAApply(t *testing.T) {
	a := rle.Encode(types.A{types.S8("a"), types.S8("b"), types.S8("c"), types.S8("d"), types.S8("e")})

	x := a.Apply(nil)
	if !reflect.DeepEqual(x, a) {
		t.Error("Unexpected Apply.nil", x)
	}

	x = a.Apply(changes.Replace{a, changes.Nil})
	if x != changes.Nil {
		t.Error("Unexpeted Apply.Replace-Delete", x)
	}

	x = a.Apply(changes.Replace{a, types.S16("OK")})
	if x != types.S16("OK") {
		t.Error("Unexpected Apply.Replace", x)
	}

	x = a.Apply(changes.Splice{1, a.Slice(1, 3), rle.Encode(types.A{types.S8("-")})})
	expected := rle.Encode(types.A{types.S8("a"), types.S8("-"), types.S8("e")})
	if !reflect.DeepEqual(x, expected) {
		t.Error("Unexpected Apply.Splice", x)
	}

	x = a.Apply(changes.Move{2, 2, -1})
	expected = rle.Encode(types.A{types.S8("a"), types.S8("c"), types.S8("d"), types.S8("b"), types.S8("e")})
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
	a = rle.Encode(types.A{types.S8("**a"), types.S8("b"), types.S8("c"), types.S8("d"), types.S8("e")})
	if !reflect.DeepEqual(x, a) {
		t.Error("Unexpected Apply.PathChange", x)
	}

	// validate that nil values can be replaced
	a = rle.Encode(types.A{nil})
	rep := changes.Replace{changes.Nil, types.S8("OK")}
	x = a.Apply(changes.PathChange{[]interface{}{0}, rep})
	if !reflect.DeepEqual(x, rle.Encode(types.A{types.S8("OK")})) {
		t.Error("Unexpected apply with nil element", x)
	}

	remove := changes.Replace{types.S8("OK"), changes.Nil}
	x = x.Apply(changes.PathChange{[]interface{}{0}, remove})
	if !reflect.DeepEqual(x, a) {
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
		(rle.A{}).Apply(poorlyDefinedChange{})
	})

	mustPanic(func() {
		(rle.A{}).Apply(changes.ChangeSet{changes.Move{7, 3, -1}})
	})

	mustPanic(func() {
		(rle.A{}).Apply(changes.PathChange{[]interface{}{5}, nil})
	})
}

// this implements Change but not CustomChange
type poorlyDefinedChange struct{}

func (p poorlyDefinedChange) Merge(o changes.Change) (changes.Change, changes.Change) {
	return o, nil
}

func (p poorlyDefinedChange) Revert() changes.Change {
	return p
}

type ratio [2]int

func (r ratio) IsEqual(o changes.Value) bool {
	if rx, ok := o.(ratio); ok {
		return r[0]*rx[1] == r[1]*rx[0]
	}
	return false
}

func (r ratio) Count() int {
	panic("should not be called")
}

func (r ratio) Slice(offset, count int) changes.Value {
	panic("should not be called")
}

func (r ratio) Apply(c changes.Change) changes.Value {
	v := (changes.Atomic{r}).Apply(c)
	if ax, ok := v.(changes.Atomic); ok {
		if rx, ok := ax.Value.(ratio); ok {
			return rx
		}
	}
	return v
}
