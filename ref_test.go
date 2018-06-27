// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
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
