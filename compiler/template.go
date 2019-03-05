// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package compiler

import (
	"regexp"
	"text/template"
)

var headerTpl = newTemplate(`
// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
//

package Package

import (
  "github.com/dotchain/dot/changes"
  "github.com/dotchain/dot/streams"
  "github.com/dotchain/dot/refs"
  {{range $import := .Imports}}{{index $import 0}} "{{index $import 1}}"
{{end}}
)
`)

var contextTpl = newTemplate(`
type {{.ContextType}} struct {
	streams.Cache
	{{range $ps := .PkgSubcomps}}{{if not (eq $ps.Pkg "")}}{{$ps.Pkg}} struct { {{end}}
		{{range $ps.Subcomps}}{{.}}
	{{end}} {{if not (eq $ps.Pkg "")}} } {{end}}
	{{end}}
}
{{ $c := index .Args 0 }}
func ({{$c.Name}} {{$c.Type}}) areArgsSame({{.NonContextArgsDecl}}) bool {
	{{range $i, $a := .Args}}{{if not (eq $i 0)}}
	{{if $a.IsArray}}
        if len({{$a.Name}}) != len({{$c.Name}}.memoized.{{$a.Name}}) {
		return false
	}
	for {{$a.Name}}Idx := range {{$a.Name}} {
		if {{$a.Name}}[{{$a.Name}}Idx] != {{$c.Name}}.memoized.{{$a.Name}}[{{$a.Name}}.Idx] {
			return false
		}
	}
	{{if $a.IsLast}}return true{{end}}
	{{else}}
		{{if $a.IsLast}}return {{$a.Name}} == {{$c.Name}}.memoized.{{$a.Name}}
		{{else}}if {{$a.Name}} != {{$c.Name}}.memoized.{{$a.Name}} { return false }
	{{end}}{{end}}{{end}}{{end}}
}

func ({{$c.Name}} {{$c.Type}}) refreshIfNeeded({{.NonContextArgsDecl}}) ({{.ResultsDecl}}) {
	if !{{$c.Name}}.initialized || !{{$c.Name}}.areArgsSame({{.NonContextArgs}}) {
		return {{$c.Name}}.refresh({{.NonContextArgs}})
	}
	return {{.MemoizedNonStateResults}}
}

func ({{$c.Name}} {{$c.Type}}) refresh({{.NonContextArgsDecl}}) ({{.ResultsDecl}}) {
	{{$c.Name}}.initialized = true
	{{$c.Name}}.stateHandler.Handle = func() {
		{{$c.Name}}.refresh({{.NonContextArgs}})
	}
	{{range $sa := .StateArgs}}
	if {{$c.Name}}.memoized.{{$sa.Name}} != nil {
		{{$c.Name}}.memoized.{{$sa.Name}} = {{$c.Name}}.memoized.{{$sa.Name}}.Latest()
	}{{end}}
	{{.MemoizedNonContextArgs}} = {{.NonContextArgs}}

	{{$c.Name}}.Cache.Begin()
	defer {{$c.Name}}.Cache.End()

	{{- range $i := .Subcomponents}}
	{{$c.Name}}.{{$i}}.Begin()
	defer {{$c.Name}}.{{$i}}.End()
	{{end -}}

	{{.MemoizedResults}} = {{.Function}}({{.AllArgs}})

	{{range $sa := .StateArgs -}}
	if {{$c.Name}}.memoized.{{$sa.Name}} != {{$c.Name}}.memoized.{{$sa.ResultName}} {
		if {{$c.Name}}.memoized.{{$sa.Name}} != nil {
			{{$c.Name}}.memoized.{{$sa.Name}}.Off(&{{$c.Name}}.stateHandler)
		}
		if {{$c.Name}}.memoized.{{$sa.ResultName}} != nil {
			{{$c.Name}}.memoized.{{$sa.ResultName}}.On(&{{$c.Name}}.stateHandler)
		}
		{{$c.Name}}.memoized.{{$sa.Name}} = {{$c.Name}}.memoized.{{$sa.ResultName}}
	}
	{{end -}}
	return {{.MemoizedNonStateResults}}
}

func ({{$c.Name}} {{$c.Type}}) close() {
	{{$c.Name}}.Cache.Begin()
	defer {{$c.Name}}.Cache.End()

	{{- range $i := .Subcomponents}}
	{{$c.Name}}.{{$i}}.Begin()
	defer {{$c.Name}}.{{$i}}.End()
	{{end -}}

	{{range $sa := .StateArgs -}}
	if {{$c.Name}}.memoized.{{$sa.Name}} != nil {
		{{$c.Name}}.memoized.{{$sa.Name}}.Off(&{{$c.Name}}.stateHandler)
	}
	{{end -}}
}

{{.ComponentComments}}
type {{.Component}} struct {
	old, current map[interface{}]{{$c.Type}}
}

// Begin starts a round
func (c *{{.Component}}) Begin() {
	c.old, c.current  = c.current, map[interface{}]{{$c.Type}}{}
}

// End finishes the round cleaning up any unused components
func (c *{{.Component}}) End() {
	for _, ctx := range c.old {
		ctx.close()
	}
	c.old = nil
}

{{.MethodComments}}
func (c *{{.Component}}) {{.Method}}({{.MethodDecl}}) ({{.ResultsDecl}}) {
	{{$c.Name}}Old, ok := c.old[{{$c.Name}}Key]
	if ok {
		delete(c.old, {{$c.Name}}Key)
	} else {
		{{$c.Name}}Old = &{{.ContextType}}{}
	}
	c.current[{{$c.Name}}Key] =  {{$c.Name}}Old
	return {{$c.Name}}Old.refreshIfNeeded({{.NonContextArgs}})
}

`)

