// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux"
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
	Root ux.Element

	styles ux.Styles
	cache  TaskEditCache

	Tasks *TasksStream
}

// NewTasksView creates a new TasksView
func NewTasksView(styles ux.Styles, showDone bool, showNotDone bool, tasks Tasks) *TasksView {
	view := &TasksView{
		ux.NewElement(ux.Props{Tag: "div", Styles: styles}),
		styles,
		TaskEditCache{},
		NewTasksStream(nil),
	}
	view.Update(styles, showDone, showNotDone, tasks)
	return view
}

// Update updates the TasksView
func (view *TasksView) Update(styles ux.Styles, showDone bool, showNotDone bool, tasks Tasks) {
	if view.styles != styles {
		view.styles = styles
		view.Root.SetProp("Styles", styles)
	}

	view.cache.Reset()
	defer view.cache.Cleanup()

	if !view.areTasksSame(view.Tasks.Value, tasks) {
		view.Tasks = view.Tasks.Update(nil, tasks)
	}

	children := []ux.Element{}
	for _, task := range tasks {
		if task.Done && !showDone || !task.Done && !showNotDone {
			continue
		}

		taskEdit, exists := view.cache.Get(task.ID, ux.Styles{}, task)
		if !exists {
			taskEdit.Task.On(&streams.Handler{view.on})
		}
		children = append(children, taskEdit.Root)
	}
	ux.UpdateChildren(view.Root, children)
}

func (view *TasksView) on() {
	// TODO: propagate change properly instead of simply recalculating

	result := append(Tasks(nil), view.Tasks.Value...)
	for kk, task := range result {
		if edit, ok := view.cache.current[task.ID]; ok {
			result[kk] = edit.Task.Value
		}
	}
	view.Tasks = view.Tasks.Update(nil, result)
	view.Tasks.Notify()
}

func (view *TasksView) areTasksSame(t1, t2 Tasks) bool {
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
