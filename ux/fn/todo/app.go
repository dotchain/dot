// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package todo is an example TODO task list app
package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

// App is a thin wrapper on top of TasksView with checkboxes for ShowDone and ShowUndone
//
// codegen: pure
func App(c *appCtx, styles core.Styles, tasks *TasksStream) core.Element {
	done := getAppStateStream(c, "done", styles, tasks)
	notDone := getAppStateStream(c, "notDone", styles, tasks)

	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		c.fn.Checkbox("done", core.Styles{}, done),
		c.fn.Checkbox("notDone", core.Styles{}, notDone),
		c.TasksView("tasks", core.Styles{}, done, notDone, tasks),
	)
}

func getAppStateStream(c *appCtx, name string, styles core.Styles, tasks *TasksStream) *streams.BoolStream {
	var result *streams.BoolStream
	var handler *streams.Handler
	if f, h, ok := c.Cache.GetSubstream(nil, "done"); ok {
		result, handler = f.(*streams.BoolStream), h
	} else {
		result = streams.NewBoolStream(true)
		handler = &streams.Handler{nil}
		result.On(handler)
	}
	handler.Handle = func() { c.refresh(styles, tasks) }
	result = result.Latest()
	c.Cache.SetSubstream(nil, name, result, handler, func() { result.Off(handler) })
	return result
}

//go:generate  go  run ../codegen.go - $GOFILE
