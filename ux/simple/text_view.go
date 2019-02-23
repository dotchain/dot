// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package simple

import "github.com/dotchain/dot/ux/core"

// TextView implements a simple inline text view control.
type TextView struct {
	// Element exposes Root
	Element
}

// NewTextView creates a new text view control.
func NewTextView(styles core.Styles, text string) *TextView {
	v := &TextView{}
	v.Update(styles, text)
	return v
}

// Update updates the text or styles of the text view control.
func (v *TextView) Update(styles core.Styles, text string) {
	v.Declare(core.Props{
		Tag:         "span",
		TextContent: text,
		Styles:      styles,
	})
}
