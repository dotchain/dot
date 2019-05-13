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

// Object is the required interface for all objects in fred
type Object interface {
	changes.Value
	Eval(*DirStream) Object
	Diff(old, next *DirStream, c changes.Change) changes.Change
}

func replace(before, after changes.Value) changes.Change {
	if before == after {
		return nil
	}
	return changes.Replace{Before: before, After: after}
}
