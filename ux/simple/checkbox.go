// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package simple

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

// Checkbox implements a checkbox control.
type Checkbox struct {
	Element

	// Consumers of Checkbox can get the latest value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on this field.
	Checked *streams.BoolStream

	// persist this
	onChangeHandler core.EventHandler
}

// NewCheckbox creates a new checkbox control.
func NewCheckbox(styles core.Styles, checked bool) *Checkbox {
	c := &Checkbox{}
	c.onChangeHandler = core.EventHandler{c.onChange}
	c.Checked = &streams.BoolStream{&streams.Notifier{}, checked, nil, nil}

	c.render(styles)
	return c
}

// Update updates the control to use the provided styles and checked value
func (c *Checkbox) Update(styles core.Styles, checked bool) {
	if c.Checked.Value != checked {
		c.Checked = c.Checked.Update(nil, checked)
	}
	c.render(styles)
}

func (c *Checkbox) render(styles core.Styles) {
	c.Declare(core.Props{
		Tag:      "input",
		Type:     "checkbox",
		Checked:  c.Checked.Value,
		Styles:   styles,
		OnChange: &c.onChangeHandler,
	})
}

func (c *Checkbox) onChange(e core.Event) {
	c.Checked = c.Checked.Update(nil, c.Root.Value() == "on")
	c.Checked.Notify()
}

// generate CheckboxCache

//go:generate go run ../templates/gen.go ../templates/cache.template Package=simple Base=Checkbox BaseType=Checkbox "Args=styles, checked" "ArgsDef=styles core.Styles, checked bool" Constructor=NewCheckbox out=checkbox_cache.go
