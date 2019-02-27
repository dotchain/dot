// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux/core"
	_ "github.com/dotchain/dot/ux/fn" // blank because it is only used by codegen
	"github.com/dotchain/dot/ux/streams"
)

// Task represents an item in the TODO list.
type Task struct {
	ID          string
	Done        bool
	Description string
}

// TaskEdit is a control that displays a task as well as allowing it
// to be edited. The current value of the data is available in the
// Task field (which is a stream and so supports On/Off methods).
//
// codegen: pure
func TaskEdit(c *taskEditCtx, styles core.Styles, task *TaskStream) core.Element {
	done := streams.NewBoolStream(task.Value.Done)
	text := streams.NewTextStream(task.Value.Description)

	onChange := func() {
		task = task.Update(nil, Task{task.Value.ID, done.Value, text.Value})
		task.Notify()
	}
	c.On(done.Notifier, onChange)
	c.On(text.Notifier, onChange)

	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		c.fn.Checkbox("cb", core.Styles{}, done),
		c.fn.TextEdit("textedit", core.Styles{}, text),
	)
}

// generate TaskStream
//go:generate go run ../../templates/gen.go ../../templates/streams.template Package=todo Base=Task BaseType=Task out=task_stream.go

//go:generate go run ../codegen.go - $GOFILE
