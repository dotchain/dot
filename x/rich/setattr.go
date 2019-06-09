// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rich

import (
	"github.com/dotchain/dot/changes"
)

type setattr struct {
	Offset        int
	Name          string
	Before, After values
}

func (s setattr) Revert() changes.Change {
	s.Before, s.After = s.After, s.Before
	return s
}

func (s setattr) ApplyTo(ctx changes.Context, c changes.Value) changes.Value {
	return c.(Text).applySetAttr(s)
}

func (s setattr) Merge(o changes.Change) (ox, sx changes.Change) {
	return s.merge(o, false)
}

func (s setattr) ReverseMerge(o changes.Change) (ox, sx changes.Change) {
	return s.merge(o, true)
}

func (s setattr) merge(o changes.Change, reverse bool) (ox, sx changes.Change) {
	switch o := o.(type) {
	case nil:
		return nil, s
	case setattr:
		return s.mergeSetAttr(o, reverse)
	case changes.Replace:
		o.Before = o.Before.Apply(nil, s)
		return o, nil
	case changes.Splice:
		return s.mergeSplice(o, reverse)
	case changes.Move:
		return s.mergeMove(o, false)
	case changes.PathChange:
		if len(o.Path) == 0 {
			return s.merge(o.Change, reverse)
		}
		return s.mergePath(o, reverse)
	}

	if reverse {
		l, r := o.Merge(s)
		return r, l
	}

	l, r := o.(changes.Custom).ReverseMerge(s)
	return r, l
}

func (s setattr) mergeSplice(o changes.Splice, reverse bool) (ox, sx changes.Change) {
	return nil, nil
}

func (s setattr) mergeMove(o changes.Move, reverse bool) (ox, sx changes.Change) {
	return nil, nil
}

func (s setattr) mergePath(o changes.PathChange, reverse bool) (ox, sx changes.Change) {
	idx := o.Path[0].(int)
	if idx < s.Offset || idx >= s.Offset+s.Before.count() {
		return o, s
	}
	non, overlap := s.split(idx, idx+1)
	replacement := changes.PathChange{
		Path: []interface{}{idx, s.Name},
		Change: changes.Replace{
			Before: overlap.(setattr).Before[0].Value,
			After:  overlap.(setattr).After[0].Value,
		},
	}
	own := changes.ChangeSet{non, replacement}
	if reverse {
		return own.ReverseMerge(o)
	}
	return own.Merge(o)
}

func (s setattr) mergeSetAttr(o setattr, reverse bool) (ox, sx changes.Change) {
	if s.Name != o.Name {
		return o, s
	}
	left, lx := s.split(o.Offset, o.Offset+o.Before.count())
	right, rx := o.split(s.Offset, s.Offset+s.Before.count())

	if lx == nil && rx == nil {
		return o, s
	}

	if reverse {
		lxx := lx.(setattr)
		lxx.Before = rx.(setattr).After
		return right, (changes.ChangeSet{lxx, left}).Simplify()
	}
	rxx := rx.(setattr)
	rxx.Before = lx.(setattr).After
	return (changes.ChangeSet{rxx, right}).Simplify(), left
}

func (s setattr) split(start, end int) (nonOverlap, overlap changes.Change) {
	send := s.Offset + s.Before.count()
	if s.Offset > start {
		start = s.Offset
	}
	if send < end {
		end = send
	}
	if start >= end {
		return s, nil
	}
	left := s.slice(0, start-s.Offset)
	mid := s.slice(start-s.Offset, end-start)
	right := s.slice(end-s.Offset, send-end)
	return (changes.ChangeSet{left, right}).Simplify(), mid
}

func (s setattr) slice(offset, count int) changes.Change {
	if count <= 0 {
		return nil
	}

	before := s.Before.slice(offset, count)
	after := s.After.slice(offset, count)
	return setattr{s.Offset + offset, s.Name, before, after}
}
