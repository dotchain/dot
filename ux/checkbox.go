// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

// Checkbox implements a checkbox control.
type Checkbox struct {
	// Root is the root dom element of this control
	Root Element

	// private state
	styles Styles

	// Consumers of Checkbox can get the latest value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on this field.
	Checked *BoolStream
}

// NewCheckbox creates a new checkbox control.
func NewCheckbox(styles Styles, checked bool) *Checkbox {
	c := &Checkbox{nil, styles, &BoolStream{&Notifier{}, checked, nil, nil}}
	c.Root = NewElement(Props{
		Tag:      "input",
		Type:     "checkbox",
		Checked:  checked,
		Styles:   styles,
		OnChange: c.onChange,
	})
	return c
}

// Update updates the value and styles of the checkbox.
func (c *Checkbox) Update(styles Styles, checked bool) {
	if c.Checked.Value != checked {
		c.Root.SetProp("Checked", checked)
		c.Checked = c.Checked.Update(nil, checked)
	}
	if c.styles != styles {
		c.styles = styles
		c.Root.SetProp("Styles", styles)
	}
}

func (c *Checkbox) onChange(e Event) {
	c.Checked = c.Checked.Update(nil, c.Root.Value() == "on")
	c.Checked.Notify()
}
