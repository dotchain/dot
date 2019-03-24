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

var structApply = template.Must(template.New("struct_apply").Parse(`
{{ $r := .Recv }}
func ({{$r}} {{.Type}}) get(key interface{}) changes.Value {
	switch key {
	{{range .Fields}}
	case "{{.Key}}": return {{.WrapR $r}}{{end}}
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
	{{$r}}Replace := changes.Replace{ {{.WrapR $r}}, {{.Wrap "value"}}}
	{{$r}}Change := changes.PathChange{[]interface{}{"{{.Key}}"}, {{$r}}Replace}
	return {{$r}}.Apply(nil, {{$r}}Change).({{$type}})
}
{{end -}}
`))

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

var sliceApply = template.Must(template.New("slice_apply").Parse(`
func ({{.Recv}} {{.Type}}) get(key interface{}) changes.Value {
	return {{.WrapR .Recv}}
}

func ({{.Recv}} {{.Type}}) set(key interface{}, v changes.Value) changes.Value {
	{{.Recv}}Clone := append({{.Type}}(nil), {{.Recv}}...)
	{{.Recv}}Clone[key.(int)] = {{.Unwrap}}
	return {{.Recv}}Clone
}

func ({{.Recv}} {{.Type}}) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	return append(append({{.Recv}}[:offset:offset], after.({{.Type}})...), {{.Recv}}[end:]...)
}

// Slice implements changes.Collection Slice() method
func ({{.Recv}} {{.Type}}) Slice(offset, count int) changes.Collection {
	return {{.Recv}}[offset:offset+count]
}

// Count implements changes.Collection Count() method
func ({{.Recv}} {{.Type}}) Count() int {
	return len({{.Recv}})
}

func ({{.Recv}} {{.Type}}) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: {{.Recv}}.get, Set: {{.Recv}}.set, Splice: {{.Recv}}.splice}).Apply(ctx, c, {{.Recv}})
}

func ({{.Recv}} {{.Type}}) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: {{.Recv}}.get, Set: {{.Recv}}.set, Splice: {{.Recv}}.splice}).ApplyCollection(ctx, c, {{.Recv}})
}

`))

var sliceSetters = template.Must(template.New("slice_setter").Parse(`
func ({{.Recv}} {{.Type}}) Splice(offset, count int, insert ...{{.ElemType}}) {{.Type}} {
	end := offset + count
	return append(append({{.Recv}}[:offset:offset], insert...), {{.Recv}}[end:]...)
}

`))
