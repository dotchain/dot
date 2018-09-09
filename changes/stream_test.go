// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"github.com/dotchain/dot/changes"
	"reflect"
	"testing"
)

func TestStream(t *testing.T) {
	initial := S("")
	var latest *changes.Stream
	v := changes.Value(initial)

	ev := func(c changes.Change, _ interface{}, l *changes.Stream) {
		v = v.Apply(c)
		latest = l
	}
	s := changes.NewStream().On("boo", ev)
	defer s.On("boo", nil)

	s1 := s.Apply(changes.Replace{changes.Nil, S("Hello World")}, nil)

	c1_1 := s1.Apply(changes.Splice{0, S(""), S("A ")}, nil)
	c1_2 := c1_1.Apply(changes.Splice{2, S(""), S("B ")}, nil)

	c2_1 := s1.Apply(changes.Splice{0, S(""), S("X ")}, nil)
	c2_1_merged := latest

	c2_2 := c2_1.Apply(changes.Splice{2, S(""), S("Y ")}, nil)
	var c2_2_with_c1_1 *changes.Stream
	c2_2.On("boque", func(_ changes.Change, _ interface{}, l *changes.Stream) {
		if c2_2_with_c1_1 == nil {
			c2_2_with_c1_1 = l
		}
	})
	c2_2.On("boque", nil)

	c1_3 := c2_1_merged.Apply(changes.Splice{6, S(""), S("C ")}, nil)
	c2_3 := c2_2_with_c1_1.Apply(changes.Splice{6, S(""), S("Z ")}, nil)

	if !reflect.DeepEqual(v, S("A B X Y C Z Hello World")) {
		t.Error("Merge failed: ", v)
		t.Error("changes", c1_1, c1_2, c1_3)
		t.Error("changes", c2_1, c2_2, c2_3)
	}
}

func TestBranch(t *testing.T) {
	initial := S("")
	v := changes.Value(initial)

	ev := func(c changes.Change, _ interface{}, l *changes.Stream) {
		v = v.Apply(c)
	}
	s := changes.NewStream().On("boo", ev)
	defer s.On("boo", nil)
	s = s.Apply(changes.Replace{changes.Nil, S("Hello World")}, nil)

	child := changes.NewStream()
	branch := changes.Branch{s, child}
	child1 := child.Apply(changes.Splice{0, S(""), S("OK ")}, nil)
	if v != S("Hello World") {
		t.Fatal("Unexpected branch updated", v)
	}
	s = s.Apply(changes.Splice{0, S(""), S("Oh ")}, nil)
	s.Apply(changes.Splice{len("Oh Hello World"), S(""), S("!")}, nil)

	branch.Push()
	if v != S("Oh OK Hello World!") {
		t.Fatal("Unexpected branch updated", v)
	}

	child1.Apply(changes.Splice{len("OK Hello World"), S(""), S("**")}, nil)
	branch.Pull()
	v = changes.Value(S("Hello World"))
	child.On("boq", ev)
	child.On("boq", nil)
	if v != S("Oh OK Hello World!**") {
		t.Fatal("Unexpected branch updated", v)
	}
}

func TestStreamNilChange(t *testing.T) {
	initial := S("")
	v := changes.Value(initial)

	ev := func(c changes.Change, _ interface{}, l *changes.Stream) {
		v = v.Apply(c)
	}
	s := changes.NewStream().On("boo", ev)
	defer s.On("boo", nil)

	child := changes.NewStream()
	branch := changes.Branch{s, child}

	s.Apply(nil, nil)
	child.Apply(nil, nil)
	branch.Merge()

	if v != S("") {
		t.Fatal("Failed merging nil changes", v)
	}
}
