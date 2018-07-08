// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

func (t Transformer) mergeSetSet(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && l == len(c2.Path) {
		if c1.Set.Key != c2.Set.Key {
			// no conflict
			return []Change{c2}, []Change{c1}
		}

		// same path.  last writer wins, assume the left op is the new op
		// TODO:  would be nice if the Change came with a timestamp itself but the timestamp
		// is on the envelop right now :(
		alteredC1Set := &SetInfo{Key: c1.Set.Key, Before: c2.Set.After, After: c1.Set.After}
		alteredC1 := Change{Path: c1.Path, Set: alteredC1Set}
		return []Change{}, []Change{alteredC1}
	}

	if l == len(c1.Path) && c1.Set.Key == c2.Path[l] {
		return t.mergeSetSubPath(c1, c2)
	}

	if l == len(c2.Path) && c2.Set.Key == c1.Path[l] {
		// same as previous case but with the orders inverted
		return t.swap(t.mergeSetSet(c2, c1))
	}

	// if it is not any of the above the operations do not conflict!
	return []Change{c2}, []Change{c1}
}

// first arg is a set, second arg can be any op.  first arg path is prefix of second
func (t Transformer) mergeSetSubPath(set, subPathOp Change) ([]Change, []Change) {
	c1 := set
	c2 := subPathOp
	l := len(c1.Path)

	if c2.Path[l] != c1.Set.Key {
		// no conflicts
		return []Change{c2}, []Change{c1}
	}

	// c1 wins but we need to rewrite c1.Set.Before to take into account effect of c2
	alteredC2Op := c2
	alteredC2Op.Path = c2.Path[l+1:]
	alteredC1Before := Utils(t).Apply(c1.Set.Before, []Change{alteredC2Op})

	alteredC1Set := &SetInfo{Key: c1.Set.Key, Before: alteredC1Before, After: c1.Set.After}
	alteredC1 := Change{Path: c1.Path, Set: alteredC1Set}
	return []Change{}, []Change{alteredC1}
}

func (t Transformer) mergeSetSplice(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && c2.Path[l] == c1.Set.Key {
		return t.mergeSetSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeSpliceSubPath(c2, c1))
	}

	// no conflicts
	return []Change{c2}, []Change{c1}
}

func (t Transformer) mergeSetMove(c1, c2 Change) ([]Change, []Change) {
	l := commonPathLength(c1.Path, c2.Path)
	if l == len(c1.Path) && c2.Path[l] == c1.Set.Key {
		return t.mergeSetSubPath(c1, c2)
	}

	if l == len(c2.Path) {
		return t.swap(t.mergeMoveSubPath(c2, c1))
	}

	// no conflicts
	return []Change{c2}, []Change{c1}
}
