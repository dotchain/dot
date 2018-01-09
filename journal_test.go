// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"encoding/json"
	datalib "github.com/dotchain/dataset/tools/lib"
	"github.com/dotchain/dot"
	"go/build"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"testing"
)

func mustNotFail(err error) {
	if err != nil {
		panic(err)
	}
}

func TestJournalSuite(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Failed:", r, string(debug.Stack()))
		}
	}()

	fname := build.Default.GOPATH + "/src/github.com/dotchain/dataset/json/journal_suite.json"
	bytes, err := ioutil.ReadFile(fname)
	mustNotFail(err)

	var data map[string]interface{}
	mustNotFail(json.Unmarshal(bytes, &data))

	data = data["test"].(map[string]interface{})
	for test := range data {
		testData := data[test].(map[string]interface{})
		journal := testData["journal"].([]interface{})
		rebased := testData["rebased"].([]interface{})
		mergeChains := testData["mergeChains"].([]interface{})
		t.Run(test, journalTestCase{journal, rebased, mergeChains}.Run)
	}
}

type journalTestCase struct {
	journal, rebased, mergeChains []interface{}
}

func (tc journalTestCase) Run(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Failed:", r, string(debug.Stack()))
		}
	}()

	ops := []dot.Operation{}
	inputs := []string{}
	outputs := []string{}

	for _, elt := range tc.journal {
		parts := elt.([]interface{})
		id := parts[0].(string)
		basisID := parts[1].(string)
		parentID := parts[2].(string)
		inp, ch := datalib.Compact{}.Decode(parts[3].(string))
		inputs = append(inputs, inp)
		outputs = append(outputs, datalib.Compact{}.Apply(inp, ch))
		changes := []dot.Change{ch}
		parents := []string{basisID, parentID}
		op := dot.Operation{ID: id, Parents: parents, Changes: changes}
		ops = append(ops, op)
	}

	log, variations := tc.getAllVariations(inputs[0], ops)
	tc.validateLog(t, log, inputs, outputs)
	final := tc.applyOperations(inputs[0], log.Rebased)
	for kk, variation := range variations {
		actual := tc.applyOperations(inputs[0], variation)
		if actual != final {
			t.Error("Variation", kk, "produced", actual, ". Initial variation produced", final)
			t.Error("Variation is", tc.toCompactForm(inputs[0], variation))
			t.Error("Rebased is", tc.toCompactForm(inputs[0], log.Rebased))
		}
	}
}

func (tc journalTestCase) getAllVariations(input string, ops []dot.Operation) (*dot.Log, [][]dot.Operation) {
	variations := [][]dot.Operation{}

	l := &dot.Log{}
	for _, op := range ops {
		mustNotFail(l.AppendOperation(op))
	}

	variations = append(variations, tc.getObserverVariations(input, ops, l)...)
	variations = append(variations, tc.getEagerClientVariations(input, ops, l)...)
	variations = append(variations, tc.getReconnectVariations(input, ops, l)...)
	return l, variations
}

// getObserverVariations only returns the variations where a client
// is simply initializing up to a part of the log and then reconciling
// with the rest.  This basically should produce just the exact same
// sequence as the rebased operations in the log and is a very trivial
// set of variations
func (tc journalTestCase) getObserverVariations(input string, ops []dot.Operation, l *dot.Log) [][]dot.Operation {
	variations := [][]dot.Operation{}

	for kk := range ops {
		slog := &dot.Log{}
		basisID := ""
		for jj := 0; jj < kk; jj++ {
			slog.AppendOperation(ops[jj])
			basisID = ops[jj].ID
		}
		clog, rebased, clientRebased, err := dot.BootstrapClientLog(slog, nil)
		mustNotFail(err)

		rest1, err1 := clog.Reconcile(l)
		mustNotFail(err1)
		_, rest2, err2 := dot.ReconnectClientLog(l, nil, basisID, "")
		mustNotFail(err2)

		var1 := append(append(append([]dot.Operation{}, rebased...), clientRebased...), rest1...)
		var2 := append(append(append([]dot.Operation{}, rebased...), clientRebased...), rest2...)
		variations = append(variations, var1, var2)
	}

	return variations
}

func (tc journalTestCase) getReconnectVariations(input string, ops []dot.Operation, l *dot.Log) [][]dot.Operation {
	result := [][]dot.Operation{}

	seq := tc.getEagerClientOrder(ops)
	slog := &dot.Log{}

	// for each state there is a client log
	clog, _, _, _ := dot.BootstrapClientLog(slog, nil)
	applied := []dot.Operation{}

	// all seen client operations from the journal
	cops := []dot.Operation{}

	basisID := ""
	parentID := ""
	basisIndex := 0

	variation := tc.getReconnectVariation(l, cops, basisID, parentID)
	result = append(result, variation)
	for _, op := range seq {
		if !strings.HasPrefix(op.ID, "c") {
			for ; basisIndex < len(ops); basisIndex++ {
				mustNotFail(slog.AppendOperation(ops[basisIndex]))
				more, err := clog.Reconcile(slog)
				mustNotFail(err)
				applied = append(applied, more...)
				if ops[basisIndex].ID == op.ID {
					break
				}
			}
			basisID = op.ID
		} else {
			cops = append(cops, op)
			applied = append(applied, op)
			more, err := clog.AppendClientOperation(slog, op)
			mustNotFail(err)
			if len(more) > 0 {
				panic("Unexpected more growth")
			}
			parentID = op.ID
		}

		variation := tc.getReconnectVariation(l, cops, basisID, parentID)
		variation = append(append([]dot.Operation{}, applied...), variation...)
		result = append(result, variation)
	}
	return result
}

