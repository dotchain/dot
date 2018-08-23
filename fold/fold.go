// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package fold implements folded data structures
//
// Folded data structures contain locally applied transient
// changes. A typical "Text editor" example is when a region of text
// is "hidden" temporarily. The effect is only maintained on local
// copies and further changes on top of folded data structures are
// correctly applied on the "Remote" value.  At any point, any of the
// changes can be "unfolded".
//
// The package implements "Folding" as a simple type to help manage
// the set of "hidden" changes  and translate between local and remote
// changes.
//
// A higher level primitive is the Foldable type.  This allows adding
// and removing folds (via Fold and Unfold respectively) as well as
// applying changes locally via the Folded type and remotely via the
// Unfolded type. The local and remote values are expected to be
// fetched via the Local and Remote properties respectively.
//
// Note that all the types here are immutable and while this implies
// they are thread-safe, the changes made on one copy will not be
// reflected on the other in any fashion.  The ver package
// (https://godoc.org/github.com/dotchain/ver) implements a scheme
// that is thread-safe as well as automerged.  At some point, this
// package will play well with that or the functionality will get
// exposed into that library.
package fold

import (
	"github.com/dotchain/dot"
	"github.com/dotchain/dot/encoding"
)

// Folding is the basic structure which holds a set of changes that
// are in a holding pattern: changes done on the local version should
// be transformed so that it can be applied on the remote side and
// changes on the remote side should be transformed so it can be
// applied on the local copy.
//
// Folding is an immutable type that provides the required methods to
// manage this.
type Folding struct {
	dot.Transformer
	Changes []dot.Change
}

// TransformLocal takes a set of local changes and find the
// transformed set of changes that should be applied on the remote if
// the "folding" changes were not present.
func (f Folding) TransformLocal(changes []dot.Change) (updated Folding, remote []dot.Change) {
	updated = f
	updated.Changes, remote = f.Transformer.ReorderChanges(f.Changes, changes)
	return updated, remote
}

// TransformRemote takes a set of remote changes and figures out the
// transformed version which can be applied on top of the folding
// changes.
func (f Folding) TransformRemote(changes []dot.Change) (updated Folding, local []dot.Change) {
	updated = f
	updated.Changes, local = f.Transformer.MergeChanges(changes, f.Changes)
	return updated, local
}

// Unfold allows any changes (described by offset and count) in the
// Folding structure to be removed.  This should have no effect on the
// remote side but the local side may now have changes that need to be
// applied to effectively "revert" the previously folded
// changes.
func (f Folding) Unfold(offset, count int) (updated Folding, reverts []dot.Change) {
	updated = f
	removed := f.Changes[offset : offset+count : offset+count]
	rest := f.Changes[offset+count : len(f.Changes) : len(f.Changes)]
	removed, rest = f.Transformer.ReorderChanges(removed, rest)
	reverts = dot.Operation{Changes: removed}.Undo().Changes
	updated.Changes = append(f.Changes[:offset], rest...)
	return updated, reverts
}

// Foldable is a higher level version of Folding that keeps track of
// both the local and remote values along with the changes that are a
// part of the folding.
//
// The local and remote values should be normalized before use via the
// LocalValue() and RemoteValue() functions. To operate on the folded
// structure (to make local changes), simple use Folded:
//
//       f = Folded(f).Apply(...)
//
// To operate on the unfolded structure (to apply remote changes), use
// Unfolded the same way:
//
//       f = Unfolded(f).Apply(remoteChanges)
//
//
// Local and Remote can be string, []interface{},
// map[string]interface{} or richer types that implement
// encoding.UniversalEncoding (see
// https://godoc.org/github.com/dotchain/dot/encoding).
type Foldable struct {
	dot.Utils
	Changes       []dot.Change
	Local, Remote interface{}
}

// LocalValue returns the normalized local value. 
func (f Foldable) LocalValue() interface{} {
	return encoding.Normalize(f.Local)
}

// RemoteValue returns the normalized remote value
func (f Foldable) RemoteValue() interface{} {
	return encoding.Normalize(f.Remote)
}

func (f Foldable) update(changes, local, remote []dot.Change) Foldable {
	f.Changes = changes[:len(changes):len(changes)]
	f.Local = f.Utils.Apply(f.Local, local)
	f.Remote = f.Utils.Apply(f.Remote, remote)
	return f
}

// Fold adds a change to the folded changes list.  The Local value
// will be updated to reflect this.
func (f Foldable) Fold(changes []dot.Change) Foldable {
	return f.update(append(f.Changes, changes...), changes, nil)
}

// Unfold removes the change specified by the offset and count (which
// point to the Changes array).  The Local value is updated to reflect
// the effect of removing those changes.
func (f Foldable) Unfold(offset, count int) Foldable {
	folding, reverts := Folding{dot.Transformer(f.Utils), f.Changes}.Unfold(offset, count)
	return f.update(folding.Changes, reverts, nil)
}

// Folded is a simple type to manage applying changes on the folded
// (or "local") value.
type Folded Foldable

// Apply modifies the local value by the provided changes and also
// propagates the transformed version on the Remote value.
func (f Folded) Apply(local []dot.Change) Foldable {
	folding, remote := Folding{dot.Transformer(f.Utils), f.Changes}.TransformLocal(local)
	return Foldable(f).update(folding.Changes, local, remote)
}

// Unfolded is a simple type to manage applying changes on the
// unfolded (or "remote") value.
type Unfolded Foldable

// Apply modifies the remote value by the provided changes and also
// propagates the transformed version on the local value.
func (f Unfolded) Apply(remote []dot.Change) Foldable {
	folding, local := Folding{dot.Transformer(f.Utils), f.Changes}.TransformRemote(remote)
	return Foldable(f).update(folding.Changes, local, remote)
}
