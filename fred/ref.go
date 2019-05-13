// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// NoRef is returned by Eval when a reference id does not exist
var NoRef = Error("No value")

// Ref refers to another object whose id is stored here
type Ref string

// Apply implements changes.Value
func (r Ref) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return r
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, r)
}

// Eval returns the value pointed to by the Ref
func (r Ref) Eval(dir *DirStream) Object {
	if v, ok := dir.Cache[r]; ok {
		return v
	}
	var result Object = NoRef
	if v, ok := dir.Value[string(r)]; ok {
		result = v.Eval(dir)
	}
	dir.Cache[r] = result
	return result
}

// Diff returns the difference between old and new
func (r Ref) Diff(old, next *DirStream, c changes.Change) changes.Change {
	if cx, ok := old.Changes[r]; ok {
		return cx
	}

	after, ok := next.Value[string(r)]
	if !ok {
		after = NoRef
	}

	cx := after.Diff(old, next, c)
	if before, ok := old.Value[string(r)]; !ok || before != after {
		cx = replace(r.Eval(old), r.Eval(next))
	}
	// TODO: is this needed? it may make things a bit more efficient
	// next.Cache[r] = r.Eval(old).Apply(nil, cx).(Object)
	old.Changes[r] = cx
	return cx
}
