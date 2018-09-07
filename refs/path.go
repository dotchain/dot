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
// should implement a MergePath method which will then get invoked
// whenever this method is invoked:
//
//     MergePath(path refs.Path) (refs.Ref, changes.Change)
//
// If no such method is implemented by the change, the change is
// ignored as if it has no side-effects.  If such a method is
// implemented, it should only return InvalidRef or another Path.
type Path []interface{}

// Merge implements Ref.Merge
func (p Path) Merge(c changes.Change) (Ref, changes.Change) {
	switch c := c.(type) {
	case changes.Replace:
		return InvalidRef, nil
	case changes.Splice:
		return p.mergeSplice(c)
	case changes.Move:
		return p.mergeMove(c)
	case changes.PathChange:
		return p.mergePathChange(c)
	case changes.ChangeSet:
		return mergeChangeSet(p, c)
	case pathMerger:
		return c.MergePath(p)
	}
	return p, nil
}

func (p Path) mergeSplice(c changes.Splice) (Ref, changes.Change) {
	if len(p) == 0 {
		return p, c
	}
	idx := p[0].(int)
	switch {
	case idx >= c.Offset && idx < c.Offset+c.Before.Count():
		return InvalidRef, nil
	case idx > c.Offset:
		idx += c.After.Count() - c.Before.Count()
		return Path(append([]interface{}{idx}, p[1:]...)), nil
	}
	return p, nil
}

func (p Path) mergeMove(c changes.Move) (Ref, changes.Change) {
	if len(p) == 0 {
		return p, c
	}
	idx := p[0].(int)
	switch {
	case idx >= c.Offset && idx < c.Offset+c.Count:
		idx += c.Distance
	case idx >= c.Offset+c.Distance && idx < c.Offset:
		idx += c.Count
	case idx >= c.Offset+c.Count && idx < c.Offset+c.Count+c.Distance:
		idx -= c.Count
	default:
		return p, nil
	}
	return Path(append([]interface{}{idx}, p[1:]...)), nil
}

func (p Path) mergePathChange(c changes.PathChange) (Ref, changes.Change) {
	p1, p2 := p, c.Path
	for len(p1) > 0 && len(p2) > 0 && p1[0] == p2[0] {
		p1, p2 = p1[1:], p2[1:]
	}
	switch {
	case len(p1) == 0 && len(p2) == 0:
		return p, c.Change
	case len(p1) == 0:
		return p, changes.PathChange{p2, c.Change}
	case len(p2) == 0:
		left := p[: len(p)-len(p1) : len(p)-len(p1)]
		p1, cx := p1.Merge(c.Change)
		if p1 == InvalidRef {
			return p1, nil
		}
		return append(left, p1.(Path)...), cx
	}
	return p, nil
}

type pathMerger interface {
	MergePath(p Path) (Ref, changes.Change)
}
