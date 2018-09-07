// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/types"
	"reflect"
	"testing"
)

type A = types.A
type S = types.S8

// The random set of tests here are more targeted to cover some rarely
// used codepaths.

func TestNilChangeSet(t *testing.T) {
	_, cx := changes.ChangeSet(nil).Merge(changes.Move{0, 5, 1})
	if cx != nil {
		t.Error("Unexpected non nil cx", cx)
	}

	_, cx = changes.ChangeSet{nil}.Merge(changes.Move{0, 5, 1})
	if cx != nil {
		t.Error("Unexpected non nil cx", cx)
	}

	if x := changes.ChangeSet(nil).Revert(); x != nil {
		t.Error("Unexpected non nil x", x)
	}

	if x := (changes.ChangeSet{nil}).Revert(); x != nil {
		t.Error("Unexpected non nil x", x)
	}
}

func TestSingleChangeSet(t *testing.T) {
	splice1 := changes.Splice{0, S("ab"), S("cd")}
	splice2 := changes.Splice{5, S("ef"), S("gh")}
	cx := changes.ChangeSet{splice1}
	if _, x := cx.Merge(splice2); x != splice1 {
		t.Error("Unexpected single merge behavior", x)
	}

	if cx.Revert() != splice1.Revert() {
		t.Error("Unexpected revert behavior", cx.Revert())
	}

	if x, _ := splice2.Merge(cx); x != splice1 {
		t.Error("Unexpected single merge behavior", x)
	}
}

func TestMultiChangeSet(t *testing.T) {
	splice1 := changes.Splice{0, S("ab"), S("cd")}
	splice2 := changes.Splice{5, S("ef"), S("gh")}
	splice3 := changes.Splice{10, S("a"), S("z")}
	cx := changes.ChangeSet{splice1, splice2}
	if x, _ := splice3.Merge(cx); !reflect.DeepEqual(x, cx) {
		t.Error("Unexpected multi merge", x)
	}
}

func TestChangeMethod(t *testing.T) {
	r := &changes.Replace{S(""), S("a")}
	s := &changes.Splice{0, S(""), S("a")}
	m := &changes.Move{0, 5, 1}
	if r.Change() != *r || s.Change() != *s || m.Change() != *m {
		t.Error("Unexpected change failure")
	}
}

func TestSpliceMapIndex(t *testing.T) {
	s := changes.Splice{5, types.S8("12"), types.S8("1234")}
	idx, ok := s.MapIndex(4)
	if idx != 4 || ok {
		t.Error("Unexpected MapIndex", idx, ok)
	}

	idx, ok = s.MapIndex(5)
	if idx != 5 || !ok {
		t.Error("Unexpected MapIndex", idx, ok)
	}

	idx, ok = s.MapIndex(6)
	if idx != 5 || !ok {
		t.Error("Unexpected MapIndex", idx, ok)
	}

	idx, ok = s.MapIndex(7)
	if idx != 9 || ok {
		t.Error("Unexpected MapIndex", idx, ok)
	}
}

func TestMoveMapIndex(t *testing.T) {
	m1 := changes.Move{2, 3, 4}
	m2 := changes.Move{5, 4, -3}

	mapped := map[int]int{1: 1, 2: 6, 4: 8, 5: 2, 8: 5, 9: 9}
	for before, after := range mapped {
		idx1 := m1.MapIndex(before)
		idx2 := m2.MapIndex(before)
		if idx1 != after || idx2 != after {
			t.Error("MapIndex failed", before, idx1, idx2, after)
		}
	}
}

func TestEmptyAtomicAndUnexpectedChange(t *testing.T) {
	expectPanic := func(msg string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("Failed to panic", msg)
			}
		}()
		fn()
	}
	expectPanic("Nil.Slice", func() { changes.Nil.Slice(0, 0) })
	expectPanic("Nil.Count", func() { changes.Nil.Count() })
	expectPanic("Nil.Replace.Delete", func() { changes.Nil.Apply(changes.Replace{types.S8(""), changes.Nil}) })
	expectPanic("Nil.Apply", func() { changes.Nil.Apply(changes.Move{0, 5, 1}) })

	if x := changes.Nil.Apply(changes.ChangeSet{changes.PathChange{}}); x != changes.Nil {
		t.Error("Unexpected apply(...)", x)
	}

	x := changes.Nil.Apply(changes.Replace{changes.Nil, S("a")})
	if x != S("a") {
		t.Error("Unexpected replace", x)
	}

	a := changes.Atomic{nil}
	expectPanic("Atomic.Slice", func() { a.Slice(0, 0) })
	expectPanic("Atomic.Count", func() { a.Count() })
	expectPanic("Atomic.Apply", func() { a.Apply(changes.Move{0, 5, 1}) })
	expectPanic("Atomic.Create", func() { a.Apply(changes.Replace{changes.Nil, types.S8("")}) })

	if x := a.Apply(changes.ChangeSet{changes.PathChange{}}); x != a {
		t.Error("Unexpected apply(...)", x)
	}

	x = a.Apply(changes.Replace{a, S("a")})
	if x != S("a") {
		t.Error("Unexpected replace", x)
	}

	z := myChange{}
	expectPanic("myChange1", func() { (changes.Replace{S(""), S("a")}).Merge(z) })
	expectPanic("myChange2", func() { (changes.Splice{0, S(""), S("a")}).Merge(z) })
	expectPanic("myChange3", func() { (changes.Move{0, 5, 1}).Merge(z) })

	expectPanic("ApplyTo", func() {
		p := changes.PathChange{[]interface{}{"OK"}, nil}
		p.ApplyTo(S(""))
	})
	p := changes.PathChange{nil, changes.Move{2, 2, 2}}
	if x := p.ApplyTo(S("abcdef")); x != S("abefcd") {
		t.Error("PathChange.ApplyTo failed", x)
	}
}

type myChange struct{}

func (m myChange) Merge(other changes.Change) (ox, cx changes.Change) {
	return nil, nil
}

func (m myChange) Revert() changes.Change {
	return nil
}
