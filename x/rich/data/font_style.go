// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/dotchain/dot/changes"
)

// FontStyle is CSS font-style
type FontStyle string

// FontStyle values
const (
	FontStyleNormal  FontStyle = "normal"
	FontStyleItalic  FontStyle = "italic"
	FontStyleOblique FontStyle = "oblique"
)

// Name is the key to use within rich.Attrs
func (f FontStyle) Name() string {
	return "FontStyle"
}

// Apply only accepts one type of change: one that Replace's the
// value.
func (f FontStyle) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return f
	case changes.Replace:
		return c.After
	}
	return c.(changes.Custom).ApplyTo(ctx, f)
}
