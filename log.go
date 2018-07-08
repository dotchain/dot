// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import "github.com/pkg/errors"

// ErrMissingParentOrBasis error indicates an operation has a basis
// ID or parent ID that is not present in the journal before it.
var ErrMissingParentOrBasis = errors.New("could not find parent or basis of operation")

// ErrInvalidOperation is used if the operation is badly formatted.
var ErrInvalidOperation = errors.New("invalid operation")

// ErrLogNeedsBackfilling is returned by AppendOperation if the
// log does not have some operations that is being referred to
// by the operation that is being appended.
var ErrLogNeedsBackfilling = errors.New("log needs backfilling")

// ModelImage represents a snapshot of the model after all operations
// including and before the BasisID have been applied
type ModelImage struct {
	Model   interface{}
	BasisID string
}

// Log represents a transformed sequence of operations for a
// model that can be applied in sequence to reconstruct the
// model.  The transformed operations are in the Rebased field.
//
// If a client submitted an operation and other operations
// (that the client was unaware of) were also submitted in
// parallel, the client will have to apply a series of
// compensating operations to merge with those changes that
// sneaked into the journal ahead of the client operation.
// This is represented by the MergeChains array for each op.
// This field is required for the calculation of the log and so
// also functions as an internal book-keeping device.
//
// The Log structure can be sparse by essentially only storing
// the MergeChains and the Rebased arrays for the most recent
// operations.  To keep book-keeping simple, these arrays are
// actually maintained sparsely -- i.e. they have zero values
// for all but the tail of the operations.  The MinIndex field
// records the minimum index value for which these arrays have
// a valid value.   Note that the IDToIndexMap is not sparse,
//
// When working with a sparse array, one can receive a
// ErrLogNeedsBackfilling error indicating that an operation
// was received that refers to operations earlier than MinIndex
// and os the log needs to be backfilled to continue.
//
// This struct is not thread-safe.
type Log struct {
	Transformer
	MinIndex int
	*ModelImage
	Rebased      []Operation
	MergeChains  [][]Operation
	IDToIndexMap map[string]int
}

// TrimMergeChain trims the merge chain of all operations that
// match the provided BasisID or appearered earlier in the log.
//
// This makes a frequent appearance in code because the particular
// state already has the effects of everything until the BasisID
// factored in (by definition of BasisID) and so this method
// removes such operations to prevent double applying of an op.
//
// The code assumes that the mergeChain is properly ordered in
// journal order.
func (l *Log) TrimMergeChain(mergeChain []Operation, basisID string) []Operation {
	if basisID == "" {
		return mergeChain
	}
	basisIndex := l.IDToIndexMap[basisID]
	for len(mergeChain) > 0 && l.IDToIndexMap[mergeChain[0].ID] <= basisIndex {
		mergeChain = mergeChain[1:]
	}
	return mergeChain
}

// getMergeTarget returns the sequence of known operations that
// happened after the provided basis and which has factored in
// the parent operation.
// This is just the concatenation of operations in the merge chain
// of the parent operation and all operations that followed the
// parent -- and then trimmed against the basis.
func (l *Log) getMergeTarget(parentID, basisID string, parentIndex, basisIndex int) []Operation {
	if parentID == "" && basisID == "" {
		return l.Rebased
	}
	if parentID == "" || basisIndex >= parentIndex {
		return l.Rebased[basisIndex+1:]
	}
	against := l.joinOperation(l.MergeChains[parentIndex], l.Rebased[parentIndex+1:])
	return l.TrimMergeChain(against, basisID)
}

// TransformOperation takes a raw journal operation; it transforms
// that operation into the Rebased and MergeChain aspects using the
// rebased and merge chains of its own parents.
//
// It does not modify the log but returns the Rebased/MergeChain for
// the operation.
func (l *Log) TransformOperation(op Operation) ([]Operation, []Operation, error) {
	// ignore duplicates
	if _, ok := l.IDToIndexMap[op.ID]; ok {
		return nil, nil, nil
	}

	if err := l.validateOp(op); err != nil {
		return nil, nil, err
	}

	bIndex := l.IDToIndexMap[op.BasisID()]
	pIndex := l.IDToIndexMap[op.ParentID()]
	against := l.getMergeTarget(op.ParentID(), op.BasisID(), pIndex, bIndex)

	left, right, ok := l.TryMergeOperations(against, []Operation{op})
	if !ok {
		return nil, nil, ErrInvalidOperation
	}
	return left, right, nil
}

func (l *Log) validateOp(op Operation) error {
	// calculate parent and basis indices
	basisID, parentID := op.BasisID(), op.ParentID()
	basisIndex, basisExists := l.IDToIndexMap[basisID]
	_, parentExists := l.IDToIndexMap[parentID]

	// validate parent and basis IDs exist in the map
	if basisID != "" && !basisExists || parentID != "" && !parentExists {
		return ErrMissingParentOrBasis
	}

	// ensure basis is loaded.
	// do not need to check that parent is loaded because if parent was earlier
	// than basis, we ignore the parent anyway and if parent is later than basis
	// and basis is loaded, parent is effectively loaded
	if (basisID == "" && l.MinIndex > 0) || (basisID != "" && basisIndex < l.MinIndex) {
		return ErrLogNeedsBackfilling
	}

	return nil
}

// AppendOperation takes a raw operations, transforms it appropriately
// and stores the rebased operation in Rebased.  It also updates
// IDToIndexMap and MergeChains as they are used for the actual
// transformation process.
//
// AppendOperation is setup such that a sequence of operations can be
// appended in the Journal order to create a log.
//
// It returns either ErrLogNeedsBackfilling or ErrMissingParentOrBasis.
func (l *Log) AppendOperation(op Operation) error {
	left, right, err := l.TransformOperation(op)

	if err != nil || left == nil {
		return err
	}

	l.Rebased = append(l.Rebased, left...)
	l.MergeChains = append(l.MergeChains, right)
	if l.IDToIndexMap == nil {
		l.IDToIndexMap = map[string]int{}
	}
	l.IDToIndexMap[op.ID] = len(l.Rebased) - 1
	return nil
}
