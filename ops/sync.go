// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"
	"github.com/dotchain/dot/changes"
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
func NewSync(transformed Store, version int, local changes.Stream, newID func() string) *Sync {
	s := &Sync{tx: transformed, ver: version, local: local}
	local.Nextf(s, func(c changes.Change, updated changes.Stream) {
		if s.mergingID != "" {
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
	local      changes.Stream
	IDs        []string
	lastSentID string
	mergingID  string
}

// Close terminates synchronization
func (s *Sync) Close() {
	s.local.Nextf(s, nil)
	*s = Sync{}
}

// Fetch synchronously fetches the operations from a store and merges
// them with the local stream.
func (s *Sync) Fetch(ctx context.Context, limit int) error {
	ops, err := s.tx.GetSince(ctx, s.ver+1, limit)
	if err == nil {
		for _, op := range ops {
			s.ver = op.Version()
			if len(s.IDs) > 0 && s.IDs[0] == op.ID() {
				if s.lastSentID == s.IDs[0] {
					s.lastSentID = ""
				}
				s.IDs = s.IDs[1:]
				_, s.local = s.local.Next()
			} else {
				s.mergingID = op.ID().(string)
				s.local = s.local.ReverseAppend(op.Changes())
			}
		}
		s.mergingID = ""
	}
	return err
}
