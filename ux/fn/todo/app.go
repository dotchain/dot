// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package todo is an example TODO task list app
package todo

import "github.com/dotchain/dot/ux/core"

// App is a thin wrapper on top of TasksView with checkboxes for ShowDone and ShowUndone
//
// codegen: pure
func App(c *appCtx, styles core.Styles, tasks *TasksStream) core.Element {
	showDone := c.newBoolStream("done", true).Latest()
	showNotDone := c.newBoolStream("notDone", true).Latest()
	refresh := func() {
		c.refresh(styles, tasks)
	}
	c.On(showDone.Notifier, refresh)
	c.On(showNotDone.Notifier, refresh)

	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		c.fn.Checkbox("done", core.Styles{}, showDone),
		c.fn.Checkbox("notDone", core.Styles{}, showNotDone),
		c.TasksView("tasks", core.Styles{}, showDone, showNotDone, tasks),
	)
}

//go:generate  go  run ../codegen.go - $GOFILE
