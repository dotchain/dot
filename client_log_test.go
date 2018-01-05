// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"fmt"
	"github.com/dotchain/dot"
	"reflect"
	"testing"
)

func TestClientLog_simple_log(t *testing.T) {
	insert := func(o int, s string) []dot.Change {
		return []dot.Change{{Splice: &dot.SpliceInfo{o, "", s}}}
	}

	p := func(basisID, parentID string) []string {
		return []string{basisID, parentID}
	}

	journal := []dot.Operation{
		{ID: "initial", Changes: insert(0, "The fox jumped over the fence")},

		{ID: "0-0", Parents: p("initial", ""), Changes: insert(4, "red ")},
		{ID: "0-1", Parents: p("initial", "0-0"), Changes: insert(4, "beautiful ")},
		{ID: "0-2", Parents: p("initial", "0-1"), Changes: insert(4, "big ")},

		{ID: "1-0", Parents: p("initial", ""), Changes: insert(24, "yellow ")},
		{ID: "1-1", Parents: p("initial", "1-0"), Changes: insert(24, "large ")},
		{ID: "1-2", Parents: p("initial", "1-1"), Changes: insert(24, "frightening ")},
	}

	// iterate over all valid permutations of this journal
	for _, perm := range GetPermutations(journal) {
		test := &clog_test{T: t, journal: perm, initial: ""}
		test.Validate()
	}
}

