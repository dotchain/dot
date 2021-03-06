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

type subsuite bool

func TestSubstream(t *testing.T) {
	subsuite(true).run(t)
	subsuite(false).run(t)
}

func (s subsuite) run(t *testing.T) {
	reverse := ""
	if s {
		reverse = "Reverse"
	}
	t.Run("FieldAppend"+reverse, s.FieldAppend)
	t.Run("ParentAppendOther"+reverse, s.ParentAppendOther)
	t.Run("ParentAppendOwn"+reverse, s.ParentAppendOwn)
	t.Run("InvalidRef"+reverse, s.InvalidRef)
}

func (s subsuite) ChangingIndex(t *testing.T) {
	parent := streams.New()
	child := streams.Substream(parent, 5)
	parent = parent.Append(changes.Splice{Offset: 2, Before: types.A{}, After: types.A{types.S16("ok"), types.S16("boo")}})
	child, cx := child.Next()
	if cx != nil {
		t.Fatal("got unexpected change", cx)
	}

	// now child should affect index 7
	c := changes.Replace{Before: types.S16("b"), After: types.S16("a")}
	child.Append(c)
	parent, cx = parent.Next()
	if !reflect.DeepEqual(parent, changes.PathChange{Path: []interface{}{7}, Change: c}) {
		t.Error("Index did not change", cx)
	}
}

func (s subsuite) InvalidRef(t *testing.T) {
	parent := streams.New()
	child := streams.Substream(parent, "boo")
	c := changes.Replace{Before: types.S16("yoo"), After: changes.Nil}
	parent = parent.Append(c)
	nn, cc := child.Next()
	if nn == nil || cc != nil {
		t.Fatal("Unexpected response", nn, cc)
	}

	c2 := changes.Replace{Before: changes.Nil, After: types.S16("goo")}
	parent.Append(changes.PathChange{Path: []interface{}{"boo"}, Change: c2})

	nn2, cc2 := nn.Next()
	if nn2 != nil || cc2 != nil {
		t.Error("Invalid ref didn't do its thing", nn2, cc2)
	}

	if nn2 := nn.Append(c); nn2 != nn {
		t.Error("Append() on invalidRef", nn2)
	}

	if nn2 := nn.ReverseAppend(c); nn2 != nn {
		t.Error("Append() on invalidRef", nn2)
	}
}

func (s subsuite) ParentAppendOther(t *testing.T) {
	parent := streams.New()
	child := streams.Substream(parent, "boo")
	if x, c := child.Next(); x != nil || c != nil {
		t.Fatal("Next() on empty stream", x, c)
	}

	c := changes.Replace{Before: changes.Nil, After: types.S16("yoo")}
	pc := changes.PathChange{Path: []interface{}{"goo"}, Change: c}
	if s {
		parent.ReverseAppend(pc)
	} else {
		parent.Append(pc)
	}
	_, nextc := child.Next()
	if !reflect.DeepEqual(nextc, nil) {
		t.Error("append parent didn't do the expected", nextc)
	}
}

func (s subsuite) ParentAppendOwn(t *testing.T) {
	parent := streams.New()
	child := streams.Substream(parent, "boo")
	if x, c := child.Next(); x != nil || c != nil {
		t.Fatal("Next() on empty stream", x, c)
	}

	c := changes.Replace{Before: changes.Nil, After: types.S16("yoo")}
	pc := changes.PathChange{Path: []interface{}{"boo"}, Change: c}
	if s {
		parent.ReverseAppend(pc)
	} else {
		parent.Append(pc)
	}
	_, nextc := child.Next()
	if !reflect.DeepEqual(changes.Simplify(nextc), c) {
		t.Error("append parent didn't do the expected", nextc)
	}

}

func (s subsuite) FieldAppend(t *testing.T) {
	parent := streams.New()
	child := streams.Substream(parent, "boo")
	c := changes.Replace{Before: changes.Nil, After: types.S16("yoo")}
	if s {
		child.ReverseAppend(c)
	} else {
		child.Append(c)
	}
	_, nextc := parent.Next()
	if !reflect.DeepEqual(nextc, changes.PathChange{Path: []interface{}{"boo"}, Change: c}) {
		t.Error("append parent didn't do the expected", nextc)
	}
}
