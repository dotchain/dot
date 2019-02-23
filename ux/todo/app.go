// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/simple"
	"github.com/dotchain/dot/ux/streams"
)

// App is a thin wrapper on top of TasksView with checkboxes for ShowDone and ShowUndone
type App struct {
	simple.Element

	styles     core.Styles
	toggles    simple.CheckboxCache
	tasksViews TasksViewCache
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
	done, notDone := app.toggles.Item("done"), app.toggles.Item("notDone")
	if done != nil {
		showDone = done.Checked.Value
	}
	if notDone != nil {
		showNotDone = notDone.Checked.Value
	}

	app.toggles.Begin()
	app.tasksViews.Begin()
	defer app.toggles.End()
	defer app.tasksViews.End()

	// do the actual rendering
	app.Declare(
		core.Props{Tag: "div", Styles: styles},
		app.listenToggle(app.toggles.TryGet("done", core.Styles{}, showDone)),
		app.listenToggle(app.toggles.TryGet("notDone", core.Styles{}, showNotDone)),
		app.listenTasks(app.tasksViews.TryGet("tasks", core.Styles{}, showDone, showNotDone, tasks)),
	)
}

// Tasks returns the current stream instance of tasks
func (app *App) Tasks() *TasksStream {
	return app.tasksViews.Item("tasks").Tasks
}

func (app *App) listenToggle(cb *simple.Checkbox, exists bool) core.Element {
	if !exists {
		cb.Checked.On(&streams.Handler{app.on})
	}
	return cb.Root
}

func (app *App) listenTasks(tasksView *TasksView, exists bool) core.Element {
	if !exists {
		tasksView.Tasks.On(&streams.Handler{app.on})
	}
	return tasksView.Root
}

func (app *App) on() {
	app.Update(app.styles, app.Tasks().Value)
}
