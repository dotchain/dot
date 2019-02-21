// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux_test

import "github.com/dotchain/dot/ux"

// TodoTaskCache can be generated from a template.
type TodoTaskCache struct {
	old, current map[string]*TodoTask
}

func (c *TodoTaskCache) Reset() {
	c.old = c.current
	c.current = map[string]*TodoTask{}
}

func (c *TodoTaskCache) Cleanup() {
	// if TodoTask had a Close() method all the old left-over items
	// can be cleaned up via that call
	c.old = nil
}

func (c *TodoTaskCache) Get(id string, styles ux.Styles, data TaskData) (*TodoTask, bool) {
	exists := false
	if t, ok := c.old[id]; !ok {
		c.current[id] = NewTodoTask(styles, data)
	} else {
		delete(c.old, id)
		t.Update(styles, data)
		c.current[id] = t
		exists = true
	}

	return c.current[id], exists
}
