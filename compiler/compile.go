// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package compiler has a set of code generation tools used in DOT
package compiler

import (
	"bytes"
	"go/format"
	"golang.org/x/tools/imports"
)

// Info contians all the info needed to generate code
type Info struct {
	Package string
	Imports [][2]string
	Streams []StreamInfo
}

// StreamInfo holds the information to generate a single stream type
type StreamInfo struct {
	StreamType      string
	ValueType       string
	Fields          []FieldInfo
	EntryStreamType string
}

// Compile generates the code needed to deal with a stream
func (s *StreamInfo) Compile() string {
	var result bytes.Buffer
	must(streamTpl.Execute(&result, s))
	for _, f := range s.Fields {
		var data struct {
			*StreamInfo
			*FieldInfo
		}
		data.StreamInfo = s
		data.FieldInfo = &f
		must(fieldTpl.Execute(&result, data))
	}

	if s.EntryStreamType != "" {
		must(entryTpl.Execute(&result, s))
	}

	return result.String()
}

// FieldInfo holds info on individual substream fields of the base stream
type FieldInfo struct {
	Field           string
	FieldType       string
	FieldStreamType string
}

// Compile returns the source code generated from the provided info
func Compile(info Info) string {
	var result bytes.Buffer
	must(headerTpl.Execute(&result, info))
	r := result.String()
	for _, s := range info.Streams {
		r += "\n" + s.Compile()
	}

	p, err := format.Source([]byte(r))
	must(err)

	p, err = imports.Process("compiled.go", p, nil)
	must(err)

	return string(p)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
