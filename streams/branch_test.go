// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func TestBranchUndoRedo(t *testing.T) {
	s := streams.New()
	b := streams.Branch(s)
	b.Undo()
	b.Redo()
}

func TestBranch(t *testing.T) {
	initial := types.S8("")
	s := streams.New()
	s1 := s.Append(changes.Replace{Before: changes.Nil, After: types.S8("Hello World")})

	child := streams.Branch(s1)
	cInitial := types.S8("Hello World")

	child1 := child.Append(changes.Splice{Before: types.S8(""), After: types.S8("OK ")})

	_, c := streams.Latest(s)
	if v := initial.Apply(nil, c); v != types.S8("Hello World") {
		t.Fatal("Unexpected branch updated", v)
	}

	s2 := s1.Append(changes.Splice{Offset: 0, Before: S(""), After: S("Oh ")})
	s2.Append(changes.Splice{Offset: len("Oh Hello World"), Before: S(""), After: S("!")})

	if err := child.Push(); err != nil {
		t.Error("Push failed", err)
	}

	_, c = streams.Latest(s)
	if v := initial.Apply(nil, c); v != types.S8("Oh OK Hello World!") {
		t.Fatal("Unexpected branch updated", v)
	}

	child1.Append(changes.Splice{Offset: len("OK Hello World"), Before: S(""), After: S("**")})
	if err := child.Pull(); err != nil {
		t.Error("Pull failed", err)
	}

	_, c = streams.Latest(child)
	if v := cInitial.Apply(nil, c); v != types.S8("Oh OK Hello World!**") {
		t.Fatal("Unexpected branch updated", v)
	}
}

func TestBranchReverseAppend(t *testing.T) {
	master := streams.New()
	child := streams.Branch(master)

	child2 := child.Append(changes.Move{Offset: 2, Count: 2, Distance: 2})
	replace := changes.Replace{Before: types.S8("1234567"), After: types.S8("boo")}
	child.ReverseAppend(replace)

	_, c1 := child2.Next()
	_, expected := replace.Merge(changes.Move{Offset: 2, Count: 2, Distance: 2})
	if c1 != expected {
		t.Error("Unexpected branch behavior", c1, expected)
	}
}

func TestDoubleBranches(t *testing.T) {
	master := streams.New()
	child := streams.Branch(master)
	grandChild := streams.Branch(child)

	master.Append(changes.Move{Offset: 2, Count: 2, Distance: 2})

	if x, _ := child.Next(); x != nil {
		t.Error("Branch merged too soon", x)
	}

	if err := child.Pull(); err != nil {
		t.Error("Pull failed", err)
	}

	if _, c := child.Next(); c != (changes.Move{Offset: 2, Count: 2, Distance: 2}) {
		t.Error("Branch move unexpected change", c)
	}

	if err := grandChild.Pull(); err != nil {
		t.Error("pull failed", err)
	}

	if _, c := grandChild.Next(); c != (changes.Move{Offset: 2, Count: 2, Distance: 2}) {
		t.Error("Branch move unexpected change", c)
	}
}

func TestBranchNilChange(t *testing.T) {
	s := streams.New()
	child := streams.Branch(s)

	s.Append(nil)
	child.Append(nil)
	if err := child.Push(); err != nil {
		t.Error("push", err)
	}
	if err := child.Pull(); err != nil {
		t.Error("pull", err)
	}

	_, c := streams.Latest(s)
	if v := types.S8("").Apply(nil, c); v != types.S8("") {
		t.Fatal("Failed merging nil changes", v)
	}
}
