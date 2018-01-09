// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

// ClientLog is a helper struct that provides the functionality
// needed by clients to deal with operations from other clients
// that get merged into the journal.  It maintains the state
// needed to calculate the "compensating" operations that a
// client can apply to get to the same converged state that
// would result if its inflight operations were merged into the
// journal.  Note that this is a moving target -- as more operations
// get added to the journal and the local client keeps adding
// more operations of its own, the compensating operations need
// to track both and yield a final converged state that mirrors
// what silent observer would end up with if it only tracked
// the server log.
//
// Please read https://github.com/dotchain/dot/docs/IntroductionToOperationalTransforms.md
// for a detailed description of how reconciliation works.
//
// The initialization of the client log involves a few different
// cases.
//
// 1. Client starts from scratch, no model at all but might have
// operations from its previous session on the device that were
// in flight (maybe in the journal, may not)
//
// In this case, the client should use #BootstrapClientLog to
// bootstrap its model.
//
// 2. Client has restarted a session with a cached model at a
// particular basis and potentially some client operations that
// were in flight before.
//
// In this case, the client should use #ReconnectClientLog to continue
// the reconciliation process
//
// Please see
// https://github.com/dotchain/site/blob/master/Protocol.md
// for a better understanding of the use of the ParentID and
// BasisID when sending them to the server.  #AppendClientOperation
// expects these to be properly setup with BasisID being the
// last value in server log when the client applied the operation
// and ParentID being the last client operation applied in the
// current session (or carried over from a previous session)
type ClientLog struct {
	Transformer

	// the following two numbers are 1+ index in server log
	// with the 1 being there because of making it easy to
	// initialize the log with zeroes.

	// 1 + index of last known operation from server log
	// that has been factored into the client log so far
	ServerIndex int

	// Rebased maintains the rebased client operations that
	// have yet to appear in the server log
	Rebased []Operation

	// MergeChain is the sequence of operations to apply after
	// the last rebased operation to get the model into a
	// converged state.  This is empty if rebased is empty
	MergeChain []Operation
}

// Reconcile takes a server log and if there are any operations
// there that have not been added to the client log, it updates
// the client log.  It returns the set of compensating operations
// to apply to the client model to get to the converged state.
func (c *ClientLog) Reconcile(l *Log) ([]Operation, error) {
	var ok bool

	if c.ServerIndex+1 <= l.MinIndex {
		return nil, ErrLogNeedsBackfilling
	}

	rebased, merge := c.Rebased, []Operation{}
	serverIndex := c.ServerIndex
	for _, op := range l.Rebased[c.ServerIndex:] {
		serverIndex++
		if len(rebased) > 0 && rebased[0].ID == op.ID {
			rebased = rebased[1:]
			continue
		}
		var m []Operation
		rebased, m, ok = c.TryMergeOperations([]Operation{op}, rebased)
		if !ok {
			return nil, ErrInvalidOperation
		}
		merge = append(merge, m...)
	}

	c.ServerIndex = len(l.Rebased)
	if len(rebased) > 0 {
		c.Rebased = append([]Operation{}, rebased...)
		c.MergeChain = append(c.MergeChain, merge...)
	} else {
		c.Rebased = nil
		c.MergeChain = nil
	}
	return merge, nil
}

