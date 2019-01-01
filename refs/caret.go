// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs

import "github.com/dotchain/dot/changes"

// Caret is a reference to a specific position in an array-like
// object.
//
// This is an immutable type
//
// This only handles the standard set of changes. Custom changes
// should implement a MergeCaret method:
//
//    MergeCaret(caret refs.Caret) (refs.Ref)
//
// Note that this is in addition to the MergePath method which is
// called first to transform the path and then the MergeCaret is
// called  on the updated Caret (based on the path returned by
// MergePath).
//
// The IsLeft flag controls whether the position sticks with the
// element to the left in case of insertions happening at the
// index. The default is for the reference to stick to the element on
// the right.
type Caret struct {
	Path
	Index  int
	IsLeft bool
}

// Merge updates the caret index based on the change.  Note that it
// always returns a nil change as there is no way for a change to
// affect the caret.
func (caret Caret) Merge(c changes.Change) (Ref, changes.Change) {
	px, cx := caret.Path.Merge(c)
	if px == InvalidRef {
		return px, cx
	}
	return caret.updateIndex(px.(Path), caret.Index, cx), nil
}

func (caret Caret) updateMoveIndex(path Path, idx int, cx changes.Move) int {
	dest := cx.Offset + cx.Distance
	if cx.Distance > 0 {
		dest += cx.Count
	}
	switch {
	case dest != idx:
		idx = cx.MapIndex(idx)
	case caret.IsLeft && cx.Distance > 0:
		idx -= cx.Count
	case !caret.IsLeft && cx.Distance < 0:
		idx += cx.Count
	}
	return idx
}

func (caret Caret) updateIndex(path Path, idx int, cx changes.Change) Ref {
	switch cx := cx.(type) {
	case changes.Replace:
		return InvalidRef
	case changes.Splice:
		if cx.Offset != idx || !caret.IsLeft {
			idx, _ = cx.MapIndex(idx)
		}
	case changes.Move:
		idx = caret.updateMoveIndex(path, idx, cx)
	case changes.PathChange:
		if len(cx.Path) == 0 {
			return caret.updateIndex(path, idx, cx.Change)
		}
	case changes.ChangeSet:
		for _, c := range cx {
			ref := caret.updateIndex(path, idx, c)
			if ref == InvalidRef {
				return ref
			}
			idx = ref.(Caret).Index
		}
	case caretMerger:
		return cx.MergeCaret(Caret{path, idx, caret.IsLeft})
	}
	return Caret{path, idx, caret.IsLeft}
}

// Equal impements Ref.Equal
func (caret Caret) Equal(other Ref) bool {
	o, ok := other.(Caret)
	return ok && caret.Path.Equal(o.Path) && caret.Index == o.Index && caret.IsLeft == o.IsLeft
}

type caretMerger interface {
	MergeCaret(caret Caret) Ref
}
