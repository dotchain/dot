// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/fn/todo"
)

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
