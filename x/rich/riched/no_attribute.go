// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched

import "github.com/dotchain/dot/changes"

// NoAttribute specified that a specific attribute should be
// nullified. This is used with Editor.SetAttribute to specify that
// the override is actively removing an attribute
//
// The name of the attribute is the underlying string.
type NoAttribute string

// Name just returns the underlying string
func (n NoAttribute) Name() string {
	return string(n)
}

// Apply implements changes.Value
func (n NoAttribute) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return n
	case changes.Replace:
		return c.After
	}
	return c.(changes.Custom).ApplyTo(ctx, n)
}
