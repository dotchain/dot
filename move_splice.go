// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

func (t Transformer) mergeMoveSplice(c1, c2 Change) ([]Change, []Change) {
	if c1.Move.Count == 0 || c1.Move.Distance == 0 {
		return []Change{c2}, []Change{}
	}

	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		return t.mergeMoveSpliceSamePath(c1, c2)
	}

	if l == len(c1.Path) {
		return t.mergeMoveSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeSpliceSubPath(c2, c1))
	}

	// no conflict
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeMoveSpliceSamePath(c1, c2 Change) ([]Change, []Change) {
	offset1, count1, distance1, dest1 := c1.Move.Offset, c1.Move.Count, c1.Move.Distance, c1.Move.Dest()
	offset2, count2, inserted2 := c2.Splice.Offset, t.toArray(c2.Splice.Before).Count(), t.toArray(c2.Splice.After).Count()

	newSpliceOffset := offset2
	newSpliceSize := count2
	newMoveOffset := offset1
	newMoveCount := count1
	newDistance := distance1

	switch {
	case offset1 <= offset2 && offset1+count1 >= offset2+count2 && count2 > 0:
		// splice is within move
		newMoveCount = count1 + inserted2 - count2
		newSpliceOffset = offset2 + distance1

	case offset1+count1 <= offset2:
		// move is fully left of splice
		if dest1 <= offset2 {
			return []Change{c2}, []Change{c1}
		}
		if dest1 < offset2+count2 {
			// split the splice around dest1 and recurse
			splices := t.splitSplice(c2, dest1-offset2)
			return t.MergeChanges([]Change{c1}, splices)
		}
		newSpliceOffset = offset2 - count1
		newDistance = distance1 + inserted2 - count2

	case offset1 < offset2 && offset1+count1 > offset2:
		// move is left-of + overlapping splice: split splice
		moves := t.splitMoveByOffsets(c1, []int{offset2, offset2 + count2})
		return t.MergeChanges(moves, []Change{c2})

	case offset1 == offset2 && count1 < count2:
		// too complex to have part of the splice in and out.
		// replace splice by a simple splice for the left part and make the right part a delete
		splices := t.splitSplice(c2, count1)
		return t.MergeChanges([]Change{c1}, splices)

	case offset1 < offset2+count2:
		// move fully inside splice or has overlap: split splice
		splices := t.splitSplice(c2, offset1-offset2)
		return t.MergeChanges([]Change{c1}, splices)

	default: // offset2 + count2 <= offset1
		if dest1 <= offset2 {
			// dest1 <= offset2 < offset1
			newSpliceOffset = offset2 + count1
			newMoveOffset = offset1 + inserted2 - count2
			newDistance = distance1 - (inserted2 - count2)
		} else if dest1 < offset2+count2 {
			// ugh.  lets split the splice into two parts and do this one by one
			splices := t.splitSplice(c2, dest1-offset2)
			return t.MergeChanges([]Change{c1}, splices)
		} else {
			// dest1 >= offset2 + count2: just need to shift move offset but everything else is ok
			newMoveOffset = offset1 + inserted2 - count2
		}
	}

	newSpliceInfo := &SpliceInfo{
		Offset: newSpliceOffset,
		Before: t.toArray(c2.Splice.Before).Slice(0, newSpliceSize),
		After:  c2.Splice.After,
	}
	newMoveInfo := &MoveInfo{Offset: newMoveOffset, Count: newMoveCount, Distance: newDistance}
	result1 := Change{Path: c2.Path, Splice: newSpliceInfo}
	result2 := Change{Path: c1.Path, Move: newMoveInfo}
	return []Change{result1}, []Change{result2}
}

func (t Transformer) splitSplice(change Change, size int) []Change {
	return []Change{
		{
			Path: change.Path,
			Splice: &SpliceInfo{
				Offset: change.Splice.Offset,
				Before: t.toArray(change.Splice.Before).Slice(0, size),
				After:  change.Splice.After,
			},
		},
		{
			Path: change.Path,
			Splice: &SpliceInfo{
				Offset: change.Splice.Offset + t.toArray(change.Splice.After).Count(),
				Before: t.toArray(change.Splice.Before).Slice(size, t.toArray(change.Splice.Before).Count()-size),
				After:  nil,
			},
		},
	}
}
