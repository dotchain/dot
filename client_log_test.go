// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"encoding/json"
	"fmt"
	datalib "github.com/dotchain/dataset/tools/lib"
	"github.com/dotchain/dot"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

type clientLogTestCase struct {
	input, final string
	cops, sops   []string
}

func TestClientLog_suite(t *testing.T) {
	tests := []clientLogTestCase{
		{
			input: "abcd",
			final: "axYyXd",
			cops:  []string{"a(=x)bcd", "axb(c=y)d"},
			sops:  []string{"abc(=X)d", "a(b=Y)cXd"},
		},
		{
			input: "abcd",
			final: "axYyXd",
			cops:  []string{"abc(=X)d", "a(b=Y)cXd"},
			sops:  []string{"a(=x)bcd", "axb(c=y)d"},
		},
		{
			input: "abc",
			final: "y",
			cops:  []string{"a(b=XY)c", "aX(Yc=)"},
			sops:  []string{"a(b=xy)c", "(ax=)yc"},
		},
	}

	c := clientLogTestSuite{}
	for _, test := range tests {
		c.Run(t, test)
	}
}

type clientLogTestSuite struct{}

func (c clientLogTestSuite) Run(t *testing.T, testCase clientLogTestCase) {
	err := c.doublePairsTest(testCase.input, testCase.final, testCase.cops, testCase.sops)
	if err != nil {
		t.Error("Failed", err)
	}
}

// doublePairsTest works with a pair of client ops and a pair of server ops
//
// The client and server both start with the initial string.  The server
// applies sop1, sop2 followed by cop1, cop2 while the client follows
//    cop1, cop2, sop1, sop2
//
// All the operations are encoded using the Compact form as defined in
// inner compact form of http://github.com/dotchain/dataset/CompactJSON.md
func (c clientLogTestSuite) doublePairsTest(input, final string, cx, sx []string) error {
	sops := []dot.Operation{c.decode("s1", nil, sx[0], input)}
	for kk := range sx {
		if kk == 0 {
			continue
		}
		id := fmt.Sprintf("s%d", kk+1)
		prev := fmt.Sprintf("s%d", kk)
		parents := []string{prev}
		sops = append(sops, c.decode(id, parents, sx[kk], ""))
	}
	cops := []dot.Operation{c.decode("c1", nil, cx[0], input)}
	for kk := range cx {
		if kk == 0 {
			continue
		}
		id := fmt.Sprintf("c%d", kk+1)
		prev := fmt.Sprintf("c%d", kk)
		parents := []string{"", prev}
		cops = append(cops, c.decode(id, parents, cx[kk], ""))
	}

	copsVariations, sops1, err := c.getDoublePairOps(input, cops, sops)
	if err != nil {
		return err
	}
	for _, variation := range copsVariations {
		if err := c.validateOps(input, final, variation); err != nil {
			return errors.Wrap(err, "client validation")
		}
	}

	if err := c.validateOps(input, final, sops1); err != nil {
		return errors.Wrap(err, "server validation")
	}
	return nil
}

func (c clientLogTestSuite) getDoublePairOps(input interface{}, cops, sops []dot.Operation) ([][]dot.Operation, []dot.Operation, error) {
	all := append(append([]dot.Operation{}, sops...), cops...)

	// all client variations where the client "reconnects" at various points
	variation, err := c.getReconnectOps(all, "", "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "variation empty, empty")
	}
	variations := [][]dot.Operation{variation}

	slog := &dot.Log{}
	clog, rx, cx, err := dot.BootstrapClientLog(slog, cops)
	if err != nil {
		return nil, nil, errors.Wrap(err, "bootstrap")
	}
	cApply := append(append([]dot.Operation{}, rx...), cx...)
	lastParentID := cops[len(cops)-1].ID

	variation, err = c.getReconnectOps(all, "", lastParentID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "variation empty, last")
	}
	variations = append(variations, append(append([]dot.Operation{}, cApply...), variation...))

	lastBasisID := ""
	for _, op := range all {
		if err := slog.AppendOperation(op); err != nil {
			return nil, nil, errors.Wrap(err, "slog.AppendOperation")
		}
		xop, err := clog.Reconcile(slog)
		if err != nil {
			return nil, nil, errors.Wrap(err, "clog.Reconcile")
		}
		cApply = append(cApply, xop...)

		if len(xop) > 0 {
			lastBasisID = xop[len(xop)-1].ID
		}
		variation, err = c.getReconnectOps(all, lastBasisID, lastParentID)
		if err != nil {
			return nil, nil, errors.Wrap(err, fmt.Sprintf("variation %s, %s", lastBasisID, lastParentID))
		}
		variation = append(append([]dot.Operation{}, cApply...), variation...)
		variations = append(variations, variation)
	}

	variations = append(variations, cApply)
	return variations, slog.Rebased, nil
}

func (clientLogTestSuite) getReconnectOps(all []dot.Operation, basisID, parentID string) ([]dot.Operation, error) {
	slog := &dot.Log{}
	for _, op := range all {
		if err := slog.AppendOperation(op); err != nil {
			return nil, errors.Wrap(err, "slog.AppendOperation")
		}
	}

	_, rx, err := dot.ReconnectClientLog(slog, nil, basisID, parentID)
	if err != nil {
		return nil, errors.Wrap(err, "reconnect")
	}

	return rx, nil
}

func (clientLogTestSuite) decode(id string, parents []string, op, input string) dot.Operation {
	input1, ch := datalib.Compact{}.Decode(op)
	if input != "" && input1 != input {
		panic(fmt.Sprintf("Expdected %v but got %v", input, input1))
	}
	return dot.Operation{ID: id, Parents: parents, Changes: []dot.Change{ch}}
}

func (clientLogTestSuite) toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (c clientLogTestSuite) validateOps(input, final interface{}, ops []dot.Operation) error {
	actual := input
	for _, op := range ops {
		actual = applyMany(actual, op.Changes)
	}

	if !dot.Utils(dot.Transformer{}).AreSame(actual, final) {
		return errors.Errorf("Mismatched: Expected %s, got %s", c.toJSON(final), c.toJSON(actual))
	}
	return nil
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
	fourthMergeChain := []dot.Operation{{ID: "third", Parents: []string{"two"}}}
	fifthMergeChain := []dot.Operation{{ID: "third"}}
	l := &dot.Log{
		MinIndex:    2,
		Rebased:     []dot.Operation{{}, {}, {ID: "third"}, {ID: "fourth"}, {ID: "fifth"}, {ID: "sixth"}},
		MergeChains: [][]dot.Operation{nil, nil, nil, fourthMergeChain, fifthMergeChain, nil},
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
		ServerIndex: 2,
		Rebased:     []dot.Operation{{ID: "fourth", Parents: []string{"second"}}},
		MergeChain:  fourthMergeChain,
	}

	op := dot.Operation{ID: "fifth", Parents: []string{"second", "fourth"}}
	merge, err := c.AppendClientOperation(l, op)

	if err != nil {
		t.Error("Unexpected AppendClientOperation response", err)
	}
	if !reflect.DeepEqual(merge, []dot.Operation{{ID: "third"}, {ID: "sixth"}}) {
		t.Error("Unexpected AppendClientOperation merge", merge)
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