var streamTpl = newTemplate(`
// StreamType is a stream of ValueType values.
type StreamType struct {
	// Notifier provides On/Off/Notify support. New instances of
	// StreamType created via the AppendLocal or AppendRemote
	// share the same Notifier value.
	*streams.Notifier
	
	// Value holds the current value. The latest value may be
	// fetched via the Latest() method.
	Value ValueType
	
	// Change tracks the change that leads to the next value.
	Change changes.Change

	// Next tracks the next value in the stream.
	Next *StreamType
}


// NewStreamType creates a new ValueType stream 
func NewStreamType(value ValueType) *StreamType {
	return &StreamType{&streams.Notifier{}, value, nil, nil}
}

// Latest returns the latest value in the stream
func (s *StreamType) Latest() *StreamType {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Append appends a local change. isLocal identifies if the caller is
// local or remote. It returns the updated stream whose value matches
// the provided value and whose Latest() converges to the latest of
// the stream.  
func (s *StreamType) Append(c changes.Change, value ValueType, isLocal bool) *StreamType {
	if c == nil {
		c = changes.Replace{Before: s.wrapValue(s.Value), After: s.wrapValue(value)}
	}

	// return value: after is correctly set to provided value
	result := &StreamType{Notifier: s.Notifier, Value: value}

	// before tracks s, after tracks result, v tracks latest value
	// of after chain
	before := s
	var v changes.Value = changes.Atomic{value}

	// walk the chain of Next and find corresponding values to
	// add to after so that both s annd after converge
	after := result
	for ; before.Next != nil; before = before.Next {
		var afterChange changes.Change

		if isLocal {
			c, afterChange = before.Change.Merge(c)
		} else {
			afterChange, c = c.Merge(before.Change)
		}
		
		if c == nil {
			// the convergence point is before.Next
			after.Change, after.Next = afterChange, before.Next
			return result
		}

		if afterChange == nil {
			continue
		}
		
		// append this to after and continue with that
		v = v.Apply(nil, afterChange)
		after.Change  = afterChange
		after.Next = &StreamType{Notifier: s.Notifier, Value: s.unwrapValue(v)}
		after = after.Next
	}

	// append the residual change (c) to converge to wherever
	// after has landed. Notify since s.Latest() has now changed
	before.Change, before.Next = c, after
	s.Notify()
	return result
}

func (s *StreamType) wrapValue(i interface{}) changes.Value {
	if x, ok := i.(changes.Value); ok {
		return x
	}
	return changes.Atomic{i}
}

func (s *StreamType) unwrapValue(v changes.Value) ValueType {
	if x, ok := v.(interface{}).(ValueType); ok {
		return x
	}
	return v.(changes.Atomic).Value.(ValueType)
}
`)

