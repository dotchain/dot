// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

// Slice has the type information of a slice for code generation of
// the Apply(), ApplyCollection() and Splice() methods
type Slice struct {
	Recv, Type, ElemType string
	Atomic               bool
}

// RawType returns the non-pointer inner type if Type is a pointer
func (s Slice) RawType() string {
	if s.Pointer() {
		return s.Type[1:]
	}
	return s.Type
}

// Elem returns a field version of the element
func (s Slice) Elem() Field {
	return Field{Type: s.ElemType, Atomic: s.Atomic}
}

// Item returns s.Recv[index]
func (s Slice) Item(index string) string {
	if s.Pointer() {
		return "(*" + s.Recv + ")[" + index + "]"
	}
	return s.Recv + "[" + index + "]"
}

// Pointer checks if the slice type is a pointer
func (s Slice) Pointer() bool {
	return s.Type[0] == '*'
}

// GenerateApply generates the code for the changes.Value Apply() method
// and the ApplyCollection() method
func (s Slice) GenerateApply(w io.Writer) error {
	return sliceApply.Execute(w, s)
}

// GenerateSetters generates the code for the field setters
func (s Slice) GenerateSetters(w io.Writer) error {
	return sliceSetters.Execute(w, s)
}

var sliceApply = template.Must(template.New("slice_apply").Parse(`
func ({{.Recv}} {{.Type}}) get(key interface{}) changes.Value {
	return {{.Elem.ToValue (.Item "key.(int)") ""}}
}

func ({{.Recv}} {{.Type}}) set(key interface{}, v changes.Value) changes.Value {
	{{.Recv}}Clone := {{.RawType}}(append([]{{.ElemType}}(nil), ({{if .Pointer}}*{{end}}{{.Recv}})...))
	{{.Recv}}Clone[key.(int)] = {{.Elem.FromValue "v" ""}}
	return {{if .Pointer}}&{{end}}{{.Recv}}Clone
}

func ({{.Recv}} {{.Type}}) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	{{.Recv}}Val := {{if .Pointer}}*{{end}}{{.Recv}}
	afterVal := {{if .Pointer}}*{{end}}(after.({{.Type}}))
	{{.Recv}}New := append(append({{.Recv}}Val[:offset:offset], afterVal...), {{.Recv}}Val[end:]...)
	return {{if .Pointer}}&{{end}}{{.Recv}}New
}

// Slice implements changes.Collection Slice() method
func ({{.Recv}} {{.Type}}) Slice(offset, count int) changes.Collection {
	{{.Recv}}Slice := ({{if .Pointer}}*{{end}}{{.Recv}})[offset:offset+count]
	return {{if .Pointer}}&{{end}}{{.Recv}}Slice
}

// Count implements changes.Collection Count() method
func ({{.Recv}} {{.Type}}) Count() int {
	return len({{if .Pointer}}*{{end}}{{.Recv}})
}

func ({{.Recv}} {{.Type}}) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: {{.Recv}}.get, Set: {{.Recv}}.set, Splice: {{.Recv}}.splice}).Apply(ctx, c, {{.Recv}})
}

func ({{.Recv}} {{.Type}}) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: {{.Recv}}.get, Set: {{.Recv}}.set, Splice: {{.Recv}}.splice}).ApplyCollection(ctx, c, {{.Recv}})
}

`))

var sliceSetters = template.Must(template.New("slice_setter").Parse(`
// Splice replaces [offset:offset+count] with insert...
func ({{.Recv}} {{.Type}}) Splice(offset, count int, insert ...{{.ElemType}}) {{.Type}} {
	{{.Recv}}Insert := {{.RawType}}(insert)
	return {{.Recv}}.splice(offset, count, {{if .Pointer}}&{{end}}{{.Recv}}Insert).({{.Type}})
}

// Move shuffles [offset:offset+count] by distance.
func ({{.Recv}} {{.Type}}) Move(offset, count, distance int) {{.Type}} {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return {{.Recv}}.ApplyCollection(nil, c).({{.Type}})
}

`))
