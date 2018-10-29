// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package changes implements the core mutation types for OT.
//
// The three basic types are Replace, Splice and Move. Replace
// replaces a value altogether while Splice replaces a sub-sequence in
// an array-like object and Move shuffles a sub-sequence.
//
// Both Slice and Move work on strings as well as strings are just a
// form of arrays as far as OT is concerned. The actual representation
// of strings is abstracted away.  See
// https://godoc.org/github.com/dotchain/dot/x/types#S8 for the
// implementation of an OT-compatible string type.
//
// ChangeSet allows a set of mutations to be grouped together.
//
// PathChange allows a mutation to refer to a "path".  This is useful
// when mutating a JSON-like object composed of arrays and maps. The
// path is simply how a particular node is to be traversed, with keys
// and indices in the order in which they appear.
//
// Custom change types can be defined. They should implement the
// Custom interface. See
// https://godoc.org/github.com/dotchain/dot/changes/x/rt#Run for an
// example custom change type.
//
// Replace and Splice change both expect the Before and After fields
// to be non-nil Value implementations. Replace can use changes.Nil
// to represent empty values for the case where a value is being
// deleted or created. Slices must make sure that the Before
// and After use the "empty" representations of the respective types.
//
// Slices also should generally make sure that the Before and After
// types are compatible -- i.e. inserting a number within a string is
// not permitted.
//
// Any custom Value implementation should implement the Value
// interface.  See https://godoc.org/github.com/dotchain/dot/x/types
// for a set of custom value types such as string, arrays and
// counters.
//
// See https://godoc.org/github.com/dotchain/dot/x/rt for a custom
// type that has a specific custom change associated with it.
package changes

// Change represents an OT-compatible mutation of the virtual JSON.
//
// The methods provided here are the core methods.  Custom changes
// should implement the Custom interface in addition to this.
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
// changes.  Array like values should also implement Collection
type Value interface {
	// Apply applies the specified change on the object and
	// returns the updated value.
	Apply(c Change) Value
}

// Collection represents array-like values
type Collection interface {
	// must also implement Value
	Value

	// ApplyCollection is just strongly typed Apply
	ApplyCollection(c Change) Collection

	// Slice should only be called on collection-like objects such
	// as the Before/After fields of a Splice. Note that unlike
	// Go's slice notation, the arguments are offset and count.
	Slice(offset, count int) Collection

	// Count should noly be called on collection-like objects such
	// as the Before/After fields of a Splice. It returns the size
	// of the collection.
	Count() int
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

// ApplyTo simply walks through the individual changes and applies
// them to the value.
func (c ChangeSet) ApplyTo(v Value) Value {
	for _, cx := range c {
		v = v.Apply(cx)

	}
	return v
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

// Custom is the interface that custom change types should implement.
// This allows the "known" types to interact with custom types.
type Custom interface {
	Change

	// ReverseMerge is used when a "known" change type is merged
	// with a custom type:
	//
	//       move.Merge(myType)
	//
	// The known type has no way of figuring out how to merge. It
	// calls ReverseMerge on the custom type to get the custom
	// type to deal with this. This is separate from the regular
	// Merge method because calling "myType.Merge(move)" may not
	// be the same:  the Merge() call is not required to be
	// symmetric. A good example of a non-symmetric situation is
	// when the left change and  the right change both are
	// "inserting" into the same array at the same point -- the
	// changes will have to be ordered so that one of them ends up
	// before the other.
	//
	// Basically, if:
	//
	//       ax, bx := a.Merge(b)
	//
	// Then:
	//
	//       bx, ax := b.ReverseMerge(a)
	//
	ReverseMerge(c Change) (Change, Change)

	// ApplyTo allows custom change types to implement a method to
	// apply their changes onto known "values".  This allows Value
	// implementations to be written without awareness of all
	// possible Change implementations
	ApplyTo(v Value) Value
}
