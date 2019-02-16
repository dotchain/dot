// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom

// Checkbox implements a checkbox.
//
// Use NewCheckbox to create a checkbox element. It can then  be
// updated with calls to Update().
type Checkbox struct {
	raw
	styles Styles

	// Consumers of Checkbox can get the current value by
	// inspecting this field.  Changes can be subscribed by
	// calling On on the stream. Calling Update on this field
	// will not update the UI though. Instead call Update on the
	// Checkbox itself.
	Checked *BoolStream
}

// NewCheckbox creates a new checkbox with the provided styles and
// checked value.
func NewCheckbox(styles Styles, checked bool) *Checkbox {
	c := &CheckBox{Styles: styles, Checked: &BoolStream{&Notifier{}, checked, nil, nil}}
	on := func(MouseEvent) {
		c.raw.SetProp("Checked", !c.Checked.Value)
		c.Checked = c.Checked.Update(nil, !c.Checked.Value)
		c.Checked.Notify()
	}

	c.raw = driver.Raw(Props{Checked: checked, Styles: styles, OnClick: on})
	return c
}

// Update updates the value or styles of the checkbox.
func (c *Checkbox) Update(styles Styles, checked bool) {
	if c.Checked.Value != checked {
		c.raw.SetProp("Checked", checked)
		c.Checked = c.Checked.Update(nil, checked)
	}
	if c.Styles != styles {
		c.Styles = styles
		c.raw.SetProp("Styles", styles)
	}
}
