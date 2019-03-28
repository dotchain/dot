// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

// UnionStream implements stream functionality for unions
type UnionStream Union

// StreamType provides the stream type of the struct
func (s UnionStream) StreamType() string {
	return (Field{Type: s.Type}).ToStreamType()
}

// Pointer specifies if the struct type is itself a pointer
func (s UnionStream) Pointer() bool {
	return s.Type[0] == '*'
}

// GenerateStream generates the stream implementation
func (s UnionStream) GenerateStream(w io.Writer) error {
	return unionStreamImpl.Execute(w, s)
}

// GenerateStreamTests generates the stream tests
func (s UnionStream) GenerateStreamTests(w io.Writer) error {
	return unionStreamTests.Execute(w, s)
}

var unionStreamImpl = template.Must(template.New("union_stream_impl").Parse(`
// {{.StreamType}} implements a stream of {{.Type}} values
type {{.StreamType}} struct {
	Stream streams.Stream
	Value {{.Type}}
}

// Next returns the next entry in the stream if there is one
func (s *{{.StreamType}}) Next() (*{{.StreamType}}, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).({{.Type}}); ok {
		return &{{.StreamType}}{Stream: next, Value: nextVal}, nextc
	}
	return &{{.StreamType}}{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *{{.StreamType}}) Latest() *{{.StreamType}} {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *{{.StreamType}}) Update(val {{.Type}}) *{{.StreamType}} {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &{{.StreamType}}{Stream: nexts, Value: val}
	}
	return s
}

func (s *{{.StreamType}}) transformer() func(changes.Change) changes.Change {
	h := ({{if .Pointer}}*{{end}}s.Value).activeKeyHeap
	var xform func(changes.Change) changes.Change
	p := []interface{}{"_heap_"}

	maxRank := func() int {
		result := -1
		h.Iterate(func(_ interface{}, r int) bool {
			result = r
			return false
		})
		return result
	}

	xform = func(c changes.Change) changes.Change {
	switch c := c.(type) {
	case changes.ChangeSet:
		result := make(changes.ChangeSet, 0, len(c))
		for _, cx := range c {
			if cx = xform(cx); cx != nil {
				result = append(result, cx)
			}
		}
		return result
	case changes.PathChange:
		if len(c.Path) == 0 {
			return xform(c.Change)
		}
		if c.Path[0] != p[0] {
			cx := h.UpdateChange(c.Path[0], maxRank()+1)
			h = h.Update(c.Path[0], maxRank()+1)
			return changes.ChangeSet{changes.PathChange{Path: p, Change: cx}, c}
		}
	}
	return c
	}
	return xform
}


{{ $stype := .StreamType}}
{{range .Fields -}}
func (s *{{$stype}}) {{.Name}}() *{{.ToStreamType}} {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &{{.ToStreamType}}{Stream: streams.Substream(stream, "{{.Key}}"), Value: {{.FromStreamValue "s.Value" .Name}} }
}
{{end -}}
`))

var unionStreamTests = template.Must(template.New("union_stream_test").Parse(`
func TestStream{{.StreamType}}(t *testing.T) {
	s := streams.New()
	values := valuesFor{{.StreamType}}()
	strong := &{{.StreamType}}{Stream: s, Value: values[0]}

	strong = strong.Update(values[1])
	if !reflect.DeepEqual(strong.Value, values[1]) {
		t.Error("Update did not change value", strong.Value)
	}

	s, c := s.Next()
	if !reflect.DeepEqual(c, changes.Replace{Before: values[0], After: values[1]}) {
		t.Error("Unexpected change", c)
	}

	c = changes.Replace{Before: values[1], After: values[2]}
	s = s.Append(c)
	c = changes.Replace{Before: values[2], After: values[3]}
	s = s.Append(c)
	strong = strong.Latest()

	if !reflect.DeepEqual(strong.Value, values[3]) {
		t.Error("Unexpected value", strong.Value)
	}

	_, c = strong.Next()
	if c != nil {
		t.Error("Unexpected change on stream", c)
	}

	s = s.Append(changes.Replace{Before: values[3], After: changes.Nil})
	if strong, c = strong.Next();  c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next();  c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

{{ $stype := .StreamType}}
{{range .Fields -}}
func TestStream{{$stype}}{{.Name}}(t *testing.T) {
	s := streams.New()
	values := valuesFor{{$stype}}()
	strong := &{{$stype}}{Stream: s, Value: values[0]}
	expected := {{.FromStreamValue "strong.Value" .Name}}
	if !reflect.DeepEqual(expected, strong.{{.Name}}().Value) {
		t.Error("Substream returned unexpected value", strong.{{.Name}}().Value)
	}

	child := strong.{{.Name}}()
	for kk := range values {
		child = child.Update({{.FromStreamValue "values[kk]" .Name}})
		strong = strong.Latest()
		if !reflect.DeepEqual(child.Value, {{.FromStreamValue "values[kk]" .Name}}) {
			t.Error("updating child didn't  take effect", child.Value)
		}
		if !reflect.DeepEqual(child.Value, {{.FromStreamValue "strong.Value" .Name}}) {
			t.Error("updating child didn't  take effect", child.Value)
		}

		{{$current := .}}
		{{range $f := $.Fields}}
		if _, ok := strong.Value.{{.Getter}}(); {{if eq $current.Name $f.Name}}!ok{{else}}ok{{end}} {
			t.Error("Getter failed")
		}
		{{end}}
	}

	v := strong.Value.{{.Setter}}(values[0].{{.Name}})
	if !reflect.DeepEqual(v.{{.Name}}, values[0].{{.Name}}) {
		t.Error("Could not update", "{{.Setter}}")
	}
}
{{end -}}
`))
