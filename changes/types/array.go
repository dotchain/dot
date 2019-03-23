// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import "github.com/dotchain/dot/changes"

// The A type represents a slice of arbitrary values. It implements
// the changes.Value interface. The actual elements can be nil (unlike
// the regular requirement for values to be non-nil). Nil values are
// treated as if they were changes.Nil
type A []changes.Value

// Slice implements changes.Collection.Slice
func (a A) Slice(offset, count int) changes.Collection {
	return a[offset : offset+count]
}

// Count returns size of the array
func (a A) Count() int {
	return len(a)
}

func (a A) get(key interface{}) changes.Value {
	if v := a[key.(int)]; v != nil {
		return v
	}
	return changes.Nil
}

func (a A) set(key interface{}, v changes.Value) changes.Value {
	clone := append(A(nil), a...)
	if v != changes.Nil {
		clone[key.(int)] = v
	} else {
		clone[key.(int)] = nil
	}
	return clone
}

func (a A) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	return append(append(a[:offset:offset], after.(A)...), a[end:]...)
}

// ApplyCollection implements changes.Collection
func (a A) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (Generic{Get: a.get, Set: a.set, Splice: a.splice}).ApplyCollection(ctx, c, a)
}

// Apply applies the change and returns the updated value
//
// Note: deleting an element via changes.Replace simply replaces it
// with nil.  It does not actually remove the element -- that needs a
// changes.Splice.
func (a A) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (Generic{Get: a.get, Set: a.set, Splice: a.splice}).Apply(ctx, c, a)
}
