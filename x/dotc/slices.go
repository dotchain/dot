// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"io"
	"text/template"
)

// SliceStream implements code generation for streams of slices
type SliceStream Slice

// GenerateStream generates the stream implementation
func (s SliceStream) GenerateStream(w io.Writer) error {
	if !isExported(s.Type) {
		return nil
	}

	return sliceStreamImpl.Execute(w, Slice(s))
}

// GenerateStreamTests generates the stream tests
func (s SliceStream) GenerateStreamTests(w io.Writer) error {
	if !isExported(s.Type) {
		return nil
	}

	return sliceStreamTests.Execute(w, Slice(s))
}

var sliceStreamImpl = template.Must(template.New("slice_stream_impl").Funcs(streamFns).Parse(`
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

// Item returns the sub item stream
func (s *{{stream .Type}}) Item(index int) *{{stream .ElemType}} {
	return &{{stream .ElemType}}{Stream: streams.Substream(s.Stream, index), Value: ({{if .Pointer}}*{{end}}s.Value)[index]}
}

// Splice splices the items
func (s *{{stream .Type}}) Splice(offset, count int, replacement ...{{.ElemType}}) *{{stream .Type}} {
	after := {{.RawType}}(replacement)
	c := changes.Replace{Before: s.Value.Slice(offset, count), After: {{if .Pointer}}&{{end}}after}
	str := s.Stream.Append(c)
	return &{{stream .Type}}{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}
`))

var sliceStreamTests = template.Must(template.New("slice_stream_tests").Funcs(streamFns).Parse(`
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

func Test{{stream .Type}}Splice(t *testing.T) {
	s := streams.New()
	values := valuesFor{{stream .Type}}()
	strong := &{{stream .Type}}{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), {{if .Pointer}}*{{end}}values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func Test{{stream .Type}}Item(t *testing.T) {
	s := streams.New()
	values := valuesFor{{stream .Type}}()
	strong := &{{stream .Type}}{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, ({{if .Pointer}}*{{end}}values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}

`))
