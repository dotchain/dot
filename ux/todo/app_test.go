// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux"
	"github.com/dotchain/dot/ux/todo"
)

func Example_renderApp() {
	tasks := todo.Tasks{
		{"one", false, "first task"},
		{"two", true, "second task"},
	}
	app := todo.NewApp(ux.Styles{}, tasks)
	fmt.Println(app.Root)

	// set "ShowDone" to false which should filter out the second task
	app.Root.(*element).children[0].(*element).setValue("off")
	fmt.Println(app.Root)

	// Output:
	// div[]( input[type:checkbox checked]() input[type:checkbox checked]() div[]( div[]( input[type:checkbox]() input[type:text](first task)) div[]( input[type:checkbox checked]() input[type:text](second task))))
	// div[]( input[type:checkbox]() input[type:checkbox checked]() div[]( div[]( input[type:checkbox]() input[type:text](first task))))
}
