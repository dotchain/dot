// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package diff compares two values and returns the changes
package diff

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Differ is the general Diff interface for computing changes in
// values as a changes.Change (which when applied  to the old value is
// guaranteed to result in the new value)
type Differ interface {
	// Diff computes the diff between old and new changes.Value
	//
	// The provided differ is meant to be used for any children or
	// inner values whose specific types may not be known to a
	// particular differ.
	Diff(d Differ, old, new changes.Value) changes.Change
}

// Std implements diffs for standard types
type Std struct{}

// Diff implemnets Differ.Diff for standard types
func (s Std) Diff(d Differ, old, new changes.Value) changes.Change {
	switch old := old.(type) {
	case types.S8:
		if _, ok := new.(types.S8); ok {
			return S8(d, old, new)
		}
	case types.S16:
		if _, ok := new.(types.S16); ok {
			return S16(d, old, new)
		}
	}
	return changes.Replace{Before: old, After: new}
}
