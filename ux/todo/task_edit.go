// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import "github.com/dotchain/dot/ux"
import "github.com/dotchain/dot/ux/simple"

// Task represents an item in the TODO list.
type Task struct {
	ID          string
	Done        bool
	Description string
}

// generate TaskStream
//go:generate go run ../templates/gen.go ../templates/streams.template Package=todo Base=Task BaseType=Task out=task_stream.go

// TaskEdit is a control that displays a task as well as allowing it
// to be edited. The current value of the data is available in the
// Task field (which is a stream and so supports On/Off methods).
type TaskEdit struct {
	Root ux.Element

	styles      ux.Styles
	cb          *simple.Checkbox
	description *ux.TextEdit

	Task *TaskStream
}

// generate the TaskEditCache for any consumers who want it

//go:generate go run ../templates/gen.go ../templates/cache.template Package=todo Base=TaskEdit BaseType=TaskEdit "Args=styles, task" "ArgsDef=styles ux.Styles, task Task" Constructor=NewTaskEdit out=task_edit_cache.go

// NewTaskEdit is the constructor for creating a TaskEdit control
func NewTaskEdit(styles ux.Styles, task Task) *TaskEdit {
	cb := simple.NewCheckbox(ux.Styles{}, task.Done)
	desc := ux.NewTextEdit(ux.Styles{}, task.Description)
	t := &TaskEdit{
		ux.NewElement(ux.Props{Tag: "div", Styles: styles}, cb.Root, desc.Root),
		styles,
		cb,
		desc,
		&TaskStream{&ux.Notifier{}, task, nil, nil},
	}
	cb.Checked.On(&ux.Handler{t.on})
	desc.Text.On(&ux.Handler{t.on})
	return t
}

// Update updates the style or the content of the task
func (t *TaskEdit) Update(styles ux.Styles, task Task) {
	if styles != t.styles {
		t.styles = styles
		t.Root.SetProp("Styles", styles)
	}

	t.Task = t.Task.Update(nil, task)
	t.cb.Update(ux.Styles{}, task.Done)
	t.description.Update(ux.Styles{}, task.Description)
}

// on is called whenever either Done or Description is modified by
// child controls
func (t *TaskEdit) on() {
	id := t.Task.Value.ID
	data := Task{id, t.cb.Checked.Value, t.description.Text.Value}
	t.Task = t.Task.Update(nil, data)
	t.Task.Notify()
}
