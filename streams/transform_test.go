// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

type xformsuite struct{}

func TestTransform(t *testing.T) {
	x := xformsuite{}
	t.Run("ParentAppend", x.ParentAppend)
	t.Run("Append", x.Append)
	t.Run("ReverseAppend", x.ReverseAppend)
	t.Run("Next", x.Next)
}

func (_ xformsuite) ParentAppend(t *testing.T) {
	parent := streams.New()
	child := streams.Transform(parent, toInt, nil)

	c := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: 2.1}}

	parent.Append(c)

	_, c1 := parent.Next()
	if !reflect.DeepEqual(c1, c) {
		t.Fatal("Unexpected stream implementation", c1)
	}

	_, c2 := child.Next()
	if !reflect.DeepEqual(c2, c) {
		t.Error("Append yielded", c2)
	}
}

func (_ xformsuite) Append(t *testing.T) {
	parent := streams.New()
	child := streams.Transform(parent, toInt, nil)

	c := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: 2.1}}
	expected := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: int(2)}}

	child.Append(c)

	_, c1 := parent.Next()
	if !reflect.DeepEqual(c1, expected) {
		t.Fatal("Unexpected stream implementation", c1)
	}

	_, c2 := child.Next()
	if !reflect.DeepEqual(c2, expected) {
		t.Error("Append yielded", c2)
	}
}

func (_ xformsuite) ReverseAppend(t *testing.T) {
	parent := streams.New()
	child := streams.Transform(parent, toInt, nil)

	c := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: 2.1}}
	child.ReverseAppend(c)

	_, c1 := parent.Next()
	if !reflect.DeepEqual(c1, c) {
		t.Fatal("Unexpected stream implementation", c1)
	}

	_, c2 := child.Next()
	if !reflect.DeepEqual(c2, c) {
		t.Error("Append yielded", c2)
	}
}

func (_ xformsuite) Next(t *testing.T) {
	parent := streams.New()
	child := streams.Transform(parent, nil, toInt)

	if x, c := child.Next(); x != nil || c != nil {
		t.Error("Unexpected next", x, c)
	}

	c := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: 2.1}}
	expected := changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: int(2)}}
	parent.Append(c)

	_, c1 := parent.Next()
	if !reflect.DeepEqual(c1, c) {
		t.Fatal("Unexpected stream implementation", c1)
	}

	_, c2 := child.Next()
	if !reflect.DeepEqual(c2, expected) {
		t.Error("Append yielded", c2)
	}

}

func toInt(c changes.Change) changes.Change {
	convert := func(v changes.Value) changes.Value {
		a, ok := v.(changes.Atomic)
		if !ok {
			return v
		}
		f, ok := a.Value.(float64)
		if !ok {
			return v
		}
		return changes.Atomic{Value: int(f)}
	}

	if r, ok := c.(changes.Replace); ok {
		return changes.Replace{Before: convert(r.Before), After: convert(r.After)}
	}
	return nil
}
