// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom

// Checkbox implements a checkbox.
//
// Use NewCheckbox to create a checkbox element. It can then  be
// updated with calls to Update().
type Checkbox struct {
	// Element is the dom element associated with this widget
	Element Element

	// styles is private state, cached for use in Update()
	styles Styles

	// Consumers of Checkbox can get the current value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on this field.
	Checked *BoolStream
}

// NewCheckbox creates a new checkbox with the provided styles and
// checked value.
func NewCheckbox(styles Styles, checked bool) *Checkbox {
	c := &Checkbox{nil, styles, &BoolStream{&Notifier{}, checked, nil, nil}}
	on := func(MouseEvent) {
		c.Element.SetProp("Checked", !c.Checked.Value)
		c.Checked = c.Checked.Update(nil, !c.Checked.Value)
		c.Checked.Notify()
	}

	props := Props{Checked: checked, Styles: styles, OnClick: on}
	c.Element = driver.NewElement(props, nil)
	return c
}

// Update updates the value or styles of the checkbox.
func (c *Checkbox) Update(styles Styles, checked bool) {
	if c.Checked.Value != checked {
		c.Element.SetProp("Checked", checked)
		c.Checked = c.Checked.Update(nil, checked)
	}
	if c.styles != styles {
		c.styles = styles
		c.Element.SetProp("Styles", styles)
	}
}
