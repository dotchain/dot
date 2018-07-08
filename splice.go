// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import "strconv"

func (t Transformer) mergeSpliceSplice(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		return t.mergeSpliceSpliceSamePath(c1, c2)
	}

	if l == len(c1.Path) {
		return t.mergeSpliceSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeSpliceSubPath(c2, c1))
	}

	// no conflicts because paths are different
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeSpliceSpliceSamePath(c1, c2 Change) ([]Change, []Change) {
	offset1 := c1.Splice.Offset
	offset2 := c2.Splice.Offset

	if offset1 <= offset2 {
		deleted1 := t.toArray(c1.Splice.Before).Count()
		deleted2 := t.toArray(c2.Splice.Before).Count()
		inserted1 := t.toArray(c1.Splice.After).Count()
		inserted2 := t.toArray(c2.Splice.After).Count()

		if deleted1 == 0 || offset1+deleted1 <= offset2 || deleted2 == 0 && offset1+deleted1 == offset2 {
			// no conflict but need to update c2 with new offset
			left := withSpliceOffset(c2, offset2+inserted1-deleted1)
			return []Change{left}, []Change{c1}
		} else if deleted2 == 0 && offset1 == offset2 {
			// also no conflict
			right := withSpliceOffset(c1, offset1+inserted2-deleted2)
			return []Change{c2}, []Change{right}
		} else if offset1+deleted1 >= offset2+deleted2 {
			// c2 is fully subsumed, has no effect.  Just need to update c1 with new Before
			before := Utils(t).applySplice(c1.Splice.Before, offset2-offset1, c2.Splice.Before, c2.Splice.After)
			splice := &SpliceInfo{Offset: offset1, Before: before, After: c1.Splice.After}
			return []Change{}, []Change{{Path: c1.Path, Splice: splice}}
		}
		// c2 is partially chopped.  just need to update offset on the left
		offset := offset1 + deleted1 - offset2
		alteredC2Splice := &SpliceInfo{
			Offset: offset1 + inserted1,
			Before: t.toArray(c2.Splice.Before).Slice(offset, deleted2-offset),
			After:  c2.Splice.After,
		}
		alteredC1Splice := &SpliceInfo{
			Offset: offset1,
			Before: t.toArray(c1.Splice.Before).Slice(0, offset2-offset1),
			After:  c1.Splice.After,
		}
		left := Change{Path: c1.Path, Splice: alteredC2Splice}
		right := Change{Path: c2.Path, Splice: alteredC1Splice}
		return []Change{left}, []Change{right}
	}

	// this part is symmetric.  just invert the order
	return t.swap(t.mergeSpliceSpliceSamePath(c2, c1))
}

// c1 is a splice.  c2 can be anything. first arg path is prefix of second
func (t Transformer) mergeSpliceSubPath(c1, c2 Change) ([]Change, []Change) {
	offset := c1.Splice.Offset
	deleted := t.toArray(c1.Splice.Before).Count()
	inserted := t.toArray(c1.Splice.After).Count()

	if index, err := strconv.Atoi(c2.Path[len(c1.Path)]); err != nil || index < offset {
		// No effective conflict.
		// TODO(rameshvk): log warning for Atoi error?
		return []Change{c2}, []Change{c1}
	} else if index >= offset && index < offset+deleted {
		// Transformed c1 should have the effect of c2
		opPath := append([]string{strconv.Itoa(index - offset)}, c2.Path[len(c1.Path)+1:]...)
		op := c2
		op.Path = opPath
		alteredC1Splice := &SpliceInfo{
			Offset: offset,
			Before: Utils(t).Apply(c1.Splice.Before, []Change{op}),
			After:  c1.Splice.After,
		}
		alteredC1 := Change{Path: c1.Path, Splice: alteredC1Splice}
		return []Change{}, []Change{alteredC1}
	} else {
		// just need to rewrite the path to be adjusted due to index change
		return c2.withUpdatedIndex(len(c1.Path), index-deleted+inserted), []Change{c1}
	}
}

func withSpliceOffset(change Change, offset int) Change {
	before := change.Splice.Before
	after := change.Splice.After
	return Change{Path: change.Path, Splice: &SpliceInfo{Offset: offset, Before: before, After: after}}
}
