// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/simple"
	"github.com/dotchain/dot/ux/streams"
)

// Task represents an item in the TODO list.
type Task struct {
	ID          string
	Done        bool
	Description string
}

// TaskEdit is a control that displays a task as well as allowing it
// to be edited. The current value of the data is available in the
// Task field (which is a stream and so supports On/Off methods).
type TaskEdit struct {
	simple.Element

	Task *TaskStream

	cb   simple.CheckboxCache
	desc simple.TextEditCache
}

// NewTaskEdit is the constructor for creating a TaskEdit control.
func NewTaskEdit(styles core.Styles, task Task) *TaskEdit {
	e := &TaskEdit{}
	e.Task = NewTaskStream(task)
	e.Update(styles, task)
	return e
}

// Update updates the styles or task forthis control.
func (e *TaskEdit) Update(styles core.Styles, task Task) {
	if e.Task.Value != task {
		e.Task = e.Task.Update(nil, task)
	}

	e.cb.Begin()
	e.desc.Begin()
	defer e.cb.End()
	defer e.desc.End()

	e.Declare(
		core.Props{Tag: "div", Styles: styles},
		e.checkbox(e.cb.TryGet("cb", core.Styles{}, task.Done)),
		e.description(e.desc.TryGet("desc", core.Styles{}, task.Description)),
	)
}

func (e *TaskEdit) checkbox(cb *simple.Checkbox, exists bool) core.Element {
	if !exists {
		cb.Checked.On(&streams.Handler{e.on})
	}
	return cb.Root
}

func (e *TaskEdit) description(edit *simple.TextEdit, exists bool) core.Element {
	if !exists {
		edit.Text.On(&streams.Handler{e.on})
	}
	return edit.Root
}

// on is called whenever either Done or Description is modified by
// child controls
func (e *TaskEdit) on() {
	cb, desc := e.cb.Item("cb"), e.desc.Item("desc")
	data := Task{e.Task.Value.ID, cb.Checked.Value, desc.Text.Value}
	e.Task = e.Task.Update(nil, data)
	e.Task.Notify()
}

// generate TaskStream
//go:generate go run ../templates/gen.go ../templates/streams.template Package=todo Base=Task BaseType=Task out=task_stream.go

// generate the TaskEditCache for any consumers who want it

//go:generate go run ../templates/gen.go ../templates/cache.template Package=todo Base=TaskEdit BaseType=TaskEdit "Args=styles, task" "ArgsDef=styles core.Styles, task Task" Constructor=NewTaskEdit out=task_edit_cache.go
