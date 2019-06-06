// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/crdt"
	"github.com/dotchain/dot/changes/types"
)

func TestSeq(t *testing.T) {
	s := crdt.Seq{}
	_, s = s.Splice(0, 0, []interface{}{"hello", "new", "world"})
	if x := s.Items(); !reflect.DeepEqual(x, []interface{}{"hello", "new", "world"}) {
		t.Fatal("Splice failed", x)
	}

	c1, s1 := s.Splice(0, 1, []interface{}{"Hello"})
	if x := s1.Items(); !reflect.DeepEqual(x, []interface{}{"Hello", "new", "world"}) {
		t.Fatal("Splice failed", x)
	}

	s2 := s1.Apply(nil, c1.Revert()).(crdt.Seq)
	if x := s2.Items(); !reflect.DeepEqual(x, s.Items()) {
		t.Fatal("Undo Splice failed", x)
	}

	_, s2 = s.Splice(1, 1, nil)
	if x := s2.Items(); !reflect.DeepEqual(x, []interface{}{"hello", "world"}) {
		t.Fatal("Splice failed", x)
	}

	_, s2 = s.Splice(1, 2, nil)
	if x := s2.Items(); !reflect.DeepEqual(x, []interface{}{"hello"}) {
		t.Fatal("Splice failed", x)
	}
}

func TestSeqMove(t *testing.T) {
	s := crdt.Seq{}
	_, s = s.Splice(0, 0, []interface{}{"hello", "new", "world"})
	if x := s.Items(); !reflect.DeepEqual(x, []interface{}{"hello", "new", "world"}) {
		t.Fatal("Move failed", x)
	}

	c1, s1 := s.Move(1, 1, -1)
	if x := s1.Items(); !reflect.DeepEqual(x, []interface{}{"new", "hello", "world"}) {
		t.Fatal("Move failed", x)
	}

	s2 := s.Apply(nil, c1).(crdt.Seq)
	if x := s2.Items(); !reflect.DeepEqual(x, s1.Items()) {
		t.Fatal("Move failed", x)
	}

}

func TestSeqMove2(t *testing.T) {
	s := crdt.Seq{}
	_, s = s.Splice(0, 0, []interface{}{"hello", "new", "world"})
	if x := s.Items(); !reflect.DeepEqual(x, []interface{}{"hello", "new", "world"}) {
		t.Fatal("Move failed", x)
	}

	c1, s1 := s.Move(0, 1, 1)
	if x := s1.Items(); !reflect.DeepEqual(x, []interface{}{"new", "hello", "world"}) {
		t.Fatal("Move failed", x)
	}

	s2 := s.Apply(nil, c1).(crdt.Seq)
	if x := s2.Items(); !reflect.DeepEqual(x, s1.Items()) {
		t.Fatal("Move failed", x)
	}

}

func TestSeqUpdate(t *testing.T) {
	s := crdt.Seq{}
	_, s = s.Splice(0, 0, []interface{}{types.S16("hello")})

	c1, s1 := s.Update(0, changes.Splice{Before: types.S16("h"), After: types.S16("H")})
	if x := s1.Items()[0]; x != types.S16("Hello") {
		t.Error("inner splice didn't work", x)
	}

	if x := s.Apply(nil, c1).(crdt.Seq); !reflect.DeepEqual(x.Items(), s1.Items()) {
		t.Error("Apply diverged", x)
	}

	s2 := s1.Apply(nil, c1.Revert()).(crdt.Seq)
	if x := s2.Items(); !reflect.DeepEqual(x, s.Items()) {
		t.Error("Apply(revert) diverged", x)
	}

	// test merges of c1 with c1.revert
	c2 := c1.Revert()
	x1, x2 := c1.Merge(c2)
	if !reflect.DeepEqual(x2, c1) || !reflect.DeepEqual(x1, c2) {
		t.Error("Merge had unexpected behavior", c1, c2)
	}

	x1, x2 = c1.(changes.Custom).ReverseMerge(c2)
	if !reflect.DeepEqual(x2, c1) || !reflect.DeepEqual(x1, c2) {
		t.Error("Merge had unexpected behavior", c1, c2)
	}
}
