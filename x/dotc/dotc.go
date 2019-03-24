// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dotc implements code-generation tools for dot.changes
package dotc

import (
	"bytes"
	"go/format"
	"io"

	"golang.org/x/tools/imports"
)

// Info tracks all information used for code generation
type Info struct {
	Package string
	Imports [][2]string
	Structs []Struct
}

// Generate implements the helper methods for the provided types
func (info Info) Generate() (string, error) {
	var buf bytes.Buffer

	info.Imports = append(
		[][2]string{
			{"", "github.com/dotchain/dot/changes"},
			{"", "github.com/dotchain/dot/changes/types"},
		}, info.Imports...)

	if err := infoTpl.Execute(&buf, info); err != nil {
		return "", err
	}

	for _, s := range info.Structs {
		if err := s.GenerateApply(&buf); err != nil {
			return "", err
		}
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), err
	}

	p2, err := imports.Process("generated.go", p, nil)
	if err != nil {
		return string(p), err
	}

	return string(p2), nil
}

var basic = map[string]bool{
	"bool": true,
	"int":  true,
}

// Field holds info for a struct field
type Field struct {
	Name, Key, Type string
	Atomic          bool
}

// Wrap returns a string form of the field that maps to a changes.Value type
func (f Field) Wrap(recv string) string {
	if f.Atomic || basic[f.Type] {
		return "changes.Atomic{" + recv + "." + f.Name + "}"
	}
	star := ""
	if f.Type == "string" {
		return "types.S16(" + star + recv + "." + f.Name + ")"
	}
	return star + recv + "." + f.Name
}

// Unwrap converts "v" to the type of the field
func (f Field) Unwrap() string {
	if f.Atomic || basic[f.Type] {
		return "v.(changes.Atomic).Value.(" + f.Type + ")"
	}

	if f.Type == "string" {
		return "string(v.(types.S16))"
	}
	return "v.(" + f.Type + ")"
}

// Struct has the type information of a struct for code generation of
// the Apply() and SetField(..) methods
type Struct struct {
	Recv, Type string
	Fields     []Field
}

// Pointer specifies if the struct type is itself a pointer
func (s Struct) Pointer() bool {
	return s.Type[0] == '*'
}

// GenerateApply generates the code for the changes.Value Apply() method
func (s Struct) GenerateApply(w io.Writer) error {
	return structApply.Execute(w, s)
}

// Slice  has the type information of a slice type for code generation
// of the Apply, ApplyCollection and splice methods
type Slice struct {
	Name, ElementName string
	Atomic            bool
}
