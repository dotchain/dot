// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux_test

import (
	"fmt"
	"github.com/dotchain/dot/ux"
)

type TaskList struct {
	ShowDone, ShowUndone bool
	Tasks                []TaskData
}

type TodoTaskList struct {
	Root       ux.Element
	styles     ux.Styles
	tasksCache TodoTaskCache
}

func NewTodoTaskList(styles ux.Styles, data TaskList) *TodoTaskList {
	t := &TodoTaskList{
		ux.NewElement(ux.Props{Tag: "div", Styles: styles}),
		styles,
		TodoTaskCache{},
	}
	t.Update(styles, data)
	return t
}

func (t *TodoTaskList) Update(styles ux.Styles, data TaskList) {
	if t.styles != styles {
		t.styles = styles
		t.Root.SetProp("Styles", styles)
	}

	t.tasksCache.Reset()
	defer t.tasksCache.Cleanup()

	children := []ux.Element{}
	for _, td := range data.Tasks {
		if td.Done && !data.ShowDone || !td.Done && !data.ShowUndone {
			continue
		}

		task, exists := t.tasksCache.Get(td.ID, ux.Styles{}, td)
		if !exists {
			task.TaskData.On(&ux.Handler{t.on})
		}
		children = append(children, task.Root)
	}
	ux.UpdateChildren(t.Root, children)
}

func (t *TodoTaskList) on() {
	// need to aggregate all received changes and generate the updated changes
}

func Example_renderTaskList() {
	tasks := []TaskData{
		{"one", false, "first task"},
		{"two", true, "second task"},
	}
	list := TaskList{true, false, tasks}
	t := NewTodoTaskList(ux.Styles{}, list)
	fmt.Println(t.Root)

	list = TaskList{false, false, tasks}
	t.Update(ux.Styles{Color: "red"}, list)
	fmt.Println(t.Root)

	list = TaskList{true, false, tasks}
	t.Update(ux.Styles{}, list)
	list = TaskList{true, true, tasks}
	t.Update(ux.Styles{}, list)
	fmt.Println(t.Root)

	// Output:
	// div[]( div[]( input[type:checkbox checked]() input[type:text](second task)))
	// div[styles:{red}]()
	// div[]( div[]( input[type:checkbox]() input[type:text](first task)) div[]( input[type:checkbox checked]() input[type:text](second task)))
}
