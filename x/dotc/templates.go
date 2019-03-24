// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import "text/template"

var infoTpl = template.Must(template.New("imports").Parse(`
// Generated.  DO NOT EDIT.
package {{.Package}}

import (
	{{range .Imports}}{{index . 0}} "{{index . 1}}"
	{{end -}}
)
`))

var structApply = template.Must(template.New("struct").Parse(`
{{ $r := .Recv }}
func ({{$r}} {{.Type}}) get(key interface{}) changes.Value {
	switch key {
	{{- range .Fields}}
	case "{{.Key}}":
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
