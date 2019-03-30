// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import "context"

// Store is the interface to talk to a Op store.
//
// An Op store is an append-only store which guarantees unique order
// of operations. It does not make any guarantees on the specific
// order beyond the requirement that all operations provided to a
// single Append call will always be together in that order and that a
// sequential calls to Append on a single client/goroutine will be
// stored in that order.
//
// The store does make guarantee that operations will not be
// duplicated. If an operation is appended with an ID that already
// exists, it will silently be dropped.
//
// See https://godoc.org/github.com/dotchain/dot/ops/pg for an example
// implementiation (for Postgres 9.5+)
type Store interface {
	// Append a sequence of operations.  If the operation IDs
	// already exist, those operations are ignored but do not
	// generate an error.
	Append(ctx context.Context, ops []Op) error

	// GetSince returns all operations with version atleast equal
	// to the specified parameter. If the number of operations is
	// larger than the limit, it is truncated.
	//
	// It is not an error if the version does not exist -- an
	// empty result is returned in that case.   Poll can be used
	// to wait for a version
	GetSince(ctx context.Context, version, limit int) ([]Op, error)

	// Poll waits for an operation at the specified version or
	// higher.  Timeouts can be specified via the context and
	// timeouts will pass back the context error.
	//
	// Poll is an opportunistic performance mechanism -- it can
	// have both false positives and false negatives. A trivial
	// implementation can simply use timeouts. A slightly more
	// interesting implementation uses long polls while a more
	// sophisticated implementation can use a dedicated
	// notification channel (websockets or such).
	Poll(ctx context.Context, version int) error

	// Close releases all resources. Any ongoing calls should not
	// be canceled unless the caller cancels them via the context.
	Close()
}
