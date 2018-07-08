// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"reflect"
	"strings"
	"testing"
)

var refPathTable = [][]string{
	// Splice replace op, output RefIndex Start + End
	// The output values use "|" to indicate new index location
	{"-a*[1|b]c-", "-a|*bc-", "-a|*bc-"},
	{"-a[*1|bc]d-", "-a|bcd-", "-a|bcd-"},
	{"-a[1*2|bc]d-", "-abc|d-", "-a|bcd-"},
	{"-a[1*|bcd]e-", "-abcd|e-", "-a|bcde-"},
	{"-a[1|bcd]*e-", "-abcd|*e-", "-abcd|*e-"},
	{"-a[1|bcd]e*-", "-abcde|*-", "-abcde|*-"},

	// Splice deletes
	{"-a*[1|]c-", "-a|*c-", "-a|*c-"},
	{"-a[*1|]d-", "-a|d-", "-a|d-"},
	{"-a[1*2|]d-", "-a|d-", "-a|d-"},
	{"-a[1*|]e-", "-a|e-", "-a|e-"},
	{"-a[1|]*e-", "-a|*e-", "-a|*e-"},
	{"-a[1|]e*-", "-ae|*-", "-ae|*-"},

	// Splice inserts
	{"-a*[|b]c-", "-a|*bc-", "-a|*bc-"},
	{"-a[|b]*c-", "-ab|*c-", "-a|b*c-"},
	{"-a[|b]c*-", "-abc|*-", "-abc|*-"},

	// Move to right tests
	{"-*a[b]c|-", "-|*acb-", "-|*acb-"},
	{"-a[*b]c|-", "-ac|*b-", "-ac|*b-"},
	{"-a[b]*c|-", "-a|*cb-", "-a|*cb-"},
	{"-a[b]c*|-", "-ac|*b-", "-ac|*b-"},
	{"-a[b]c|*-", "-acb|*-", "-acb|*-"},

	// Move to left tests
	{"-*a|b[c]-", "-|*acb-", "-|*acb-"},
	{"-a|*b[c]-", "-ac|*b-", "-ac|*b-"},
	{"-a|b[*c]-", "-a|*cb-", "-a|*cb-"},
	{"-a|b[c*]-", "-ac|*b-", "-ac|*b-"},
	{"-a|b[c]*-", "-acb|*-", "-acb|*-"},
}

func TestRefPathApply(t *testing.T) {
	rootPath := dot.NewRefPath([]string{"root"})
	for _, entry := range refPathTable {
		ix, output, c := parseMutation(entry[0], false)
		c.Path = []string{"root"}
		input := ix.(string)
		index := strings.Index(input, "*")
		expectedIndex := strings.Index(output.(string), "*")

		// RefIndexPointer test
		ri := &dot.RefIndex{Type: dot.RefIndexPointer, Index: index}
		r := rootPath.Append("", ri).Append("plus", nil)

		result, ok := r.Apply([]dot.Change{c})
		if expectedIndex == -1 {
			if ok {
				t.Fatal("Failed", input, result.Encode())
			}
		} else if expectedIndex == index && result != r {
			t.Fatal("Changed", input, result.Encode())
		} else {
			ri = &dot.RefIndex{Type: ri.Type, Index: expectedIndex}
			expected := rootPath.Append("", ri).Append("plus", nil).Encode()
			if !reflect.DeepEqual(result.Encode(), expected) {
				t.Fatal("Differed", input, result.Encode(), expected)
			}
		}

		// RefIndexStart entry[1] test
		expectedIndex = strings.Index(entry[1], "|")
		ri = &dot.RefIndex{Type: dot.RefIndexStart, Index: index}
		r = rootPath.Append("", ri).Append("plus", nil)

		result, _ = r.Apply([]dot.Change{c})
		ri = &dot.RefIndex{Type: ri.Type, Index: expectedIndex}
		expected := rootPath.Append("", ri).Append("plus", nil)
		if !reflect.DeepEqual(result.Encode(), expected.Encode()) {
			t.Fatal("Differed start", entry[0], entry[1], result.Encode(), expected.Encode())
		}

		// RefIndexEnd entry[2] test
		expectedIndex = strings.Index(entry[2], "|")
		ri = &dot.RefIndex{Type: dot.RefIndexEnd, Index: index}
		r = rootPath.Append("", ri).Append("plus", nil)

		result, _ = r.Apply([]dot.Change{c})
		ri = &dot.RefIndex{Type: ri.Type, Index: expectedIndex}
		expected = rootPath.Append("", ri).Append("plus", nil)
		if !reflect.DeepEqual(result.Encode(), expected.Encode()) {
			t.Fatal("Differed end", entry[0], entry[2], result.Encode(), expected.Encode())
		}
	}
}

