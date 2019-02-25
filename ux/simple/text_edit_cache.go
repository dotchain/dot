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

// Item fetches the item at the specific key
func (c *TextEditCache) Item(key interface{}) *TextEdit {
	return c.current[key]
}

// TextEdit fetches a TextEdit from the cache (updating it)
// or creates a new TextEdit
func (c *TextEditCache) TextEdit(key interface{}, styles core.Styles, text string) *TextEdit {
	if item, ok := c.old[key]; !ok {
		c.current[key] = NewTextEdit(styles, text)
	} else {
		delete(c.old, key)
		item.Update(styles, text)
		c.current[key] = item
	}

	return c.current[key]
}
