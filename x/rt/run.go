// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rt

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
)

// Run implements a custom change type which applies a provided
// inner change to a range of items in an array.  This is particularly
// useful for rich text operations
type Run struct {
	Offset, Count int
	changes.Change
}

// ApplyTo just converts the method into a set of path changes
func (r Run) ApplyTo(v changes.Value) changes.Value {
	for kk := r.Offset; kk < r.Offset+r.Count; kk++ {
		v = v.Apply(changes.PathChange{[]interface{}{kk}, r.Change})
	}
	return v
}

// Revert undoes the effect of the Run
func (r Run) Revert() changes.Change {
	if r.Change == nil {
		return nil
	}
	return Run{r.Offset, r.Count, r.Change.Revert()}
}

// Merge implements the main merge routing for a change
func (r Run) Merge(o changes.Change) (changes.Change, changes.Change) {
	return r.merge(o, false)
}

// ReverseMerge is like Merge except the args and receiver are
// inverted. Basically if someone calls "ch.Merge(r)" and ch does not
// know how to implement merge with r, it calls r.ReverseMerge(ch).
func (r Run) ReverseMerge(o changes.Change) (changes.Change, changes.Change) {
	return r.merge(o, true)
}

// MergePath implements the method needed to work with refs.Merge
func (r Run) MergePath(p []interface{}) *refs.MergeResult {
	idx := p[0].(int)
	if idx < r.Offset || idx >= r.Offset+r.Count {
		return &refs.MergeResult{P: p, Unaffected: r}
	}

	return refs.Merge(p[1:], r.Change).Prefix(p[:1])
}

func (r Run) merge(o changes.Change, reverse bool) (changes.Change, changes.Change) {
	if r.Change == nil {
		return o, nil
	}

	switch o := o.(type) {
	case nil:
		return nil, r
	case changes.Replace:
		o.Before = o.Before.Apply(r)
		return o, nil
	case changes.Splice:
		return r.mergeSplice(o)
	case changes.Move:
		return r.mergeMove(o)
	case Run:
		return r.mergeRun(o, reverse)
	case changes.PathChange:
		return r.mergePathChange(o, reverse)
	}

	if reverse {
		return swap(o.Merge(r))
	}
	return swap(o.(revMerge).ReverseMerge(r))
}

func (r Run) mergeSplice(o changes.Splice) (changes.Change, changes.Change) {
	oEnd := o.Offset + o.Before.Count()
	switch {
	case r.Offset >= oEnd:
		r.Offset += o.After.Count() - o.Before.Count()
	case r.Offset+r.Count <= o.Offset:
	case r.Offset >= o.Offset && r.Offset+r.Count <= oEnd:
		r.Offset -= o.Offset
		o.Before = o.Before.Apply(r)
		return o, nil
	case r.Offset <= o.Offset && r.Offset+r.Count >= oEnd:
		o.Before = o.Before.Apply(Run{0, o.Before.Count(), r.Change})
		left := Run{r.Offset, o.Offset - r.Offset, r.Change}
		right := Run{o.Offset + o.After.Count(), r.Offset + r.Count - oEnd, r.Change}
		return o, changes.ChangeSet{left, right}
	case r.Offset < o.Offset && o.Offset < r.Offset+r.Count:
		o.Before = o.Before.Apply(Run{0, r.Offset + r.Count - o.Offset, r.Change})
		r.Count = o.Offset - r.Offset
	case r.Offset > o.Offset && r.Offset < oEnd:
		o.Before = o.Before.Apply(Run{r.Offset - o.Offset, oEnd - r.Offset, r.Change})
		r.Count = r.Offset + r.Count - oEnd
		r.Offset = o.Offset + o.After.Count()
	}
	return o, r
}

func (r Run) mergeMove(o changes.Move) (changes.Change, changes.Change) {
	rEnd, oEnd := r.Offset+r.Count, o.Offset+o.Count
	oDest := oEnd + o.Distance
	if o.Distance < 0 {
		oDest = o.Offset + o.Distance
	}
	switch {
	case rEnd <= o.Offset && rEnd <= oDest:
	case r.Offset >= oDest && rEnd <= o.Offset:
		r.Offset += o.Count
	case r.Offset >= o.Offset && rEnd <= oEnd:
		r.Offset += o.Distance
	case r.Offset >= oEnd && rEnd <= oDest:
		r.Offset -= o.Count
	case r.Offset >= oEnd && rEnd >= oDest:
	default:
		return r.split3(oDest, o)
	}
	return o, r
}

func (r Run) mergeRun(o Run, reverse bool) (changes.Change, changes.Change) {
	rEnd, oEnd := r.Offset+r.Count, o.Offset+o.Count
	switch {
	case rEnd <= o.Offset || oEnd <= r.Offset:
		return o, r
	case r.Offset == o.Offset && rEnd == oEnd:
		var ox, rx changes.Change
		if reverse && o.Change != nil {
			rx, ox = o.Change.Merge(r.Change)
		} else {
			ox, rx = r.Change.Merge(o.Change)
		}
		return Run{o.Offset, o.Count, ox}, Run{r.Offset, r.Count, rx}
	}
	left := r.splitRuns([]changes.Change{r}, o.Offset)
	left = r.splitRuns(left, o.Offset+o.Count)
	right := r.splitRuns([]changes.Change{o}, r.Offset)
	right = r.splitRuns(right, r.Offset+r.Count)
	lx := changes.ChangeSet(left)
	rx := changes.ChangeSet(right)

	if reverse {
		x, y := rx.Merge(lx)
		return y, x
	}
	return lx.Merge(rx)
}

func (r Run) mergePathChange(o changes.PathChange, reverse bool) (changes.Change, changes.Change) {
	if len(o.Path) == 0 {
		return r.merge(o.Change, reverse)
	}
	idx := o.Path[0].(int)
	if idx < r.Offset || idx >= r.Offset+r.Count {
		return o, r
	}
	var left, right changes.Change
	if idx > r.Offset {
		left = Run{r.Offset, idx - r.Offset, r.Change}
	}
	if idx+1 < r.Offset+r.Count {
		right = Run{idx + 1, r.Offset + r.Count - idx - 1, r.Change}
	}
	other := changes.PathChange{o.Path[1:], o.Change}
	var ox, mid changes.Change
	if reverse {
		mid, ox = other.Merge(r.Change)
	} else {
		ox, mid = r.Change.Merge(other)
	}
	ox = changes.PathChange{o.Path[:1], ox}
	mid = changes.PathChange{o.Path[:1], mid}
	return ox, changes.ChangeSet{left, mid, right}
}

func (r Run) splitRuns(runs []changes.Change, idx int) []changes.Change {
	result := make([]changes.Change, 0, len(runs)+1)
	for _, rx := range runs {
		run := rx.(Run)
		if idx > run.Offset && idx < run.Offset+run.Count {
			left := Run{run.Offset, idx - run.Offset, run.Change}
			right := Run{idx, run.Count + run.Offset - idx, run.Change}
			result = append(result, left, right)
		} else {
			result = append(result, run)
		}
	}
	return result
}

func (r Run) split3(dest int, o changes.Move) (changes.Change, changes.Change) {
	c := r.splitRuns([]changes.Change{r}, dest)
	c = r.splitRuns(c, o.Offset)
	c = r.splitRuns(c, o.Offset+o.Count)
	return changes.ChangeSet(c).Merge(o)
}

type revMerge interface {
	ReverseMerge(changes.Change) (changes.Change, changes.Change)
}

func swap(x, y changes.Change) (xx, yy changes.Change) {
	return y, x
}
