// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package refs implements reference paths, carets and selections.
//
// A reference path, caret or selection refers to an item or a
// position in an array-like object or a set of items in an array-like
// object.  As changes are applied, the path may be affected as well
// as items that the path refers to. This package provides the
// mechanism to deal with these.
package refs

import "github.com/dotchain/dot/changes"

// Ref represents the core Reference type
type Ref interface {
	// Merge takes a change and returns an Ref that reflects the
	// effect of the change.  If the change affects the item
	// specified by the Ref, it returns a modified version of the
	// change that can be applied to the value at the Ref.
	Merge(c changes.Change) (Ref, changes.Change)

	// Equal tests if two refs are the same.
	Equal(other Ref) bool
}

// InvalidRef refers to a ref that no longer exists.
var InvalidRef = invalidRef{}

type invalidRef struct{}

func (r invalidRef) Merge(c changes.Change) (Ref, changes.Change) {
	return r, nil
}

func (r invalidRef) Equal(other Ref) bool {
	_, ok := other.(invalidRef)
	return ok
}

// PathMerger is the interface that custom Change types should
// implement.
type PathMerger interface {
	MergePath(p []interface{}) *MergeResult
}

// MergeResult contains the result of calling Merge on a path. If the
// path is invalidated by the change, the whole result is nil.
// Otherwise the field P specifies the updated path.  Affected returns
// a version of the change that can be applied on the object at the
// path (and is set to nil if the chagne does not have any local
// effect). Unaffected contains any changes that do not affect the
// path.  Unaffected+Affected should be the equivalent of the original
// change.
type MergeResult struct {
	P          []interface{}
	Scoped     changes.Change
	Affected   changes.Change
	Unaffected changes.Change
}

func (p *MergeResult) join(o *MergeResult) *MergeResult {
	if o == nil {
		return nil
	}
	p.P = o.P
	p.Scoped = p.joinChanges(p.Scoped, o.Scoped)
	p.Affected = p.joinChanges(p.Affected, o.Affected)
	p.Unaffected = p.joinChanges(p.Unaffected, o.Unaffected)
	return p
}

func (p *MergeResult) joinChanges(c1, c2 changes.Change) changes.Change {
	switch {
	case c1 == nil:
		return c2
	case c2 == nil:
		return c1
	}
	if c1x, ok := c1.(changes.ChangeSet); ok {
		if c2x, ok := c2.(changes.ChangeSet); ok {
			return append(c1x, c2x...)
		}
		return append(c1x, c2)
	}
	if c2x, ok := c2.(changes.ChangeSet); ok {
		return append(changes.ChangeSet{c1}, c2x...)
	}
	return changes.ChangeSet{c1, c2}
}

// Prefix updates the merge result to include the prefix to the path.
// It does not update the Scoped field.
func (p *MergeResult) Prefix(other []interface{}) *MergeResult {
	if p != nil {
		p.P = append(append([]interface{}(nil), other...), p.P...)
		p.Affected = p.pc(other, p.Affected)
		p.Unaffected = p.pc(other, p.Unaffected)
	}
	return p
}

func (p *MergeResult) addPathPrefix(other []interface{}) *MergeResult {
	return p.Prefix(other)
}

func (p *MergeResult) pc(path []interface{}, c changes.Change) changes.Change {
	if c == nil {
		return nil
	}
	return changes.PathChange{Path: path, Change: c}
}
