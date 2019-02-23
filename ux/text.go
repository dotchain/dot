// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

import "github.com/dotchain/dot/ux/streams"

// TextStream alias
type TextStream = streams.TextStream

// Notifier alias
type Notifier = streams.Notifier

// TextSpan implements a simple text control.
type TextSpan struct {
	// Root is the root dom element of this control
	Root Element

	// private state
	styles Styles
	text   string
}

// NewTextSpan creates a new text control.
func NewTextSpan(styles Styles, text string) *TextSpan {
	root := NewElement(Props{
		Tag:         "span",
		TextContent: text,
		Styles:      styles,
	})
	return &TextSpan{root, styles, text}
}

// Update updates the text or styles of the checkbox.
func (s *TextSpan) Update(styles Styles, text string) {
	if s.text != text {
		s.Root.SetProp("TextContent", text)
		s.text = text
	}
	if s.styles != styles {
		s.styles = styles
		s.Root.SetProp("Styles", styles)
	}
}

// TextEdit implements a simple text edit control.
type TextEdit struct {
	// Root is the root dom element of this control
	Root Element

	// private state
	styles Styles

	// Consumers of TextEdit can get the latest value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on this field.
	Text *TextStream
}

// NewTextEdit creates a new text edit control.
func NewTextEdit(styles Styles, text string) *TextEdit {
	n := &Notifier{}
	t := &TextEdit{nil, styles, &TextStream{n, text, nil, nil}}
	t.Root = NewElement(Props{
		Tag:         "input",
		Type:        "text",
		TextContent: text,
		Styles:      styles,
		OnChange:    &EventHandler{t.onChange},
	})
	return t
}

// Update updates the value and style for the text edit widget.
func (t *TextEdit) Update(styles Styles, text string) {
	if t.Text.Value != text {
		t.Root.SetProp("TextContent", text)
		t.Text = t.Text.Update(nil, text)
	}

	if t.styles != styles {
		t.styles = styles
		t.Root.SetProp("Styles", styles)
	}
}

func (t *TextEdit) onChange(e Event) {
	t.Text = t.Text.Update(nil, t.Root.Value())
	t.Text.Notify()
}
