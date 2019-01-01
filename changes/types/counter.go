// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import "github.com/dotchain/dot/changes"

// Counter implements a 32-bit counter that can be incremented,
// decremented or replaced. Counter implements the changes.Value
// interface.
type Counter int32

// Slice implements changes.Value.Slice but it is not expected to ever
// by used for Counters
func (c Counter) Slice(offset, count int) changes.Collection {
	panic("Slice call not expected on counter")
}

// Count always returns 1 for counters
func (c Counter) Count() int {
	if c != 0 {
		return 1
	}
	return 0
}

// ApplyCollection implements Collection interface
func (c Counter) ApplyCollection(ctx changes.Context, cx changes.Change) changes.Collection {
	switch cx := cx.(type) {
	case changes.Splice:
		after := cx.After.(Counter)
		return c + after
	}
	panic("Unexpected change on Apply")
}

// Apply only supports Replace and Inserts
func (c Counter) Apply(ctx changes.Context, cx changes.Change) changes.Value {
	switch cx := cx.(type) {
	case nil:
		return c
	case changes.Replace:
		if cx.IsDelete() {
			return changes.Nil
		}
		return cx.After
	case changes.Custom:
		return cx.ApplyTo(ctx, c)
	}
	return c.ApplyCollection(ctx, cx)
}

// Increment returns a change which implements the increment
// operation.
func (c Counter) Increment(by int32) changes.Change {
	return changes.Splice{0, Counter(0), Counter(by)}
}

// Set returns a change which implements updating the value
func (c Counter) Set(v int32) changes.Change {
	return changes.Replace{Before: c, After: Counter(v)}
}
