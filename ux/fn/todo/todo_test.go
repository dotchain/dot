// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/fn/todo"
	"github.com/dotchain/dot/ux/streams"
)

func Example_renderApp() {
	tasks := todo.Tasks{
		{"one", false, "first task"},
		{"two", true, "second task"},
	}
	cache := todo.AppCache{}

	cache.Begin()
	root := cache.App("root", core.Styles{}, todo.NewTasksStream(tasks))
	cache.End()

	fmt.Println(root)

	// set "ShowDone" to false which should filter out the second task
	root.(*element).children[0].(*element).setValue("off")
	fmt.Println(root)

	// Output:
	// div[]( input[type:checkbox checked]() input[type:checkbox checked]() div[]( div[]( input[type:checkbox]() input[type:text](first task)) div[]( input[type:checkbox checked]() input[type:text](second task))))
	// div[]( input[type:checkbox]() input[type:checkbox checked]() div[]( div[]( input[type:checkbox]() input[type:text](first task))))
}

func Example_renderTasks() {
	cache := todo.TasksViewCache{}

	tasks := todo.Tasks{
		{"one", false, "first task"},
		{"two", true, "second task"},
	}
	s := todo.NewTasksStream(tasks)
	showDone, showNotDone := streams.NewBoolStream(true), streams.NewBoolStream(false)
	cache.Begin()
	root := cache.TasksView("root", core.Styles{}, showDone, showNotDone, s)
	cache.End()
	fmt.Println(root)

	showDone = showDone.Append(nil, false, true)
	cache.Begin()
	root = cache.TasksView("root", core.Styles{Color: "red"}, showDone, showNotDone, s)
	cache.End()
	fmt.Println(root)

	showDone = showDone.Append(nil, true, true)
	cache.Begin()
	_ = cache.TasksView("root", core.Styles{}, showDone, showNotDone, s)
	cache.End()
	showNotDone = showNotDone.Append(nil, true, true)
	cache.Begin()
	root = cache.TasksView("root", core.Styles{}, showDone, showNotDone, s)
	cache.End()
	fmt.Println(root)

	// Output:
	// div[]( div[]( input[type:checkbox checked]() input[type:text](second task)))
	// div[styles:{red}]()
	// div[]( div[]( input[type:checkbox]() input[type:text](first task)) div[]( input[type:checkbox checked]() input[type:text](second task)))
}

func Example_renderTask() {
	task := todo.NewTaskStream(todo.Task{Done: false, Description: "first task"})
	cache := todo.TaskEditCache{}
	cache.Begin()
	root := cache.TaskEdit("root", core.Styles{Color: "blue"}, task)
	cache.End()
	fmt.Println(root)

	cache.Begin()
	root = cache.TaskEdit("root", core.Styles{Color: "red"}, task)
	cache.End()
	fmt.Println(root)

	next := task.Value
	next.Done = true
	task = task.Append(nil, next, true)
	cache.Begin()
	root = cache.TaskEdit("root", core.Styles{Color: "red"}, task)
	cache.End()
	fmt.Println(root)

	// Output:
	// div[styles:{blue}]( input[type:checkbox]() input[type:text](first task))
	// div[styles:{red}]( input[type:checkbox]() input[type:text](first task))
	// div[styles:{red}]( input[type:checkbox checked]() input[type:text](first task))
}
