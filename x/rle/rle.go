// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rle

import "github.com/dotchain/dot/changes"

// Run represents a virtual array of the same value
type Run struct {
	changes.Value
	Count int
}

// The A type represents a slice of arbitrary values but uses a
// more compact encoding to remove duplicate types.  The values
// provided to it should either be comparable or implement IsEqual
type A []Run

// Slice implements changes.Value.Slice
func (a A) Slice(offset, count int) changes.Collection {
	return a.slice(offset, count)
}

func (a A) slice(offset, count int) A {
	result := A{}
	seen := 0
	for _, run := range a {
		start, end := seen, seen+run.Count
		seen = end
		if offset > start {
			start = offset
		}
		if offset+count < end {
			end = offset + count
		}
		if end > start {
			result = append(result, Run{run.Value, end - start})
		}
	}
	if offset+count > seen {
		panic("invalid index")
	}
	return result
}

// Count returns size of the array
func (a A) Count() int {
	seen := 0
	for _, run := range a {
		seen += run.Count
	}
	return seen
}

// IsEqual tests if the two encodings match
func (a A) IsEqual(o changes.Value) bool {
	if ox, ok := o.(A); ok {
		if len(a) != len(ox) {
			return false
		}

		for kk, v := range a {
			if v.Count != ox[kk].Count || !isEqual(v.Value, ox[kk].Value) {
				return false
			}
		}
		return true
	}
	return false
}

// ApplyCollection implements collection interface
func (a A) ApplyCollection(c changes.Change) changes.Collection {
	switch c := c.(type) {
	case changes.Splice:
		remove := c.Before.Count()
		right := a.Count() - c.Offset - remove
		return a.slice(0, c.Offset).
			append(c.After.(A)).
			append(a.slice(c.Offset+remove, right))
	case changes.Move:
		c = c.Normalize()
		ox, cx, dx := c.Offset, c.Count, c.Distance
		slice1 := a.slice(0, ox)
		slice2 := a.slice(ox, cx)
		slice3 := a.slice(ox+cx, dx)
		slice4 := a.slice(ox+cx+dx, a.Count()-ox-cx-dx)
		return slice1.append(slice3).append(slice2).append(slice4)
	case changes.PathChange:
		idx := c.Path[0].(int)
		left, right, v := a.split(idx)
		if v == nil {
			v = changes.Nil
		}
		v = v.Apply(changes.PathChange{c.Path[1:], c.Change})
		if v == changes.Nil {
			v = nil
		}
		return left.append(A{{v, 1}}).append(right)
	}
	panic("Unexpected change")
}

// Apply applies the change and returns the updated value
//
// Note: deleting an element via changes.Replace simply replaces it
// with nil.  It does not actually remove the element -- that needs a
// changes.Splice.
func (a A) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return a
	case changes.Replace:
		if c.IsDelete() {
			return changes.Nil
		}
		return c.After
	case changes.PathChange:
		if len(c.Path) == 0 {
			return a.Apply(c.Change)
		}
	case changes.Custom:
		return c.ApplyTo(a)
	}
	return a.ApplyCollection(c)
}

func (a A) append(o A) A {
	if len(a) == 0 {
		return o
	}
	if len(o) == 0 {
		return a
	}

	clone := append(A(nil), a...)
	if isEqual(a[len(a)-1].Value, o[0].Value) {
		clone[len(a)-1].Count += o[0].Count
		o = o[1:]
	}
	return append(clone, o...)
}

func (a A) split(idx int) (left, right A, v changes.Value) {
	seen := 0
	for kk, elt := range a {
		if seen <= idx && seen+elt.Count > idx {
			left := a[:kk:kk]
			right := a[kk+1:]
			if seen < idx {
				left = left.append(A{{elt.Value, idx - seen}})
			}
			if idx+1 < seen+elt.Count {
				run := Run{elt.Value, seen + elt.Count - idx - 1}
				right = (A{run}).append(right)
			}
			return left, right, elt.Value
		}
		seen += elt.Count
	}
	panic("invalid index")
}

// Encode takes a regular value array and returns a run-length encoded
// version of it.
func Encode(v []changes.Value) A {
	var last changes.Value
	result := A{}
	for _, elt := range v {
		if len(result) > 0 && isEqual(last, elt) {
			result[len(result)-1].Count++
		} else {
			result = append(result, Run{elt, 1})
		}
		last = elt
	}
	return result
}

func isEqual(v1, v2 changes.Value) bool {
	if ee, ok := v1.(IsEqual); ok {
		return ee.IsEqual(v2)
	}
	return v1 == v2
}

// IsEqual is the interface that values stored within an encoded array
// can implement. This is useful for deduping complex objects for run
// length encoding.
type IsEqual interface {
	IsEqual(o changes.Value) bool
}