// AppendClientOperation appends a client operation to the client log.
//
// It returns an error if the server log needs backfilling. The returned
// set of compensating operations can be used by the client to update
// its state to factor in the effect of any unaccounted ops in the log.
func (c *ClientLog) AppendClientOperation(l *Log, op Operation) ([]Operation, error) {
	if index, ok := l.IDToIndexMap[op.ID]; ok {
		if index < l.MinIndex {
			return nil, ErrLogNeedsBackfilling
		}

		c.Rebased = nil
		c.MergeChain = l.joinOperation(l.MergeChains[index], l.Rebased[index+1:])
		c.ServerIndex = len(l.Rebased)
		return c.MergeChain, nil
	}

	basisID, parentID := op.BasisID(), op.ParentID()
	bIndex, bExists := l.IDToIndexMap[basisID]
	pIndex, pExists := l.IDToIndexMap[parentID]

	if basisID != "" && !bExists {
		return nil, ErrMissingParentOrBasis
	}

	if bIndex < l.MinIndex {
		// TODO: real check should be
		// !pExists && bIndex < l.MinIndex ||
		// pExists && pIndex < l.MinIndex
		return nil, ErrLogNeedsBackfilling
	}

	var merge []Operation
	if len(c.Rebased) == 0 || (pExists && pIndex >= bIndex) {
		if parentID != "" && !pExists {
			return nil, ErrMissingParentOrBasis

		}

		merge = l.getMergeTarget(parentID, basisID, pIndex, bIndex)
	} else {
		merge = l.TrimMergeChain(c.MergeChain, basisID)
	}
	rebased, merged, ok := l.TryMergeOperations(merge, []Operation{op})
	if !ok {
		return nil, ErrInvalidOperation
	}

	c.Rebased = append(c.Rebased, rebased...)
	c.MergeChain = append([]Operation{}, merged...)
	c.ServerIndex = len(l.Rebased)
	return merged, nil
}

// BootstrapClientLog creates a new client log for a client that does
// not have a model.
//
// Errors: It returns ErrMissingParentOrBasis if the log
// has not advanced enough for the clientOps.  It can return
// ErrLogNeedsBackfilling if the log is not backfilled enough for the
// operation to complete.
//
// It returns the client log and a pair of operation collections. The
// first operation collection is the set of rebased server operations
// a client can apply to get to a good state and the second operations
// collection is the set of rebased client operations which is meant
// to be applied on top of the server rebased.
func BootstrapClientLog(l *Log, clientOps []Operation) (*ClientLog, []Operation, []Operation, error) {
	clog, _, err := newClientLog(l, clientOps, "", "")
	if err != nil {
		return nil, nil, nil, err
	}
	c := append([]Operation{}, l.Rebased...)
	r := append([]Operation{}, clog.Rebased...)
	return clog, c, r, nil
}

// ReconnectClientLog creates a new client log for a client that has
// an existing model (with the provided parentID and basisID). Note
// that if client operations are provided, the parentID will be
// ignored and the last op in that list will be used.
//
// Errors: It returns ErrMissingParentOrBasis if the log
// has not advanced enough for the clientOps.  It can return
// ErrLogNeedsBackfilling if the log is not backfilled enough for the
// operation to complete.
//
// It also returns a set of operations that the client can apply to
// get it back to a mainline state
func ReconnectClientLog(l *Log, clientOps []Operation, basisID, parentID string) (*ClientLog, []Operation, error) {
	return newClientLog(l, clientOps, basisID, parentID)
}

func newClientLog(l *Log, clientOps []Operation, basisID, parentID string) (*ClientLog, []Operation, error) {
	clog := &ClientLog{Transformer: l.Transformer}

	if _, ok := l.IDToIndexMap[basisID]; !ok && basisID != "" {
		return nil, nil, ErrMissingParentOrBasis
	}

	if len(clientOps) > 0 {
		parentID = clientOps[len(clientOps)-1].ID
	} else if _, ok := l.IDToIndexMap[parentID]; !ok && parentID != "" {
		return nil, nil, ErrMissingParentOrBasis
	}

	var merge []Operation
	for _, op := range clientOps {
		merged, err := clog.AppendClientOperation(l, op)
		if err != nil {
			return nil, nil, err
		}
		merge = merged
	}

	clog.ServerIndex = len(l.Rebased)
	if merge != nil {
		merge = l.TrimMergeChain(merge, basisID)
		return clog, merge, nil
	}

	bIndex := l.IDToIndexMap[basisID]
	pIndex := l.IDToIndexMap[parentID]
	merge = l.getMergeTarget(parentID, basisID, pIndex, bIndex)
	return clog, merge, nil
}
