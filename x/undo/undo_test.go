// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package undo_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/types"
	"github.com/dotchain/dot/x/undo"
	"testing"
)

func TestNextf(t *testing.T) {
	orig := changes.NewStream()
	downstream, stack := undo.New(orig)
	defer stack.Close()

	count := 0
	downstream.Nextf("key", func(c changes.Change, _ changes.Stream) {
		count++
	})
	orig.Append(changes.Move{1, 2, 3})
	orig.Append(changes.Move{2, 3, 4})
	downstream.Nextf("key", nil)
	orig.Append(changes.Move{4, 5, 6})
	if count != 2 {
		t.Fatal("Nextf did not proxy as expected", count)
	}
}

func TestSimpleUndoRedo(t *testing.T) {
	upstream := changes.NewStream()
	downstream, stack := undo.New(changes.NewStream())
	b := &changes.Branch{upstream, downstream}

	downstream.Append(changes.Splice{10, types.S8(""), types.S8("hello")})
	upstream.Append(changes.Splice{0, types.S8(""), types.S8("abcde")})
	b.Merge()

	// now undo should rewrite downstream to remove at index 15
	stack.Undo()
	c, _ := b.Local.Next()
	expected := changes.Splice{15, types.S8("hello"), types.S8("")}
	if c != expected {
		t.Fatal("Undo failed", c)
	}

	// now sneak in another upstream op increasing the offset again
	upstream.Append(changes.Splice{0, types.S8(""), types.S8("abcde")})
	b.Merge()

	// now redo and confirm that the redo offset is bumped up by 5more
	stack.Redo()
	c, _ = b.Local.Next()
	expected = changes.Splice{20, types.S8(""), types.S8("hello")}
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

func testUndo(t *testing.T, test string) {
	upstream := changes.NewStream()
	downstream, stack := undo.New(changes.NewStream())
	defer stack.Close()
	b := &changes.Branch{upstream, downstream}
	expected, _ := prepareBranch(b, stack, test)
	stack.Undo()
	c, _ := b.Local.Next()
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
	upstream := changes.NewStream()
	downstream, stack := undo.New(changes.NewStream())
	defer stack.Close()
	b := &changes.Branch{upstream, downstream}
	_, expected := prepareBranch(b, stack, test)
	stack.Redo()

	c, _ := b.Local.Next()
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

func prepareBranch(b *changes.Branch, stack undo.Stack, test string) (string, string) {
	letters := "abcdefghijklmnopqrstuvwxyz"
	ops := []string{}
	for kk, c := range test {
		next := letters[kk : kk+1]
		splice := changes.Splice{0, types.S8(""), types.S8(next)}
		switch c {
		case 'C':
			b.Local.Append(splice)
			ops = append(ops, next)
		case 'U':
			stack.Undo()
			ops = ops[:len(ops)-1]
		case 'R':
			stack.Redo()
			ops = ops[:len(ops)+1]
		case 'S':
			b.Master.Append(splice)
		}
		b.Merge()
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