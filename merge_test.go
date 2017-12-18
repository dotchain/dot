// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"fmt"
	"github.com/dotchain/dot"
	"math/rand"
	"reflect"
	"testing"
)

func generateSplices(input, insert []interface{}) []dot.Change {
	changes := []dot.Change{}
	for offset := 0; offset < len(input); offset++ {
		for count := 0; count < len(input)-offset; count++ {
			// first append a deletion op
			before := input[offset : offset+count]
			deletion := &dot.SpliceInfo{Offset: offset, Before: before}
			insertion := &dot.SpliceInfo{Offset: offset, Before: before, After: insert}
			changes = append(changes, dot.Change{Splice: deletion}, dot.Change{Splice: insertion})
		}
	}
	return changes
}

func generateMoves(input []interface{}) []dot.Change {
	changes := []dot.Change{}
	for offset := 0; offset < len(input); offset++ {
		for count := 1; count < len(input)-1-offset; count++ {
			for dest := 0; dest < len(input); dest++ {
				distance := 0
				if dest < offset {
					distance = dest - offset
				} else if distance > offset+count {
					distance = dest - offset - count
				}
				if distance != 0 {
					move := &dot.MoveInfo{Offset: offset, Count: count, Distance: distance}
					changes = append(changes, dot.Change{Move: move})
				}
			}
		}
	}
	return changes
}

func generateSets(input map[string]interface{}, newVal interface{}) []dot.Change {
	changes := []dot.Change{{Set: &dot.SetInfo{Key: "New", After: newVal}}}
	for key, val := range input {
		deletion := &dot.SetInfo{Key: key, Before: val}
		update := &dot.SetInfo{Key: key, Before: val, After: newVal}
		changes = append(changes, dot.Change{Set: deletion}, dot.Change{Set: update})
	}
	return changes
}

func generateRanges(input []interface{}, changes []dot.Change) []dot.Change {
	ops := []dot.Change{}
	for offset := 0; offset < len(input); offset++ {
		for end := offset + 1; end < len(input); end++ {
			ops = append(ops, dot.Change{
				Range: &dot.RangeInfo{Offset: offset, Count: end - offset, Changes: changes},
			})
		}
	}
	return ops
}

// note input is of form "map of sequence of map of sequence"
func createTestChanges(input map[string]interface{}) []dot.Change {
	changes := []dot.Change{}

	// top level changes are only sets.
	topLevelChanges := generateSets(input, "Hello World")
	changes = append(changes, topLevelChanges...)

	seq := 0

	// next level changes are array changes.
	for key, val := range input {
		outer := val.([]interface{})
		path := []interface{}{key}
		outerLevelChanges := setPath(path, generateSplices(outer, []interface{}{key}))
		changes = append(changes, outerLevelChanges...)
		outerLevelChanges = setPath(path, generateMoves(outer))
		changes = append(changes, outerLevelChanges...)

		// outer range generation is made tricky by the fact that we need to be able
		// to generate three for each range.  We generate three changes:
		// a splice in inner1, a move in inner2, an explicit set of inner3 to a monotonic seq
		rangeSplice := &dot.SpliceInfo{Offset: 0, After: []interface{}{42}}
		rangeMove := &dot.MoveInfo{Offset: 0, Count: 1, Distance: 2}
		rangeSet := &dot.SetInfo{Key: "inner3", After: fmt.Sprintf("%#v", seq)}
		seq++
		rangeChanges := []dot.Change{
			{Path: []string{"inner1"}, Splice: rangeSplice},
			{Path: []string{"inner2"}, Move: rangeMove},
			{Set: rangeSet},
		}
		outerLevelChanges = setPath(path, generateRanges(outer, rangeChanges))
		changes = append(changes, outerLevelChanges...)

		// for each outer level change, there is also the inner level changes
		for index, val := range outer {
			innerm, _ := val.(map[string]interface{})
			path := []interface{}{key, index}
			innermChanges := setPath(path, generateSets(innerm, index))
			changes = append(changes, innermChanges...)

			// each innerm has in-turn got two arrays but this is too much
			// so we only test against the first key in innerm and we also
			// only pick a random set of changes in the inner moves
			// this is ok since the exact move operation or splice operation
			// should not affect the outcome
			inner := innerm["inner1"].([]interface{})
			path = []interface{}{key, index, "inner1"}
			innerChanges := pickSome(setPath(path, generateSplices(inner, []interface{}{"yo"})))
			changes = append(changes, innerChanges...)
			innerChanges = pickSome(setPath(path, generateMoves(inner)))
			changes = append(changes, innerChanges...)

			innerSplice := &dot.SpliceInfo{Offset: 0, After: []interface{}{92}}
			innerChanges = pickSome(setPath(path, generateRanges(inner, []dot.Change{{Splice: innerSplice}})))
			changes = append(changes, innerChanges...)
		}
	}

	return changes
}

