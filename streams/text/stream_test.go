// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package text_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/text"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

func TestStream(t *testing.T) {
	t.Run("Use16=false", streamSuite(false).Run)
	t.Run("Use16=true", streamSuite(true).Run)
}

type streamSuite bool

func (suite streamSuite) Run(t *testing.T) {
	t.Run("Append", suite.testAppend)
	t.Run("ReverseAppend", suite.testReverseAppend)
	t.Run("Scheduler", suite.testScheduler)
	t.Run("CollapsedSelection", suite.testCollapsedSelection)
	t.Run("NonCollapsedSelection", suite.testNonCollapsedSelection)
}

func (suite streamSuite) testAppend(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	change := changes.PathChange{[]interface{}{"Value"}, changes.Move{0, 1, 1}}
	after := s.Append(change)
	suite.validate(t, s, after.(*text.Stream))

	_, sx := s.Next()
	if _, x := sx.Next(); x != nil {
		t.Error("Unexpected non-nil next", x)
	}

	after = sx.Append(changes.Replace{s.E, types.S8("okok")})
	vs, ok := after.(*streams.ValueStream)
	if !ok || vs.Value != types.S8("okok") {
		t.Error("Unexpected replace result", after)
	}
	if _, x := sx.Next(); !reflect.DeepEqual(after, x) {
		t.Error("Unexpected divergence", x)
	}
}

func (suite streamSuite) testReverseAppend(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	change := changes.PathChange{[]interface{}{"Value"}, changes.Move{0, 1, 1}}
	after := s.ReverseAppend(change)
	suite.validate(t, s, after.(*text.Stream))

	_, sx := s.Next()
	if _, x := sx.Next(); x != nil {
		t.Error("Unexpected non-nil next", x)
	}

	after = sx.ReverseAppend(changes.Replace{s.E, types.S8("okok")})
	vs, ok := after.(*streams.ValueStream)
	if !ok || vs.Value != types.S8("okok") {
		t.Error("Unexpected replace result", after)
	}
	if _, x := sx.Next(); !reflect.DeepEqual(after, x) {
		t.Error("Unexpected divergence", x)
	}
}

func (suite streamSuite) testScheduler(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))
	async := &streams.AsyncScheduler{}
	if x := s.WithScheduler(async).Scheduler(); x != async {
		t.Error("Scheduler change didn't take", x)
	}
}

func (suite streamSuite) testCollapsedSelection(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))

	// test caret
	after := s.SetSelection(3, 3, false)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 3 || left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 3 || left {
		t.Error("Unexpected end", idx, left)
	}

	s = after
	after = s.SetSelection(3, 3, true)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 3 || !left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 3 || !left {
		t.Error("Unexpected end", idx, left)
	}
}

func (suite streamSuite) testNonCollapsedSelection(t *testing.T) {
	s := text.StreamFromString("Hello", bool(suite))

	after := s.SetSelection(3, 5, true)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 3 || left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 5 || !left {
		t.Error("Unexpected end", idx, left)
	}

	s = after
	after = s.SetSelection(5, 3, true)
	suite.validate(t, s, after)
	if idx, left := after.E.Start(); idx != 5 || !left {
		t.Error("Unexpected start", idx, left)
	}
	if idx, left := after.E.End(); idx != 3 || left {
		t.Error("Unexpected end", idx, left)
	}
}

func (suite streamSuite) validate(t *testing.T, before, after *text.Stream) {
	if _, next := before.Next(); !reflect.DeepEqual(next, after) {
		t.Error("Divergent change", next.(*text.Stream).E, "x", after.E)
	}
	var next streams.Stream
	before.Nextf("validate", func(_ changes.Change, str streams.Stream) {
		before.Nextf("validate", nil)
		next = str
	})
	if !reflect.DeepEqual(next, after) {
		t.Error("Divergent change")
	}
}
