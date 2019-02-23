// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/todo"
)

func Example_renderTask() {
	task := todo.Task{Done: false, Description: "first task"}
	t := todo.NewTaskEdit(core.Styles{Color: "blue"}, task)
	fmt.Println(t.Root)

	t.Update(core.Styles{Color: "red"}, task)
	fmt.Println(t.Root)

	task.Done = true
	t.Update(core.Styles{Color: "red"}, task)
	fmt.Println(t.Root)

	// Output:
	// div[styles:{blue}]( input[type:checkbox]() input[type:text](first task))
	// div[styles:{red}]( input[type:checkbox]() input[type:text](first task))
	// div[styles:{red}]( input[type:checkbox checked]() input[type:text](first task))
}