func (tc journalTestCase) getReconnectVariation(l *dot.Log, cops []dot.Operation, basisID, parentID string) []dot.Operation {
	_, rest, err := dot.ReconnectClientLog(l, cops, basisID, parentID)
	mustNotFail(err)
	return rest
}

// getEagerClientVariation attempts to execute the operations in
// the order of a client -- for each client operation, all previous
// operation up to the basis of the current client operation are
// appended into server log but no more -- then the next client operation
// is added etc with a final reconcile
func (tc journalTestCase) getEagerClientVariations(input string, ops []dot.Operation, l *dot.Log) [][]dot.Operation {
	seq := tc.getEagerClientOrder(ops)
	slog := &dot.Log{}

	// for each state there is a client log
	clog, _, _, _ := dot.BootstrapClientLog(slog, nil)
	clientLogs := []*dot.ClientLog{clog}
	clientApplied := [][]dot.Operation{nil}

	// all seen client operations from the journal
	cops := []dot.Operation{}

	for _, op := range seq {
		if !strings.HasPrefix(op.ID, "c") {
			for _, oo := range ops {
				mustNotFail(slog.AppendOperation(oo))
				if oo.ID == op.ID {
					break
				}
			}
			continue
		}

		cops = append(cops, op)
		for jj, clog := range clientLogs {
			catchup, err := clog.Reconcile(slog)
			mustNotFail(err)
			clientApplied[jj] = append(clientApplied[jj], catchup...)
			more, err := clog.AppendClientOperation(slog, op)
			mustNotFail(err)
			clientApplied[jj] = append(clientApplied[jj], op)
			clientApplied[jj] = append(clientApplied[jj], more...)
		}

		// start a new client at this point
		clog, r, c, err := dot.BootstrapClientLog(slog, cops)
		mustNotFail(err)
		clientLogs = append(clientLogs, clog)
		clientApplied = append(clientApplied, append(append([]dot.Operation{}, r...), c...))
	}

	// reconcile the rest
	for jj, clog := range clientLogs {
		more, err := clog.Reconcile(slog)
		mustNotFail(err)
		clientApplied[jj] = append(clientApplied[jj], more...)
	}
	return clientApplied
}

func (tc journalTestCase) getEagerClientOrder(ops []dot.Operation) []dot.Operation {
	result, pending := []dot.Operation{}, []dot.Operation{}
	basisID := ""
	for _, op := range ops {
		if !strings.HasPrefix(op.ID, "c") {
			pending = append(pending, op)
			continue
		}
		if op.BasisID() != basisID {
			basisID = op.BasisID()
			for _, p := range pending {
				result = append(result, p)
				pending = pending[1:]
				if p.ID == basisID {
					break
				}
			}
		}
		result = append(result, op)
	}
	return append(result, pending...)
}

func (tc journalTestCase) getBootstrappedVariation(c string, slog *dot.Log, clog *dot.ClientLog, cops, rest []dot.Operation) []dot.Operation {
	for _, op := range rest {
		if strings.HasPrefix(op.ID, c) {
			more, err := clog.AppendClientOperation(slog, op)
			mustNotFail(err)
			cops = append(cops, more...)
		} else {
			mustNotFail(slog.AppendOperation(op))
		}
	}
	more, err := clog.Reconcile(slog)
	mustNotFail(err)
	return append(cops, more...)
}

func (tc journalTestCase) applyOperations(input string, ops []dot.Operation) string {
	result := input
	for _, op := range ops {
		result = datalib.Compact{}.ApplyMany(result, op.Changes)
	}
	return result
}

func (tc journalTestCase) toCompactForm(input string, ops []dot.Operation) []string {
	changes := []dot.Change{}
	for _, op := range ops {
		if len(op.Changes) == 0 {
			changes = append(changes, dot.Change{})
		} else {
			changes = append(changes, op.Changes...)
		}
	}
	return datalib.Compact{}.Encode(input, changes)
}

func (tc journalTestCase) validateLog(t *testing.T, l *dot.Log, inputs, outputs []string) {
	b, err := json.Marshal(tc.toCompactForm(inputs[0], l.Rebased))
	mustNotFail(err)

	actual := string(b)
	b, err = json.Marshal(tc.rebased)
	mustNotFail(err)

	expected := string(b)
	if actual != expected {
		t.Error("Rebased diverged. Expected", expected, "\nActual", actual)
	}

	mergeChains := [][]string{}
	for kk, output := range outputs {
		chain := tc.toCompactForm(output, l.MergeChains[kk])
		mergeChains = append(mergeChains, chain)
	}

	b, err = json.Marshal(mergeChains)
	mustNotFail(err)

	actual = string(b)
	b, err = json.Marshal(tc.mergeChains)
	mustNotFail(err)

	expected = string(b)
	if actual != expected {
		t.Error("MergeChains diverged. Expected", expected, "\nActual", actual)
	}
}
