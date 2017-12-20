// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"github.com/dotchain/dot/encoding"
	"github.com/pkg/errors"
)

// ErrMergeWithSelf is used by MergeOperations if it finds two operations
// with same ID being merged against each other
var ErrMergeWithSelf = errors.New("cannot merge an operation with itself")

// ModelBuilder represents a function that can incrementally build a model
type ModelBuilder func(oldModel interface{}, rebased []Operation) interface{}

// Transformer provides the basic functionality to transform
// Change and Operation values.
type Transformer struct {
	C encoding.Catalog
}

// left is Splice
func (t Transformer) mergeSpliceAny(left, right Change) ([]Change, []Change) {
	switch {
	case right.Splice != nil:
		return t.mergeSpliceSplice(left, right)
	case right.Set != nil:
		return t.swap(t.mergeSetSplice(right, left))
	case right.Move != nil:
		return t.swap(t.mergeMoveSplice(right, left))
	case right.Range != nil:
		return t.swap(t.mergeRangeSplice(right, left))
	}
	return []Change{right}, []Change{left}
}

// left is Set
func (t Transformer) mergeSetAny(left, right Change) ([]Change, []Change) {
	switch {
	case right.Splice != nil:
		return t.mergeSetSplice(left, right)
	case right.Set != nil:
		return t.mergeSetSet(left, right)
	case right.Move != nil:
		return t.mergeSetMove(left, right)
	case right.Range != nil:
		return t.swap(t.mergeRangeSet(right, left))
	}
	return []Change{right}, []Change{left}
}

// left is Move
func (t Transformer) mergeMoveAny(left, right Change) ([]Change, []Change) {
	switch {
	case right.Splice != nil:
		return t.mergeMoveSplice(left, right)
	case right.Set != nil:
		return t.swap(t.mergeSetMove(right, left))
	case right.Move != nil:
		return t.mergeMoveMove(left, right)
	case right.Range != nil:
		return t.swap(t.mergeRangeMove(right, left))
	}
	return []Change{right}, []Change{left}
}

// left is Range
func (t Transformer) mergeRangeAny(left, right Change) ([]Change, []Change) {
	switch {
	case right.Splice != nil:
		return t.mergeRangeSplice(left, right)
	case right.Set != nil:
		return t.mergeRangeSet(left, right)
	case right.Move != nil:
		return t.mergeRangeMove(left, right)
	case right.Range != nil:
		return t.mergeRangeRange(left, right)
	}
	return []Change{right}, []Change{left}
}

func (t Transformer) mergeChange(left, right Change) ([]Change, []Change) {
	switch {
	case left.Splice != nil:
		return t.mergeSpliceAny(left, right)
	case left.Set != nil:
		return t.mergeSetAny(left, right)
	case left.Move != nil:
		return t.mergeMoveAny(left, right)
	case left.Range != nil:
		return t.mergeRangeAny(left, right)
	}
	return []Change{right}, []Change{left}
}

// MergeChanges merges two sets of operation infos (based on the same initial state) and returns
// "compensation" Changes, which when applied after the corresponding set will guarantee
// convergence.
//
// For example:
//
//   // Assume initial state = initialState.  Two separate strand of mutations = left, right
//   transformer := ot.Transformer{}
//   leftCompensation, rightCompensation := transformer.MergeChanges(left, right)
//   finalLeftState := initialState.apply(left).apply(leftCompensation)
//   finalRightState := initialState.apply(right).apply(rightCompensation)
//   assert(finalLeftState.isEqual(finalRightState)
//
//
// For those familiar with `git`, another way to look at this is that transforming left
// against right returns `MERGE` and `REBASE` where `MERGE` provides the set of changes
// that would merge the `right` branch into `left` and `REBASE` provides the set of changes
// which would rebase `left` onto `right`.
//
// Yet another way to look at this is to consider the result of tranforming `left` and
// `right` is a pair of operations that have the effects of right and left (respectively)
// and which can be applied on `left` and `right` (respectively) to yield a common state.
//
func (t Transformer) MergeChanges(left, right []Change) ([]Change, []Change) {
	llen, rlen := len(left), len(right)

	if llen == 0 || rlen == 0 {
		return right, left
	}

	if llen == 1 && rlen == 1 {
		return t.mergeChange(left[0], right[0])
	}

	if llen > 1 {
		left1, right1 := t.MergeChanges(left[:llen-1], right)
		finalLeft, rightRemainder := t.MergeChanges(left[llen-1:], left1)
		return finalLeft, t.join(right1, rightRemainder)
	}

	// At this point rlen > 1 and llen == 1
	left1, right1 := t.MergeChanges(left, right[:rlen-1])
	leftRemainder, finalRight := t.MergeChanges(right1, right[rlen-1:])
	return t.join(left1, leftRemainder), finalRight
}