func TestAllMerges(t *testing.T) {
	input := createTestInput()
	changes := createTestChanges(input)

	x := dot.Transformer{}
	for leftIndex, left := range changes {
		for rightIndex, right := range changes {
			t.Run(fmt.Sprintf("Test %#v %#v", leftIndex, rightIndex), func(t *testing.T) {
				//				lj, _ := json.Marshal(left)
				//				rj, _ := json.Marshal(right)
				//				fmt.Printf("Testing %s\n%s\n", string(lj), string(rj))
				leftAll, rightAll := []dot.Change{left}, []dot.Change{right}
				leftRemainder, rightRemainder := x.MergeChanges(leftAll, rightAll)
				//				lx, _ := json.Marshal(leftRemainder)
				//				rx, _ := json.Marshal(rightRemainder)
				//				fmt.Printf("Remainders %s\n%s\n", string(lx), string(rx))
				leftAll, rightAll = append(leftAll, leftRemainder...), append(rightAll, rightRemainder...)
				validateMerge(t, input, leftAll, rightAll)
			})
		}
	}
}

func validateMerge(t *testing.T, input interface{}, leftAll, rightAll []dot.Change) {
	resultLeft, resultRight := applyMany(input, leftAll), applyMany(input, rightAll)
	if !dot.Utils(x).AreSame(resultLeft, resultRight) {
		t.Errorf("Merge of %#v %#v failed with outputs %v %v\n", leftAll[0], rightAll[0], resultLeft, resultRight)
	}
}

func setPath(path []interface{}, changes []dot.Change) []dot.Change {
	spath := []string{}
	for _, part := range path {
		switch part := part.(type) {
		case int:
			spath = append(spath, fmt.Sprintf("%#v", part))
		case string:
			spath = append(spath, part)
		}
	}
	result := []dot.Change{}
	for _, change := range changes {
		change.Path = spath
		result = append(result, change)
	}
	return result
}

// input = "set of sequence of set of sequence" so that we can test all combinations of operations
func createTestInput() map[string]interface{} {
	leafIndex := 0
	leaf := func() string {
		leafIndex++
		return fmt.Sprintf("%#v", leafIndex)
	}

	input := map[string]interface{}{"outer1": []interface{}{}}
	seqCount := 8

	createInner := func() map[string]interface{} {
		inner := map[string]interface{}{"inner1": []interface{}{}, "inner2": []interface{}{}}
		inner1 := []interface{}{}
		inner2 := []interface{}{}
		for i := 0; i < seqCount; i++ {
			inner1 = append(inner1, []interface{}{leaf()})
			inner2 = append(inner2, []interface{}{leaf()})
		}
		inner["inner1"] = inner1
		inner["inner2"] = inner2
		return inner
	}

	outer1 := []interface{}{}
	for i := 0; i < seqCount; i++ {
		outer1 = append(outer1, createInner())
	}
	input["outer1"] = outer1
	return input
}

func pickSome(changes []dot.Change) []dot.Change {
	// we pick 4 tests to be representative
	count := 4
	result := []dot.Change{}
	for i := 0; i < count; i++ {
		result = append(result, changes[rand.Int()%len(changes)])
	}
	return result
}

func TestEmptyMerges(t *testing.T) {
	splice := dot.Change{Splice: &dot.SpliceInfo{}}
	move := dot.Change{Move: &dot.MoveInfo{}}
	rangex := dot.Change{Range: &dot.RangeInfo{}}
	set := dot.Change{Set: &dot.SetInfo{}}

	x := dot.Transformer{}
	for _, ch := range []dot.Change{splice, move, rangex, set} {
		l, r := []dot.Change{ch}, []dot.Change{{}}
		lx, rx := x.MergeChanges(l, r)
		if !reflect.DeepEqual(lx, r) || !reflect.DeepEqual(rx, l) {
			t.Error("Mismatch", l, r, lx, rx)
		}
		rx, lx = x.MergeChanges(r, l)
		if !reflect.DeepEqual(lx, r) || !reflect.DeepEqual(rx, l) {
			t.Error("Mismatch", l, r, lx, rx)
		}
	}
}

