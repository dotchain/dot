// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/fn/todo"
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
