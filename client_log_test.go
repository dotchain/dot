// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"testing"
)

func TestClientLog_Reconcile_needs_backfilling(t *testing.T) {
	l := &dot.Log{
		MinIndex:    2,
		Rebased:     []dot.Operation{{}, {}, {ID: "third"}},
		MergeChains: [][]dot.Operation{nil, nil, nil},
		IDToIndexMap: map[string]int{
			"first":  0,
			"second": 1,
			"third":  2,
		},
	}

	c := &dot.ClientLog{}

	// now attempting to reconcile c with l should barf since C needs items from
	// the very start
	if _, err := c.Reconcile(l); err != dot.ErrLogNeedsBackfilling {
		t.Error("Unexpected Reconcile response", err)
	}
}

func TestClientLog_InitializeFromJournal_needs_backfilling(t *testing.T) {
	l := &dot.Log{
		MinIndex:    2,
		Rebased:     []dot.Operation{{}, {}, {ID: "third"}},
		MergeChains: [][]dot.Operation{nil, nil, nil},
		IDToIndexMap: map[string]int{
			"first":  0,
			"second": 1,
			"third":  2,
		},
	}

	c := &dot.ClientLog{}

	op := dot.Operation{ID: "second"}
	if _, err := c.AppendClientOperation(l, op); err != dot.ErrLogNeedsBackfilling {
		t.Error("Unexpected Reconcile response", err)
	}
}

func TestClientLog_AppendClientOperation_missing_parent_basis(t *testing.T) {
	l := &dot.Log{}
	c := &dot.ClientLog{}

	op := dot.Operation{ID: "one", Parents: []string{"something"}}
	if _, err := c.AppendClientOperation(l, op); err != dot.ErrMissingParentOrBasis {
		t.Error("Unexpected Reconcile response", err)
	}

	op = dot.Operation{ID: "one", Parents: []string{"", "something"}}
	if _, err := c.AppendClientOperation(l, op); err != dot.ErrMissingParentOrBasis {
		t.Error("Unexpected Reconcile response", err)
	}
}

func TestClientLog_AppendClientOperation_needs_backfilling(t *testing.T) {
	l := &dot.Log{
		MinIndex:    2,
		Rebased:     []dot.Operation{{}, {}, {ID: "third"}},
		MergeChains: [][]dot.Operation{nil, nil, nil},
		IDToIndexMap: map[string]int{
			"first":  0,
			"second": 1,
			"third":  2,
		},
	}

	c := &dot.ClientLog{}
	op := dot.Operation{ID: "four", Parents: []string{"first"}}
	if _, err := c.AppendClientOperation(l, op); err != dot.ErrLogNeedsBackfilling {
		t.Error("Unexpected Reconcile response", err)
	}

	// the following should actually succeed because if parent is earlier than
	// basis, parent can be disregarded
	op = dot.Operation{ID: "four", Parents: []string{"third", "first"}}
	if _, err := c.AppendClientOperation(l, op); err != nil {
		t.Error("Unexpected Reconcile response", err)
	}
}

func TestClientLog_AppendClientOperation_second_op_invalid_basis(t *testing.T) {
	l := &dot.Log{}
	c := &dot.ClientLog{}

	initial := []dot.Operation{{ID: "one"}}
	for _, op := range initial {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("AppendOperation failed", err)
		}
	}

	if _, err := c.Reconcile(l); err != nil {
		t.Fatal("Reconcile failed", err)
	}

	validOp := dot.Operation{ID: "two", Parents: []string{"one"}}
	if _, err := c.AppendClientOperation(l, validOp); err != nil {
		t.Fatal("AppendClientOperation failed", err)
	}

	invalidOp := dot.Operation{ID: "three", Parents: []string{"blah"}}

	// now attempting to add this op should fail
	if _, err := c.AppendClientOperation(l, invalidOp); err != dot.ErrMissingParentOrBasis {
		t.Error("Unexpected Reconcile response", err)
	}
}

