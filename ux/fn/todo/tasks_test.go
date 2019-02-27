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

	showDone = showDone.Update(nil, false)
	cache.Begin()
	root = cache.TasksView("root", core.Styles{Color: "red"}, showDone, showNotDone, s)
	cache.End()
	fmt.Println(root)

	showDone = showDone.Update(nil, true)
	cache.Begin()
	_ = cache.TasksView("root", core.Styles{}, showDone, showNotDone, s)
	cache.End()
	showNotDone = showNotDone.Update(nil, true)
	cache.Begin()
	root = cache.TasksView("root", core.Styles{}, showDone, showNotDone, s)
	cache.End()
	fmt.Println(root)

	// Output:
	// div[]( div[]( input[type:checkbox checked]() input[type:text](second task)))
	// div[styles:{red}]()
	// div[]( div[]( input[type:checkbox]() input[type:text](first task)) div[]( input[type:checkbox checked]() input[type:text](second task)))
}
