// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"strconv"
	"testing"
)

func generateMoves(input interface{}) []Change {
	ops := []Change{}
	for offset := 0; offset < arraySize(input)-1; offset++ {
		for count := 1; count < arraySize(input)-1-offset; count++ {
			for dest := 0; dest < arraySize(input); dest++ {
				if dest < offset {
					move := &MoveInfo{Offset: offset, Count: count, Distance: dest - offset}
					ops = append(ops, Change{Path: []string{}, Move: move})
				} else if dest > offset+count {
					move := &MoveInfo{Offset: offset, Count: count, Distance: dest - offset - count}
					ops = append(ops, Change{Path: []string{}, Move: move})
				}
			}
		}
	}
	return ops
}

func moveConflicts(left, right Change) bool {
	offsets := []int{left.Move.Offset, right.Move.Offset,
		left.Move.Offset + left.Move.Count,
		right.Move.Offset + right.Move.Count,
		left.Move.Dest(), right.Move.Dest(),
	}
	return len(x.splitMoveByOffsets(left, offsets)) > 1 || len(x.splitMoveByOffsets(right, offsets)) > 1
}

func TestMergeMoveMoveSamePathNoConflicts(t *testing.T) {
	input := "Hello yo"
	ops := generateMoves(input)
	testMergeFiltered(t, input, ops, ops, moveConflicts)
}

func TestMergeMoveMoveSamePathWithConflicts(t *testing.T) {
	input := "Hello yo"
	ops := generateMoves(input)
	onlyConflict := func(left, right Change) bool { return !moveConflicts(left, right) }
	testMergeFiltered(t, input, ops, ops, onlyConflict)
}

func TestMergeMoveMoveSubPath(t *testing.T) {
	input := "Hello World"
	inputOuter := []interface{}{input, input, input, input, input, input}

	outerMoves := []Change{
		{Move: &MoveInfo{Offset: 2, Count: 2, Distance: -1}},
		{Move: &MoveInfo{Offset: 2, Count: 2, Distance: 1}},
	}

	innerMoves := []Change{}
	for index := 0; index < len(inputOuter); index++ {
		innerMoves = append(innerMoves, Change{
			Path: []string{strconv.Itoa(index)},
			Move: &MoveInfo{Offset: 2, Count: 3, Distance: 1},
		})
	}

	testMerge(t, inputOuter, outerMoves, innerMoves)
}

func TestMergeMoveSplice(t *testing.T) {
	input := "Hello yo"
	testMerge(t, input, generateMoves(input), generateSplices(input, "yo"))
}
