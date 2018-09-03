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
	r := &changes.Replace{Before: S(""), After: S("a")}
	s := &changes.Splice{0, S(""), S("a")}
	m := &changes.Move{0, 5, 1}
	if r.Change() != *r || s.Change() != *s || m.Change() != *m {
		t.Error("Unexpected change failure")
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
	expectPanic("Nil.Count", func() { changes.Nil.Apply(changes.Replace{IsDelete: true}) })
	expectPanic("Nil.Apply", func() { changes.Nil.Apply(changes.Move{0, 5, 1}) })

	if x := changes.Nil.Apply(changes.ChangeSet{changes.PathChange{}}); x != changes.Nil {
		t.Error("Unexpected apply(...)", x)
	}

	x := changes.Nil.Apply(changes.Replace{IsInsert: true, Before: changes.Nil, After: S("a")})
	if x != S("a") {
		t.Error("Unexpected replace", x)
	}

	a := changes.Atomic{nil}
	expectPanic("Atomic.Slice", func() { a.Slice(0, 0) })
	expectPanic("Atomic.Count", func() { a.Count() })
	expectPanic("Atomic.Count", func() { a.Apply(changes.Replace{IsInsert: true}) })
	expectPanic("Atomic.Apply", func() { a.Apply(changes.Move{0, 5, 1}) })

	if x := a.Apply(changes.ChangeSet{changes.PathChange{}}); x != a {
		t.Error("Unexpected apply(...)", x)
	}

	x = a.Apply(changes.Replace{Before: a, After: S("a")})
	if x != S("a") {
		t.Error("Unexpected replace", x)
	}

	z := myChange{}
	expectPanic("myChange1", func() { (changes.Replace{Before: S(""), After: S("a")}).Merge(z) })
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
