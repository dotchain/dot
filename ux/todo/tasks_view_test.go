// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux"
	"github.com/dotchain/dot/ux/todo"
)

func Example_renderTasks() {
	tasks := todo.Tasks{
		{"one", false, "first task"},
		{"two", true, "second task"},
	}
	t := todo.NewTasksView(ux.Styles{}, true, false, tasks)
	fmt.Println(t.Root)

	t.Update(ux.Styles{Color: "red"}, false, false, tasks)
	fmt.Println(t.Root)

	t.Update(ux.Styles{}, true, false, tasks)
	t.Update(ux.Styles{}, true, true, tasks)
	fmt.Println(t.Root)

	// Output:
	// div[]( div[]( input[type:checkbox checked]() input[type:text](second task)))
	// div[styles:{red}]()
	// div[]( div[]( input[type:checkbox]() input[type:text](first task)) div[]( input[type:checkbox checked]() input[type:text](second task)))
}
