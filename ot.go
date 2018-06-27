// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dot implements conflict free stateless transforms
// that form the basis for a larger platform
//
// Please read http://github.com/dotchain/site/blob/master/IntroductionToOperationalTransforms.md
// for an overview of Operational Transforms as used here.
//
// Please read https://github.com/dotchain/site/blob/master/Manifesto.md
// for the overall goals of the project.
//
// The DOT Architecture is a pub-sub system which relies on a
// backend service to provide a consistent order of untrasnformed
// operations submitted by multiple clients.
//
// This ordered sequence of untransformed operations is a "journal"
// and can be accessed by a simple protocol defined in
// https://github.com/dotchain/site/blob/master/Protocol.md.
//
// An implementation of this journal protocol is available here:
// https://godoc.org/github.com/dotchain/dots/journal
//
// The raw sequence of operations in the journal can be transformed
// using the #Log struct.  If the client has local changes, it can
// calculate the "compensating actions" that can be safely applied
// to get the converged state.  The #ClientLog struct helps with
// this.
//
// It is more practical for clients to simply use the "log" protocol
// which implements the reconciliation procedure transforming all
// the operations internally and providing the client with the
// sequence of "compensating actions".  This makes it easy to implement
// clients as the clients no longer even need to understand or implement
// operational transforms.
//
// This log protocol is also documented at
// https://github.com/dotchain/site/blob/master/Protocol.md and
// a reference implementation is available at
// https://godoc.org/github.com/dotchain/dots/log
//
// Note that reconcilers can also be implemented on the client side
// for added responsiveness in situations where clients are almost
// constantly editing collaboratively.
//
// This package implements the core transformation engine via:
// the Transformer (which implements the actual Merge actions),
// the Log (which implements the mechanism which converts journal
// operations to rebased operations clients can apply) and the ClientLog
// (which implements the reconciliation mechanisms and provides the
// compensating actions for clients to apply to their local models)
//
// Virtual JSON model
//
// The unit of granularity for syncing is a Model. The DOT engine
// interprets each Model as a logical JSON entity (i.e. a tree of
// JSON objects and arrays). The real representation can be anything
// but the changes operate on the logical model.  When clients actually
// apply the operations they would need to figure out the right
// way to apply them based on the real type/schema of the model.
//
// The smallest granularity mutation is represented via a #Change
// struct.  This struct has the path where the mutation is happening.
// The actual Change structure is a union of #SpliceInfo, #MoveInfo,
// #RangeInfo or #SetInfo types.  Since Go does not have native
// union types, this is modeled with a struct where only one of the
// four values are set.
//
// In reality, most models will have a much richer mutation semantics
// than just these three operations.  For example, Rich Text may
// support "bolding" of a region of text. Please see
// https://github.com/dotchain/site/blob/master/ComposableOperations.md
// for a discussion of how this can be achieved in DOT.
//
// Most of these individual changes are weakly typed.  #SpliceInfo,
// for instance, can represent the splice action on a string or an array
// or a run length encoded array. #SetInfo similarly works
// with string keys but arbitrary interface types for the value.
// This makes it convenient when working with a JSON schema but does not
// preclude applying the operations to a strongly typed state object.
//
// #Transformer
//
// For those not familar with OT, please read
// https://github.com/dotchain/dot/blob/master/IntroductionToOperationalTransforms
// for a gentle introduction to the topic as it pertains to DOT
//
// The heart of the logic of OT is a set of low-level merge functions
// that can take any two primitive #Change mutations (that were applied
// to the same initial model independently) and find equivalent
// "compensating" Change values, which when applied on top of
// the corresponding initial mutations will converge to the exact same state.
// These standard functions are implemented by the #Transformer struct.
//
// This provides the backbone for implementing a richer construct -- if
// two different clients have independently made a sequence of mutations
// (not just single changes), can we "converge" these two clients
// to the same state reliably?
//
// The #Transformer.MergeChanges method does exactly that by using the
// single mutation merge function as the basis to repeatedly transform
// operations. The return value is a matching pair of "compensating"
// mutations that converge the two states.
//
// An advanced version of the transformation is the #Log.AppendOperation
// function which transforms an operation with arbitrary basis (i.e. does
// not have to be directly on top of the last operation) against the
// journal (which is the authoritative ordered sequence of operations)
// into a "rebased" operation.  The rebased operation can be applied on
// top of the journal safely to construct the cumulative model state.
//
// The #ClientLog component uses the Log structure and its helper methods
// to find the set of compensating actions that a client can take to
// merge changes from the server into the client model.
package dot

