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
	"strings"
	"testing"
)

func TestJournalSuite(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Failed:", r)
		}
	}()

	fname := build.Default.GOPATH + "/src/github.com/dotchain/dataset/json/journal_suite.json"
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}

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
			t.Error("Failed:", r)
		}
	}()

	ops := []dot.Operation{}
	input := ""

	for kk, elt := range tc.journal {
		parts := elt.([]interface{})
		id := parts[0].(string)
		basisID := parts[1].(string)
		parentID := parts[2].(string)
		inp, ch := datalib.Compact{}.Decode(parts[3].(string))
		if kk == 0 {
			input = inp
		}
		changes := []dot.Change{ch}
		parents := []string{basisID, parentID}
		op := dot.Operation{ID: id, Parents: parents, Changes: changes}
		ops = append(ops, op)
	}

	log, variations := tc.getAllVariations(input, ops)
	tc.validateLog(t, log, input)
	final := tc.applyOperations(input, log.Rebased)
	for kk, variation := range variations {
		actual := tc.applyOperations(input, variation)
		if actual != final {
			t.Error("Variation", kk, "produced", actual, ". Initial variation produced", final)
		}
	}
}

func (tc journalTestCase) getAllVariations(input string, ops []dot.Operation) (*dot.Log, [][]dot.Operation) {
	variations := [][]dot.Operation{}

	l := &dot.Log{}
	for _, op := range ops {
		if err := l.AppendOperation(op); err != nil {
			panic(err)
		}
	}

	variations = append(variations, tc.getObserverVariations(input, ops, l)...)
	variations = append(variations, tc.getClientVariations(input, "c", ops, l)...)
	return l, variations
}

func (tc journalTestCase) getReconnectVariations(input, c string, ops, cops []dot.Operation, l, slog *dot.Log, clog *dot.ClientLog) ([]dot.Operation, [][]dot.Operation) {
	variations := [][]dot.Operation{}

	for jj := 0; jj < len(ops); jj++ {
		if !strings.HasPrefix(ops[jj].ID, c) {
			continue
		}

		basisID := ops[jj].BasisID()
		parentID := ops[jj].ID

		// before attempting the client op, ensure all server ops
		// up to basis have been applied
		if basisID != "" {
			for qq := range ops {
				slog.AppendOperation(ops[qq])
				if ops[qq].ID == basisID {
					break
				}
			}
		}
		more, err := clog.Reconcile(slog)
		if err != nil {
			panic(err)
		}
		cops = append(cops, more...)
		cops = append(cops, ops[jj])
		more, err = clog.AppendClientOperation(slog, ops[jj])
		if err != nil || len(more) > 0 {
			panic(err)
		}

		_, rest, err := dot.ReconnectClientLog(l, []dot.Operation{ops[jj]}, basisID, parentID)
		if err != nil {
			panic(err)
		}
		variation := append(append([]dot.Operation{}, cops...), rest...)
		variations = append(variations, variation)
	}
	return cops, variations
}

func (tc journalTestCase) getClientVariations(input, c string, ops []dot.Operation, l *dot.Log) [][]dot.Operation {
	var own, others []dot.Operation

	variations := [][]dot.Operation{}

	for kk, op := range ops {
		if !strings.HasPrefix(op.ID, c) {
			own, others = []dot.Operation{}, ops[:kk+1]
			continue
		}

		own = append(own, op)
		slog := &dot.Log{}
		basisID := ""
		for _, oo := range others {
			slog.AppendOperation(oo)
			basisID = op.ID
		}
		clog, reb, clientreb, err := dot.BootstrapClientLog(slog, own)
		if err != nil {
			panic(err)
		}

		cops := append(append([]dot.Operation{}, reb...), clientreb...)
		parentID := op.ID

		// attempt a full reconnect in the current state
		// and add it as a variation
		_, rest, err := dot.ReconnectClientLog(l, own, basisID, parentID)
		if err != nil {
			panic(err)
		}
		variations = append(variations, append(append([]dot.Operation{}, cops...), rest...))
		cops, inner := tc.getReconnectVariations(input, c, ops[kk+1:], cops, l, slog, clog)

		variations = append(variations, inner...)
		more, err := clog.Reconcile(l)
		if err != nil {
			panic(err)
		}
		variations = append(variations, append(cops, more...))
	}

	return variations
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
		if err != nil {
			panic(err)
		}
		rest1, err1 := clog.Reconcile(l)
		_, rest2, err2 := dot.ReconnectClientLog(l, nil, basisID, "")
		if err1 != nil {
			panic(err1)
		}
		if err2 != nil {
			panic(err2)
		}

		var1 := append(append(append([]dot.Operation{}, rebased...), clientRebased...), rest1...)
		var2 := append(append(append([]dot.Operation{}, rebased...), clientRebased...), rest2...)
		variations = append(variations, var1, var2)
	}

	return variations
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

func (tc journalTestCase) validateLog(t *testing.T, l *dot.Log, input string) {
	b, err := json.Marshal(tc.toCompactForm(input, l.Rebased))
	if err != nil {
		panic(err)
	}
	actual := string(b)
	b, err = json.Marshal(tc.rebased)
	if err != nil {
		panic(err)
	}
	expected := string(b)
	if actual != expected {
		t.Error("Rebased diverged. Expected", expected, "\nActual", actual)
	}
}
