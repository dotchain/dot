// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	_ "github.com/dotchain/dot/encoding/sparse"
	"testing"
)

func TestUtilsBuildImage_string(t *testing.T) {
	splice := &dot.SpliceInfo{After: "Hello"}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Splice: splice}}},
	}
	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, "Hello") {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_iterates(t *testing.T) {
	splice1 := &dot.SpliceInfo{After: "Hello"}
	splice2 := &dot.SpliceInfo{Offset: 5, After: " World"}
	ops := []dot.Operation{
		{ID: "one", Changes: []dot.Change{{Splice: splice1}}},
		{ID: "two", Changes: []dot.Change{{Splice: splice2}}},
	}
	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops[:1])

	if !u.AreSame(result.Model, "Hello") {
		t.Error("Unexpected output of splice", result.Model)
	}
	if result.BasisID != "one" {
		t.Error("Unexpected basis ID", result.BasisID)
	}

	result = u.BuildImage(result, ops[1:])
	if !u.AreSame(result.Model, "Hello World") {
		t.Error("Unexpected output of splice", result.Model)
	}
	if result.BasisID != "two" {
		t.Error("Unexpected basis ID", result.BasisID)
	}
}

func TestUtilsBuildImage_array(t *testing.T) {
	splice := &dot.SpliceInfo{After: []interface{}{"q", 42.0}}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Splice: splice}}},
	}
	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, splice.After) {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_sparse_array(t *testing.T) {
	sparse := map[string]interface{}{
		"dot:encoding": "SparseArray",
		"dot:encoded":  []interface{}{5, 122},
	}
	splice := &dot.SpliceInfo{After: sparse}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Splice: splice}}},
	}
	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, sparse) {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_map(t *testing.T) {
	set := &dot.SetInfo{Key: "hello", After: "world"}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Set: set}}},
	}
	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, map[string]interface{}{"hello": "world"}) {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_ignores_non_empty_path(t *testing.T) {
	badSplice := &dot.SpliceInfo{After: "qqq"}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Path: []string{"hello"}, Splice: badSplice}}},
		{Changes: []dot.Change{{Path: []string{}, Splice: &dot.SpliceInfo{After: "good"}}}},
	}

	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, "good") {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_ignores_bad_splice(t *testing.T) {
	badSplice := &dot.SpliceInfo{After: 42}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Splice: badSplice}}},
		{Changes: []dot.Change{{Splice: &dot.SpliceInfo{After: "good"}}}},
	}

	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, "good") {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_ignores_bad_move(t *testing.T) {
	badMove := &dot.MoveInfo{Count: 1, Distance: 1, Offset: 2}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Move: badMove}}},
		{Changes: []dot.Change{{Splice: &dot.SpliceInfo{After: "good"}}}},
	}

	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, "good") {
		t.Error("Unexpected output of splice", result.Model)
	}
}

func TestUtilsBuildImage_ignores_bad_range(t *testing.T) {
	badRange := &dot.RangeInfo{Count: 1, Offset: 2}
	ops := []dot.Operation{
		{Changes: []dot.Change{{Range: badRange}}},
		{Changes: []dot.Change{{Splice: &dot.SpliceInfo{After: "good"}}}},
	}

	u := dot.Utils(dot.Transformer{})
	result := u.BuildImage(nil, ops)

	if !u.AreSame(result.Model, "good") {
		t.Error("Unexpected output of splice", result.Model)
	}
}
