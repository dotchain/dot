// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import "github.com/dotchain/dot/changes"

// Op represents a collection of changes, like a git-commit.
//
// The ID of the op uniquely identifies it.
//
// All operations are effectively ordered sequentually in a "master"
// branch.  The index in this store is the "Version" of the operation.
//
// Every operation has a "parent" operation that it is based on. That
// operation is identified by the ParentID. Each operation can be
// thought of as creating a branch.  But a client can also "pull" in
// changes from the "master" branch. When this happens, further
// operations will effectively be based upon both the "Parent" and the
// merged in master version.  The Basis tracks the lasts version of
// the master branch that was pulled in and so must always be
// non-negative.
//
// Type Operation implements this interface though the rest of the
// package makes no assumptions about the concrete implementation of
// this interface.
type Op interface {
	// ID is the unique identifier for the operation. It should be
	// of any type that can be compared for equality
	ID() interface{}

	// Version is the index of the operation in the sequential
	// store. For operations that have not yet been stored (such
	// as client operations), this must be negative. Once an
	// operation is stored, its version will never change.
	Version() int

	// WithVersion returns a new op with the specified version.
	WithVersion(int) Op

	// Parent is the ID of the local operation this particular
	// operation is based upon.  If this operaiton is based upon
	// an operation that has a version at this point (i.e. it has
	// been stored), then ParentID must be nil.
	Parent() interface{}

	// Basis is the last "stored" operation (i.e. with
	// non-negative version) that has been factored in before the
	// current operation.
	//
	// Example:
	//
	// If a client makes a sequence of operations (A, B, C), then
	// all of them will share the same Basis.  If, at this point,
	// the client receives a new operation from the server and
	// merges it in, its basis for future operations would be set
	// to that version but its parent would remain as the last
	// local operation.
	//
	// If the client receives its own operation back from the
	// server, there is no specific action the client has to take
	// but the next future operation will have the acknowledged
	// operation as its Basis and an empty Parent.
	Basis() int

	// Changes is the changes associated with this operation. It
	// can be nil.
	Changes() changes.Change

	// WithChanges creates a new operation with the same metadata
	// as the current but with a different set of changes.
	WithChanges(changes.Change) Op
}

// Operation holds the basic info needed for Op with string IDs
type Operation struct {
	OpID, ParentID interface{}
	VerID, BasisID int
	changes.Change
}

// ID implements Op.ID
func (o Operation) ID() interface{} {
	return o.OpID
}

// Version implements Op.Version
func (o Operation) Version() int {
	return o.VerID
}

// WithVersion implements Op.WithVersion
func (o Operation) WithVersion(v int) Op {
	o.VerID = v
	return o
}

// Parent implements Op.Parent
func (o Operation) Parent() interface{} {
	return o.ParentID
}

// Basis implements Op.Basis
func (o Operation) Basis() int {
	return o.BasisID
}

// Changes implements Op.Changes
func (o Operation) Changes() changes.Change {
	return o.Change
}

// WithChanges implements Op.WithChanges
func (o Operation) WithChanges(c changes.Change) Op {
	o.Change = c
	return o
}
