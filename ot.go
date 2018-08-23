// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dot implements core operational transforms for composable
// JSON arrays and objects
//
// Clients and servers
//
// This package is a low level transformation library.  It does
// not include any client or server implementations.  Please see
// https://godoc.org/github.com/dotchain/ver for an usable client
// library. Please ssee https://github.com/dotchain/dots for server
// implementations.
//
// Quick Introduction to Operational Transforms
//
// Operational transformation is a technique for conflict-free
// synchronization of data that bears similarity to git. Each change
// of the data is represented as a compact data structure that is
// streamed to all other clients.  A transformation procedure is used
// to "reconcile" two changes made on the same initial state on
// separate clients, to bring them both into convergence.
//
// For example, if the initial data is the string "!hello world" and
// two independent changes were made: move the exclamation mark to the
// end and insert a comma after the "hello".  These can be encoded as
// changes like so:
//
//     move from offset 0 count 1 to the right by 11
//     insert at offset 6 ","
//
// Now if each of those client directly applied the change from the
// other client, the results would not converge: "hello ,world!" and
// "hello, worl!d".  The trick with operational transformation is to
// transform the two operations against each other such that applying
// the transformed operations converges:
//
//     move from offset 0 count 1 to the right by 11
//        + insert at offset 5 ","
//
//     insert at offset 5 ","//
//        + move from offset 0 count 1 to the right by 12
//
//
// Stateless transforms and chaining
//
// As described above, the individual changes can be represented in a
// compac way without reference to the actual data structure. In
// addition, by careful choice, one can set up the changes so that the
// changes can be transformed against one another without reference to
// the current state of the data.  All changes and transformations in
// this package are stateless.
//
// Note that it is only legal to transform two changes if they were
// both based on the same initial state. If there are multiple
// operations off the same initial state, one can recursively
// transform the list and this is the central premise.
//
// Changes and operations
//
// This package defines four Change types: Set, Splice, Move and
// Range. Set is used to set (and delete) a key in a dictionary while
// the other three operations apply to collections. It is possible to
// implement quite a few interesting types on top of these
// operations. For example, a simple integer counter can be
// implemented as an array of integers -- each increment simply an
// insertion into this array.  This might seem like a wasteful
// approach as the array would keep growing in size -- but note that
// only the change needs to be represented in such fashion.  The
// actual client can choose to implement the underlying data structure
// as an integer.
//
// The change also includes a Path to deal with composition.  For
// example, if the core data structure was an array of arrays, editing
// the inner array would cause the path to be set to the index in the
// outer array.  Similarly, for objects, the path is set to the key of
// the item.
//
// Operations are collections of changes with an ID to uniquely
// identify the operation.  In addition, each operation contains a
// Parent field where the first element represents the **BasisID** and
// the second represents the **ParentID**.  The core model of this
// package is the assumption that there is a central server where the
// operations are persisted in a consistent order.  The BasisID refers
// to the ID of the last operation that a client has applied locally
// (after transformations) that it received from this persisted
// log. In addition, the client may have done a sequence of local
// operation.  Each such local operation would set the ParentID to be
// the last local operation.  See the DOT protocol documentation for
// further details on these IDs
// (https://github.com/dotchain/site/blob/master/Protocol.md).
//
// Transformer, Log and ClientLog
//
// The main export of this package is the Transformer type which
// provides the functionality to transform any set of changes or
// operations based on the same initial state.
//
// The Log type provides additional support for taking a raw sequence
// of operations in the order that they are persisted (with their
// varying BasisID and ParentID) and transforming them into version
// that can be applied to derive the model needed.
//
// References
//
// A particularly useful construct is the idea of a reference within
// the JSON object.  A path into the JSON object would need to change
// as operations are applied because an item within an array may
// change index due to an insertion before that item.  The Ref type
// helps manage this process.
//
// Undos
//
// Each change is represented in a way that its undo can be calculated
// without reference to the original state. There are some
// complications with undos when remote changes can intervene local
// changes. While the transforms in this package guarantee
// convergence, there will be some unexpected effects if remote
// operations intervene that affect the same region. The UndoStack
// type handles transforming undos and redos when there are
// intervening operations.
//
// Custom types
//
// The ver package (https://godoc.org/github.com/dotchain/dot)
// provides client-side to work with collections, maps and to
// customize them to a degree.  Occasionally, an app would need to
// build a custom type of collections (say, one which also builds an
// index along with it).  In these cases, one can write custom
// "encodings".  See
// https://godoc.org/github.com/dotchain/dot/encoding/richtext for a
// rich text encoding which looks like a regular array but is actually
// stored and transmitted differently on the wire.
//
// Another interesting example is the case of "counters" -- integers
// that can be incremented or decremented.  These provide the
// impression of arrays but are basically just stored as integers.
// See https://godoc.org/github.com/dotchain/dot/encoding/counters for
// an implementation of counters.  Here the wire protocol for mutating
// a counter does not actually look different -- the encoding only
// plays a role in the construction of the counter type (for example,
// when an object field is initialized to a counter).
//
package dot

import "github.com/dotchain/dot/conv"

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
	prefix, suffix := path[:offset:offset], path[offset+1:]
	updated := append(append(prefix, conv.FromIndex(index)), suffix...)
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