func TestRefUpdateEmptyLog(t *testing.T) {
	l := &dot.Log{}
	ref := dot.Ref{Path: dot.NewRefPath([]string{"5"})}
	u, err := ref.Update(l)
	if err != nil || !reflect.DeepEqual(ref, u) {
		t.Fatal("Unexpected update", u, err)
	}

	ref.ParentID = "something"
	u, err = ref.Update(l)
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Unexpected success with invalid parent", err, u)
	}
	ref.BasisID = "something"
	ref.ParentID = ""
	u, err = ref.Update(l)
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Unexpected success with invalid basis", err, u)
	}
}

func TestRefUpdateLog(t *testing.T) {
	changes1 := []dot.Change{{Splice: &dot.SpliceInfo{After: []interface{}{nil}}}}
	changes2 := []dot.Change{{Splice: &dot.SpliceInfo{After: []interface{}{nil}}}}
	changes3 := []dot.Change{{Splice: &dot.SpliceInfo{After: []interface{}{nil}}}}
	ops := []dot.Operation{
		{ID: "one", Changes: changes1},
		{ID: "two", Parents: []string{"one"}, Changes: changes2},
		{ID: "three", Parents: []string{"two"}, Changes: changes3},
	}

	l := &dot.Log{}
	for _, op := range ops {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("Append failed", err, op)
		}
	}
	ref := dot.Ref{Path: dot.NewRefPath([]string{"5"})}
	u, err := ref.Update(l)
	if err != nil {
		t.Fatal("Update failed", err)
	}
	if !reflect.DeepEqual(u.Path.Encode(), []string{"8"}) {
		t.Fatal("Unexpected new path", u.Path.Encode())
	}
	if u.BasisID != "three" {
		t.Fatal("Unexpected basis ID", u.BasisID)
	}
}

func TestRefUpdateInvalidated(t *testing.T) {
	changes := []dot.Change{{Splice: &dot.SpliceInfo{Before: []interface{}{nil}}}}
	l := &dot.Log{}
	op := dot.Operation{ID: "one", Changes: changes}
	if err := l.AppendOperation(op); err != nil {
		t.Fatal("Append failed", err, op)
	}

	ref := dot.Ref{Path: dot.NewRefPath([]string{"0"})}
	_, err := ref.Update(l)
	if err != dot.ErrPathInvalidated {
		t.Fatal("Update unexpected", err)
	}
}

func TestRefUpdateClientLog(t *testing.T) {
	clog := &dot.ClientLog{}
	deleteFirst := dot.Change{Splice: &dot.SpliceInfo{Before: []interface{}{nil}}}
	op := dot.Operation{ID: "local", Changes: []dot.Change{deleteFirst}}
	_, err := clog.AppendClientOperation(&dot.Log{}, op)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	// first attempt to transform [1] -- it should now be [0]
	ref := dot.Ref{Path: dot.NewRefPath([]string{"1"})}
	updated, err := ref.UpdateClient(clog)
	if err != nil || !reflect.DeepEqual(updated.Encode(), []string{"0"}) {
		t.Fatal("Unexpected", err, updated.Encode())
	}

	// now attempt to transform [0], it should fail
	ref = dot.Ref{Path: dot.NewRefPath([]string{"0"})}
	updated, err = ref.UpdateClient(clog)
	if err != dot.ErrPathInvalidated || updated != nil {
		t.Fatal("Unexpected..", err, updated.Encode())
	}
}
