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

// GenerateStream generates the stream implementation
func (s UnionStream) GenerateStream(w io.Writer) error {
	return unionStreamImpl.Execute(w, StructStream(Union(s)))
}

// GenerateStreamTests generates the stream tests
func (s UnionStream) GenerateStreamTests(w io.Writer) error {
	return unionStreamTests.Execute(w, StructStream(Union(s)))
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

{{ $stype := .StreamType}}
{{range .Fields -}}
func (s *{{$stype}}) {{.Name}}() *{{.ToStreamType}} {
	return &{{.ToStreamType}}{Stream: streams.Substream(s.Stream, "{{.Key}}"), Value: {{.FromStreamValue "s.Value" .Name}} }
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
}
{{end -}}
`))
