// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import "strconv"

func (t Transformer) mergeRangeRange(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		return t.mergeRangeRangeSamePath(c1, c2)
	}

	if l == len(c1.Path) {
		return t.mergeRangeSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeRangeSubPath(c2, c1))
	}

	// no conflicts because paths are different
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeRangeRangeSamePath(c1, c2 Change) ([]Change, []Change) {
	start := max(c1.Range.Offset, c2.Range.Offset)
	end := min(c1.Range.Offset+c1.Range.Count, c2.Range.Offset+c2.Range.Count)
	if end <= start {
		return []Change{c2}, []Change{c1}
	}

	changes1, changes2 := c1.Range.Changes, c2.Range.Changes
	merged1, merged2 := t.MergeChanges(changes1, changes2)
	replacement1 := RangeInfo{Offset: start, Count: end - start, Changes: merged1}
	replacement2 := RangeInfo{Offset: start, Count: end - start, Changes: merged2}
	return t.replaceConflictingRange(c2, replacement1), t.replaceConflictingRange(c1, replacement2)
}

func (t Transformer) replaceConflictingRange(c1 Change, r RangeInfo) []Change {
	replacement := Change{Path: c1.Path, Range: &r}
	result := []Change{replacement}

	if c1.Range.Offset < r.Offset {
		left := *c1.Range
		left.Count = r.Offset - left.Offset
		if left.Count > 0 {
			result = []Change{{Path: c1.Path, Range: &left}, replacement}
		}
	}

	offset, count := r.Offset+r.Count, c1.Range.Offset+c1.Range.Count-r.Offset-r.Count
	if count > 0 {
		right := RangeInfo{Offset: offset, Count: count, Changes: c1.Range.Changes}
		result = append(result, Change{Path: c1.Path, Range: &right})
	}
	return result
}

// mergeRangeSubPath expects c1 to be a Range mutation and c2 to be any mutation with
// the path referring to within an element in the array that c1 is working on.
// in particular, c2 does not have to be a range mutation
func (t Transformer) mergeRangeSubPath(c1, c2 Change) ([]Change, []Change) {
	index, err := strconv.Atoi(c2.Path[len(c1.Path)])
	if err != nil {
		panic(err)
	}

	if index < c1.Range.Offset || index >= c1.Range.Offset+c1.Range.Count {
		return []Change{c2}, []Change{c1}
	}

	// convert c2 into a Range mutation with same path as c1 and then
	// we can simply use mergeRangeRangeSamePath!

	c2.Path = c2.Path[len(c1.Path)+1:]
	c2Range := RangeInfo{Offset: index, Count: 1, Changes: []Change{c2}}
	return t.mergeRangeRangeSamePath(c1, Change{Path: c1.Path, Range: &c2Range})
}

func (t Transformer) mergeRangeSplice(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		return t.mergeRangeSpliceSamePath(c1, c2)
	}

	if l == len(c1.Path) {
		return t.mergeRangeSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeSpliceSubPath(c2, c1))
	}

	// no conflicts because paths are different
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeRangeSpliceSamePath(c1, c2 Change) ([]Change, []Change) {
	r, splice := c1.Range, c2.Splice
	// Case 1: splice is to the right, no conflicts
	if splice.Offset > r.Offset+r.Count {
		return []Change{c2}, []Change{c1}
	}

	changes := r.Changes
	beforeCount, afterCount := t.toArray(splice.Before).Count(), t.toArray(splice.After).Count()

	// Case 2: splice is to the left, just need to update offset of range
	if splice.Offset+beforeCount <= r.Offset {
		// no conflict but may need to update Offset
		if beforeCount != afterCount {
			newRange := *r
			newRange.Offset += afterCount - beforeCount
			c1.Range = &newRange
		}
		return []Change{c2}, []Change{c1}
	}

	// Case 3: splice covers range, ignore range but update before of splice
	if splice.Offset <= r.Offset && splice.Offset+beforeCount >= r.Offset+r.Count {
		newSplice := *splice
		newSplice.Before = Utils(t).applyRange(splice.Before, r.Offset-splice.Offset, r.Count, changes)
		c2.Splice = &newSplice
		return []Change{c2}, []Change{}
	}

	// Case 4: range covers splice, update range size as well as "After" of splice
	if splice.Offset >= r.Offset && splice.Offset+beforeCount <= r.Offset+r.Count {
		newSplice := *splice
		before, ok1 := Utils(t).tryApplyRange(splice.Before, 0, beforeCount, changes)
		after, ok2 := Utils(t).tryApplyRange(splice.After, 0, afterCount, changes)

		if ok1 && ok2 {
			newSplice.Before = before
			newSplice.After = after
			c2.Splice = &newSplice

			newRange := *r
			newRange.Count = r.Count + afterCount - beforeCount
			c1.Range = &newRange
			return []Change{c2}, []Change{c1}
		}
		// if there was an error with the apply, the range mutation is incompatible
		// with the final state of the splice.  To make things converge consistently,
		// we treat this case similar to case 5 and simply ignore the range mutation
		// for the intersection alone.
		// This is a spec issue.
	}

	// Case 5: Partial overlap of range and splice, left or right truncate range
	// Follow-through from case-4 can also be the case where splice covers range

	start, end := max(splice.Offset, r.Offset), min(splice.Offset+beforeCount, r.Offset+r.Count)

	out1 := []Change{}

	left := *r
	left.Count = start - r.Offset
	if left.Count > 0 {
		out1 = append(out1, Change{Path: c1.Path, Range: &left})
	}

	right := *r
	right.Count = r.Offset + r.Count - end
	if right.Count > 0 {
		right.Offset = end + afterCount - beforeCount
		out1 = append(out1, Change{Path: c1.Path, Range: &right})
	}

	// new splice calculation is only about updating "before" array
	newSplice := *splice
	newSplice.Before = Utils(t).applyRange(splice.Before, start-splice.Offset, end-start, changes)
	c2.Splice = &newSplice

	return []Change{c2}, out1
}

func (t Transformer) mergeRangeSet(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		panic("cannot apply range and set with the same path")
	}

	if l == len(c1.Path) {
		return t.mergeRangeSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeSetSubPath(c2, c1))
	}

	// no conflicts because paths are different
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeRangeMove(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		return t.mergeRangeMoveSamePath(c1, c2)
	}

	if l == len(c1.Path) {
		return t.mergeRangeSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeMoveSubPath(c2, c1))
	}

	// no conflicts because paths are different
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeRangeMoveSamePath(c1, c2 Change) ([]Change, []Change) {
	intervals := c2.Move.TransformInterval(c1.Range.Offset, c1.Range.Count)
	if intervals == nil {
		return []Change{c2}, []Change{c1}
	}
	out1 := []Change{}
	for _, interval := range intervals {
		r := RangeInfo{Offset: interval[0], Count: interval[1], Changes: c1.Range.Changes}
		out1 = append(out1, Change{Path: c1.Path, Range: &r})
	}
	return []Change{c2}, out1
}
