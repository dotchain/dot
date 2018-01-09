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