// MergeOperations merges two sets of operations (based on the same initial state) and returns
// "compensation" Operation sets, which when applied after the corresponding set will guarantee
// convergence.
//
// For example:
//
//   // Assume initial state = initialState.  Two separate strand of mutations = left, right
//   transformer := ot.Transformer{}
//   leftCompensation, rightCompensation := transformer.MergeOperations(left, right)
//   finalLeftState := initialState.apply(left).apply(leftCompensation)
//   finalRightState := initialState.apply(right).apply(rightCompensation)
//   assert(finalLeftState.isEqual(finalRightState)
//
// Please see MergeChange as this is a very thin wrapper around MergeChanges.
//
// Note that MergeOprations is not symmetric.  If the order of arguments are inverted,
// the transformed versions may be different and more importantly, the converged document
// will be different in some cases.  For predictable behavior, the order of the arguments
// should map to some prior fixed order (such as the order in a Journal or Log, if one
// exists).  For example, Log:AppendOperation ensures that "left" is always an earlier
// operation in the log
//
func (t Transformer) MergeOperations(left, right []Operation) ([]Operation, []Operation) {
	llen, rlen := len(left), len(right)
	if llen == 0 || rlen == 0 {
		return right, left
	}

	if llen == 1 && rlen == 1 {
		if left[0].ID != "" && left[0].ID == right[0].ID {
			panic(ErrMergeWithSelf)
		}
		linfo, rinfo := t.MergeChanges(left[0].Changes, right[0].Changes)
		return []Operation{right[0].withChanges(linfo)}, []Operation{left[0].withChanges(rinfo)}
	}

	if llen > 1 {
		left1, right1 := t.MergeOperations(left[:llen-1], right)
		finalLeft, rightRemainder := t.MergeOperations(left[llen-1:], left1)
		return finalLeft, t.joinOperation(right1, rightRemainder)
	}

	// At this point rlen > 1 and llen == 1
	left1, right1 := t.MergeOperations(left, right[:rlen-1])
	leftRemainder, finalRight := t.MergeOperations(right1, right[rlen-1:])
	return t.joinOperation(left1, leftRemainder), finalRight
}

// TryMergeOperations is same as #MergeOperations except that it returns
// success or failure.  MergeOperation can panic if the operations
// are badly formatted and effectively invalid.  TryOperation catches
// the panic and recovers but sets status to false
func (t Transformer) TryMergeOperations(left, right []Operation) (l []Operation, r []Operation, ok bool) {
	// it is safe to recover since all the methods called by this
	// or functions within it do not modify any state
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	l, r = t.MergeOperations(left, right)
	return l, r, true
}

// Utils

func (t Transformer) swap(info1, info2 []Change) ([]Change, []Change) {
	return info2, info1
}

func (t Transformer) join(infos ...[]Change) []Change {
	result := []Change{}
	for _, info := range infos {
		result = append(result, info...)
	}
	return result
}

func (t Transformer) joinOperation(infos ...[]Operation) []Operation {
	result := []Operation{}
	for _, info := range infos {
		result = append(result, info...)
	}
	return result
}

func (t Transformer) toArray(i interface{}) encoding.UniversalEncoding {
	if i == nil {
		return t.C.Get([]interface{}{})
	}
	return t.C.Get(i)
}
