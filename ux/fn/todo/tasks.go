// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

// Tasks represents a collection of tasks
type Tasks []Task

// generate TasksStream
//go:generate go run ../../templates/gen.go ../../templates/streams.template Package=todo Base=Tasks BaseType=Tasks out=tasks_stream.go

// TasksView is a control that renders tasks using TaskEdit.
//
// Individual tasks can be modified underneath. The current list of
// tasks is available via Tasks field which supports On/Off to receive
// notifications.
//
// codegen: pure
func TasksView(c *tasksViewCtx, styles core.Styles, showDone *streams.BoolStream, showNotDone *streams.BoolStream, tasks *TasksStream) core.Element {
	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		renderTasks(tasks.Value, func(index int, t Task) core.Element {
			if t.Done && !showDone.Value || !t.Done && !showNotDone.Value {
				return nil
			}

			return c.TaskEdit(t.ID, core.Styles{}, tasks.Substream(c.Cache, index))
		})...,
	)
}

func renderTasks(t Tasks, fn func(int, Task) core.Element) []core.Element {
	result := make([]core.Element, len(t))
	for kk, elt := range t {
		result[kk] = fn(kk, elt)
	}
	return result
}

// generate the TasksViewCache for any consumers who want it
//go:generate go run ../codegen.go - $GOFILE
