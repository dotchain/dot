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
	}
	ioutil.WriteFile("generated.go", []byte(compiler.Generate(info)), 0644)
}
