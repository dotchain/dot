// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"testing"
)

func TestCounterApply(t *testing.T) {
	c := types.Counter(5)
	if x := c.Apply(nil, nil); x != c {
		t.Error("Apply(nil, nil)", x)
	}

	if x := c.Apply(nil, changes.Replace{c, changes.Nil}); x != changes.Nil {
		t.Error("Replace(IsDelete)", x)
	}

	if x := c.Apply(nil, changes.Replace{c, types.S8("OK")}); x != types.S8("OK") {
		t.Error("Replace()", x)
	}

	if x := c.Apply(nil, c.Increment(2)); x != c+2 {
		t.Error("Increment()", x)
	}

	if x := c.Apply(nil, c.Set(42)); x != types.Counter(42) {
		t.Error("Set", x)
	}

	l := changes.ChangeSet{c.Increment(2), c.Increment(-3)}
	r := changes.ChangeSet{c.Increment(44), c.Increment(-42)}
	lx, rx := l.Merge(changes.PathChange{nil, r})
	lval := c.Apply(nil, l).Apply(nil, lx)
	rval := c.Apply(nil, r).Apply(nil, rx)
	if lval != rval || lval != c+2-3+44-42 {
		t.Error("Unexpected lval, rval", lval, rval)
	}
}

func TestCounterPanics(t *testing.T) {
	mustPanic := func(fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Failed to panic")
			}
		}()
		fn()
	}

	mustPanic(func() {
		(types.Counter(0)).Apply(nil, poorlyDefinedChange{})
	})

	mustPanic(func() {
		types.Counter(0).Slice(0, 0)
	})
}
