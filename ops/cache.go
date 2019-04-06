// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

// Cache holds transformation state for each operation.
//
// The transformed op is a version of the original op but one which
// can be applied directly on top of the previous version. The
// untransformed operation generally cannot be applied because it
// might have a Basis() or Parent() that is much older than the
// previous version as it may not have merged with those.
//
// The mergeChain is the counterpart: the client which sent the
// original op can apply this to get to the same state as applying the
// transformed op on the previous version.
type Cache interface {
	// Load returns nil if the items are not present in the cache.
	Load(version int) (transformed Op, mergeChain []Op)

	// Store saves the transformed op and its merge chain.
	Store(version int, transformed Op, mergeChain []Op)
}
