// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

import "github.com/dotchain/dot/ux/core"

// Checkbox implements a checkbox control.
type Checkbox struct {
	// Component is embedded. Root element is exposed through this.
	Component

	// persist this
	onChangeHandler core.EventHandler

	// Consumers of Checkbox can get the latest value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on this field.
	Checked *BoolStream
}

// NewCheckbox creates a new checkbox control.
func NewCheckbox(styles Styles, checked bool) *Checkbox {
	c := &Checkbox{}
	c.onChangeHandler = core.EventHandler{c.onChange}
	c.Checked = &BoolStream{&Notifier{}, checked, nil, nil}

	c.render(styles)
	return c
}

// Update updates the control to use the provided styles and checked value
func (c *Checkbox) Update(styles Styles, checked bool) {
	if c.Checked.Value != checked {
		c.Checked = c.Checked.Update(nil, checked)
	}
	c.render(styles)
}

func (c *Checkbox) render(styles Styles) {
	c.Declare(core.Props{
		Tag:      "input",
		Type:     "checkbox",
		Checked:  c.Checked.Value,
		Styles:   styles,
		OnChange: &c.onChangeHandler,
	})
}

func (c *Checkbox) onChange(e Event) {
	c.Checked = c.Checked.Update(nil, c.Root.Value() == "on")
	c.Checked.Notify()
}
