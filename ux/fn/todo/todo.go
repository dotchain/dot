// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package todo demonstrates a simple todo mvc app
package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

// Task represents an item in the TODO list.
type Task struct {
	ID          string
	Done        bool
	Description string
}

// Tasks represents a collection of tasks
type Tasks []Task

// TaskEdit is a control that displays a task as well as allowing it
// to be edited. The current value of the data is available in the
// Task field (which is a stream and so supports On/Off methods).
func TaskEdit(c *taskEditCtx, styles core.Styles, task *TaskStream) core.Element {
	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		c.fn.Checkbox("cb", core.Styles{}, task.DoneSubstream(c.Cache)),
		c.fn.TextEdit("textedit", core.Styles{}, task.DescriptionSubstream(c.Cache)),
	)
}

// TasksView is a control that renders tasks using TaskEdit.
//
// Individual tasks can be modified underneath. The current list of
// tasks is available via Tasks field which supports On/Off to receive
// notifications.
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

// AppWithState is a thin wrapper on top of TasksView with checkboxes for ShowDone and ShowUndone
//
func AppWithState(c *appCtx, styles core.Styles, tasks *TasksStream, done *streams.BoolStream, notDone *streams.BoolStream) core.Element {
	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		c.fn.Checkbox("done", core.Styles{}, done),
		c.fn.Checkbox("notDone", core.Styles{}, notDone),
		c.TasksView("tasks", core.Styles{}, done, notDone, tasks),
	)
}
