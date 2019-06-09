// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package changes implements the core mutation types for OT.
//
// The three basic types are Replace, Splice and Move. Replace
// replaces a value altogether while Splice replaces a sub-sequence in
// an array-like object and Move shuffles a sub-sequence.
//
// Both Slice and Move work on strings as well. The actual type for
// Slice and Move is represented by the Collection interface (while
// Replace uses a much more lax Value interface)
//
// See
// https://godoc.org/github.com/dotchain/dot/changes/types#S16 for an
// implementation of an OT-compatible string type.
//
// Composition
//
// ChangeSet allows a set of mutations to be grouped together.
//
// PathChange allows a mutation to refer to a "path".  For example, a
// field in a "struct type" can be thought of as having the path of
// "field name".  An element of a collection can be thought  of as
// having the path of the index in that array.
//
// The type of the elements in the path is not specified but it is
// assumed that they are comparable for equality.  Collections are
// required to use the index of type int for the path elements.
//
// Custom Changes
//
// Custom change types can be defined. They should implement the
// Custom interface. See
// https://godoc.org/github.com/dotchain/dot/changes/x/rt#Run for an
// example custom change type.
//
// The general asssumption underlying OT is that the Merge method
// produces convergence:
//
//    if: c1, c2 = changes on top of "initial"
//    and: c1x, c2x := c1.Merge(c2)
//    then: initial + c1 + c1x == initial + c2 + c2x
//
// Notes
//
// Replace and Splice change both expect the Before and After fields
// to be non-nil Value implementations. Replace can use changes.Nil
// to represent empty values for the case where a value is being
// deleted or created. Slices must make sure that the Before
// and After use the "empty" representations of the respective types.
//
// Slices also should generally make sure that the Before and After
// types are compatible -- i.e. each should be able to be spliced
// within the other.
//
//
// Value Interface
//
//
// Any custom Value implementation should implement the Value
// interface.  See https://godoc.org/github.com/dotchain/dot/changes/types
// for a set of custom value types such as string, arrays and
// counters.
//
// See https://godoc.org/github.com/dotchain/dot/x/rt for a custom
// type that has a specific custom change associated with it.
//
// It is common to have a value type (say *Node) that is meant as an
// atomic value. In that case, one can use the Atomic{} type to hold
// such values.
package changes

// Change represents an OT-compatible mutation
//
// The methods provided here are the core methods.  Custom changes
// should implement the Custom interface in addition to this.
// Note that it is legal for a change to be nil (meaning the value
// isnt change at all)
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

// Custom is the interface that custom change types should implement.
// This allows the "known" types to interact with custom types.
//
// Custom changes might also need to implement the refs.PathMerger
// interface (see https://godoc.org/github.com/dotchain/dot/refs)
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
	ApplyTo(ctx Context, v Value) Value
}

// Value represents an immutable JSON object that can apply
// changes.  Array like values should also implement Collection
type Value interface {
	// Apply applies the specified change on the object and
	// returns the updated value.
	Apply(ctx Context, c Change) Value
}

// Collection represents an immutable array-like value
type Collection interface {
	// must also implement Value
	Value

	// ApplyCollection is just strongly typed Apply
	ApplyCollection(ctx Context, c Change) Collection

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
	idx, results := 0, make(ChangeSet, len(c))
	for _, elt := range c {
		other, results[idx] = Merge(elt, other)
		if results[idx] != nil {
			idx++
		}
	}
	return other, results.Simplify()
}

// ReverseMerge is like merge except with receiver and args inverted
func (c ChangeSet) ReverseMerge(other Change) (otherx, cx Change) {
	idx, results := 0, make(ChangeSet, len(c))
	for _, elt := range c {
		results[idx], other = Merge(other, elt)
		if results[idx] != nil {
			idx++
		}
	}

	return other, results.Simplify()
}

// Revert implements Change.Revert.
func (c ChangeSet) Revert() Change {
	idx, results := 0, make(ChangeSet, len(c))
	for kk := range c {
		if elt := c[len(c)-kk-1]; elt != nil {
			results[idx] = elt.Revert()
			idx++
		}
	}
	return results.Simplify()
}

// ApplyTo simply walks through the individual changes and applies
// them to the value.
func (c ChangeSet) ApplyTo(ctx Context, v Value) Value {
	for _, cx := range c {
		v = v.Apply(ctx, cx)

	}
	return v
}

// Simplify converts an empty or single element change-set
// into a simpler version
func (c ChangeSet) Simplify() Change {
	result := ChangeSet{}
	for _, cx := range c {
		if cx = Simplify(cx); cx != nil {
			result = append(result, cx)
		}
	}
	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	}
	return result
}

// Context defines the context in which a change is being
// applied. This is useful to capture data such as the "current user"
// or "virtual time" etc.  For true convergence, the context itself
// should be derived from the change -- say via Meta.
//
// Note that this interface is a subset of the standard golang
// "context.Context"
type Context interface {
	Value(key interface{}) interface{}
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

// Merge is effectively c1.Merge(c2) except that c1 can be nil.
//
// As with individual merge implementations, applying c1+c1x is
// effectively the same as applying c2+c2x.
func Merge(c1, c2 Change) (c1x, c2x Change) {
	if c1 == nil {
		return c2, c1
	}
	return c1.Merge(c2)
}

// Simplify converts a change to a simpler form if possible
func Simplify(c Change) Change {
	if simp, ok := c.(simplifier); ok {
		return simp.Simplify()
	}
	return c
}

type simplifier interface {
	Simplify() Change
}
