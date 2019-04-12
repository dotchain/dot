// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

type S = types.S8

func TestStream(t *testing.T) {
	initial := S("")
	var latest streams.Stream
	v := changes.Value(initial)

	ev := func() {
		var c changes.Change
		latest, c = latest.Next()
		v = v.Apply(nil, c)
	}
	s := streams.New()
	latest = s
	s.Nextf("boo", ev)
	defer s.Nextf("boo", nil)

	s1 := s.Append(changes.Replace{Before: changes.Nil, After: S("Hello World")})

	c1_1 := s1.Append(changes.Splice{Offset: 0, Before: S(""), After: S("A ")})
	c1_2 := c1_1.Append(changes.Splice{Offset: 2, Before: S(""), After: S("B ")})

	c2_1 := s1.Append(changes.Splice{Offset: 0, Before: S(""), After: S("X ")})
	c2_1_merged := latest

	c2_2 := c2_1.Append(changes.Splice{Offset: 2, Before: S(""), After: S("Y ")})
	c2_2_with_c1_1, _ := c2_2.Next()

	c1_3 := c2_1_merged.Append(changes.Splice{Offset: 6, Before: S(""), After: S("C ")})
	c2_3 := c2_2_with_c1_1.Append(changes.Splice{Offset: 6, Before: S(""), After: S("Z ")})

	if !reflect.DeepEqual(v, S("A B X Y C Z Hello World")) {
		t.Error("Merge failed: ", v)
		t.Error("changes", c1_1, c1_2, c1_3)
		t.Error("changes", c2_1, c2_2, c2_3)
	}
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

	streams.Push(child)

	_, c = streams.Latest(s)
	if v := initial.Apply(nil, c); v != types.S8("Oh OK Hello World!") {
		t.Fatal("Unexpected branch updated", v)
	}

	child1.Append(changes.Splice{Offset: len("OK Hello World"), Before: S(""), After: S("**")})
	streams.Pull(child)

	_, c = streams.Latest(child)
	if v := cInitial.Apply(nil, c); v != types.S8("Oh OK Hello World!**") {
		t.Fatal("Unexpected branch updated", v)
	}
}

func TestConnectedBranches(t *testing.T) {
	var master changes.Value = S("")
	var local changes.Value = S("")

	bm := streams.New()
	bl := streams.New()
	bm.Nextf("key", func() {
		var c changes.Change
		bm, c = bm.Next()
		master = master.Apply(nil, c)
	})
	bl.Nextf("key", func() {
		var c changes.Change
		bl, c = bl.Next()
		local = local.Apply(nil, c)
	})

	streams.Connect(bm, bl)
	bl.Append(changes.Splice{Offset: 0, Before: S(""), After: S("OK")})
	if master != S("OK") || local != S("OK") {
		t.Fatal("Unexpected master, local", master, local)
	}

	bm.Append(changes.Splice{Offset: 2, Before: S(""), After: S(" Computer")})
	if master != S("OK Computer") || local != S("OK Computer") {
		t.Fatal("Unexpected master, local", master, local)
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

	streams.Pull(child)
	if _, c := child.Next(); c != (changes.Move{Offset: 2, Count: 2, Distance: 2}) {
		t.Error("Branch move unexpected change", c)
	}

	streams.Pull(grandChild)
	if _, c := grandChild.Next(); c != (changes.Move{Offset: 2, Count: 2, Distance: 2}) {
		t.Error("Branch move unexpected change", c)
	}
}

func TestStreamNilChange(t *testing.T) {
	initial := S("")
	v := changes.Value(initial)

	var latest streams.Stream
	ev := func() {
		var c changes.Change
		latest, c = latest.Next()
		v = v.Apply(nil, c)
	}
	s := streams.New()
	latest = s
	s.Nextf("boo", ev)
	defer s.Nextf("boo", nil)

	child := streams.Branch(s)

	s.Append(nil)
	child.Append(nil)
	streams.Push(child)
	streams.Pull(child)

	if v != S("") {
		t.Fatal("Failed merging nil changes", v)
	}
}