var fieldTpl = newTemplate(`

// SetField updates the field with a new value
func (s *StreamType) SetField(v FieldType) *StreamType {
	c := changes.Replace{s.wrapValue(s.Value.Field), s.wrapValue(v)}
	value := s.Value
	value.Field = v
	key := []interface{}{"Field"}
	return s.Append(changes.PathChange{key, c}, value, true)
}

// FieldSubstream returns a stream for Field that is automatically
// connected to the current StreamType instance.  Updates to one
// automatically update the other.
func (s *StreamType) FieldSubstream(cache streams.Cache) (field *FieldStrType) {
	n := s.Notifier
	handler := &streams.Handler{nil}
	if f, h, ok := cache.GetSubstream(n, "Field"); ok {
		field, handler  = f.(*FieldStrType), h
	} else {
		field = NewFieldStrType(s.Value.Field)
		parent, merging, path := s, false, []interface{}{"Field"}
		handler.Handle = func() {
			if merging {
				return
			}

			merging = true
			for ;field.Next != nil; field = field.Next {
				v := parent.Value
				v.Field = field.Next.Value
				c := changes.PathChange{path, field.Change}
				parent = parent.Append(c, v, true)
			}

			for ;parent.Next != nil; parent = parent.Next {
				result := refs.Merge(path, parent.Change)
				if result ==  nil {
					field = field.Append(nil, parent.Next.Value.Field, true)
				}  else {
					field = field.Append(result.Affected, parent.Next.Value.Field, true)
				}
			}
			merging = false
		}
		field.On(handler)
		parent.On(handler)
	}

	handler.Handle()
	field = field.Latest()
	n2 := field.Notifier
	close := func() { n.Off(handler); n2.Off(handler); }
	cache.SetSubstream(n, "Field", field, handler, close)
	return field
}
`)

var entryTpl = newTemplate(`

// Substream returns a stream for the specified index that is
// automatically connected to the current StreamType instance.  Updates to
// one automatically update the other.
func (s *StreamType) Substream(cache streams.Cache, index int) (entry *EntryStrType) {
	n := s.Notifier
	handler := &streams.Handler{nil}
	if f, h, ok := cache.GetSubstream(n, index); ok {
		entry, handler  = f.(*EntryStrType), h
	} else {
		entry = NewEntryStrType(s.Value[index])
		parent, merging, path := s, false, []interface{}{index}
		handler.Handle = func() {
			if merging {
				return
			}

			merging = true
			for ;entry.Next != nil; entry = entry.Next {
				v := append(ValueType(nil), parent.Value...)
				v[index] = entry.Next.Value
				c := changes.PathChange{path, entry.Change}
				parent = parent.Append(c, v, true)
			}

			for ;parent.Next != nil; parent = parent.Next {
				result := refs.Merge(path, parent.Change)
				var c changes.Change
				if result !=  nil {
					index = result.P[0].(int)
                                        // TODO: if the index changed fix up
                                        // the key in the cache
					c = result.Affected
				}
				entry = entry.Append(c, parent.Next.Value[index], true)
			}
			merging = false
		}
		entry.On(handler)
		parent.On(handler)
	}

	handler.Handle()
	entry = entry.Latest()
	n2 := entry.Notifier
	close := func() { n.Off(handler); n2.Off(handler); }
	cache.SetSubstream(n, index, entry, handler, close)
	return entry
}
`)

func newTemplate(s string) *template.Template {
	replacements := [][2]string{
		{"Package", "{{$.Package}}"},
		{"StreamType", "{{$.StreamType}}"},
		{"ValueType", "{{$.ValueType}}"},
		{"EntryStrType", "{{$.EntryStreamType}}"},
		{"NewEntryStrType", "{{$.EntryConstructor}}"},
		{"NewFieldStrType", "{{$.FieldConstructor}}"},
		{"FieldStrType", "{{$.FieldStreamType}}"},
		{"FieldSubstream", "{{$.FieldSubstream}}"},
		{"FieldType", "{{$.FieldType}}"},
	}

	for _, rr := range replacements {
		s = regexp.MustCompile(rr[0]).ReplaceAllString(s, rr[1])
	}
	repl := func(s string) string {
		return "{{$.Field}}" + s[len(s)-1:]
	}
	s = regexp.MustCompile("Field[^CST]").ReplaceAllStringFunc(s, repl)

	return template.Must(template.New("code").Parse(s))
}
