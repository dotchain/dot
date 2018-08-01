// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"errors"
	"github.com/dotchain/dot/conv"
)

func (t Transformer) mergeMoveMove(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		return t.mergeMoveMoveSamePath(c1, c2)
	}

	if l == len(c1.Path) {
		return t.mergeMoveSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeMoveSubPath(c2, c1))
	}

	// no conflict
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeMoveMoveSamePath(c1, c2 Change) ([]Change, []Change) {
	offset1, count1, distance1, dest1 := c1.Move.Offset, c1.Move.Count, c1.Move.Distance, c1.Move.Dest()
	offset2, count2, distance2, dest2 := c2.Move.Offset, c2.Move.Count, c2.Move.Distance, c2.Move.Dest()

	if count1 == 0 || distance1 == 0 || count2 == 0 || distance2 == 0 {
		return []Change{c2}, []Change{c1}
	}

	// check if there are any overlaps in the ranges of any kind..
	offsets := []int{offset1, offset2, offset1 + count1, offset2 + count2, dest1, dest2}
	c1Split := t.splitMoveByOffsets(c1, offsets)
	c2Split := t.splitMoveByOffsets(c2, offsets)
	if len(c1Split) > 1 || len(c2Split) > 1 {
		return t.MergeChanges(c1Split, c2Split)
	}

	if offset1 == offset2 && count1 == count2 && distance1 == distance2 {
		// the two ops are the same, no conflicts
		return []Change{}, []Change{}
	}

	// exact same ranges..
	if offset1 == offset2 {
		// both are moving the same block.  let c1 win
		c2Mod := c2
		c2Mod.Move = &MoveInfo{Offset: offset2 + distance2, Count: count2, Distance: distance1 - distance2}
		return []Change{}, []Change{c2Mod}
	}

	return t.mergeMoveMoveSamePathSimple(c1, c2)
}

// mergeMoveMoveSamePathSimple merges two move Change structures that do not conflict
// i.e. no overlapping ranges, destination is not within any ranges etc.
// The algorithm here is just a blind enumeration of cases and the proof that it works is
// based on the fact that the tests work.
func (t Transformer) mergeMoveMoveSamePathSimple(c1, c2 Change) ([]Change, []Change) {
	if c1.Move.Offset > c2.Move.Offset {
		return t.swap(t.mergeMoveMoveSamePathSimple(c2, c1))
	}

	result1 := Change{Path: c2.Path, Move: &MoveInfo{}}
	*result1.Move = *c2.Move
	result2 := Change{Path: c1.Path, Move: &MoveInfo{}}
	*result2.Move = *c1.Move

	offset1, count1, dest1 := c1.Move.Offset, c1.Move.Count, c1.Move.Dest()
	offset2, count2, dest2 := c2.Move.Offset, c2.Move.Count, c2.Move.Dest()

	if dest1 < dest2 {
		if dest1 > offset2 {
			result1.Move.Offset -= count1
			result1.Move.Distance += count1
			result2.Move.Distance -= count2
		}
		if dest2 <= offset1 {
			result1.Move.Distance += count1
			result2.Move.Offset += count2
			result2.Move.Distance -= count2
		}
	} else {
		if dest1 > offset2 {
			result1.Move.Offset -= count1
		}
		if offset1 < dest2 && dest1 <= offset2 {
			result1.Move.Distance -= count1
		} else if dest2 <= offset1 && offset2 < dest1 {
			result1.Move.Distance += count1
		}

		if dest2 <= offset1 {
			result2.Move.Offset += count2
		}
		if dest2 <= offset1 && c2.Move.Offset < dest1 {
			result2.Move.Distance -= count2
		} else if offset1 < dest2 && dest1 <= offset2 {
			result2.Move.Distance += count2
		}
	}

	return []Change{result1}, []Change{result2}
}

// splitMoveByOffsets takes a move Change and returns a sequence of moves that
// have the same effect.  Any offsets provided in the input are guaranteed to not
// be within the range of any individaul move. The method sorts the unsortedOffsets arg
func (t Transformer) splitMoveByOffsets(change Change, unsortedOffsets []int) []Change {
	result := []Change{}

	path, offset, count, distance := change.Path, change.Move.Offset, change.Move.Count, change.Move.Distance
	sortedOffsets := append([]int{}, unsortedOffsets...)
	sortInts(sortedOffsets)

	if distance > 0 {
		// apply operations from the right so that we don't have to rebase offsets against prior moves
		lastOffset := offset + count
		for ii := len(sortedOffsets) - 1; ii >= 0; ii-- {
			off := sortedOffsets[ii]
			if off < offset+count && off > offset {
				move := &MoveInfo{Offset: off, Count: lastOffset - off, Distance: distance}
				result = append(result, Change{Path: path, Move: move})
				lastOffset = off
			}
		}
		move := &MoveInfo{Offset: offset, Count: lastOffset - offset, Distance: distance}
		result = append(result, Change{Path: path, Move: move})
	} else {
		lastOffset := offset
		for ii := 0; ii < len(sortedOffsets); ii++ {
			off := sortedOffsets[ii]
			if off < offset+count && off > offset {
				move := &MoveInfo{Offset: lastOffset, Count: off - lastOffset, Distance: distance}
				result = append(result, Change{Path: path, Move: move})
				lastOffset = off
			}
		}
		move := &MoveInfo{Offset: lastOffset, Count: offset + count - lastOffset, Distance: distance}
		result = append(result, Change{Path: path, Move: move})
	}
	return result
}

// mergeMoveSubPath is called with a move operation and any other operation.  The move operation
// must have a path that is a prefix of the other operation
func (t Transformer) mergeMoveSubPath(move, otherSubPath Change) ([]Change, []Change) {
	l := len(move.Path)
	offset, count, dest := move.Move.Offset, move.Move.Count, move.Move.Dest()

	if !conv.IsIndex(otherSubPath.Path[l]) {
		panic(errors.New("invalid array key, not a number"))
	}
	index := conv.ToIndex(otherSubPath.Path[l])

	// case 0: offset <= index < offset + count
	if offset <= index && index < offset+count {
		return otherSubPath.withUpdatedIndex(l, index+move.Move.Distance), []Change{move}
	}

	// case 1: index < dest < offset - no changes
	// case 2: dest <= index < offset - shift index right by count
	// case 3: dest < offset < offset + count <= index -- no changes
	// case 4: index < offset < dest - no changes
	// case 5: offset + count <= index < dest: shift index left by count
	// case 6: offset + count < dest <= index - no changes

	if dest <= index && index < offset {
		return otherSubPath.withUpdatedIndex(l, index+count), []Change{move}
	} else if offset+count <= index && index < dest {
		return otherSubPath.withUpdatedIndex(l, index-count), []Change{move}
	} else {
		return []Change{otherSubPath}, []Change{move}
	}
}

// simple insertion sort
func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
