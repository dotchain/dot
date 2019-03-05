// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build ignore

package main

import (
	"github.com/dotchain/dot/compiler"
	"io/ioutil"
)

func main() {
	info := compiler.Info{
		Package: "todo",
		Imports: [][2]string{
			{"uxstreams", "github.com/dotchain/dot/ux/streams"},
		},
		Streams: []compiler.StreamInfo{
			{
				StreamType: "TaskStream",
				ValueType:  "Task",
				Fields: []compiler.FieldInfo{{
					Field:            "Done",
					FieldType:        "bool",
					FieldStreamType:  "uxstreams.BoolStream",
					FieldSubstream:   "DoneSubstream",
					FieldConstructor: "uxstreams.NewBoolStream",
				}, {
					Field:            "Description",
					FieldType:        "string",
					FieldStreamType:  "uxstreams.TextStream",
					FieldSubstream:   "DescriptionSubstream",
					FieldConstructor: "uxstreams.NewTextStream",
				}},
				EntryStreamType: "",
			},
			{
				StreamType:       "TasksStream",
				ValueType:        "Tasks",
				EntryStreamType:  "TaskStream",
				EntryConstructor: "NewTaskStream",
			},
		},
		Contexts: []compiler.ContextInfo{
			{
				ContextType: "taskEditCtx",

				Function:      "TaskEdit",
				Subcomponents: []string{"fn.CheckboxCache", "fn.ElementCache", "fn.TextEditCache"},
				Params: []compiler.ParamInfo{
					{Name: "ctx", Type: "*taskEditCtx"},
					{Name: "styles", Type: "core.Styles"},
					{Name: "task", Type: "*TaskStream"},
				},
				Results: []compiler.ResultInfo{{Name: "", Type: "core.Element"}},

				Component:         "TaskEditCache",
				ComponentComments: "// TaskEditCache is good",
				Method:            "TaskEdit",
				MethodComments:    "// TaskEdit is good",
			},
			{
				ContextType: "tasksViewCtx",

				Function:      "TasksView",
				Subcomponents: []string{"TaskEditCache", "fn.ElementCache"},
				Params: []compiler.ParamInfo{
					{Name: "ctx", Type: "*tasksViewCtx"},
					{Name: "styles", Type: "core.Styles"},
					{Name: "showDone", Type: "*uxstreams.BoolStream"},
					{Name: "showNotDone", Type: "*uxstreams.BoolStream"},
					{Name: "tasks", Type: "*TasksStream"},
				},
				Results: []compiler.ResultInfo{{Name: "", Type: "core.Element"}},

				Component:         "TasksViewCache",
				ComponentComments: "// TasksViewCache is good",
				Method:            "TasksView",
				MethodComments:    "// TasksView is good",
			},
			{
				ContextType: "appCtx",

				Function:      "App",
				Subcomponents: []string{"TasksViewCache", "fn.ElementCache", "fn.CheckboxCache"},
				Params: []compiler.ParamInfo{
					{"ctx", "*appCtx"},
					{"styles", "core.Styles"},
					{"tasks", "*TasksStream"},
					{"doneState", "*uxstreams.BoolStream"},
					{"notDoneState", "*uxstreams.BoolStream"},
				},
				Results: []compiler.ResultInfo{
					{"", "core.Element"},
					{"", "*uxstreams.BoolStream"},
					{"", "*uxstreams.BoolStream"},
				},

				Component:         "AppCache",
				ComponentComments: "// AppCache is good",
				Method:            "App",
				MethodComments:    "// App is good",
			},
		},
	}
	ioutil.WriteFile("generated.go", []byte(compiler.Generate(info)), 0644)
}
