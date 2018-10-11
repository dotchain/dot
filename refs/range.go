// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs

import "github.com/dotchain/dot/changes"

// Range is a reference to a specific selection of elements in an array-like
// object.  Ranges can be "collapsed" when they behave the same as Caret.
//
// This is an immutable type
//
// This only handles the standard set of changes. Custom changes
// should implement a MergeCaret method:
//
//    MergeCaret(caret refs.Caret) (refs.Ref)
//
// Range is implemented on top of Caret and doing the above is sufficient.
//
// Note that this is in addition to the MergePath method which is
// called first to transform the path and then the MergeCaret is
// called  on the updated Caret (based on the path returned by
// MergePath).
//
// Note that the paths for the Start and End are expected to be the same.
type Range struct {
	Start, End Caret
}

// Merge updates the range based on the change.  Note that it
// always returns a nil change as there are cases where there isn't
// sufficient information to return the correct changes.Change.
func (r Range) Merge(c changes.Change) (Ref, changes.Change) {
	sx, _ := r.Start.Merge(c)
	ex, _ := r.End.Merge(c)
	sy, ok1 := sx.(Caret)
	ey, ok2 := ex.(Caret)
	if !ok1 || !ok2 {
		return InvalidRef, nil
	}

	// special case: collapsing a range should make sure
	// IsLeft values match.
	if r.Start.Index != r.End.Index && sy.Index == ey.Index {
		ey.IsLeft = sy.IsLeft
	}

	return Range{sy, ey}, nil
}

// Equal implements Ref.Equal
func (r Range) Equal(other Ref) bool {
	o, ok := other.(Range)
	return ok && r.Start.Equal(o.Start) && r.End.Equal(o.End)
}
