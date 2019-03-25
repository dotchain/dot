// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

// StructStream has the type information of a struct for code generation of
// the corresponding Stream
type StructStream Struct

// PublicFields lists all the exported fields of the underlying struct
func (s StructStream) PublicFields() []Field {
	result := []Field{}
	for _, f := range s.Fields {
		if isExported(f.Name) {
			result = append(result, f)
		}
	}
	return result
}

// GenerateStream generates the stream implementation
func (s StructStream) GenerateStream(w io.Writer) error {
	return structStreamImpl.Execute(w, s)
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
	return &{{stream .Type}}{Stream: streams.Substream(s.Stream, "{{.Key}}"), Value: s.Value.{{.Name}}}
}
{{end -}}
`))
