// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package undo_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/undo"
)

func TestSimpleUndoRedo(t *testing.T) {
	upstream := streams.New()
	downstream := undo.New(upstream)

	downstream.Append(changes.Splice{Offset: 10, Before: types.S8(""), After: types.S8("hello")})
	upstream.Append(changes.Splice{Offset: 0, Before: types.S8(""), After: types.S8("abcde")})

	// now undo should rewrite downstream to remove at index 15
	downstream = latest(downstream)
	downstream.Undo()
	downstream, c := downstream.Next()
	expected := changes.Splice{Offset: 15, Before: types.S8("hello"), After: types.S8("")}
	if c != expected {
		t.Fatal("Undo failed", c)
	}

	// now sneak in another upstream op increasing the offset again
	upstream.Append(changes.Splice{Offset: 0, Before: types.S8(""), After: types.S8("abcde")})

	// now redo and confirm that the redo offset is bumped up by 5more
	downstream = latest(downstream)
	downstream.Redo()
	_, c = downstream.Next()
	expected = changes.Splice{Offset: 20, Before: types.S8(""), After: types.S8("hello")}
	if c != expected {
		t.Fatal("Redo failed", c)
	}
}

func TestUndo(t *testing.T) {
	// To make the tests readable, the undo log would consist of
	// letters C, S, U and R to represent local(client), remote(server)
	// undo and redo operations.  A star represents the correct operation
	// is the one that follows it.
	// The lack of a star in the input implies there is no available undo
	// operation
	tests := []string{
		"*C",
		"C*C",
		"S",
		"SS",
		"SS*CSS",
		"CS*CSS",
		"*CCU",
		"S*CCSSSSU",
		"CU*R", // note that redo should be picked!
		"CU*C",
		"CU*CU*R",
		"CSUC*C",
		// no undo possible here
		"CCCUUU",
		"CSCSCSUSUSU",
		"SSSS",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			testUndo(t, name)
		})
	}
}

func TestRedo(t *testing.T) {
	// To make the tests readable, the undo log would consist of
	// letters C, S, U and R to represent local(client), remote(server)
	// undo and redo operations.  A star represents the correct operation
	// is the one that follows it.
	// The lack of a star in the input implies there is no available undo
	// operation
	tests := []string{
		"C*U",
		"C*USSSSSS",
		"CCC*USSSSS",
		"CCCUR*US",
		"CSCSCSUSRSS*USS",
		"CCCUU*U",
		"CCSS*UUSR",
		// No redo possible
		"CUR",
		"CCCCUC",
		"CSUSRS",
		"CSCSCSUSCS",
		"SSSS",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			testRedo(t, name)
		})
	}
}

func TestPushPull(t *testing.T) {
	upstream := streams.New()
	downstream := undo.New(upstream)

	if err := downstream.Push(); err != nil {
		t.Error("Unexpected push error", err)
	}
	if err := downstream.Pull(); err != nil {
		t.Error("Unexpected pull error", err)
	}

}

func testUndo(t *testing.T, test string) {
	upstream := streams.New()
	downstream := undo.New(upstream)

	expected, _ := prepareBranch(upstream, downstream, test)

	downstream = latest(downstream)
	downstream.Undo()
	_, c := downstream.Next()
	if expected == "" {
		if c != nil {
			t.Error("Unexpected non-nil response", c)
		}
		return
	}

	splice, ok := c.(changes.Splice)
	if !ok {
		t.Fatal("Unexpected change type", c)
	}

	if splice.Before != types.S8(expected) || splice.After != types.S8("") {
		t.Error("Failed test", splice, "\nExpected", expected)
	}
}

func testRedo(t *testing.T, test string) {
	upstream := streams.New()
	downstream := undo.New(upstream)
	_, expected := prepareBranch(upstream, downstream, test)

	downstream = latest(downstream)
	downstream.Redo()

	_, c := downstream.Next()
	if expected == "" {
		if c != nil {
			t.Error("Unexpected non-nil response", c)
		}
		return
	}

	splice, ok := c.(changes.Splice)
	if !ok {
		t.Fatal("Unexpected change type", c)
	}

	if splice.After != types.S8(expected) || splice.Before != types.S8("") {
		t.Error("Failed test", splice, "\nExpected", expected)
	}
}

func prepareBranch(upstream, downstream streams.Stream, test string) (string, string) {
	letters := "abcdefghijklmnopqrstuvwxyz"
	ops := []string{}
	for kk, c := range test {
		next := letters[kk : kk+1]
		splice := changes.Splice{Offset: 0, Before: types.S8(""), After: types.S8(next)}
		switch c {
		case 'C':
			latest(downstream).Append(splice)
			ops = append(ops, next)
		case 'U':
			downstream.Undo()
			ops = ops[:len(ops)-1]
		case 'R':
			downstream.Redo()
			ops = ops[:len(ops)+1]
		case 'S':
			latest(upstream).Append(splice)
		}
	}
	last, next := "", ""
	if len(ops) > 0 {
		last = ops[len(ops)-1]
	}
	if cap(ops) > len(ops) {
		ops = ops[:len(ops)+1]
		next = ops[len(ops)-1]
	}

	return last, next
}

func latest(s streams.Stream) streams.Stream {
	for v, _ := s.Next(); v != nil; v, _ = s.Next() {
		s = v
	}
	return s
}
