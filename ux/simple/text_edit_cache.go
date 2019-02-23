// This file is generated by:
//    github.com/dotchain/dot/ux/templates/cache.template
//
// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package simple

import "github.com/dotchain/dot/ux/core"

// TextEditCache holds a cache of TextEdit controls.
//
// Controls that have manage a bunch of TextEdit controls
// should maintain a cache created like so:
//
//     cache := &TextEditCache{}
//
// When updating, the cache can be used to reuse controls:
//
//     cache.Begin()
//     defer cache.End()
//
//     ... for each TextEdit control needed do:
//     cache.Get(key, styles, text)
//
// This allows the cache to reuse the control if the key exists.
// Otherwise a new control is created via NewTextEdit(styles, text)
//
// When a control is reused, it is also automatically updated.
type TextEditCache struct {
	old, current map[interface{}]*TextEdit
}

// Begin should be called before the start of a round
func (c *TextEditCache) Begin() {
	c.old = c.current
	c.current = map[interface{}]*TextEdit{}
}

// End should be called at the end of a round
func (c *TextEditCache) End() {
	// if components had a Close() method all the old left-over items
	// can be cleaned up via that call
	c.old = nil
}

// TryGet fetches a TextEdit from the cache (updating it)
// or creates a new TextEdit
//
// It returns the TextEdit but also whether the control existed.
// This can be used to conditionally setup listeners.
func (c *TextEditCache) TryGet(key interface{}, styles core.Styles, text string) (*TextEdit, bool) {
	exists := false
	if item, ok := c.old[key]; !ok {
		c.current[key] = NewTextEdit(styles, text)
	} else {
		delete(c.old, key)
		item.Update(styles, text)
		c.current[key] = item
		exists = true
	}

	return c.current[key], exists
}

// Item fetches the item at the specific key
func (c *TextEditCache) Item(key interface{}) *TextEdit {
	return c.current[key]
}

// Get fetches a TextEdit from the cache (updating it)
// or creates a new TextEdit
//
// Use TryGet to also fetch whether the control from last round was reused
func (c *TextEditCache) Get(key interface{}, styles core.Styles, text string) *TextEdit {
	v, _ := c.TryGet(key, styles, text)
	return v
}
