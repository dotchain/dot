// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
)

func validateMerge(t *testing.T, initial changes.Value, left, right changes.Change) {
	leftx, rightx := left.Merge(right)
	lval := initial.Apply(nil, changes.ChangeSet{left, leftx})
	rval := initial.Apply(nil, changes.ChangeSet{right, rightx})
	if lval != rval {
		t.Error("Diverged", lval, rval, left, right)
	}
}

func TestConvergenceNonEmptyInitial(t *testing.T) {
	ForEachChange(S("xyz"), func(initial changes.Value, left changes.Change) {
		ForEachChange(S("ab"), func(_ changes.Value, right changes.Change) {
			validateMerge(t, initial, left, right)
			validateMerge(t, initial, right, left)
		})
	})
}

func TestConvergenceEmptyInitial(t *testing.T) {
	initial := changes.Nil
	cx := []changes.Change{
		changes.Replace{changes.Nil, S("hello")},
		changes.Replace{changes.Nil, S("world")},
	}
	for _, left := range cx {
		for _, right := range cx {
			validateMerge(t, initial, left, right)
		}
	}
}

func TestConvergenceChangeSet(t *testing.T) {
	initial := S("hello")
	left := changes.ChangeSet{
		changes.Replace{initial, changes.Nil},
		changes.Replace{changes.Nil, S("World")},
		changes.Splice{5, S(""), S("!")},
	}
	right := changes.ChangeSet{
		changes.Splice{0, S("h"), S("j")},
		changes.Splice{5, S(""), S(" shots!")},
	}
	validateMerge(t, initial, left, right)
}

func ForEachChange(replacement changes.Collection, fn func(initial changes.Value, c changes.Change)) {
	initial := S("abcdef")
	fn(initial, changes.ChangeSet{nil})
	fn(initial, changes.Replace{initial, replacement})
	fn(initial, changes.Replace{initial, changes.Nil})
	fn(initial, changes.Move{0, 5, 0})
	fn(initial, changes.Move{5, 0, 1})
	for offset := 0; offset <= initial.Count(); offset++ {
		for count := 0; count <= initial.Count()-offset; count++ {
			before := initial.Slice(offset, count)
			if count > 0 {
				fn(initial, changes.Splice{offset, before, S("")})
			}
			fn(initial, changes.Splice{offset, before, replacement})
		}
	}

	count := replacement.Count()
	for offset := 0; offset <= initial.Count()-count; offset++ {
		for dest := 0; dest <= initial.Count(); dest++ {
			if dest < offset {
				fn(initial, changes.Move{offset, count, dest - offset})
			} else if dest > offset+count {
				fn(initial, changes.Move{offset, count, dest - offset - count})
			}
		}
	}
}