func TestClientLog_staggered_basis(t *testing.T) {
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

	p := func(basis, parent string) []string {
		return []string{basis, parent}
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
	journal := []dot.Operation{
		{ID: "0-0", Parents: p("", ""), Changes: insert(0, "hello world")},
		{ID: "0-1", Parents: p("", "0-0"), Changes: insert(6, "beautiful ")},
		{ID: "0-2", Parents: p("", "0-1"), Changes: insert(6+len("beautiful "), "crazy ")},
		{ID: "1-0", Parents: p("0-0", ""), Changes: insert(11, "!")},
		{ID: "1-1", Parents: p("0-0", "1-0"), Changes: insert(0, "This ")},
		{ID: "1-2", Parents: p("0-0", "1-1"), Changes: insert(5, "is ")},
		{ID: "0-3", Parents: p("1-1", "0-2"), Changes: remove(len("This hello beautiful cra"), "z")},
	}

	// iterate over all valid permutations of this journal
	for _, perm := range GetPermutations(journal) {
		test := &clog_test{T: t, journal: perm, initial: ""}
		test.Validate()
	}
}

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

func TestClientLog_AppendClientOperation_autoreconciles(t *testing.T) {
	l := &dot.Log{
		MinIndex:    2,
		Rebased:     []dot.Operation{{}, {}, {ID: "third"}, {ID: "fourth"}, {ID: "fifth"}, {ID: "sixth"}},
		MergeChains: [][]dot.Operation{nil, nil, nil},
		IDToIndexMap: map[string]int{
			"first":  0,
			"second": 1,
			"third":  2,
			"fourth": 3,
			"fifth":  4,
			"sixth":  5,
		},
	}

	c := &dot.ClientLog{
		ClientIndex: 2,
		ServerIndex: 2,
		Rebased:     []dot.Operation{{ID: "fourth", Parents: []string{"two"}}},
		MergeChain:  []dot.Operation{{ID: "third", Parents: []string{"two"}}},
	}

	op := dot.Operation{ID: "fifth", Parents: []string{"two", "fourth"}}
	merge, err := c.AppendClientOperation(l, op)

	if err != nil {
		t.Error("Unexpected AppendClientOperation response", err)
	}
	if !reflect.DeepEqual(merge, []dot.Operation{{ID: "third"}, {ID: "sixth"}}) {
		t.Error("Unexpected AppendClientOperation merge", merge)
	}
	if c.ServerIndex != len(l.Rebased) || len(c.Rebased) > 0 || len(c.MergeChain) > 0 {
		t.Error("Unexpected clog", c)
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

type clog_test struct {
	*testing.T
	journal []dot.Operation
	initial interface{}
}

func (test *clog_test) GetClientOperations() [][]dot.Operation {
	clients := [][]dot.Operation{}

	// separate ops into different clients
	for _, op := range test.journal {
		if c, seq, ok := test.parseID(op.ID); ok {
			for len(clients) <= c {
				clients = append(clients, nil)
			}
			for len(clients[c]) <= seq {
				clients[c] = append(clients[c], dot.Operation{})
			}
			clients[c][seq] = op
		}
	}
	return clients
}

func (test *clog_test) GetFinalServerState() interface{} {
	t := test.T
	l := &dot.Log{}

	for _, op := range test.journal {
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("Append operation failed with error", err)
		}
	}
	return test.applyOps(test.initial, l.Rebased)
}

func (test *clog_test) GetClientStateForOperation(slog *dot.Log, op dot.Operation) (*dot.ClientLog, interface{}) {
	clog := &dot.ClientLog{}
	basisIndex := slog.IDToIndexMap[op.BasisID()]
	if op.BasisID() == "" {
		basisIndex = -1
	}

	m, err := clog.AppendClientOperation(slog, op)

	if err == nil {
		client := test.applyOps(test.initial, slog.Rebased[:basisIndex+1])
		client = test.applyOps(client, []dot.Operation{op})
		return clog, test.applyOps(client, m)
	}
	test.Fatal("Could not start a new client log with an intermediate op", err)
	return nil, nil
}

func (test *clog_test) validateFinalStateFromFullJournal(op dot.Operation) {
	slog := &dot.Log{}
	for _, opx := range test.journal {
		if err := slog.AppendOperation(opx); err != nil {
			test.Fatal("Could not append", opx)
		}
	}
	serverState := test.applyOps(test.initial, slog.Rebased)
	_, finalState := test.GetClientStateForOperation(slog, op)
	if !reflect.DeepEqual(serverState, finalState) {
		test.Fatal("Failed to converge for op", op, serverState, "<--->", finalState)
	}
}

func (test *clog_test) sameIntermediateClientState(clog1, clog2 *dot.ClientLog, result1, result2 interface{}) {
	if !reflect.DeepEqual(clog1, clog2) {
		test.Fatal("Client logs differ", clog1, "<--->", clog2)
	}
	if !reflect.DeepEqual(result1, result2) {
		test.Fatal("Client states differ", result1, "<--->", result2)
	}
}

func (test *clog_test) createLog(ops, clientOps []dot.Operation) *dot.Log {
	result := &dot.Log{}
	for kk := range ops {
		if err := result.AppendOperation(test.journal[kk]); err != nil {
			test.Fatal("Unexpected append operation fail", err)
		}
	}

	rawOps := []dot.Operation{}
	for _, cop := range clientOps {
		for _, op := range test.journal {
			if cop.ID == op.ID {
				rawOps = append(rawOps, op)
				break
			}
		}
	}

	for _, op := range rawOps {
		if err := result.AppendOperation(op); err != nil {
			test.Fatal("Unexpected append operation fail", err)
		}
	}

	return result
}

func (test *clog_test) GetFinalClientState(ops []dot.Operation) interface{} {
	// slog and clog maintiain the client's view of its server and itself
	slog, clog := &dot.Log{}, &dot.ClientLog{}

	// client state
	client := test.initial

	// merged operations for client to execute
	for _, op := range ops {
		// update server log up to op.BasisID() and reconcile clog with it
		client = test.applyOps(client, test.UpdateLogAndReconcile(slog, clog, &op))

		if len(clog.Rebased) > 1 {
			// create a copy of slog and add all but the last Rebased
			// operation to slog
			rebased := clog.Rebased[:len(clog.Rebased)-1]
			dupe := test.createLog(test.journal[:len(slog.Rebased)], rebased)
			clog2 := &dot.ClientLog{}
			if _, err := clog2.AppendClientOperation(dupe, op); err != nil {
				test.Fatal("Unexpected append operation error", err)
			}
			if !reflect.DeepEqual(rebased, dupe.Rebased[len(slog.Rebased):]) {
				test.Fatal("Differing rebased values!", rebased, dupe.Rebased[len(slog.Rebased)+1:])
			}
			continue
		}

		// also create a temporary client log with this op and check
		// if it ends up with the same state.  But we can do this only if
		// there is exactly one operations in clog.Rebased

		clog2, client2 := test.GetClientStateForOperation(slog, op)
		test.sameIntermediateClientState(clog, clog2, client, client2)

		// validate that intermeidate fetching form the op leads to same
		// final state
		test.validateFinalStateFromFullJournal(op)
	}

	// update server log one last time in case there is still some stuff left over
	return test.applyOps(client, test.UpdateLogAndReconcile(slog, clog, nil))
}

func (test *clog_test) Validate() {
	clients := test.GetClientOperations()
	serverState := test.GetFinalServerState()

	// reconstruct client actions and through that the client state
	for kk, c := range clients {
		clientState := test.GetFinalClientState(c)

		// validate client model has converged to server
		if !reflect.DeepEqual(clientState, serverState) {
			test.Fatal("Client(non-greedy)", kk, "diverged.  Expected", serverState, "got", clientState)
		}
	}
}

func (test *clog_test) UpdateLogAndReconcile(slog *dot.Log, clog *dot.ClientLog, op *dot.Operation) []dot.Operation {
	if op == nil {
		// use the last op ID as basisID to effectively flush the full journal
		test.UpdateLog(slog, test.journal[len(test.journal)-1].ID)
	} else {
		// use the basisID of the provided op for this
		test.UpdateLog(slog, op.BasisID())
	}

	// reconcile
	merge, err := clog.Reconcile(slog)
	if err != nil {
		test.Fatal("Client failed to reconcile op", err)
	}

	result := append([]dot.Operation{}, merge...)

	// merge the op
	if op != nil {
		result = append(result, *op)
		merge, err := clog.AppendClientOperation(slog, *op)
		if err != nil {
			test.Error("Append failed", *op)
			test.Error("clog server index", clog.ServerIndex, "client index", clog.ClientIndex)
			test.Error("clog.Rebased", clog.Rebased)
			test.Error("clog.MergeChain", clog.MergeChain)
			test.Error("slog.Rebased", slog.Rebased)
			test.Error("journal", test.journal)
			test.Fatal("Client failed to append op", err)
		}
		result = append(result, merge...)
	}
	return result
}

func (test *clog_test) UpdateLog(l *dot.Log, basisID string) {
	t := test.T

	if _, ok := l.IDToIndexMap[basisID]; ok || basisID == "" {
		return
	}

	start := -1

	if len(l.Rebased) > 0 {
		lastID := l.Rebased[len(l.Rebased)-1].ID
		for start = 0; test.journal[start].ID != lastID; start++ {
		}
	}

	for kk := start + 1; kk < len(test.journal); kk++ {
		op := test.journal[kk]
		if err := l.AppendOperation(op); err != nil {
			t.Fatal("Failed to append op", kk, "when updating basis to", basisID, err)
		}
		if op.ID == basisID {
			return
		}
	}
}

func (test *clog_test) parseID(s string) (client int, seq int, ok bool) {
	n, err := fmt.Sscanf(s, "%d-%d", &client, &seq)
	return client, seq, (n == 2 && err == nil)
}

func (test *clog_test) applyOps(result interface{}, ops []dot.Operation) interface{} {
	for _, op := range ops {
		result = applyMany(result, op.Changes)
	}
	return result
}
