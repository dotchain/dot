// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package todo is an example TODO task list app
package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/simple"
	"github.com/dotchain/dot/ux/streams"
)

// App is a thin wrapper on top of TasksView with checkboxes for ShowDone and ShowUndone
type App struct {
	simple.Element

	styles core.Styles

	simple.CheckboxCache
	TasksViewCache
	streams.Subs
}

// NewApp creates the new app control
func NewApp(styles core.Styles, tasks Tasks) *App {
	app := &App{}
	app.Update(styles, tasks)
	return app
}

// Update app
func (app *App) Update(styles core.Styles, tasks Tasks) {
	app.styles = styles

	// the following ugliness is because the checkboxes are used as
	// state to drive the rendering
	showDone, showNotDone := true, true
	done, notDone := app.CheckboxCache.Item("done"), app.CheckboxCache.Item("notDone")
	if done != nil {
		showDone = done.Checked.Value
	}
	if notDone != nil {
		showNotDone = notDone.Checked.Value
	}

	app.CheckboxCache.Begin()
	app.TasksViewCache.Begin()
	app.Subs.Begin()
	defer app.CheckboxCache.End()
	defer app.TasksViewCache.End()
	defer app.Subs.End()

	done = app.Checkbox("done", core.Styles{}, showDone)
	notDone = app.Checkbox("notDone", core.Styles{}, showNotDone)
	tasksView := app.TasksView("tasks", core.Styles{}, showDone, showNotDone, tasks)

	app.On(done.Checked.Notifier, app.on)
	app.On(notDone.Checked.Notifier, app.on)
	app.On(tasksView.Tasks.Notifier, app.on)

	// do the actual rendering
	app.Declare(core.Props{Tag: "div", Styles: styles}, done.Root, notDone.Root, tasksView.Root)
}

// Tasks returns the current stream instance of tasks
func (app *App) Tasks() *TasksStream {
	return app.TasksViewCache.Item("tasks").Tasks
}

func (app *App) on() {
	app.Update(app.styles, app.Tasks().Value)
}
