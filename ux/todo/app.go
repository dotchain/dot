// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo

import "github.com/dotchain/dot/ux"

// App is a thin wrapper on top of TasksView
type App struct {
	Root ux.Element

	styles                ux.Styles
	showDone, showNotDone *ux.Checkbox
	tasksView             *TasksView

	Tasks *TasksStream
}

// NewApp creates the new app control
func NewApp(styles ux.Styles, tasks Tasks) *App {
	// TODO: need labels for these two + a container to wrap them
	showDone := ux.NewCheckbox(ux.Styles{}, true)
	showNotDone := ux.NewCheckbox(ux.Styles{}, true)

	tasksView := NewTasksView(ux.Styles{}, true, true, tasks)
	root := ux.NewElement(
		ux.Props{Tag: "div", Styles: styles},
		showDone.Root,
		showNotDone.Root,
		tasksView.Root,
	)
	app := &App{root, styles, showDone, showNotDone, tasksView, tasksView.Tasks}
	showDone.Checked.On(&ux.Handler{app.refresh})
	showNotDone.Checked.On(&ux.Handler{app.refresh})
	tasksView.Tasks.On(&ux.Handler{app.refresh})

	return app
}

// Update props
func (app *App) Update(styles ux.Styles, tasks Tasks) {
	if app.styles != styles {
		app.styles = styles
		app.Root.SetProp("Styles", styles)
	}
	app.tasksView.Update(ux.Styles{}, app.showDone.Checked.Value, app.showNotDone.Checked.Value, tasks)
}

func (app *App) refresh() {
	app.Update(app.styles, app.tasksView.Tasks.Value)
}
