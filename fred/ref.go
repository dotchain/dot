// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// NoRef is returned by Eval when a reference id does not exist
var NoRef = Error("No value")

// Ref refers to another object whose id is stored here
type Ref string

func (r Ref) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return r
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, r)
}

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

func (r Ref) Next(old, next *DirStream, c changes.Change) (Object, changes.Change) {
	cx, ok := old.Changes[r]

	if !ok {
		var v Object

		after, ok := next.Value[string(r)]
		if !ok {
			after = NoRef
		}
		v, cx = after.Next(old, next, c)
		next.Cache[r] = v.Eval(next)
		if before, ok := old.Value[string(r)]; !ok || before != after {
			cx = replace(r.Eval(old), next.Cache[r])
		}
		old.Changes[r] = cx
	}

	return r, cx
}
