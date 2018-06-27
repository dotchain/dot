// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"strconv"
	"testing"
)

func generateRanges(input, changes interface{}) []Change {
	ops := []Change{}
	for offset := 0; offset < arraySize(input); offset++ {
		for end := offset + 1; end < arraySize(input); end++ {
			ops = append(ops, Change{
				Range: &RangeInfo{Offset: offset, Count: end - offset, Changes: changes.([]Change)},
			})
		}
	}
	return ops
}

func TestMergeRangeRangeSamePathNoConflict(t *testing.T) {
	input := array(array(0, 1), array(1, 2), array(2, 3), array(3, 4), array(4, 5))
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(), After: array(42)}}
	change2 := Change{Splice: &SpliceInfo{Offset: 2, Before: array(), After: array(43)}}
	changes := []Change{change1, change2}
	ops := generateRanges(input, changes)
	testMerge(t, input, ops, ops)
}

func TestMergeRangeRangeSamePathConflicting1(t *testing.T) {
	input := array(array(0, 1), array(0, 2), array(0, 3), array(0, 4), array(0, 5))
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(0), After: array(42)}}
	changes := []Change{change1}
	ops := generateRanges(input, changes)
	testMerge(t, input, ops, ops)
}

func TestMergeRangeRangeSamePathConflicting2(t *testing.T) {
	input := array(array(0, 1), array(0, 2), array(0, 3), array(0, 4), array(0, 5))
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(0), After: array(42)}}
	change2 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(0), After: array(43)}}
	ops1 := generateRanges(input, []Change{change1})
	ops2 := generateRanges(input, []Change{change2})
	testMerge(t, input, ops1, ops2)
}

func TestMergeRangeSpliceSamePath(t *testing.T) {
	input := array(array(0, 1), array(0, 2), array(0, 3), array(0, 4), array(0, 5))
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(0), After: array(42)}}
	rangeOps := generateRanges(input, []Change{change1})
	spliceOps := generateSplices(input, array(array(0, 400), array(0, 500)))
	testMerge(t, input, rangeOps, spliceOps)
}

func TestMergeRangeMoveSamePath(t *testing.T) {
	input := array(array(0, 1), array(0, 2), array(0, 3), array(0, 4), array(0, 5))
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(0), After: array(42)}}
	rangeOps := generateRanges(input, []Change{change1})
	moveOps := generateMoves(input)
	testMerge(t, input, rangeOps, moveOps)
}

func TestMergeRangeRangeSubPath(t *testing.T) {
	input := array()
	for kk := 0; kk < 5; kk++ {
		input = append(input, array(array(0, 0), array(kk+1, kk+2), array(kk+2, kk+3)))
	}
	// input: [ [[0, 0], [1, 2], [2, 3]], [[0, 0], [2, 3], [3, 4]], + 3 more times]

	// outer range modifies the [0, 0] element into [42, 42]
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(array(0, 0)), After: array(array(42, 42))}}
	outer := generateRanges(input, []Change{change1})

	// inner range will take input [[0, 0], [kk+1, kk+2], [kk+2, kk+3]] and path [kk]
	// it will simply splice kk + 1 into 43 (always offset 1)
	inner := []Change{}
	for kk := range input {
		path := []string{strconv.Itoa(kk)}
		innerSplice := &SpliceInfo{Offset: 0, Before: array(kk + 1), After: array(43)}
		innerChanges := []Change{{Splice: innerSplice}}
		innerChange := Change{Path: path, Range: &RangeInfo{Offset: 1, Count: 1, Changes: innerChanges}}
		inner = append(inner, innerChange)
	}

	testMerge(t, input, outer, inner)
}

func TestMergeRangeSpliceSubPath(t *testing.T) {
	input := array()
	for kk := 0; kk < 5; kk++ {
		input = append(input, array(array(0, 0), array(kk+1, kk+2), array(kk+2, kk+3)))
	}
	// input: [ [[0, 0], [1, 2], [2, 3]], [[0, 0], [2, 3], [3, 4]], + 3 more times]

	// outer range modifies the [0, 0] element into [42, 42]
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(array(0, 0)), After: array(array(42, 42))}}
	outer := generateRanges(input, []Change{change1})

	inner := []Change{}

	// inner change is a splice with path = [kk] or a splice with path = [kk, 0] (which is always 0, 0)
	for kk := range input {
		path := []string{strconv.Itoa(kk)}
		splice1 := SpliceInfo{Offset: 1, Before: array(array(kk+1, kk+2)), After: array(array(43, 43))}
		change1 := Change{Path: path, Splice: &splice1}
		splice2 := SpliceInfo{Offset: 0, Before: array(0), After: array(44)}
		change2 := Change{Path: append(path, "0"), Splice: &splice2}
		inner = append(inner, change1, change2)
	}

	testMerge(t, input, outer, inner)
}

