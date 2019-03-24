// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dotc implements code-generation tools for dot.changes
package dotc

import (
	"bytes"
	"go/format"
	"io"
	"text/template"

	"golang.org/x/tools/imports"
)

// Info tracks all information used for code generation
type Info struct {
	Package string
	Imports [][2]string
	Structs []Struct
}

// Generate generates a single file with all the provided info
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

var infoTpl = template.Must(template.New("letter").Parse(`
// Generated.  DO NOT EDIT.
package {{.Package}}

import (
	{{range .Imports}}{{index . 0}} "{{index . 1}}"
	{{end -}}
)
`))

var basic = map[string]bool{
	"bool": true,
	"int":  true,
}

// Field holds info for a struct field
type Field struct {
	Name, Key, Type  string
	Atomic, Nullable bool
}

func (f Field) Wrap(recv string) string {
	if f.Atomic || basic[f.Type] {
		return "changes.Atomic{" + recv + "." + f.Name + "}"
	}
	star := ""
	if f.Nullable {
		star = "*"
	}
	if f.Type == "string" {
		return "types.S16(" + star + recv + "." + f.Name + ")"
	}
	return star + recv + "." + f.Name
}

func (f Field) Unwrap() string {
	if f.Atomic || basic[f.Type] {
		return "v.(changes.Atomic).Value.(" + f.Type + ")"
	}

	nullable := func(s string) string {
		if !f.Nullable {
			return s
		}
		return "func(x " + f.Type + ") *" + f.Type + " { return &x }(" + s + ")"
	}

	if f.Type == "string" {
		return nullable("string(v.(types.S16))")
	}
	return nullable("v.(" + f.Type + ")")
}

// Struct has the type information of a struct for code generation of
// the Apply() and SetField(..) methods
type Struct struct {
	Recv, Type string
	Fields     []Field
}

func (s Struct) Pointer() bool {
	return s.Type[0] == '*'
}

func (s Struct) GenerateApply(w io.Writer) error {
	return structApply.Execute(w, s)
}

var structApply = template.Must(template.New("letter").Parse(`
{{ $r := .Recv }}
func ({{$r}} {{.Type}}) get(key interface{}) changes.Value {
	switch key {
	{{- range .Fields}}
	case "{{.Key}}":
		{{- if .Nullable}}if {{$r}}.{{.Name}} == nil { return changes.Nil } {{end}}
		return {{.Wrap $r}}
        {{- end }}
        }
	panic(key)
}

func ({{$r}} {{.Type}}) set(key interface{}, v changes.Value) changes.Value {
	{{$r}}Clone := {{if .Pointer}}*{{end}}{{$r}}
	switch key {
	{{- range .Fields}}
	case "{{.Key}}":
		{{- if .Nullable}}if v == changes.Nil { {{$r}}Clone.{{.Name}} = nil; break } 
		{{end -}}
		{{$r}}Clone.{{.Name}} = {{.Unwrap}}
        {{- end }}
	default: 
		panic(key)
        }
	return {{if .Pointer}}&{{end}} {{.Recv}}Clone
}

func ({{$r}} {{.Type}}) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: {{$r}}.get, Set: {{$r}}.set}).Apply(ctx, c, {{$r}})
}
`))

// Union has the information of a struct used for unions for code
// generation of the Apply() and SetField() methods
type Union struct {
	Name   string
	Fields []Field
}

// Slice  has the type information of a slice type for code generation
// of the Apply, ApplyCollection and splice methods
type Slice struct {
	Name, ElementName string
	Atomic, Nullable  bool
}
