// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package changes implements the core mutation types for OT.
//
// The three basic types are Replace, Splice and Move. Replace
// replaces a value altogether while Splice replaces a sub-sequence in
// an array-like object and Move shuffles a sub-sequence.
//
// ChangeSet allows a set of mutations to be grouped together.
//
// PathChange allows a mutation to refer to a "path".  This is useful
// when mutating a JSON-like object composed of arrays and maps. The
// path is simply how a particular node is to be traversed, with keys
// and indices in the order in which they appear.
//
// Custom changes can be defined so long as it follows the Change
// interface and also implements a "ReverseMerge" method. See
// https://godoc.org/github.com/dotchain/dot/changes/x/rt#Run for a
// custom change type.
//
// The Replace and Splice change both expect the Before, After fields
// to be non-nil Value implementations.  Nil is an example empty value
// and Atomic can be used to wrap any immutable value.
//
// Any custom Value implementation should implement the Value
// interface.  See https://godoc.org/github.com/dotchain/dot/types
// for a set of custom value types
package changes

// Change represents an OT-compatible mutation of the virtual JSON.
//
// Note that it is legal for a change to be nil. It represents a
// noop.
type Change interface {
	// Merge takes the current change and another change that were
	// both applied to the same virtual JSON and returns
	// transformed versions such that:
	//
	//     c + otherx = other + cx
	//
	// Otherx captures the intent behind other and cx captures the
	// intent behind the current change.  The equation above
	// guarantees that the combined intent can be achieved by
	// applying the transformed changes on top of the local
	// change.
	//
	// Note that there is no requirement that c.Merge(other) and
	// other.Merge(c) should both yield the same transforms. This
	// requires that for correctness, there must be a clear way to
	// decide which change is the receiver and which is the param.
	Merge(other Change) (otherx, cx Change)

	// Revert returns the opposite effect of the current
	// change. If applied directly after the current change it
	// should undo the effect of the current change.
	Revert() Change
}

// Value represents an immutable JSON object that can apply
// changes.
type Value interface {
	// Slice should only be called on collection-like objects such
	// as the Before/After fields of a Splice. Note that unlike
	// Go's slice notation, the arguments are offset and count.
	Slice(offset, count int) Value

	// Count should noly be called on collection-like objects such
	// as the Before/After fields of a Splice. It returns the size
	// of the collection.
	Count() int

	// Apply applies the specified change on the object and
	// returns the updated value.
	Apply(c Change) Value
}

// ChangeSet represents a collection of changes. It implements the
// Change interface thereby allowing merging groups of changes against
// each other.
type ChangeSet []Change

// Merge implements Change.Merge.
func (c ChangeSet) Merge(other Change) (otherx, cx Change) {
	idx, results := 0, make([]Change, len(c))
	for _, elt := range c {
		if elt != nil {
			other, results[idx] = elt.Merge(other)
			if results[idx] != nil {
				idx++
			}
		}
	}
	switch idx {
	case 0:
		return other, nil
	case 1:
		return other, results[0]
	}
	return other, ChangeSet(results[:idx])
}

// ReverseMerge is like merge except with receiver and args inverted
func (c ChangeSet) ReverseMerge(other Change) (otherx, cx Change) {
	idx, results := 0, make([]Change, len(c))
	for _, elt := range c {
		l, r := elt, other
		if other != nil {
			l, r = other.Merge(elt)
		}
		results[idx], other = l, r
		if l != nil {
			idx++
		}
	}
	switch idx {
	case 0:
		return other, nil
	case 1:
		return other, results[0]
	}
	return other, ChangeSet(results[:idx])
}

// Revert implements Change.Revert.
func (c ChangeSet) Revert() Change {
	idx, results := 0, make([]Change, len(c))
	for kk := range c {
		if elt := c[len(c)-kk-1]; elt != nil {
			results[idx] = elt.Revert()
			idx++
		}
	}
	switch idx {
	case 0:
		return nil
	case 1:
		return results[0]
	}
	return ChangeSet(results[:idx])
}

type revMerge interface {
	ReverseMerge(o Change) (ox, mx Change)
}

// helper method
type changeable interface {
	Change() Change
}

func change(x, y changeable) (Change, Change) {
	return x.Change(), y.Change()
}

func swap(x, y Change) (Change, Change) {
	return y, x
}