func TestMergeRangeMoveSubPath(t *testing.T) {
	input := array()
	for kk := 0; kk < 5; kk++ {
		input = append(input, array(array(0, 0), array(kk+1, kk+2), array(kk+2, kk+3)))
	}
	// input: [ [[0, 0], [1, 2], [2, 3]], [[0, 0], [2, 3], [3, 4]], + 3 more times]

	// outer range modifies the [0, 0] element into [42, 42]
	change1 := Change{Splice: &SpliceInfo{Offset: 0, Before: array(array(0, 0)), After: array(array(42, 42))}}
	outer := generateRanges(input, []Change{change1})

	inner := []Change{}

	// inner change is a move with path = [kk] or a move with path = [kk, n]
	for kk := range input {
		path := []string{strconv.Itoa(kk)}
		move1 := MoveInfo{Offset: 1, Count: 1, Distance: 1}
		change1 := Change{Path: path, Move: &move1}

		move2 := MoveInfo{Offset: 1, Count: 1, Distance: -1}
		change2 := Change{Path: append(path, "0"), Move: &move2}
		inner = append(inner, change1, change2)
	}

	testMerge(t, input, outer, inner)
}

func TestMergeSpliceRangeSubPath(t *testing.T) {
	input := array()
	for kk := 0; kk < 5; kk++ {
		input = append(input, array(array(0, 0), array(kk+1, kk+2), array(kk+2, kk+3)))
	}
	// input: [ [[0, 0], [1, 2], [2, 3]], [[0, 0], [2, 3], [3, 4]], + 3 more times]

	outer := generateSplices(input, array(42))

	inner := []Change{}
	for kk, elt := range input {
		path := []string{strconv.Itoa(kk)}
		// each elt is of the form [[0, 0], [kk + 1, kk +2], [kk + 2, kk + 3]]
		change := Change{Splice: &SpliceInfo{Offset: 1, Before: array(), After: array(43)}}
		for _, op := range generateRanges(elt, []Change{change}) {
			op.Path = path
			inner = append(inner, op)
		}
	}

	testMerge(t, input, outer, inner)
}

func TestMergeMoveRangeSubPath(t *testing.T) {
	input := array()
	for kk := 0; kk < 5; kk++ {
		input = append(input, array(array(0, 0), array(kk+1, kk+2), array(kk+2, kk+3)))
	}
	// input: [ [[0, 0], [1, 2], [2, 3]], [[0, 0], [2, 3], [3, 4]], + 3 more times]

	outer := generateMoves(input)

	inner := []Change{}
	for kk, elt := range input {
		path := []string{strconv.Itoa(kk)}
		// each elt is of the form [[0, 0], [kk + 1, kk +2], [kk + 2, kk + 3]]
		change := Change{Splice: &SpliceInfo{Offset: 1, Before: array(), After: array(43)}}
		for _, op := range generateRanges(elt, []Change{change}) {
			op.Path = path
			inner = append(inner, op)
		}
	}

	testMerge(t, input, outer, inner)
}

func TestMergeRangeSetSubPath(t *testing.T) {
	set := map[string]interface{}{"hello": "world", "alpha": "beta"}
	input := []interface{}{}
	for kk := 0; kk < 2; kk++ {
		index := strconv.Itoa(kk)
		elt := map[string]interface{}{"index": index}
		for k, v := range set {
			elt[k] = v
		}
		input = append(input, elt)
	}

	rangeOps := generateRanges(input, generateSets(input, "blimey"))

	setOps := []Change{}
	for kk, elt := range input {
		path := []string{strconv.Itoa(kk)}
		elt := elt.(map[string]interface{})
		ops := generateSets(elt, "checkmate!")
		for _, op := range ops {
			op.Path = path
			setOps = append(setOps, op)
		}
	}

	testMerge(t, input, rangeOps, setOps)
}

func TestMergeSetRangeSubPath(t *testing.T) {
	input := map[string]interface{}{
		"hello": array(array(0), array(1), array(2), array(3)),
		"world": array(array(0), array(4), array(5), array(6)),
	}
	setOps := generateSets(input, array(42, 42, 42, 42))

	rangeOps := []Change{}
	for k, v := range input {
		elt := v.([]interface{})
		splice := SpliceInfo{Offset: 0, Before: nil, After: array(42)}
		rr := generateRanges(elt, []Change{{Splice: &splice}})
		for _, op := range rr {
			op.Path = []string{k}
			rangeOps = append(rangeOps, op)
		}
	}

	testMerge(t, input, setOps, rangeOps)
}
