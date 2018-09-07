// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package refs implements reference paths, carets and selections.
//
// A reference path, caret or selection refers to an item or a
// position in an array-like object or a set of items in an array-like
// object.  As changes are applied, the path may be affected as well
// as items that the path refers to. This package provides the
// mechanism to deal with these.
package refs

import "github.com/dotchain/dot/changes"

// Ref represents the core Reference type
type Ref interface {
	// Merge takes a change and returns an Ref that reflects the
	// effect of the change.  If the change affects the item
	// specified by the Ref, it returns a modified version of the
	// change that can be applied to the value at the Ref.
	Merge(c changes.Change) (Ref, changes.Change)
}

// InvalidRef refers to a ref that no longer exists.
var InvalidRef = invalidRef{}

type invalidRef struct{}

func (r invalidRef) Merge(c changes.Change) (Ref, changes.Change) {
	return r, nil
}

func mergeChangeSet(ref Ref, c changes.ChangeSet) (Ref, changes.Change) {
	result := make([]changes.Change, len(c))
	idx := 0
	for _, cc := range c {
		ref, cc = ref.Merge(cc)
		if cc != nil {
			result[idx] = cc
			idx++
		}
	}
	switch idx {
	case 0:
		return ref, nil
	case 1:
		return ref, result[0]
	}
	return ref, changes.ChangeSet(result)
}
