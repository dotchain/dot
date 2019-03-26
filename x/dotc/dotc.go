// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dotc implements code-generation tools for dot.changes
package dotc

import (
	"bytes"
	"go/format"

	"golang.org/x/tools/imports"
)

// Info tracks all information used for code generation
type Info struct {
	Package string
	Imports [][2]string
	Structs []Struct
	Unions  []Union
	Slices  []Slice
}

// GenerateTests generates the tests
func (info Info) GenerateTests() (result string, err error) {
	var buf bytes.Buffer

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	info.Imports = append(
		[][2]string{
			{"", "github.com/dotchain/dot/changes"},
			{"", "github.com/dotchain/dot/changes/types"},
			{"", "github.com/dotchain/dot/streams"},
			{"", "reflect"},
			{"", "testing"},
		}, info.Imports...)

	must(infoTpl.Execute(&buf, info))

	for _, s := range info.Structs {
		must(StructStream(s).GenerateStreamTests(&buf))
	}
	for _, u := range info.Unions {
		must(UnionStream(u).GenerateStreamTests(&buf))
	}
	for _, s := range info.Slices {
		must(SliceStream(s).GenerateStreamTests(&buf))
	}

	result = buf.String()
	p, err := format.Source([]byte(result))
	must(err)

	result = string(p)
	p, err = imports.Process("generated_test.go", p, nil)
	return string(p), err
}

// Generate implements the helper methods for the provided types
func (info Info) Generate() (result string, err error) {
	var buf bytes.Buffer

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	info.Imports = append(
		[][2]string{
			{"", "github.com/dotchain/dot/changes"},
			{"", "github.com/dotchain/dot/changes/types"},
			{"", "github.com/dotchain/dot/streams"},
		}, info.Imports...)

	must(infoTpl.Execute(&buf, info))

	for _, s := range info.Structs {
		must(s.GenerateApply(&buf))
		must(s.GenerateSetters(&buf))
		must(StructStream(s).GenerateStream(&buf))
	}

	for _, u := range info.Unions {
		must(u.GenerateApply(&buf))
		must(u.GenerateSetters(&buf))
		must(UnionStream(u).GenerateStream(&buf))
	}

	for _, s := range info.Slices {
		must(s.GenerateApply(&buf))
		must(s.GenerateSetters(&buf))
		must(SliceStream(s).GenerateStream(&buf))
	}

	result = buf.String()
	p, err := format.Source([]byte(result))
	must(err)

	result = string(p)
	p, err = imports.Process("generated.go", p, nil)
	return string(p), err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
