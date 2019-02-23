// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package simple

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

// TextEdit implements a simple inline text edit view control.
type TextEdit struct {
	// Element exposes Root.
	Element

	// Text tracks the current text in the control.
	Text *streams.TextStream

	changeHandler core.EventHandler
}

// NewTextEdit creates a new text edit control
func NewTextEdit(styles core.Styles, text string) *TextEdit {
	v := &TextEdit{}
	v.Text = streams.NewTextStream(text)
	v.changeHandler = core.EventHandler{func(core.Event) {
		v.Text = v.Text.Update(nil, v.Root.Value())
		v.Text.Notify()
	}}
	v.Update(styles, text)
	return v
}

// Update updates the text or styles of the text view control.
func (v *TextEdit) Update(styles core.Styles, text string) {
	v.Declare(core.Props{
		Tag:         "input",
		Type:        "text",
		TextContent: text,
		Styles:      styles,
		OnChange:    &v.changeHandler,
	})
}
