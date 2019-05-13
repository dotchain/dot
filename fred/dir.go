// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

// Dir the main fred container which holds all objects.
type Dir map[string]Object

func (d Dir) set(key interface{}, value changes.Value) changes.Value {
	clone := Dir{}
	for k, v := range d {
		if k != key {
			clone[k] = v
		}
	}

	if value != changes.Nil {
		clone[key.(string)] = value.(Object)
	}
	return clone
}

func (d Dir) get(key interface{}) changes.Value {
	if v, ok := d[key.(string)]; ok {
		return v
	}
	return changes.Nil
}

// Apply implements the changes.Value interface
func (d Dir) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: d.set, Get: d.get}).Apply(ctx, c, d)
}

// NewDirStream returns a new dir stream. The second arg is optional
func NewDirStream(dir Dir, s streams.Stream) *DirStream {
	if s == nil {
		s = streams.New()
	}
	return &DirStream{
		Value:   dir,
		Stream:  s,
		Cache:   map[interface{}]Object{},
		Changes: map[interface{}]changes.Change{},
	}
}

// DirStream implements a stream of dir versions
type DirStream struct {
	Value  Dir
	Stream streams.Stream

	// the following are really performance optimizations
	next    *DirStream
	Cache   map[interface{}]Object
	Changes map[interface{}]changes.Change
}

// Next returns the next dir stream. Any evaluated ObjectStream values
// should be advanced via their own Next calls. If they are no longer
// referred to, none of their antecedents will be computed
func (s *DirStream) Next() (*DirStream, changes.Change) {
	if s.Stream != nil {
		next, c := s.Stream.Next()
		if next != nil && s.next == nil {
			s.next = NewDirStream(s.Value.Apply(nil, c).(Dir), next)
		}
		return s.next, c
	}
	return nil, nil
}

// Eval evaluates an object with possible references and calculations
// within the directory.  All the sub-computations are cached and
// future calls to Next() on the returned object stream will use the
// cached values for incremental computation.  The cache is located on
// the DirStream instance and as such the cache itself is not
// automatically carried over to the next instance to avoid leaking
// memory.
//
// Note that the Object being evaluated could be a Func expression or
// a Ref.  The returned value is guaranteed to not have any references.
func (s *DirStream) Eval(obj Object) (Object, *ObjectStream) {
	result := &ObjectStream{s, obj}
	return result.Eval(), result
}

// ObjectStream is a holder for the evaluations from Eval.
type ObjectStream struct {
	*DirStream
	Object
}

// Next iterates through the object stream yielding the changes. The
// actual value can be recalculated (most likely simple pulled from
// cache) via Eval.
func (s *ObjectStream) Next() (*ObjectStream, changes.Change) {
	if next, c := s.DirStream.Next(); next != nil {
		cx := s.Object.Diff(s.DirStream, next, c)
		return &ObjectStream{next, s.Object}, cx
	}
	return nil, nil
}

// Eval evaluates the object value.
func (s *ObjectStream) Eval() Object {
	return s.Object.Eval(s.DirStream)
}
