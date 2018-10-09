// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs

import "github.com/dotchain/dot/changes"

// Path represents a reference to a value at a specific path. A nil
// or empty path refers to the root value.
//
// This is an immutable type -- none of the methods modify the
// provided path itself.
//
// This only handles the standard set of changes. Custom changes
// should implement the PathMerger interface.
//
// If no such method is implemented by the change, the change is
// ignored as if it has no side-effects.
type Path []interface{}

// Merge implements Ref.Merge
func (p Path) Merge(c changes.Change) (Ref, changes.Change) {
	if result := Merge(p, c); result != nil {
		return Path(result.P), result.Affected
	}
	return InvalidRef, nil
}

// Equal implements equality comparison
func (p Path) Equal(o Path) bool {
	if len(p) != len(o) {
		return false
	}
	for kk, elt := range p {
		if o[kk] != elt {
			return false
		}
	}
	return true
}

// Merge merges a path with a change. If the path is invalidated, it
// returns nil. Otherwise, it returns the updated path. The version of
// the change that can be applied to the just object at the path
// itself is in Affected.  Unaffected holds the changes that does not
// concern the provided path.
//
// Custom changes should implement the PathMerger interface or the
// change will be considered as not affecting the path in any way
//
// For most purposes, the Path type is a better fit than directly
// calling Merge.
func Merge(p []interface{}, c changes.Change) *MergeResult {
	if len(p) == 0 {
		return &MergeResult{nil, c, nil}
	}

	switch c := c.(type) {
	case changes.Replace:
		return nil
	case changes.Splice:
		return mergeSplice(p, c)
	case changes.Move:
		return mergeMove(p, c)
	case changes.PathChange:
		idx := 0
		for len(p) > idx && len(c.Path) > idx {
			if p[idx] == c.Path[idx] {
				idx++
				continue
			}
			return &MergeResult{P: p, Unaffected: c}
		}
		if len(p) == idx {
			c.Path = c.Path[idx:]
			return &MergeResult{p, c, nil}
		}

		return Merge(p[idx:], c.Change).addPathPrefix(p[:idx])
	case changes.ChangeSet:
		result := &MergeResult{P: p}
		for _, cx := range c {
			result = result.join(Merge(result.P, cx))
			if result == nil {
				return nil
			}
		}
		return result
	case PathMerger:
		return c.MergePath(p)
	}

	return &MergeResult{P: p, Unaffected: c}
}

func mergeMove(p []interface{}, c changes.Move) *MergeResult {
	idx := c.MapIndex(p[0].(int))
	return &MergeResult{
		P:          append([]interface{}{idx}, p[1:]...),
		Unaffected: c,
	}
}

func mergeSplice(p []interface{}, c changes.Splice) *MergeResult {
	idx, ok := c.MapIndex(p[0].(int))
	if ok {
		return nil
	}
	return &MergeResult{
		P:          append([]interface{}{idx}, p[1:]...),
		Unaffected: c,
	}
}
