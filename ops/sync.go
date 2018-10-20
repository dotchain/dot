// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// NewSync connects a op Store to a Stream.
//
// All changes made on the stream are sent upstream to the Store
// immediately.  Changes can be fetched from the store by an explicit
// call to Fetch. These are transformed and applied to the local
// stream.
//
// The newID function is used to create new IDs for the operation. See
// https://godoc.org/github.com/dotchain/dot/x/idgen#New for an
// example implementation
//
// The provided store must fetch transformed operations. Use
// Transformed() to convert a raw store for use with NewSync
func NewSync(transformed Store, version int, local streams.Stream, newID func() interface{}) *Sync {
	s := &Sync{tx: transformed, ver: version, local: local}
	local.Nextf(s, func() {
		var c changes.Change
		local, c = local.Next()

		if s.mergingID != nil {
			return
		}
		id := newID()
		op := Operation{OpID: id, BasisID: s.ver, VerID: -1, ParentID: s.lastSentID, Change: c}
		s.IDs = append(s.IDs, id)
		s.lastSentID = id
		if err := s.tx.Append(context.Background(), []Op{op}); err != nil {
			panic(err)
		}
	})
	return s
}

// Sync holds the state to manage two-way synchronization of a Store
// with a Stream.  See NewSync() for how to setup synchronization
type Sync struct {
	tx         Store
	ver        int
	pending    []Op
	local      streams.Stream
	IDs        []interface{}
	lastSentID interface{}
	mergingID  interface{}
}

// Version returns the latest server version
func (s *Sync) Version() int {
	return s.ver
}

// Close terminates synchronization
func (s *Sync) Close() {
	s.local.Nextf(s, nil)
	*s = Sync{}
}

// Prefetch synchronously fetches the next batch of operations. It
// returns if it found any new operations.  It does not actually
// update the local stream.  This is useful as it is a slow blocking
// network call.  The prefetched operations are stashed until a call
// to ApplyPrefetched
func (s *Sync) Prefetch(ctx context.Context, limit int) (bool, error) {
	var err error
	if len(s.pending) == 0 {
		s.pending, err = s.tx.GetSince(ctx, s.Version()+1, limit)
	}
	return len(s.pending) > 0, err
}

// ApplyPrefetched is guaranteed to not make any blocking calls. It
// only looks at any prefetched operations and applies them.
func (s *Sync) ApplyPrefetched() {
	for _, op := range s.pending {
		s.ver = op.Version()
		if len(s.IDs) > 0 && s.IDs[0] == op.ID() {
			if s.lastSentID == s.IDs[0] {
				s.lastSentID = nil
			}
			s.IDs = s.IDs[1:]
			s.local, _ = s.local.Next()
		} else {
			s.mergingID = op.ID()
			s.local = s.local.ReverseAppend(op.Changes())
		}
	}
	s.mergingID = nil
	s.pending = nil
}

// Fetch synchronously fetches the operations from a store and merges
// them with the local stream.  This is a simple combination of
// Prefetch and ApplyPrefetched
func (s *Sync) Fetch(ctx context.Context, limit int) error {
	_, err := s.Prefetch(ctx, limit)
	s.ApplyPrefetched()
	return err
}
