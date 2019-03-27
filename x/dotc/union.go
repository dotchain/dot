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
	return unionApply.Execute(w, s)
}

// GenerateSetters generates the code for the field setters
func (s Union) GenerateSetters(w io.Writer) error {
	return unionSetters.Execute(w, s)
}

var unionApply = template.Must(template.New("union_apply").Parse(`
{{ $r := .Recv }}
func ({{$r}} {{.Type}}) get(key interface{}) changes.Value {
	switch key {
	case "_heap_": return {{$r}}.activeKeyHeap
	{{range .Fields}}
	case "{{.Key}}": return {{.ToValue $r .Name}}{{end}}
        }
	panic(key)
}

func ({{$r}} {{.Type}}) set(key interface{}, v changes.Value) changes.Value {
	{{$r}}Clone := {{if .Pointer}}*{{end}}{{$r}}
	switch key {
	case "_heap_": {{$r}}Clone.activeKeyHeap = v.(heap.Heap)
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

var unionSetters = template.Must(template.New("union_setter").Parse(`
func ({{.Recv}} {{.Type}}) activeKey() string {
	result := ""
	// fetch the largest ranked key => latest
	{{.Recv}}.activeKeyHeap.Iterate(func(key interface{}, _ int) bool {
		if s, ok := key.(string); ok {
			result = s
		}
		return false
	})
	return result
}

func ({{.Recv}} {{.Type}}) maxRank() int {
	rank := -1
	// fetch the largest rank
	{{.Recv}}.activeKeyHeap.Iterate(func(_ interface{}, r int) bool {
		rank = r
		return false
	})
	return rank
}

{{- range $f := .Fields}}
func ({{$.Recv}} {{$.Type}}) {{$f.Setter}}(value {{$f.Type}}) {{$.Type}} {
	rank := {{$.Recv}}.maxRank() + 1
	h := {{$.Recv}}.activeKeyHeap.Update("{{.Key}}", rank)
	return {{$.Ctor}}{activeKeyHeap: h, {{.Name}}: value}
}

func ({{$.Recv}} {{$.Type}}) {{$f.Getter}}() ({{$f.Type}}, bool) {
	var result {{$f.Type}}
	if {{$.Recv}}.activeKey() != "{{.Key}}" {
		return result, false
	}
	return {{$.Recv}}.{{.Name}}, true
}
{{end -}}
`))
