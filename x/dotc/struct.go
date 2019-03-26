// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

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

// GenerateSetters generates the code for the field setters
func (s Struct) GenerateSetters(w io.Writer) error {
	return structSetters.Execute(w, s)
}

var structApply = template.Must(template.New("struct_apply").Parse(`
{{ $r := .Recv }}
func ({{$r}} {{.Type}}) get(key interface{}) changes.Value {
	switch key {
	{{range .Fields}}
	case "{{.Key}}": return {{.ToValue $r .Name}}{{end}}
        }
	panic(key)
}

func ({{$r}} {{.Type}}) set(key interface{}, v changes.Value) changes.Value {
	{{$r}}Clone := {{if .Pointer}}*{{end}}{{$r}}
	switch key {
	{{- range .Fields}}
	case "{{.Key}}":
		{{$r}}Clone.{{.Name}} = {{.FromValue "v" ""}}
        {{- end }}
	default: 
		panic(key)
        }
	return {{if .Pointer}}&{{end}} {{$r}}Clone
}

func ({{$r}} {{.Type}}) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: {{$r}}.get, Set: {{$r}}.set}).Apply(ctx, c, {{$r}})
}
`))

var structSetters = template.Must(template.New("struct_setter").Parse(`
{{ $r := .Recv }}
{{ $type := .Type}}
{{- range .Fields}}
func ({{$r}} {{$type}}) {{.Setter}}(value {{.Type}}) {{$type}} {
	{{$r}}Replace := changes.Replace{ {{.ToValue $r .Name}}, {{.ToValue "value" ""}}}
	{{$r}}Change := changes.PathChange{[]interface{}{"{{.Key}}"}, {{$r}}Replace}
	return {{$r}}.Apply(nil, {{$r}}Change).({{$type}})
}
{{end -}}
`))
