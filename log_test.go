// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"reflect"
	"testing"
)

func logAppend(op dot.Operation, rebased []dot.Operation, mergeChains [][]dot.Operation, idToIndexMap map[string]int) (dot.Operation, []dot.Operation, error) {
	l := dot.Log{Rebased: rebased, MergeChains: mergeChains, IDToIndexMap: idToIndexMap}
	err := l.AppendOperation(op)
	if err != nil {
		return dot.Operation{}, nil, err
	}
	return l.Rebased[len(l.Rebased)-1], l.MergeChains[len(l.MergeChains)-1], nil
}

func TestLog_emptyLogSuccess(t *testing.T) {
	op := dot.Operation{Changes: []dot.Change{
		{Splice: &dot.SpliceInfo{Offset: 0, Before: "", After: "The fox jumped over the fence"}},
	}}
	rebased, mergeChain, err := logAppend(op, nil, nil, map[string]int{})
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err.Error())
	} else if len(mergeChain) > 0 {
		t.Errorf("MergeChain expected to be zero length: %v\n", mergeChain)
	} else if !reflect.DeepEqual(op, rebased) {
		t.Errorf("Unexpected rebase result %#v\n", rebased)
	}
}

func TestLog_emptyLogMissingBasis(t *testing.T) {
	op := dot.Operation{
		Parents: []string{"hello"},
		Changes: []dot.Change{
			{Splice: &dot.SpliceInfo{Offset: 0, Before: "", After: "The fox jumped over the fence"}},
		},
	}
	_, _, err := logAppend(op, nil, nil, map[string]int{})
	if err != dot.ErrMissingParentOrBasis {
		t.Error("Unexpected err return value")
	}
}

func TestLog_fullMerge(t *testing.T) {
	insert := func(o int, s string) []dot.Change {
		return []dot.Change{{Splice: &dot.SpliceInfo{o, "", s}}}
	}

	p := func(basisID, parentID string) []string {
		return []string{basisID, parentID}
	}

	journal := []dot.Operation{
		// initial
		{ID: "initial", Changes: insert(0, "The fox jumped over the fence")},

		// client1 inserts "big beautiful red " before "fox"
		{ID: "0-0", Parents: p("initial", ""), Changes: insert(4, "red ")},
		{ID: "0-1", Parents: p("initial", "0-0"), Changes: insert(4, "beautiful ")},
		{ID: "0-2", Parents: p("initial", "0-1"), Changes: insert(4, "big ")},

		// client2 inserts "frightening large yellow " before "fence"
		{ID: "1-0", Parents: p("initial", ""), Changes: insert(24, "yellow ")},
		{ID: "1-1", Parents: p("initial", "1-0"), Changes: insert(24, "large ")},
		{ID: "1-2", Parents: p("initial", "1-1"), Changes: insert(24, "frightening ")},
	}

	// validate all combinations
	for _, perm := range GetPermutations(journal) {
		l := &dot.Log{}

		for kk, op := range perm {
			if err := l.AppendOperation(op); err != nil {
				t.Errorf("Unexpected err return value: %#v, %#v\n", err, kk)
				return
			}
		}

		var result interface{} = ""
		for _, op := range l.Rebased {
			result = applyMany(result, op.Changes)
		}
		if !dot.Utils(x).AreSame(result, "The big beautiful red fox jumped over the frightening large yellow fence") {
			t.Errorf("Unexpected result of rebasing: %#v\n", result)
		}
	}
}

func TestLog_InvalidBasis(t *testing.T) {
	initial := dot.Change{Splice: &dot.SpliceInfo{0, "", "hello world"}}
	log := []dot.Operation{
		{ID: "initial", Changes: []dot.Change{initial}},
		{ID: "fail", Parents: []string{"blah"}},
	}

	l := &dot.Log{}
	if err := l.AppendOperation(log[0]); err != nil {
		t.Error("Unexpected err", err)
	}

	err := l.AppendOperation(log[1])
	if err != dot.ErrMissingParentOrBasis {
		t.Error("Unexpected err", err)
	}
}

func TestLog_InvalidParent(t *testing.T) {
	initial := dot.Change{Splice: &dot.SpliceInfo{0, "", "hello world"}}
	log := []dot.Operation{
		{ID: "initial", Changes: []dot.Change{initial}},
		{ID: "fail", Parents: []string{"", "blah"}},
	}

	l := &dot.Log{}
	if err := l.AppendOperation(log[0]); err != nil {
		t.Error("Unexpected err", err)
	}

	err := l.AppendOperation(log[1])
	if err != dot.ErrMissingParentOrBasis {
		t.Error("Unexpected err", err)
	}
}

func TestLog_InvalidOperation(t *testing.T) {
	initial := dot.Change{Splice: &dot.SpliceInfo{1, []interface{}{5}, []interface{}{10}}}
	incompatible := dot.Change{Splice: &dot.SpliceInfo{0, "hello", "world"}}

	log := []dot.Operation{
		{ID: "initial", Changes: []dot.Change{initial}},
		{ID: "invalid", Changes: []dot.Change{incompatible}},
	}

	l := &dot.Log{}
	if err := l.AppendOperation(log[0]); err != nil {
		t.Error("Unexpected err", err)
	}

	if err := l.AppendOperation(log[1]); err != dot.ErrInvalidOperation {
		t.Error("Unexpected err", err)
	}
}

