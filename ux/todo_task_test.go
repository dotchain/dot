// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux_test

import (
	"fmt"
	"github.com/dotchain/dot/ux"
)

type TaskData struct {
	Done        bool
	Description string
}

type TaskDataStream struct {
	*ux.Notifier
	TaskData
	ux.Change
	Next *TaskDataStream
}

func (s *TaskDataStream) Latest() *TaskDataStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

func (s *TaskDataStream) Update(data TaskData) *TaskDataStream {
	s.Next = &TaskDataStream{s.Notifier, data, nil, nil}
	return s.Next
}

type TodoTask struct {
	Root     ux.Element
	styles   ux.Styles
	TaskData *TaskDataStream

	cb          *ux.Checkbox
	description *ux.TextEdit
}

func NewTodoTask(styles ux.Styles, data TaskData) *TodoTask {
	cb := ux.NewCheckbox(ux.Styles{}, data.Done)
	desc := ux.NewTextEdit(ux.Styles{}, data.Description)
	t := &TodoTask{
		ux.NewElement(ux.Props{Tag: "div", Styles: styles}, cb.Root, desc.Root),
		styles,
		&TaskDataStream{&ux.Notifier{}, data, nil, nil},
		cb,
		desc,
	}
	cb.Checked.On(&ux.Handler{t.on})
	desc.Text.On(&ux.Handler{t.on})
	return t
}

func (t *TodoTask) Update(styles ux.Styles, data TaskData) {
	if styles != t.styles {
		t.styles = styles
		t.Root.SetProp("Styles", styles)
	}

	t.TaskData = t.TaskData.Update(data)
	t.cb.Update(ux.Styles{}, data.Done)
	t.description.Update(ux.Styles{}, data.Description)
}

func (t *TodoTask) on() {
	data := TaskData{t.cb.Checked.Value, t.description.Text.Value}
	t.TaskData = t.TaskData.Update(data)
	t.TaskData.Notify()
}

func Example_renderTask() {
	data := TaskData{Done: false, Description: "first task"}
	t := NewTodoTask(ux.Styles{Color: "blue"}, data)
	fmt.Println("Task:", t.Root)

	t.Update(ux.Styles{Color: "red"}, data)
	fmt.Println("Task:", t.Root)

	data.Done = true
	t.Update(ux.Styles{Color: "red"}, data)
	fmt.Println("Task:", t.Root)

	// Output:
	// Task: Props{div false   {blue} <nil> <nil>}( Props{input false checkbox  {} <nil> 0x10f0470}() Props{input false text first task {} <nil> 0x10f04b0}())
	// Task: Props{div false   {red} <nil> <nil>}( Props{input false checkbox  {} <nil> 0x10f0470}() Props{input false text first task {} <nil> 0x10f04b0}())
	// Task: Props{div false   {red} <nil> <nil>}( Props{input true checkbox  {} <nil> 0x10f0470}() Props{input false text first task {} <nil> 0x10f04b0}())

}
