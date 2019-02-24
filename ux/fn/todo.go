// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fn

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/simple"
	"github.com/dotchain/dot/ux/todo"
)

// The following generate header generates  TasksViewCache which
// allows TasksView to be used nicely within other functional
// components.
//
// The code also generates the tasksViewCtx context struct which is
// only used from within this file and consumers are not expected to
// refer to it at all.  Instead consumers are expected to use whatever
// key they prefer to use.

//go:generate go run cmd/gen.go TasksView $GOFILE

// TasksView is a function representation of a task view component.
//
// An example consumer of TasksView  can look like this:
//
//      func TasksViewConsumer(c *tasksViewConsumerCtx, args...) core.Element {
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
func TasksView(c *tasksViewCtx, styles core.Styles, showDone bool, showNotDone bool, tasks todo.Tasks) core.Element {

	// the c.Element call here ends  up calling
	// c.ElementCache.Element(key, ...)
	return c.Element(
		"root_key",
		core.Props{Tag: "div", Styles: styles},
		mapTasks(tasks, func(t todo.Task) core.Element {
			if t.Done && !showDone || t.Done && !showNotDone {
				return nil
			}

			key := t.ID

			// the c.TaskEdit call ends up calling
			// c.TaskEditCache.TaskEdit(key,...)
			return c.TaskEdit(key, core.Styles{}, t).Root
		})...,
	).Root
}

// ElementCache is an adapter to call into the simple-style
// struct. This adaptor is not needed for any caller of TasksView as
// the corresponding cache is generated
type ElementCache struct {
	old, current map[interface{}]*simple.Element
}

// Begin starts a round
func (e *ElementCache) Begin() {
	e.old, e.current = e.current, map[interface{}]*simple.Element{}
}

// End ends a round
func (e *ElementCache) End() {
	e.old = nil
}

// Element returns a core element
func (e *ElementCache) Element(key interface{}, props core.Props, children ...core.Element) *simple.Element {
	if old, ok := e.old[key]; ok {
		e.current[key] = old
	} else {
		e.current[key] = &simple.Element{}
	}
	e.current[key].Declare(props, children...)
	return e.current[key]
}

// TaskEditCache is an adaptor to call into the simple-style struct.
// This adaptor is not needed for any caller of TasksView as that
// corresponding cache is generated.  Note how this Cache is only
// different in the name  of the function -- so potentially the simple
// scheme can be changed to generate matching signatures..
type TaskEditCache struct {
	todo.TaskEditCache
}

// TaskEdit maps to todo.TaskCache.Get.  Might be worth simply changing
// all the cache fetch names in all those caches too.
func (t *TaskEditCache) TaskEdit(key interface{}, styles core.Styles, task todo.Task) *todo.TaskEdit {
	return t.Get(key, styles, task)
}

func mapTasks(tasks todo.Tasks, fn func(todo.Task) core.Element) []core.Element {
	result := []core.Element{}
	for _, t := range tasks {
		result = append(result, fn(t))
	}
	return result
}
