// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/simple"
	"github.com/dotchain/dot/ux/streams"
)

// Tasks represents a collection of tasks
type Tasks []Task

// generate TasksStream
//go:generate go run ../templates/gen.go ../templates/streams.template Package=todo Base=Tasks BaseType=Tasks out=tasks_stream.go

// TasksView is a control that renders tasks using TaskEdit.
//
// Individual tasks can be modified underneath. The current list of
// tasks is available via Tasks field which supports On/Off to receive
// notifications.
type TasksView struct {
	simple.Element
	Tasks *TasksStream

	edits TaskEditCache
}

// NewTasksView creates a new TasksView
func NewTasksView(styles core.Styles, showDone bool, showNotDone bool, tasks Tasks) *TasksView {
	v := &TasksView{}
	v.Tasks = NewTasksStream(tasks)
	v.Update(styles, showDone, showNotDone, tasks)
	return v
}

// Update updates the TasksView
func (v *TasksView) Update(styles core.Styles, showDone bool, showNotDone bool, tasks Tasks) {
	if !v.areTasksSame(v.Tasks.Value, tasks) {
		v.Tasks = v.Tasks.Update(nil, tasks)
	}

	v.edits.Begin()
	defer v.edits.End()

	v.Declare(
		core.Props{Tag: "div", Styles: styles},
		v.renderTasks(tasks, func(t Task) core.Element {
			if t.Done && !showDone || !t.Done && !showNotDone {
				return nil
			}
			return v.addTaskListener(v.edits.TryGet(t.ID, core.Styles{}, t))
		})...,
	)
}

func (v *TasksView) addTaskListener(edit *TaskEdit, exists bool) core.Element {
	if !exists {
		edit.Task.On(&streams.Handler{v.on})
	}
	return edit.Root
}

func (v *TasksView) on() {
	// TODO: propagate change properly instead of simply recalculating

	result := append(Tasks(nil), v.Tasks.Value...)
	for kk, task := range result {
		if edit := v.edits.Item(task.ID); edit != nil {
			result[kk] = edit.Task.Value
		}
	}
	v.Tasks = v.Tasks.Update(nil, result)
	v.Tasks.Notify()
}

func (v *TasksView) areTasksSame(t1, t2 Tasks) bool {
	if len(t1) != len(t2) {
		return false
	}
	for kk, task := range t1 {
		if task != t2[kk] {
			return false
		}
	}
	return true
}

func (v *TasksView) renderTasks(t Tasks, fn func(Task) core.Element) []core.Element {
	result := make([]core.Element, len(t))
	for kk, elt := range t {
		result[kk] = fn(elt)
	}
	return result
}
