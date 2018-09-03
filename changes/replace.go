// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Replace represents replacing a particular value with another
// IsDelete is used to distinguish the case where the replace is
// deleting a value as opposed to setting it to nil
//
// Before is ignored for IsInsert and After is ignored for IsDelete.
type Replace struct {
	IsDelete      bool  `json:",omitempty"`
	IsInsert      bool  `json:",omitempty"`
	Before, After Value `json:",omitempty"`
}

// Revert inverts the effect of the replace
func (s Replace) Revert() Change {
	return Replace{s.IsInsert, s.IsDelete, s.After, s.Before}
}

// MergeReplace merges against another Replace change.  The last writer wins
// here with the receiver assumed to be the earlier change
func (s Replace) MergeReplace(other Replace) (other1, s1 *Replace) {
	if s.IsDelete && other.IsDelete {
		return nil, nil
	}

	other.IsInsert = s.IsDelete
	other.Before = s.After
	return &other, nil
}

// MergeSplice merges against a Splice change.  The replace wins
func (s Replace) MergeSplice(other Splice) (other1 *Splice, s1 *Replace) {
	s.Before = s.Before.Apply(other)
	return nil, &s
}

// MergeMove merges against a Move change.  The replace wins
func (s Replace) MergeMove(other Move) (other1 *Move, s1 *Replace) {
	s.Before = s.Before.Apply(other)
	return nil, &s
}

// Merge implements the Change.Merge method
func (s Replace) Merge(other Change) (otherx, cx Change) {
	if other == nil {
		return nil, s
	}

	switch o := other.(type) {
	case Replace:
		return change(s.MergeReplace(o))
	case Splice:
		return change(s.MergeSplice(o))
	case Move:
		return change(s.MergeMove(o))
	case Custom:
		return swap(o.ReverseMerge(s))
	}
	panic("Unexpected change")
}

// Change returns either nil or a Change
func (s *Replace) Change() Change {
	if s == nil {
		return nil
	}
	return *s
}
