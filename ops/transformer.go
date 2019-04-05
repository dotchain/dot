// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import (
	"context"

	"github.com/dotchain/dot/changes"
)

// opInfo contains the transformed op (remote master) as well as the
// sequence of merged operations (local branch).
//
// The transformed op is based on the op with the previous version and
// can be applied on top of that to get the same effect as the Op
// provided.
//
// The sequence of merge operations can be applied on top of the
// current Op to get it to the same state as one would end up applying
// the transformed op on top of all previous operations.
type opInfo struct {
	xformed Op
	merge   []Op
}

// Transformed takes a store of raw operations and converts them to a
// transformed version that can be applied in sequence.
//
// Every operation has a Basis (last server acknowledged server
// operation applied before) and a Parent (the previous client
// operation on top of which the current operation is applied)
//
// The different basis/parent values imply that raw operations cannot
// directly be applied on each other but need to be transformed. This
// function returns a store which silently transforms the operations
// underneath and returns that (so clients can apply them to get
// accumulated state)
//
// The transformed Store only modifies the GetSince method.
//
// The Cache interface matches sync.Map and is required for
// efficiency.  It does not need to be a synchronized map if all calls
// to Transformed are guaranteed to be non-concurrent which is
// generally the case for client uses (as with Connect)
func Transformed(raw Store, cache Cache) Store {
	return transformer{raw, cache}
}

type transformer struct {
	Store
	Cache
}

func (t transformer) GetSince(ctx context.Context, version, count int) ([]Op, error) {
	ops, err := t.Store.GetSince(ctx, version, count)
	if err != nil || len(ops) == 0 {
		return ops, err
	}

	result := make([]Op, len(ops))
	for kk := range ops {
		// Transform all the returned operations
		opInfo, err := t.xform(ctx, ops[kk])
		if err != nil {
			return nil, err
		}
		result[kk] = opInfo.xformed
	}
	return result, nil
}

// xform transforms a single operation, recursively calling itself on
// other operations it needs. The cache is used to "memoize" such
// prior results to avoid redoing them.
//
// The returned opInfo includes the the transformed version as well as
// the collection of merge operations. The merge operations can be
// used to get the client state (after it had applied the raw
// operation) to the converged state (i.e. by sequental application of
// all the transformed operations until this raw operation).
func (t transformer) xform(ctx context.Context, op Op) (opInfo, error) {
	basis, version := op.Basis(), op.Version()

	// if the operation is available in the cache, just return it
	if result, ok := t.Cache.Load(version); ok {
		return result.(opInfo), nil
	}

	// if this operation is based on the last operation in the
	// store, there is no transformation needed
	gap := version - basis - 1
	if gap == 0 {
		return opInfo{op, nil}, nil
	}

	// fetch all the operations since the basis
	ops, err := t.Store.GetSince(ctx, basis+1, gap)
	if err != nil {
		return opInfo{}, err
	}

	// skip all those before the parent of the current op
	for op.Parent() != nil && ops[0].ID() != op.Parent() {
		ops = ops[1:]
	}

	// The current op is on top of the parent op and so
	// the parent op should first be transformed against the
	// "merge chain" of the parent op
	result, err := t.getMergeChain(ctx, op, ops)
	if err != nil {
		return result, err
	}

	if op.Parent() != nil {
		// skip parent op
		ops = ops[1:]
	}

	// The residual op needs to be merged against all ops
	// since the "parent" or "basis".  For each of these,
	// we need the transformed version first
	for _, opx := range ops {
		info, err := t.xform(ctx, opx)
		if err != nil {
			return opInfo{}, err
		}
		result.xformed, opx = t.merge(info.xformed, result.xformed)
		result.merge = append(result.merge, opx)
	}

	// stash the result to avoid calculating this again
	t.Cache.Store(version, result)
	return result, nil
}

// getMergeChain gets the merge chain for the op and transforms the
// current op against that to get the updated current op as well as
// its merge chain
func (t transformer) getMergeChain(ctx context.Context, op Op, ops []Op) (opInfo, error) {
	parent, basis := op.Parent(), op.Basis()

	if parent == nil {
		return opInfo{op, nil}, nil
	}

	info, err := t.xform(ctx, ops[0])
	if err != nil {
		return opInfo{}, err
	}

	merge := info.merge
	for len(merge) > 0 && merge[0].Version() <= basis {
		merge = merge[1:]
	}

	mergeChain := []Op(nil)
	for _, opx := range merge {
		op, opx = t.merge(opx, op)
		mergeChain = append(mergeChain, opx)
	}

	return opInfo{op, mergeChain}, nil
}

// merge merges the changes in two operations
func (t transformer) merge(left, right Op) (Op, Op) {
	lx, rx := changes.Merge(left.Changes(), right.Changes())
	return right.WithChanges(lx), left.WithChanges(rx)
}
