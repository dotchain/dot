// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"reflect"
	"testing"
)

// TestValidateUndoOrderInOperation only validates that the order
// of performing undos in an Operation.Undo is right.  It actually
// does not test whether the individual undo is right -- that is
// covered in common_test.go with every apply operation
func TestValidateUndoOrderInOperation(t *testing.T) {
	// insert hello at 5 and then world at 10
	splice1 := dot.Change{Splice: &dot.SpliceInfo{5, "", "hello"}}
	splice2 := dot.Change{Splice: &dot.SpliceInfo{10, "", "world"}}
	op := dot.Operation{Changes: []dot.Change{splice1, splice2}}
	changes := op.Undo().Changes
	if !reflect.DeepEqual(changes, []dot.Change{splice2.Undo(), splice1.Undo()}) {
		t.Error("Expect undo to be reversed properly")
	}
}