// TODO: replace with a smaller library so it won't explode when used with GopherJS
import "strconv"

// SpliceInfo represents mutating a sequence by replacing a contiguous
// sub-sequence with the provided alternate value.  If the Before value
// is as empty sequence, this represents an insertion.  If the After value
// is an empty sequence, this represents a deletion.
//
// This is weakly typed.  The Before/After values can be arrays or
// strings, etc.  Also, they may be nil to indicate empty arrays
//
// Please note that unlike Go, the native string representation in DOT
// is UTF16 and so the offsets here should refer to UTF16 offsets.
type SpliceInfo struct {
	Offset int
	Before interface{} `json:",omitempty"`
	After  interface{} `json:",omitempty"`
}

// Undo returns a new change which nullifies the effect of the current change
func (t SpliceInfo) Undo() SpliceInfo {
	return SpliceInfo{Offset: t.Offset, Before: t.After, After: t.Before}
}

// MoveInfo represents mutating a sequence by shifting a contiguous sub-sequence
// by the amount specified. If Distance is negative, the shift is to the left while
// a positive value shifts it to the right.  A zero Distance makes this  no-op.
type MoveInfo struct {
	Offset, Count, Distance int
}

// Undo returns a new change which nullifies the effect of the current change
func (t MoveInfo) Undo() MoveInfo {
	return MoveInfo{Offset: t.Offset + t.Distance, Count: t.Count, Distance: -t.Distance}
}

// Dest calculates a useful offset in the sequence -- where the destination of the
// move is in the original sequence.
func (t MoveInfo) Dest() int {
	if t.Distance > 0 {
		return t.Offset + t.Count + t.Distance
	}
	return t.Offset + t.Distance
}

// TransformInterval considers what happens to a specific interval after
// a move. It gets transformed to zero (unchanged), one (moved) or two
// (split + moved) intervals. The return value can be nil (to indicate no
// change) or it can contain the new intervals, each represented by
// offset/count pairs.  The offset refers to the position after the
// move has been performed.
func (t MoveInfo) TransformInterval(offset, count int) [][2]int {
	// normalize the move so that it is always moving to the right (i.e distance > 0)
	o, c, d := t.Offset, t.Count, t.Distance
	if d < 0 {
		o, c, d = t.Offset+t.Distance, -t.Distance, t.Count
	}

	result := [][2]int{}

	// add interval but take the chance to coalesce if possible
	addInterval := func(ox, cx int) {
		if cx > 0 {
			last := len(result) - 1
			if last >= 0 && ox == result[last][0]+result[last][1] {
				result[last][1] += cx
			} else {
				result = append(result, [2]int{ox, cx})
			}
		}
	}

	// calculate intersection of two intervals
	intersection := func(o1, c1, o2, c2 int) (int, int) {
		start, end := max(o1, o2), min(o1+c1, o2+c2)
		return start, end - start
	}

	// consider the move as four segments: {0, o}, {o, o+c}, {o+c, o+c+d}, {o+c+d, ...}
	o1, c1 := intersection(0, o, offset, count)
	o2, c2 := intersection(o, c, offset, count)
	o3, c3 := intersection(o+c, d, offset, count)
	o4, c4 := intersection(o+c+d, offset+count-o-c-d, offset, count)

	// adding in order of o1, o3, o2, o4 for best coaslescing possibilities
	addInterval(o1, c1)
	addInterval(o3-c, c3)
	addInterval(o2+d, c2)
	addInterval(o4, c4)

	// special case for when offset and count are unchanged
	if result[0][0] == offset && result[0][1] == count {
		return nil
	}
	return result
}

