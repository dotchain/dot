// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"reflect"
	"testing"
)

func TestUndoStack_getUndoOperation(t *testing.T) {
	// this test only validates that the right operation is picked.
	// to make the tests readable, the undo log would consist of
	// letters C, S, U and R to represent local(client), remote(server)
	// undo and redo operations.  A star represents the correct operation
	// is the one that follows it.
	// The lack of a star in the input implies there is no available undo
	// operation
	tests := []string{
		"*C",
		"C*C",
		"S",
		"SS",
		"SS*CSS",
		"CS*CSS",
		"*CCU",
		"S*CCSSSSU",
		"CU*R", // note that redo should be picked!
		"CU*C",
		"CU*CU*R",
		"CSUC*C",
		// no undo possible here
		"CCCUUU",
		"CSCSCSUSUSU",
		"SSSS",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			var expected *dot.Operation
			stack := &dot.UndoStack{}
			for kk, t := range name {
				set := &dot.SetInfo{Key: name[0 : kk+1], After: true}
				op := &dot.Operation{
					ID:      name[0 : kk+1],
					Changes: []dot.Change{{Set: set}},
				}

				switch t {
				case 'C':
					stack = stack.Push(*op, true)
				case 'S':
					stack = stack.Push(*op, false)
				case 'U':
					op, stack = stack.GetUndo(op.ID)
				case 'R':
					op, stack = stack.GetRedo(op.ID)
				}

				if kk > 0 && name[kk-1] == '*' {
					expected = op
				}
			}
			result, _ := stack.GetUndo("latest")
			if result != nil && expected != nil {
				result.ID = expected.ID
				u := expected.Undo()
				expected = &u
			}

			if !reflect.DeepEqual(result, expected) {
				t.Error("Mismatch", result, expected)
			}

		})
	}

}

func TestUndoStack_getRedoOperation(t *testing.T) {
	// this test only validates that the right operation is picked.
	// to make the tests readable, the undo log would consist of
	// letters C, S, U and R to represent local(client), remote(server)
	// undo and redo operations.  A star represents the correct operation
	// is the one that follows it.
	// The lack of a star in the input implies there is no available redo
	// operation
	tests := []string{
		"C*U",
		"C*USSSSSS",
		"CCC*USSSSS",
		"CCCUR*US",
		"CSCSCSUSRSS*USS",
		"CCCUU*U",
		"CCSS*UUSR",
		// No redo possible
		"CUR",
		"CCCCUC",
		"CSUSRS",
		"CSCSCSUSCS",
		"SSSS",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			var expected *dot.Operation
			stack := &dot.UndoStack{}
			for kk, t := range name {
				set := &dot.SetInfo{Key: name[0 : kk+1]}
				op := &dot.Operation{
					ID:      name[0 : kk+1],
					Changes: []dot.Change{{Set: set}},
				}
				switch t {
				case 'C':
					stack = stack.Push(*op, true)
				case 'S':
					stack = stack.Push(*op, false)
				case 'U':
					op, stack = stack.GetUndo(op.ID)
				case 'R':
					op, stack = stack.GetRedo(op.ID)
				}

				if kk > 0 && name[kk-1] == '*' {
					expected = op
				}
			}
			result, _ := stack.GetRedo("latest")
			if result != nil {
				result.ID = expected.ID
			}

			if !reflect.DeepEqual(result, expected) {
				t.Error("Mismatch", result, expected)
			}

		})
	}

}
