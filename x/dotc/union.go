// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

// Union has the type information of a union for code generation of
// the Apply() and SetField(..) methods
type Union struct {
	Recv, Type string
	Fields     []Field
}

// Pointer specifies if the union type is itself a pointer
func (s Union) Pointer() bool {
	return s.Type[0] == '*'
}

// Ctor returns the type used to create this
func (s Union) Ctor() string {
	if s.Type[0] == '*' {
		return "&" + s.Type[1:]
	}
	return s.Type
}

// GenerateApply generates the code for the changes.Value Apply() method
func (s Union) GenerateApply(w io.Writer) error {
	// union also uses structApply
	return structApply.Execute(w, s)
}

// GenerateSetters generates the code for the field setters
func (s Union) GenerateSetters(w io.Writer) error {
	return unionSetters.Execute(w, s)
}

var unionSetters = template.Must(template.New("union_setter").Parse(`
{{ $r := .Recv }}
{{ $type := .Type}}
{{ $ctor := .Ctor}}
{{ $ptr := .Pointer}}
{{- range .Fields}}
func ({{$r}} {{$type}}) {{.Setter}}(value {{.Type}}) {{$type}} {
	return {{$ctor}}{ {{.Name}}: value}
}
{{end -}}
`))
