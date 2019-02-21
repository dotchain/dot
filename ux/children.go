// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

// UpdateChildren updates a parent node with the new set of children
// attempting to minimize calls to RemoveChild and InsertChild
func UpdateChildren(parent Element, after []Element) {
	before := parent.Children()
	for _, op := range bestDiff(before, after, 0, nil) {
		if op.insert {
			parent.InsertChild(op.index, op.elt)
		} else {
			parent.RemoveChild(op.index)
		}
	}
}

type diff struct {
	insert bool
	elt    Element
	index  int
}

func bestDiff(before, after []Element, offset int, ops []diff) []diff {
	for len(before) > 0 && len(after) > 0 && before[0] == after[0] {
		offset++
		before, after = before[1:], after[1:]
	}

	switch {
	case len(before) == 0:
		for _, e := range after {
			ops = append(ops, diff{true, e, offset})
			offset++
		}
	case len(after) == 0:
		for range before {
			ops = append(ops, diff{false, nil, offset})
		}
	default:
		ops = chooseDiff(before, after, offset, ops)
	}

	return ops
}

func chooseDiff(before, after []Element, offset int, ops []diff) []diff {
	// choice1 = clone of ops + delete first before elt
	choice1 := append(ops, diff{false, nil, offset})
	choice1 = bestDiff(before[1:], after, offset, choice1)

	index := indexOf(before[0], after)
	if index == -1 {
		return choice1
	}

	// choice2 = clone of ops + insert index after elts
	choice2 := append([]diff(nil), ops...)
	for kk := 0; kk < index+1; kk++ {
		choice2 = append(choice2, diff{true, after[kk], offset + kk})
	}
	choice2 = bestDiff(before, after[index+1:], offset+index+1, choice2)
	if len(choice1) < len(choice2) {
		return choice1
	}
	return choice2
}

func indexOf(elt Element, elts []Element) int {
	for kk, e := range elts {
		if e == elt {
			return kk
		}
	}
	return -1
}
