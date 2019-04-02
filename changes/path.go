// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// PathChange represents a change at the provided "path" which can
// consist of strings (for map-like objects) and integers for
// array-like objects. In particular, each element of the path should
// be a proper comparable value (so slices and such cannot be part of
// th path)
type PathChange struct {
	Path []interface{}
	Change
}

// Revert implements Change.Revert
func (pc PathChange) Revert() Change {
	if pc.Change == nil {
		return nil
	}
	return PathChange{pc.Path, pc.Change.Revert()}
}

// Merge implements Change.Merge
func (pc PathChange) Merge(o Change) (Change, Change) {
	opc, ok := o.(PathChange)
	if !ok {
		opc = PathChange{nil, o}
	}
	return pc.mergePathChange(opc, false)
}

// ReverseMerge implements but with receiver and arg interchanged.
func (pc PathChange) ReverseMerge(o Change) (Change, Change) {
	opc, ok := o.(PathChange)
	if !ok {
		opc = PathChange{nil, o}
	}
	return pc.mergePathChange(opc, true)
}

// ApplyTo is not relevant to PathChange.  It only works when the path
// is empty. In all other cases, it panics.
func (pc PathChange) ApplyTo(ctx Context, v Value) Value {
	if len(pc.Path) == 0 {
		return v.Apply(ctx, pc.Change)
	}
	panic("Unexpected use of PathChange.ApplyTo")
}

func (pc PathChange) mergePathChange(o PathChange, reverse bool) (Change, Change) {
	prefixLen := pc.commonPrefixLen(pc.Path, o.Path)
	switch {
	case len(pc.Path) != prefixLen && len(o.Path) != prefixLen:
		return o, pc
	case len(pc.Path) == prefixLen && len(o.Path) == prefixLen:
		return pc.prefixMerge(pc.Path, pc.Change, o.Change, reverse)
	case len(pc.Path) == prefixLen:
		return pc.mergeSubPath(o, reverse)
	}

	return swap(o.mergeSubPath(pc, !reverse))
}

func (pc PathChange) prefixMerge(prefix []interface{}, l, r Change, reverse bool) (Change, Change) {
	rev, ok := l.(Custom)
	switch {
	case ok && reverse:
		l, r = rev.ReverseMerge(r)
	case reverse && r != nil:
		r, l = r.Merge(l)
	case l != nil:
		l, r = l.Merge(r)
	case !reverse && l == nil:
		l, r = r, l
	}
	return PathChange{prefix, l}, PathChange{prefix, r}
}

func (pc PathChange) updateSubPathIndex(o PathChange, idx int) (Change, Change) {
	path := append([]interface{}(nil), o.Path...)
	path[len(pc.Path)] = idx
	return PathChange{path, o.Change}, pc
}

func (pc PathChange) mergeSubPath(o PathChange, reverse bool) (Change, Change) {
	sub := o.Path[len(pc.Path):]
	switch change := pc.Change.(type) {
	case nil:
		return o, nil
	case Replace:
		change.Before = change.Before.Apply(nil, PathChange{sub, o.Change})
		return nil, PathChange{pc.Path, change}
	case Splice:
		idx := sub[0].(int)
		beforeSize, afterSize := change.Before.Count(), change.After.Count()
		switch {
		case idx < change.Offset:
			return o, pc
		case idx >= change.Offset+beforeSize:
			return pc.updateSubPathIndex(o, idx+afterSize-beforeSize)
		}
		sub := append([]interface{}(nil), sub...)
		sub[0] = idx - change.Offset
		change.Before = change.Before.ApplyCollection(nil, PathChange{sub, o.Change})
		return nil, PathChange{pc.Path, change}
	case Move:
		idx := sub[0].(int)
		dest, end := change.dest(), change.Offset+change.Count
		switch {
		case idx >= change.Offset && idx < end:
			return pc.updateSubPathIndex(o, idx+change.Distance)
		case idx >= dest && idx < change.Offset:
			return pc.updateSubPathIndex(o, idx+change.Count)
		case idx >= end && idx < dest:
			return pc.updateSubPathIndex(o, idx-change.Count)
		}
		return o, pc
	}

	return pc.prefixMerge(pc.Path, pc.Change, PathChange{sub, o.Change}, reverse)
}

func (pc PathChange) commonPrefixLen(a, b []interface{}) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	for kk, elt := range a {
		if b[kk] != elt {
			return kk
		}
	}
	return len(a)
}