// SetInfo represents mutating a dictionary or struct by updting
// the value of the key to the provided new value.  Note that the
// Before and After values are weakly typed as they are not interpreted.
//
// Note also that After may be nil (or a zero value) to indicate
// deleting the key.  There is no difference between setting a key
// to nil vs deleting the key.
type SetInfo struct {
	Key    string
	Before interface{} `json:",omitempty"`
	After  interface{} `json:",omitempty"`
}

// Undo returns a new change which nullifies the effect of the
// current change
func (t SetInfo) Undo() SetInfo {
	return SetInfo{Key: t.Key, Before: t.After, After: t.Before}
}

// RangeInfo represents a set of bulk mutations which only differ
// in one aspect: the path refers to a contiguous sub-sequence of
// array elements.  In this case, the path is split into a primary
// path (represented in the top level Change structure) to refer
// to the array. Offset and Count fields are used to represent
// the sub-sequence portion and the rest of the path can be obtained
// via the Change field (which is weakly typed but really refers to
// the Change structure)
type RangeInfo struct {
	Offset, Count int
	// Changes is expected to be of type []Change but not
	// specifying the type here to avoid recursive
	// type declarations
	Changes []Change `json:",omitempty"`
}

// Undo returns a new change which nullifies the effect of the current change
func (t RangeInfo) Undo() RangeInfo {
	input := t.Changes
	output := make([]Change, len(input))
	for kk, ch := range input {
		// we reverse the order of the changes
		// here because that is how undos work!
		output[len(input)-kk-1] = ch.Undo()
	}
	return RangeInfo{Offset: t.Offset, Count: t.Count, Changes: output}
}

// Change represents a single mutation at the specified Path.
// Only one of Splice, Move, Range and Set must be non-nil.
type Change struct {
	Path   []string    `json:",omitempty"`
	Splice *SpliceInfo `json:",omitempty"`
	Move   *MoveInfo   `json:",omitempty"`
	Range  *RangeInfo  `json:",omitempty"`
	Set    *SetInfo    `json:",omitempty"`
}

// Undo returns a new change which nullifies the effect of the current change
func (change Change) Undo() Change {
	result := change
	if result.Splice != nil {
		ss := result.Splice.Undo()
		result.Splice = &ss
	} else if result.Move != nil {
		mm := result.Move.Undo()
		result.Move = &mm
	} else if result.Range != nil {
		rr := result.Range.Undo()
		result.Range = &rr
	} else if result.Set != nil {
		ss := result.Set.Undo()
		result.Set = &ss
	}
	return result
}

func (change Change) withUpdatedIndex(offset, index int) []Change {
	path := change.Path
	prefix, suffix := append([]string{}, path[:offset]...), path[offset+1:]
	updated := append(append(prefix, strconv.Itoa(index)), suffix...)
	newChange := change
	newChange.Path = updated
	return []Change{newChange}
}

// Operation represents an atomic batch of mutations (to be applied in sequence).
type Operation struct {
	ID      string
	Parents []string
	Changes []Change
}

// Undo returns a new operation which nullifies the effect of the current operation
func (op Operation) Undo() Operation {
	result := op
	ll := len(result.Changes)
	result.Changes = make([]Change, ll)
	for kk, ch := range op.Changes {
		// reverse the order of the changes!
		result.Changes[ll-kk-1] = ch.Undo()
	}
	return result
}

func (op Operation) withChanges(changes []Change) Operation {
	return Operation{
		ID:      op.ID,
		Parents: op.Parents,
		Changes: changes,
	}
}

// BasisID returns the ID of the last operation from the journal that is factored
// into the model on top of which this operation is to be applied.  BasisID does
// not capture local operations that may have been applied on the model.  Please
// see ParentID for that.
// This is required for all but the first operation.
func (op Operation) BasisID() string {
	if len(op.Parents) >= 1 {
		return op.Parents[0]
	}
	return ""
}

// ParentID returns the ID of the last local operation (that is yet to appear in
// the journal).  This is optional and operations without a parent are considered
// to have been applied directly on state that matches some version of the journal.
func (op Operation) ParentID() string {
	if len(op.Parents) >= 2 {
		return op.Parents[1]
	}
	return ""
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
