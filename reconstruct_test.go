// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"testing"
)

func TestUtilsReconstruct(t *testing.T) {
	guessEmptyValue := func(ops []dot.Change) (interface{}, bool) {
		if len(ops) == 0 {
			return nil, true
		}
		op := ops[0]
		u := dot.Utils(dot.Transformer{})

		if len(op.Path) > 0 {
			panic("Unexpected path with empty before")
		}

		if op.Splice != nil {
			switch {
			case op.Splice.Before != nil:
				if val, ok := u.C.TryGet(op.Splice.Before); ok {
					return val.Slice(0, 0), true
				}
				panic("Unexpected splice op")
			case op.Splice.After != nil:
				if val, ok := u.C.TryGet(op.Splice.After); ok {
					return val.Slice(0, 0), true
				}
				panic("Unexpected splice op")
			}
			return nil, false
		}
		if op.Set != nil {
			return map[string]interface{}{}, true
		}
		if op.Move != nil {
			if op.Move.Count == 0 || op.Move.Distance == 0 {
				return nil, false
			}
			panic("Unexpected move on nil input")
		}
		if op.Range != nil {
			if op.Range.Count == 0 {
				return nil, false
			}
			panic("Unexpected range on nil input")
		}
		return nil, false
	}

	validate := func(m interface{}) {
		u := dot.Utils(dot.Transformer{})
		changes := u.Reconstruct(m)
		initial, ok := guessEmptyValue(changes)
		if !ok {
			t.Error("Failed to reconstruct", m)
			return
		}

		actual := u.Apply(initial, changes)
		if !u.AreSame(m, actual) {
			t.Error("Failed to reconstruct", m, "Got:", actual)
		}
	}

	validate(nil)
	validate("")
	validate([]interface{}{})
	validate("hello")
	validate([]interface{}{"p", 42.0})
	validate(map[string]interface{}{})
	validate(map[string]interface{}{"p": 42, "q": "hello"})
	validate(map[string]interface{}{"p": []interface{}{"a", 42.0}})
}
