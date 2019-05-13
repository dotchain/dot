// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package fred implements a convergent functional reactive engine.
//
// It uses a pull-based functional reactive system (instead of push
// based setup) to make it easy to only rebuild the parts of the graph
// needed.
//
// The main entrypoint is Dir which maintains a directory of object
// defintions.  Objects are indexed by an ID (string) and values can
// be fixed or derived.  Fixed values evaluate to themselves while
// derived values evaluate to a different result.
package fred

import (
	"github.com/dotchain/dot/changes"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return e
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, e)
}

func (e Error) Eval(_ *DirStream) Object {
	return e
}

func (e Error) Next(old, next *DirStream, c changes.Change) (Object, changes.Change) {
	return e, nil
}
