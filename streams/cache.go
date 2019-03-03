// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

type cacheEntry struct {
	stream interface{}
	h      *Handler
	close  func()
}

// Cache is used to manage a set of streams
type Cache struct {
	old, current map[interface{}]*cacheEntry
}

// Begin starts a round of using the cache. Any items that were
// present in the cache last time are available. If they are not
// reused, they are closed when End iss called
func (c *Cache) Begin() {
	c.old, c.current = c.current, map[interface{}]*cacheEntry{}
}

// End calls close on any items in the cache that were not reused.
func (c *Cache) End() {
	for _, v := range c.old {
		v.close()
	}
	c.old = nil
}

// GetSubstream returns an entry from the old cache if it exists.
func (c *Cache) GetSubstream(n *Notifier, key interface{}) (interface{}, *Handler, bool) {
	key = [2]interface{}{n, key}
	if v, ok := c.old[key]; ok {
		return v.stream, v.h, true
	}
	return nil, nil, false
}

// SetSubstream updates the entry in the cache
func (c *Cache) SetSubstream(n *Notifier, key, v interface{}, h *Handler, close func()) {
	key = [2]interface{}{n, key}
	c.current[key] = &cacheEntry{v, h, close}
	delete(c.old, key)
}
