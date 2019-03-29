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

// Pointer simply proxies to Slice, needed because of template limitation
func (s SliceStream) Pointer() bool {
	return Slice(s).Pointer()
}

// RawType simply proxies to Slice, needed because of template limitation
func (s SliceStream) RawType() string {
	return Slice(s).RawType()
}

// Elem simply proxies to Slice, needed because of template limitation
func (s SliceStream) Elem() Field {
	return Slice(s).Elem()
}

// StreamType provides the stream type of the struct
func (s SliceStream) StreamType() string {
	return (Field{Type: s.Type}).ToStreamType()
}

// GenerateStream generates the stream implementation
func (s SliceStream) GenerateStream(w io.Writer) error {
	return sliceStreamImpl.Execute(w, s)
}

// GenerateStreamTests generates the stream tests
func (s SliceStream) GenerateStreamTests(w io.Writer) error {
	return sliceStreamTests.Execute(w, s)
}

var sliceStreamImpl = template.Must(template.New("slice_stream_impl").Parse(`
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

// Item returns the sub item stream
func (s *{{.StreamType}}) Item(index int) *{{.Elem.ToStreamType}} {
	return &{{.Elem.ToStreamType}}{Stream: streams.Substream(s.Stream, index), Value: ({{if .Pointer}}*{{end}}s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *{{.StreamType}}) Splice(offset, count int, replacement ...{{.ElemType}}) *{{.StreamType}} {
	after := {{.RawType}}(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: {{if .Pointer}}&{{end}}after}
	str := s.Stream.Append(c)
	return &{{.StreamType}}{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *{{.StreamType}}) Move(offset, count, distance int) *{{.StreamType}} {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &{{.StreamType}}{Stream: str, Value: s.Value.Move(offset, count, distance)}
}
`))

var sliceStreamTests = template.Must(template.New("slice_stream_tests").Parse(`
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

func TestStream{{.StreamType}}Splice(t *testing.T) {
	s := streams.New()
	values := valuesFor{{.StreamType}}()
	strong := &{{.StreamType}}{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), {{if .Pointer}}*{{end}}values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStream{{.StreamType}}Move(t *testing.T) {
	s := streams.New()
	values := valuesFor{{.StreamType}}()
	strong := &{{.StreamType}}{Stream: s, Value: values[1]}
	v2 := {{if .Pointer}}*{{end}}values[2]
	strong1 := strong.Splice(strong.Value.Count(), 0, v2[len(v2)-1])
	strong2 := strong1.Move(0, 1, 1)
	if reflect.DeepEqual(strong1.Value, strong2.Value) {
		t.Error("Move did the unexpected", strong1.Value, strong2.Value)
	}
	strong2 = strong2.Move(0, 1, 1)
	if !reflect.DeepEqual(strong1.Value, strong2.Value) {
		t.Error("Move did the unexpected", strong1.Value, strong2.Value)
	}
}

func TestStream{{.StreamType}}Item(t *testing.T) {
	s := streams.New()
	values := valuesFor{{.StreamType}}()
	strong := &{{.StreamType}}{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, ({{if .Pointer}}*{{end}}values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := {{if .Pointer}}*{{end}}values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := ({{if .Pointer}}*{{end}}strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}		
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
	}
}

`))
