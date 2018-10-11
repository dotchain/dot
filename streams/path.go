// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
)

// ChildOf returns a stream that only contains changes specific to the
// path provided. For instance, if the base stream is for an array like
// object, ChildOf(base, 5) refers to the changes applicable to the
// 5th element in the array.
func ChildOf(s Stream, keys ...interface{}) Stream {
	toChild := func(r *refs.MergeResult) changes.Change {
		return r.Scoped
	}
	toParent := func(c changes.Change, p []interface{}) changes.Change {
		return changes.PathChange{p, c}
	}
	return &xform{s, refs.Path(keys), toChild, toParent}
}

// FilterPath returns a stream which is focused only on the provided
// path. All changes that do not affect this path are ignored.  Unlike
// ChildOf, the changes themselves are not transformed, just filtered.
//
// FilterPath can cause unexpected results if used with array indices
// in the path.  For example:
//
//    derived := FilterPath(base, 5)
//    base.Append(changes.Splice{0, types.S8(""), types.S8("abc")})
//
// The splice above would effectively be filtered out because it does
// not affect the provided path. But without this splice, all other
// affecting the 5th element (now the 8th element) cannot be applied
// without transformations.
func FilterPath(s Stream, keys ...interface{}) Stream {
	toChild := func(r *refs.MergeResult) changes.Change {
		return r.Affected
	}
	toParent := func(c changes.Change, p []interface{}) changes.Change {
		return c
	}
	return &xform{s, refs.Path(keys), toChild, toParent}
}

// FilterOutPath returns a stream which is focused only on the provided
// path. All changes that do not affect this path are returned but paths
// that affect the provided path are ignored.  Unlike
// ChildOf, the changes themselves are not transformed, just filtered.
func FilterOutPath(s Stream, keys ...interface{}) Stream {
	toChild := func(r *refs.MergeResult) changes.Change {
		return r.Unaffected
	}
	toParent := func(c changes.Change, p []interface{}) changes.Change {
		return c
	}
	return &xform{s, refs.Path(keys), toChild, toParent}
}

// xform transforms changes from the parent stream to child stream and
// vice-versa. toChild converts a change in the parent stream to that
// in the child.  toParent converts a change in the child stream to
// that in the parent.
type xform struct {
	Stream
	Path     []interface{}
	toChild  func(r *refs.MergeResult) changes.Change
	toParent func(c changes.Change, p []interface{}) changes.Change
}

func (x *xform) clone(s Stream, p []interface{}) *xform {
	return &xform{s, p, x.toChild, x.toParent}
}

func (x *xform) Append(c changes.Change) Stream {
	return x.clone(x.Stream.Append(x.toParent(c, x.Path)), x.Path)
}

func (x *xform) ReverseAppend(c changes.Change) Stream {
	return x.clone(x.Stream.ReverseAppend(x.toParent(c, x.Path)), x.Path)
}

func (x *xform) Next() (changes.Change, Stream) {
	c, s := x.Stream.Next()
	r := refs.Merge(x.Path, c)
	if r == nil || s == nil {
		return nil, nil
	}

	return x.toChild(r), x.clone(s, r.P)
}

func (x *xform) Nextf(key interface{}, fn func()) {
	if fn != nil {
		last, prev := x, fn
		fn = func() {
			_, l := last.Next()
			if l != nil {
				last = l.(*xform)
				prev()
			}
		}
	}
	x.Stream.Nextf(key, fn)
}

func (x *xform) Scheduler() Scheduler {
	return x.Stream.Scheduler()
}

func (x *xform) WithScheduler(sch Scheduler) Stream {
	return x.clone(x.Stream.WithScheduler(sch), x.Path)
}