func TestClientLog_AppendClientOperation_invalid_op(t *testing.T) {
	l := &dot.Log{}
	c := &dot.ClientLog{}

	change1 := dot.Change{Splice: &dot.SpliceInfo{1, []interface{}{5}, []interface{}{10}}}
	change2 := dot.Change{Splice: &dot.SpliceInfo{0, "hello", "world"}}

	initial := []dot.Operation{{ID: "one", Changes: []dot.Change{change1}}}
	for _, op := range initial {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("AppendOperation failed", err)
		}
	}

	if _, err := c.Reconcile(l); err != nil {
		t.Fatal("Reconcile failed", err)
	}

	invalidOp := dot.Operation{ID: "three", Changes: []dot.Change{change2}}
	if _, err := c.AppendClientOperation(l, invalidOp); err != dot.ErrInvalidOperation {
		t.Error("Unexpected AppendClientOperation result", err)
	}

}

func TestClientLog_AppendClientOperation_invalid_op2(t *testing.T) {
	l := &dot.Log{}
	c := &dot.ClientLog{}

	change1 := dot.Change{Splice: &dot.SpliceInfo{1, []interface{}{5}, []interface{}{10}}}
	change2 := dot.Change{Splice: &dot.SpliceInfo{0, "hello", "world"}}

	invalidOp := dot.Operation{ID: "three", Changes: []dot.Change{change2}}
	if _, err := c.AppendClientOperation(l, invalidOp); err != nil {
		t.Error("Unexpected AppendClientOperation result", err)
	}

	initial := []dot.Operation{{ID: "one", Changes: []dot.Change{change1}}}
	for _, op := range initial {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("AppendOperation failed", err)
		}
	}

	if _, err := c.Reconcile(l); err != dot.ErrInvalidOperation {
		t.Fatal("Reconcile failed", err)
	}

}

func TestClientLog_bootstrap_invalid_bootstrap(t *testing.T) {
	l := &dot.Log{}

	change1 := dot.Change{Splice: &dot.SpliceInfo{1, []interface{}{5}, []interface{}{10}}}

	initial := []dot.Operation{{ID: "one", Changes: []dot.Change{change1}}}
	for _, op := range initial {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("AppendOperation failed", err)
		}
	}

	change2 := dot.Change{Splice: &dot.SpliceInfo{0, "hello", "world"}}
	invalidOp := dot.Operation{ID: "three", Changes: []dot.Change{change2}}
	_, _, _, err := dot.BootstrapClientLog(l, []dot.Operation{invalidOp})
	if err != dot.ErrInvalidOperation {
		t.Fatal("Bootstrap failed to fail", err)
	}
}

func TestClientLog_bootstrap_missing_parent_basis(t *testing.T) {
	l := &dot.Log{}

	invalidOp := dot.Operation{ID: "three", Parents: []string{"miss1", ""}}
	_, _, err := dot.ReconnectClientLog(l, []dot.Operation{invalidOp}, "", "")
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Bootstrap failed to fail", err)
	}

	invalidOp = dot.Operation{ID: "three", Parents: []string{"", "miss2"}}
	_, _, err = dot.ReconnectClientLog(l, []dot.Operation{invalidOp}, "", "")
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Bootstrap failed to fail", err)
	}
}

func TestClientLog_reconnect_missing_parent_basis(t *testing.T) {
	l := &dot.Log{}

	invalidOp := dot.Operation{ID: "three", Parents: []string{"miss1", ""}}
	_, _, err := dot.ReconnectClientLog(l, []dot.Operation{invalidOp}, "", "")
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Bootstrap failed to fail", err)
	}

	invalidOp = dot.Operation{ID: "three", Parents: []string{"", "miss2"}}
	_, _, err = dot.ReconnectClientLog(l, []dot.Operation{invalidOp}, "", "")
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Bootstrap failed to fail", err)
	}

	_, _, err = dot.ReconnectClientLog(l, nil, "miss3", "")
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Bootstrap failed to fail", err)
	}

	_, _, err = dot.ReconnectClientLog(l, nil, "", "miss4")
	if err != dot.ErrMissingParentOrBasis {
		t.Fatal("Bootstrap failed to fail", err)
	}
}
