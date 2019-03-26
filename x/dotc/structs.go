// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

// StructStream implements code generation for streams of structs
type StructStream Struct

// PublicFields lists all fields which support streams
func (s StructStream) PublicFields() []Field {
	result := []Field{}
	for _, f := range s.Fields {
		if !f.Atomic || streamTypes[f.Type] != "" {
			result = append(result, f)
		}
	}
	return result
}

// GenerateStream generates the stream implementation
func (s StructStream) GenerateStream(w io.Writer) error {
	return structStreamImpl.Execute(w, s)
}

// GenerateStreamTests generates the stream tests
func (s StructStream) GenerateStreamTests(w io.Writer) error {
	return structStreamTests.Execute(w, s)
}

var streamFns = template.FuncMap{"stream": streamType}

var structStreamImpl = template.Must(template.New("struct_stream_impl").Funcs(streamFns).Parse(`
// {{stream .Type}} implements a stream of {{.Type}} values
type {{stream .Type}} struct {
	Stream streams.Stream
	Value {{.Type}}
}

// Next returns the next entry in the stream if there is one
func (s *{{stream .Type}}) Next() (*{{stream .Type}}, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).({{.Type}}); ok {
		return &{{stream .Type}}{Stream: next, Value: nextVal}, nextc
	}
	return &{{stream .Type}}{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *{{stream .Type}}) Latest() *{{stream .Type}} {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *{{stream .Type}}) Update(val {{.Type}}) *{{stream .Type}} {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &{{stream .Type}}{Stream: nexts, Value: val}
	}
	return s
}

{{ $stype := stream .Type}}
{{range .PublicFields -}}
func (s *{{$stype}}) {{.Name}}() *{{stream .Type}} {
	return &{{stream .Type}}{Stream: streams.Substream(s.Stream, "{{.Key}}"), Value: {{.Unstringify}}(s.Value.{{.Name}})}
}
{{end -}}
`))

var structStreamTests = template.Must(template.New("struct_stream_test").Funcs(streamFns).Parse(`
func Test{{stream .Type}}(t *testing.T) {
	s := streams.New()
	values := valuesFor{{stream .Type}}()
	strong := &{{stream .Type}}{Stream: s, Value: values[0]}

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

{{ $stype := stream .Type}}
{{range .PublicFields -}}
func Test{{$stype}}{{.Name}}(t *testing.T) {
	s := streams.New()
	values := valuesFor{{$stype}}()
	strong := &{{$stype}}{Stream: s, Value: values[0]}
	if !reflect.DeepEqual(strong.Value.{{.Name}}, strong.{{.Name}}().Value) {
		t.Error("Substream returned unexpected value", strong.{{.Name}}().Value)
	}
}
{{end -}}
`))