func TestLog_EmptyBasis(t *testing.T) {
	initial := dot.Change{Splice: &dot.SpliceInfo{0, "", "hello world"}}
	second := dot.Change{Splice: &dot.SpliceInfo{0, "", "HELLO WORLD"}}
	log := []dot.Operation{
		{ID: "initial", Changes: []dot.Change{initial}},
		{ID: "empty", Changes: []dot.Change{second}},
	}

	l := &dot.Log{}
	if err := l.AppendOperation(log[0]); err != nil {
		t.Error("Unexpected err", err)
	}

	err := l.AppendOperation(log[1])
	if err != nil {
		t.Error("Unexpected err", err)
	}

	// expect that second gets offset by length of "hello world" when rebased
	// so mergeChain1 should simply be initial effectively
	if !reflect.DeepEqual(l.MergeChains[1][0].Changes, []dot.Change{initial}) {
		t.Error("Unexpected mergechain1", l.MergeChains[1][0].Changes[0].Splice)
	}
	// rebased should be second but with offset bumped by "hello world" length
	second.Splice.Offset += len("hello world")
	if !reflect.DeepEqual(l.Rebased[1].Changes, []dot.Change{second}) {
		t.Error("Unexpected rebased", l.Rebased[1].Changes[0].Splice)
	}
}

func TestLog_UnloadedBasis(t *testing.T) {
	l := &dot.Log{
		MinIndex:     1,
		Rebased:      []dot.Operation{{}, {ID: "two"}},
		MergeChains:  [][]dot.Operation{nil, nil},
		IDToIndexMap: map[string]int{"one": 0, "two": 1},
	}
	err := l.AppendOperation(dot.Operation{ID: "three", Parents: []string{"one", ""}})
	if err != dot.ErrLogNeedsBackfilling {
		t.Error("Unexpected err", err)
	}
}

func TestLog_DuplicateOp(t *testing.T) {
	l := &dot.Log{
		MinIndex:     1,
		Rebased:      []dot.Operation{{}, {ID: "two"}},
		MergeChains:  [][]dot.Operation{nil, nil},
		IDToIndexMap: map[string]int{"one": 0, "two": 1},
	}
	oldRebased := append([]dot.Operation{}, l.Rebased...)
	err := l.AppendOperation(l.Rebased[1])
	if err != nil {
		t.Error("Unexpected err", err)
	}
	if !reflect.DeepEqual(oldRebased, l.Rebased) {
		t.Error("Rebased modified unexpectedly", l.Rebased)
	}
}

func TestLog_staggeredBasis(t *testing.T) {
	// appending op must have basisIndex < parentIndex but greater than index of parent op.
	// Easiest way to simulate that is have op1, op2, op3 where
	// op3 is based and parented on op1 and then we append an op
	// based on op2 and parented on op3. To make things interesting,
	// we need to make op2 conflict with op3 -- so we will make both insert
	// at the same location.  In addition, we will add an op4 based on op2
	// so we can be sure that the new op is properly transformed afterwards.

	insert := func(o int, s string) []dot.Change {
		return []dot.Change{{Splice: &dot.SpliceInfo{o, "", s}}}
	}

	remove := func(o int, s string) []dot.Change {
		return []dot.Change{{Splice: &dot.SpliceInfo{o, s, ""}}}
	}

	// The actual details:
	// Op1 = insert "hello world" at 0.  Basis/Parent = nothing.
	// Op2 = insert "beautiful " at 6 (after "hello ", basis/parent = op1)
	// Op3 = insert "crazy " at 6 (after "hello ") Basis/Parent = op1
	// Op4 = insert "!" after "hello world" -- Basis/Parent = op1
	//
	// Append happens with basis = op2 but parent = op3.
	// That is, we consider the client which applied op1, op3 to get
	// "hello crazy world" and then factored in the effect of
	// Op2 to get "hello beautiful crazy world" and now does the
	// delete after "hello beautiful cra"
	// Append = delete "z" from "crazy".  Basis = Op3, Parent = Op2
	log := []dot.Operation{
		{ID: "one", Changes: insert(0, "hello world")},
		{ID: "two", Changes: insert(6, "beautiful "), Parents: []string{"one"}},
		{ID: "three", Changes: insert(6, "crazy "), Parents: []string{"one"}},
		{ID: "four", Changes: insert(11, "!"), Parents: []string{"one"}},
		{ID: "five", Changes: remove(len("hello beautiful cra"), "z"), Parents: []string{"two", "three"}},
	}

	l := &dot.Log{}
	for _, op := range log {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("Failed to apply operation", op)
		}
	}

	var result interface{} = ""
	for _, op := range l.Rebased {
		result = applyMany(result, op.Changes)
	}

	if !dot.Utils(x).AreSame(result, "hello beautiful cray world!") {
		t.Error("Got unexpected result", result)
	}

	// Note that when the last op is appended, we expected the client
	// to simply have "hello beautiful cray world" but not have the exclamation at the end
	// lets check that the last merge chain is good as well
	mergeChain := l.MergeChains[len(l.MergeChains)-1]
	result = "hello beautiful cray world"
	for _, op := range mergeChain {
		result = applyMany(result, op.Changes)
	}

	if !dot.Utils(x).AreSame(result, "hello beautiful cray world!") {
		t.Error("Got unexpected merge result", result)
	}

}
