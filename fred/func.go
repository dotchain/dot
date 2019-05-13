// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "github.com/dotchain/dot/changes"

// Functor is the interface to be implemented by user defined functions
type Functor interface {
	// Eval should evaluate the function in the context of the dir.
	//
	// The objects provided are unevaluated -- it is up to the
	// implementation to evaluate it if needed.
	//
	// The return value is likely to be cached
	// automatically. dir.Cache can be used to cache interesting
	// values but at this point there is no simple way to create
	// an appropriate cache key.  At some point, a context can be
	// provided to enable such caching
	Eval(dir *DirStream, args []Object) Object

	// Diff evaluates the chagne in value between old and new
	Diff(old, next *DirStream, c changes.Change, args []Object) changes.Change
}

// Func is the placeholder for a function.  The actual implementation
// of functions is in the Functor while this is just a wrapper to
// marshal args values
type Func struct {
	Functor
	Args interface{}
}

// Apply implements changes.Value
func (f Func) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return f
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, f)
}

// Eval returns the calculated value
func (f Func) Eval(dir *DirStream) Object {
	if v, ok := dir.Cache[f]; ok {
		return v
	}
	var result Object = f.Functor.Eval(dir, FromTuple(f.Args))
	dir.Cache[f] = result
	return result
}

// Diff returns the difference between old and new
func (f Func) Diff(old, next *DirStream, c changes.Change) changes.Change {
	if cx, ok := old.Changes[f]; ok {
		return cx
	}

	cx := f.Functor.Diff(old, next, c, FromTuple(f.Args))
	old.Changes[f] = cx
	return cx
}
