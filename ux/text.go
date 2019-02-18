// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom

// TextSpan implements a simple text segment.
type TextSpan struct {
	// Element is the dom element associated with this widget
	Element Element

	// styles and text are private state, cached for use in Update()
	styles Styles
	text   string
}

// NewTextSpaan creates a new text span with the provided styles and
// text.
func NewTextSpan(styles Styles, text string) *TextSpan {
	s := &TextSpan{nil, styles, text}
	props := Props{Tag: "span", TextContent: text, Styles: styles}
	s.Element = driver.NewElement(props)
	return s
}

// Update updates the text or styles of the checkbox.
func (s *TextSpan) Update(styles Styles, text string) {
	if s.text != text {
		s.Element.SetProp("TextContent", text)
		s.text = text
	}
	if s.styles != styles {
		s.styles = styles
		s.Element.SetProp("Styles", styles)
	}
}

// TextEdit implements a simple text edit control
type TextEdit struct {
	// Element is the dom element associated with this widget
	Element Element

	// styles is private state, cached for use in Update()
	styles Styles

	// Consumers of Checkbox can get the current value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on this field.
	Text *TextStream
}

// NewTextEdit creates a new text edit widget with the provided styles
// and text value
func NewTextEdit(styles Styles, text string) *TextEdit {
	n := &Notifier{}
	t := &TextEdit{nil, styles, &TextStream{n, text, nil, nil}}
	on := func(Event) {
		t.Text = t.Text.Update(nil, t.Element.Value())
		t.Text.Notify()
	}

	props := Props{Tag: "input", Type: "text", TextContent: text, Styles: styles, OnChange: on}
	t.Element = driver.NewElement(props)
	return t
}

// Update updates the value and style for the text edit widget.
func (t *TextEdit) Update(styles Styles, text string) {
	if t.Text.Value != text {
		t.Element.SetProp("TextContent", text)
		t.Text = t.Text.Update(nil, text)
	}

	if t.styles != styles {
		t.styles = styles
		t.Element.SetProp("Styles", styles)
	}
}
