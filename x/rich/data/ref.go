// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import "github.com/dotchain/dot/changes"

// Ref represents a refereence to an embedded object storeed
// in a Dir.  The Dir must be an ancestor of the Ref
type Ref struct {
	ID interface{}
}

// Name returns the key name for use with Attrs
func (r *Ref) Name() string {
	return "Embed"
}

// Apply implements channges.Value
func (r *Ref) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return r
	case changes.Replace:
		return c.After
	}
	return c.(changes.Custom).ApplyTo(ctx, r)
}
