// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"github.com/dotchain/dot/changes"
)

// FontWeight is CSS font-weight
type FontWeight int

// FontWeight constants
const (
	FontThin FontWeight = 100 * (iota + 1)
	FontExtraLight
	FontLight
	FontNormal
	FontMedium
	FontSemibold
	FontBold
	FontExtraBold
	FontBlack
)

// Name is the key to use within rich.Attrs
func (f FontWeight) Name() string {
	return "FontWeight"
}

// Apply only accepts one type of change: one that Replace's the
// value.
func (f FontWeight) Apply(ctx changes.Context, c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return f
	case changes.Replace:
		return c.After
	}
	return c.(changes.Custom).ApplyTo(ctx, f)
}
