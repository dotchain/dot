// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"github.com/dotchain/dot/encoding"
	"testing"
)

type applyPathSpliceTestCase struct {
	obj        interface{}
	path       []string
	start, end int
	before     interface{}
	replace    interface{}
	expected   interface{}
}

func (a applyPathSpliceTestCase) validate(t *testing.T) {
	splice := &SpliceInfo{
		Offset: a.start,
		Before: a.before,
		After:  a.replace,
	}
	actual := Utils(x).Apply(a.obj, []Change{{Path: a.path, Splice: splice}})
	ensureSame(t, a.expected, actual)
}

func (a applyPathSpliceTestCase) validateAllCases(t *testing.T) {
	// validate basic test
	a.validate(t)

	// path refers to index in an array
	applyPathSpliceTestCase{
		obj:   []interface{}{a.obj},
		path:  append([]string{"0"}, a.path...),
		start: a.start, end: a.end,
		before:   a.before,
		replace:  a.replace,
		expected: []interface{}{a.expected},
	}.validate(t)

	// path refers to keys in a map
	applyPathSpliceTestCase{
		obj:   map[string]interface{}{"key": a.obj, "key2": "yo"},
		path:  append([]string{"key"}, a.path...),
		start: a.start, end: a.end,
		before:   a.before,
		replace:  a.replace,
		expected: map[string]interface{}{"key": a.expected, "key2": "yo"},
	}.validate(t)
}

func TestApplyPathToStrings(t *testing.T) {
	applyPathSpliceTestCase{
		obj:      "hello world",
		path:     []string{},
		start:    2,
		end:      5,
		before:   "llo ",
		replace:  "LLO ",
		expected: "heLLO world",
	}.validateAllCases(t)
}

func TestApplyPathToArrays(t *testing.T) {
	applyPathSpliceTestCase{
		obj:      []interface{}{1, 2, 44, 33, 21},
		path:     []string{},
		start:    1,
		end:      3,
		before:   []interface{}{2, 44},
		replace:  []interface{}{3, 55},
		expected: []interface{}{1, 3, 55, 33, 21},
	}.validateAllCases(t)
}

func TestObjectSetOnMap(t *testing.T) {
	input := map[string]interface{}{"hello": "world"}

	// modify existing key
	output := encoding.Get(input).Set("hello", "new world")
	expected := map[string]interface{}{"hello": "new world"}
	ensureSame(t, expected, output)

	// modify a different key
	output = encoding.Get(input).Set("new", "wave")
	expected = map[string]interface{}{"hello": "world", "new": "wave"}
	ensureSame(t, expected, output)

	// delete an existing key
	output = encoding.Get(input).Set("hello", nil)
	expected = map[string]interface{}{}
	ensureSame(t, expected, output)

	// zero an existing key
	output = encoding.Get(input).Set("hello", "")
	expected = map[string]interface{}{"hello": ""}
	ensureSame(t, expected, output)
}

func ensureSame(t *testing.T, expected interface{}, actual interface{}) {
	if !Utils(x).AreSame(expected, actual) {
		t.Errorf("Mismatch.  Expected %#v but got %#v", expected, actual)
	}
}

func objectMove(input interface{}, offset, count, distance int) interface{} {
	data := encoding.Get(input)
	slice := data.Slice(offset, count)
	return data.Splice(offset, slice, nil).Splice(offset+distance, nil, slice)
}

func TestObjectMoveOnString(t *testing.T) {
	input := "hello cruel world"

	// no op
	ensureSame(t, input, objectMove(input, 6, 0, 2))
	ensureSame(t, input, objectMove(input, 6, 2, 0))

	// move to right
	expected := "hello world cruel"
	actual := objectMove(input, 5, 6, 6)
	ensureSame(t, expected, actual)

	// move to left
	expected = "cruel hello world"
	actual = objectMove(input, 6, 6, -6)
	ensureSame(t, expected, actual)
}

func TestObjectRangeApplySplice(t *testing.T) {
	input := []interface{}{[]interface{}{0, 1}, []interface{}{2, 3}, []interface{}{4, 5}, []interface{}{6, 7}}

	splice := SpliceInfo{Offset: 0, Before: nil, After: []interface{}{42}}
	change := Change{Path: []string{}, Splice: &splice}
	change = Change{Range: &RangeInfo{Offset: 1, Count: 2, Changes: []Change{change, change}}}
	actual := Utils(x).Apply(input, []Change{change})

	// expect two 42s inserted in the middle two arrays
	expected := []interface{}{[]interface{}{0, 1}, []interface{}{42, 42, 2, 3}, []interface{}{42, 42, 4, 5}, []interface{}{6, 7}}
	ensureSame(t, expected, actual)
}

func TestEmptyApply(t *testing.T) {
	input := "hello world"
	change := Change{}
	actual1 := Utils(x).Apply(input, []Change{change})
	actual2 := Utils(x).Apply(input, nil)

	// expect both to be same as input
	ensureSame(t, input, actual1)
	ensureSame(t, input, actual2)
}

func TestFailedTryApply(t *testing.T) {
	// none of these failures should panic
	failures := [][]interface{}{
		// initial, change
		{"hello", Change{Splice: &SpliceInfo{0, "q", "bee"}}},
		{map[string]interface{}{}, Change{Set: &SetInfo{"hello", "world", "new_world"}}},
		{map[string]interface{}{"hello": 42}, Change{Set: &SetInfo{"hello", "world", "new_world"}}},
		{map[string]interface{}{"hello": 42}, Change{Set: &SetInfo{"hello", nil, "new_world"}}},
		{map[string]interface{}{"hello": 42}, Change{Path: []string{"boom"}, Set: &SetInfo{"hello", "world", "new_world"}}},
		{map[string]interface{}{"hello": 42}, Change{Path: []string{"hello"}, Set: &SetInfo{"hello", "world", "new_world"}}},
		{map[string]interface{}{"hello": 42}, Change{Path: []string{"hello"}, Move: &MoveInfo{0, 1, 1}}},
		{42, Change{Path: []string{"boom"}}},
	}

	u := Utils(Transformer{})
	for _, failure := range failures {
		initial, changes := failure[0], []Change{failure[1].(Change)}
		if result, ok := u.TryApply(initial, changes); ok {
			t.Error("TryApply false positive", initial, changes[0], result)
		}
		if result := u.Apply(initial, changes); result != nil {
			t.Error("TryApply false positive", initial, changes[0], result)
		}
	}
}
