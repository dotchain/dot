// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

type S = types.S8

func TestStream(t *testing.T) {
	initial := S("")
	var latest streams.Stream
	v := changes.Value(initial)

	ev := func(c changes.Change, l streams.Stream) {
		v = v.Apply(c)
		latest = l
	}
	s := streams.New()
	s.Nextf("boo", ev)
	defer s.Nextf("boo", nil)

	s1 := s.Append(changes.Replace{changes.Nil, S("Hello World")})

	c1_1 := s1.Append(changes.Splice{0, S(""), S("A ")})
	c1_2 := c1_1.Append(changes.Splice{2, S(""), S("B ")})

	c2_1 := s1.Append(changes.Splice{0, S(""), S("X ")})
	c2_1_merged := latest

	c2_2 := c2_1.Append(changes.Splice{2, S(""), S("Y ")})
	_, c2_2_with_c1_1 := c2_2.Next()

	c1_3 := c2_1_merged.Append(changes.Splice{6, S(""), S("C ")})
	c2_3 := c2_2_with_c1_1.Append(changes.Splice{6, S(""), S("Z ")})

	if !reflect.DeepEqual(v, S("A B X Y C Z Hello World")) {
		t.Error("Merge failed: ", v)
		t.Error("changes", c1_1, c1_2, c1_3)
		t.Error("changes", c2_1, c2_2, c2_3)
	}
}

func TestBranch(t *testing.T) {
	initial := S("")
	v := changes.Value(initial)

	ev := func(c changes.Change, l streams.Stream) {
		v = v.Apply(c)
	}
	s := streams.New()
	s.Nextf("boo", ev)
	defer s.Nextf("boo", nil)
	s = s.Append(changes.Replace{changes.Nil, S("Hello World")})

	child := streams.New()
	branch := &streams.Branch{s, child}
	child1 := child.Append(changes.Splice{0, S(""), S("OK ")})
	if v != S("Hello World") {
		t.Fatal("Unexpected branch updated", v)
	}
	s = s.Append(changes.Splice{0, S(""), S("Oh ")})
	s.Append(changes.Splice{len("Oh Hello World"), S(""), S("!")})

	branch.Push()
	if v != S("Oh OK Hello World!") {
		t.Fatal("Unexpected branch updated", v)
	}

	child1.Append(changes.Splice{len("OK Hello World"), S(""), S("**")})
	branch.Pull()
	v = changes.Value(S("Hello World"))
	child.Nextf("boq", ev)
	child.Nextf("boq", nil)
	if v != S("Oh OK Hello World!**") {
		t.Fatal("Unexpected branch updated", v)
	}
}

func TestConnectedBranches(t *testing.T) {
	var master changes.Value = S("")
	var local changes.Value = S("")

	b := streams.Branch{streams.New(), streams.New()}
	b.Master.Nextf("key", func(c changes.Change, _ streams.Stream) {
		master = master.Apply(c)
	})
	b.Local.Nextf("key", func(c changes.Change, _ streams.Stream) {
		local = local.Apply(c)
	})

	b.Connect()
	b.Local.Append(changes.Splice{0, S(""), S("OK")})
	if master != S("OK") || local != S("OK") {
		t.Fatal("Unexpected master, local", master, local)
	}

	b.Master.Append(changes.Splice{2, S(""), S(" Computer")})
	if master != S("OK Computer") || local != S("OK Computer") {
		t.Fatal("Unexpected master, local", master, local)
	}

	b.Disconnect()
	b.Master.Append(changes.Splice{2, S(""), S("!")})
	b.Local.Append(changes.Splice{11, S(""), S("s")})

	if master != S("OK! Computer") || local != S("OK Computers") {
		t.Fatal("Unexpected master, local", master, local)
	}

	b.Connect()
	if master != S("OK! Computers") || local != master {
		t.Fatal("Unexpected master, local", master, local)
	}
}

func TestStreamNilChange(t *testing.T) {
	initial := S("")
	v := changes.Value(initial)

	ev := func(c changes.Change, l streams.Stream) {
		v = v.Apply(c)
	}
	s := streams.New()
	s.Nextf("boo", ev)
	defer s.Nextf("boo", nil)

	child := streams.New()
	branch := streams.Branch{s, child}

	s.Append(nil)
	child.Append(nil)
	branch.Merge()

	if v != S("") {
		t.Fatal("Failed merging nil changes", v)
	}
}
