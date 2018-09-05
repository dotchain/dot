// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

import "context"

// Transformed creates a new Store which returns transformed
// operations.
//
// The regular sequence of operations stored cannot be directly
// applied one on top of another because they may have different basis
// and parent values.
//
// The transformed store returns the same operations but modifies the
// changes so that it these change have the same effect as the
// original but can be applied in sequence.
func Transformed(raw Store) Store {
	return transformer{raw, nullCache{}}
}

// TransformedWithCache is the same as Transformed but with the
// addition of an internal cache for performance.  The cache
// interfaace is compatible with sync.Map.  The keys in the map are
// the version numbers of the ops in the store -- the caller can
// implement their eviction strategy based on this if needed (such as
// evicting old version numbers). Custom caches can be implemented for
// more sophisticated use cases.
func TransformedWithCache(raw Store, cache Cache) Store {
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
		opInfo, err := t.xform(ctx, ops[kk])
		if err != nil {
			return nil, err
		}
		result[kk] = opInfo.xformed
	}
	return result, nil
}

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

func (t *transformer) xform(ctx context.Context, op Op) (opInfo, error) {
	basis, version, parent := op.Basis(), op.Version(), op.Parent()

	if result, ok := t.Cache.Load(version); ok {
		return result.(opInfo), nil
	}

	if version == basis+1 {
		return opInfo{op, nil}, nil
	}

	ops, err := t.Store.GetSince(ctx, basis+1, version-basis-1)
	if err != nil {
		return opInfo{}, err
	}

	if parent != nil {
		for ops[0].ID() != parent {
			ops = ops[1:]
		}
		info, err := t.xform(ctx, ops[0])
		if err != nil {
			return opInfo{}, err
		}
		for len(info.merge) > 0 && info.merge[0].Version() <= basis {
			info.merge = info.merge[1:]
		}
		ops = append(info.merge, ops[1:]...)
	}

	mergeChain := make([]Op, len(ops))

	for kk, opx := range ops {
		info, err := t.xform(ctx, opx)
		if err != nil {
			return opInfo{}, err
		}
		op, mergeChain[kk] = t.merge(info.xformed, op)
	}

	result := opInfo{op, mergeChain[:len(mergeChain):len(mergeChain)]}
	t.Cache.Store(version, result)
	return result, nil
}

func (t transformer) merge(left, right Op) (Op, Op) {
	leftChanges := left.Changes()
	if leftChanges == nil {
		return right, left
	}
	lx, rx := leftChanges.Merge(right.Changes())
	return right.WithChanges(lx), left.WithChanges(rx)
}
