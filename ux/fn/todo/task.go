// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux/core"
	_ "github.com/dotchain/dot/ux/fn" // blank because it is only used by codegen
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
	return c.fn.Element(
		"root",
		core.Props{Tag: "div", Styles: styles},
		c.fn.Checkbox("cb", core.Styles{}, task.DoneSubstream(c.Cache)),
		c.fn.TextEdit("textedit", core.Styles{}, task.DescriptionSubstream(c.Cache)),
	)
}

//go:generate go run ../codegen.go - $GOFILE
