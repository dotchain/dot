// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"reflect"
	"testing"
)

func TestRefPath_Encode(t *testing.T) {
	var r *RefPath
	num1, num2, num3 := NewRefIndex("1"), NewRefIndex("2+"), NewRefIndex("3-")

	s := r.Append("root", nil).
		Append("", num1).
		Append("", num2).
		Append("", num3).
		Encode()
	expected := []string{"root", "1", "2+", "3-"}
	if !reflect.DeepEqual(s, expected) {
		t.Fatal("Unexpected r", s, expected)
	}

	s2 := NewRefPath(s).Encode()
	if !reflect.DeepEqual(s2, s) {
		t.Fatal("Unexpected", s2, s)
	}
}

func TestRefPath_Resolve(t *testing.T) {
	var r *RefPath
	num1, num2, num3 := NewRefIndex("1"), NewRefIndex("2+"), NewRefIndex("3-")

	inner1 := []interface{}{nil, nil, nil, "foo"}
	inner2 := []interface{}{nil, nil, inner1}
	obj := map[string]interface{}{"root": []interface{}{nil, inner2}}
	v, ok := r.Append("root", nil).
		Append("", num1).
		Append("", num2).
		Append("", num3).
		Resolve(obj)

	if !ok || !Utils(x).AreSame(v, "foo") {
		t.Fatal("Failed", ok, v)
	}
}

func TestRefPath_Apply_splice(t *testing.T) {
	var r *RefPath
	num1, num2, num3 := NewRefIndex("1"), NewRefIndex("2+"), NewRefIndex("3-")

	// insert before  1 so "1" becomes "2"
	splice1 := &SpliceInfo{Offset: 1, After: []interface{}{nil}}
	// insert at 2, so so becomes "4"
	splice2 := &SpliceInfo{Offset: 2, After: []interface{}{nil, nil}}
	// insert at 3, it should get ignored
	splice3 := &SpliceInfo{Offset: 3, After: []interface{}{nil, nil, nil}}

	changes := []Change{
		{Path: []string{"root"}, Splice: splice1},
		{Path: []string{"root", "2"}, Splice: splice2},
		{Path: []string{"root", "2", "4"}, Splice: splice3},
	}
	v, ok := r.Append("root", nil).
		Append("", num1).
		Append("", num2).
		Append("", num3).
		Apply(changes)

	expected := []string{"root", "2", "4+", "3-"}
	if !ok || !reflect.DeepEqual(v.Encode(), expected) {
		t.Fatal("Failed", ok, v.Encode())
	}
}
