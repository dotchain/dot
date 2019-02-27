// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fn

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/todo"
)

//go:generate go run codegen.go - $GOFILE

// The above line generates a couple of types for each function
// marked with a doc-block header "codegen: options"
//
// Such functions should have the first parameter be a "context" of a
// unique type (which is one of the generated types). The context can
// be used to access other components.  The context itself is private
// and meant only to be used within this particular component. For
// each such function, a more public Cache type is also generated. See
// the TasksView example below for actual usage.

// TasksView is a function representation of a task view component.
//
// An example consumer of TasksView  can look like this:
//
//      func TasksViewConsumer(c *tvcCtx, args...) core.Element {
//              return c.Element(
//                     <key>,
//                     core.Props{...},
//                     c.TasksView(
//                          <key>, // note: note tasksViewCtx type
//                          core.Styles{...},
//                          showDone, showNotDone, tasks,
//                     ),
//                     ... other children...
//              )
//       }
//
// Note also that this function returns a core element directly. It
// can return any time and any number of elements though that is not
// encouraged.
//
// codegen: pure
func TasksView(c *tvCtx, styles core.Styles, showDone bool, showNotDone bool, tasks *todo.TasksStream) core.Element {
	tasks = tasks.Latest()

	// the c.Element call here ends  up calling
	// c.ElementCache.Element(key, ...)
	return c.Element(
		"root_key",
		core.Props{Tag: "div", Styles: styles},
		mapTasks(tasks.Latest().Value, func(idx int, t todo.Task) core.Element {
			if t.Done && !showDone || t.Done && !showNotDone {
				return nil
			}

			key := t.ID

			// the c.todo.TaskEdit call ends up calling
			// todo.TaskEditCache.TaskEdit(key,...)
			edit := c.todo.TaskEdit(key, core.Styles{}, t)

			// pass notifications upwards appropriately
			c.On(edit.Task.Notifier, func() {
				updated := append(todo.Tasks(nil), tasks.Value...)
				updated[idx] = edit.Task.Latest().Value
				tasks.Update(nil, updated)
			})

			return edit.Root
		})...,
	)
}

func mapTasks(tasks todo.Tasks, fn func(int, todo.Task) core.Element) []core.Element {
	result := []core.Element{}
	for kk, t := range tasks {
		result = append(result, fn(kk, t))
	}
	return result
}
