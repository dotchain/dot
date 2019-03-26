// Generated.  DO NOT EDIT.
package myslice

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

func TestStreamMySliceStream(t *testing.T) {
	s := streams.New()
	values := valuesForMySliceStream()
	strong := &MySliceStream{Stream: s, Value: values[0]}

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
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

func TestStreamMySliceStreamSplice(t *testing.T) {
	s := streams.New()
	values := valuesForMySliceStream()
	strong := &MySliceStream{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStreamMySliceStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesForMySliceStream()
	strong := &MySliceStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}

func TestStreammySlice2Stream(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2Stream()
	strong := &mySlice2Stream{Stream: s, Value: values[0]}

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
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

func TestStreammySlice2StreamSplice(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2Stream()
	strong := &mySlice2Stream{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStreammySlice2StreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2Stream()
	strong := &mySlice2Stream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}

func TestStreammySlice3Stream(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3Stream()
	strong := &mySlice3Stream{Stream: s, Value: values[0]}

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
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

func TestStreammySlice3StreamSplice(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3Stream()
	strong := &mySlice3Stream{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStreammySlice3StreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3Stream()
	strong := &mySlice3Stream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}

func TestStreamMySlicePStream(t *testing.T) {
	s := streams.New()
	values := valuesForMySlicePStream()
	strong := &MySlicePStream{Stream: s, Value: values[0]}

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
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

func TestStreamMySlicePStreamSplice(t *testing.T) {
	s := streams.New()
	values := valuesForMySlicePStream()
	strong := &MySlicePStream{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), *values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStreamMySlicePStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesForMySlicePStream()
	strong := &MySlicePStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (*values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}

func TestStreammySlice2PStream(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2PStream()
	strong := &mySlice2PStream{Stream: s, Value: values[0]}

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
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

func TestStreammySlice2PStreamSplice(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2PStream()
	strong := &mySlice2PStream{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), *values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStreammySlice2PStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2PStream()
	strong := &mySlice2PStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (*values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}

func TestStreammySlice3PStream(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3PStream()
	strong := &mySlice3PStream{Stream: s, Value: values[0]}

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
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: values[3]})
	if _, c = strong.Next(); c != nil {
		t.Error("Unexpected change on terminated stream", c)
	}
}

func TestStreammySlice3PStreamSplice(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3PStream()
	strong := &mySlice3PStream{Stream: s, Value: values[1]}
	strong1 := strong.Splice(0, strong.Value.Count(), *values[2]...)
	if !reflect.DeepEqual(strong1.Value, values[2]) {
		t.Error("Splice did the unexpected", strong1.Value)
	}
}

func TestStreammySlice3PStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3PStream()
	strong := &mySlice3PStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (*values[1])[0]) {
		t.Error("Splice did the unexpected", item0.Value)
	}
}