// MergeOperation is indirectly tested via Log:AppendOperation tests
// but this particular test is only for the case that Log:AppendOperation
// does not cover -- where left is of size one and right is longer
func TestMergeOperationRightChain(t *testing.T) {
	// left = insert "hello" at offset 5
	left1 := dot.Change{Splice: &dot.SpliceInfo{5, "", "hello"}}

	// right = insert "first" at offset 5 and then insert "second" at offset 10
	right1 := dot.Change{Splice: &dot.SpliceInfo{5, "", "first"}}
	right2 := dot.Change{Splice: &dot.SpliceInfo{10, "", "second"}}

	l := []dot.Operation{{Changes: []dot.Change{left1}}}
	r := []dot.Operation{{Changes: []dot.Change{right1}}, {Changes: []dot.Change{right2}}}
	lx, rx := (dot.Transformer{}).MergeOperations(l, r)

	// expect rx to not change
	if !reflect.DeepEqual(rx[0].Changes[0], left1) {
		t.Error("Unexpected rx", rx)
	}

	// but expect lx to all be offset by 5 more
	changes := append(lx[0].Changes, lx[1].Changes...)
	right1.Splice.Offset += 5
	right2.Splice.Offset += 5
	if !reflect.DeepEqual(changes, []dot.Change{right1, right2}) {
		t.Error("Unexpected lx", lx)
	}
}

func TestMergeDuplicateOperations(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok && err == dot.ErrMergeWithSelf {
				return
			}
		}
		t.Error("MergeOperations did not panic(ErrMergeWithSelf)")
	}()

	id := "one"
	(dot.Transformer{}).MergeOperations([]dot.Operation{{ID: id}}, []dot.Operation{{ID: id}})
}

func TestInvalidMoveMoveMerge(t *testing.T) {
	defer func() {
		expected := `strconv.Atoi: parsing "non-integer": invalid syntax`
		if r := recover(); r != nil {
			if err, ok := r.(error); ok && err.Error() == expected {
				return
			}
		}
		t.Error("MergeOperations did not panic(strconv.Atoi error)")
	}()

	path1 := []string{"array1"}
	move1 := &dot.MoveInfo{5, 10, 2}
	path2 := []string{"array1", "non-integer"}
	move2 := &dot.MoveInfo{10, 10, 1}
	ops := []dot.Operation{
		{ID: "one", Changes: []dot.Change{{Path: path1, Move: move1}}},
		{ID: "two", Changes: []dot.Change{{Path: path2, Move: move2}}},
	}
	(dot.Transformer{}).MergeOperations(ops[:1], ops[1:])
}

func TestInvalidRangeSubPathMerge1(t *testing.T) {
	defer func() {
		expected := `strconv.Atoi: parsing "non-integer": invalid syntax`
		if r := recover(); r != nil {
			if err, ok := r.(error); ok && err.Error() == expected {
				return
			}
		}
		t.Error("MergeOperations did not panic(strconv.Atoi error)")
	}()

	path1 := []string{"array1"}
	move1 := &dot.MoveInfo{5, 10, 2}
	path2 := []string{"array1", "non-integer"}
	range2 := &dot.RangeInfo{10, 1, nil}
	ops := []dot.Operation{
		{ID: "one", Changes: []dot.Change{{Path: path1, Move: move1}}},
		{ID: "two", Changes: []dot.Change{{Path: path2, Range: range2}}},
	}
	(dot.Transformer{}).MergeOperations(ops[:1], ops[1:])
}

func TestInvalidRangeSubPathMerge2(t *testing.T) {
	defer func() {
		expected := `strconv.Atoi: parsing "non-integer": invalid syntax`
		if r := recover(); r != nil {
			if err, ok := r.(error); ok && err.Error() == expected {
				return
			}
		}
		t.Error("MergeOperations did not panic(strconv.Atoi error)")
	}()

	path1 := []string{"array1", "non-integer"}
	move1 := &dot.MoveInfo{5, 10, 2}
	path2 := []string{"array1"}
	range2 := &dot.RangeInfo{10, 1, nil}
	ops := []dot.Operation{
		{ID: "one", Changes: []dot.Change{{Path: path1, Move: move1}}},
		{ID: "two", Changes: []dot.Change{{Path: path2, Range: range2}}},
	}
	(dot.Transformer{}).MergeOperations(ops[:1], ops[1:])
}

func TestInvalidRangeSetMerge(t *testing.T) {
	defer func() {
		expected := "cannot apply range and set with the same path"
		if r := recover(); r != nil {
			if err, ok := r.(string); ok && err == expected {
				return
			}
		}
		t.Error("MergeOperations did not panic")
	}()

	path1 := []string{"array1"}
	range1 := &dot.RangeInfo{5, 10, nil}
	path2 := []string{"array1"}
	set2 := &dot.SetInfo{"somekey", nil, "hello"}
	ops := []dot.Operation{
		{ID: "one", Changes: []dot.Change{{Path: path1, Range: range1}}},
		{ID: "two", Changes: []dot.Change{{Path: path2, Set: set2}}},
	}
	(dot.Transformer{}).MergeOperations(ops[:1], ops[1:])
}
