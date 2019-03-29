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

func TestStreamMySliceStreamMove(t *testing.T) {
	s := streams.New()
	values := valuesForMySliceStream()
	strong := &MySliceStream{Stream: s, Value: values[1]}
	v2 := values[2]
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

func TestStreamMySliceStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesForMySliceStream()
	strong := &MySliceStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := (strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
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

func TestStreammySlice2StreamMove(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2Stream()
	strong := &mySlice2Stream{Stream: s, Value: values[1]}
	v2 := values[2]
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

func TestStreammySlice2StreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2Stream()
	strong := &mySlice2Stream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := (strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
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

func TestStreammySlice3StreamMove(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3Stream()
	strong := &mySlice3Stream{Stream: s, Value: values[1]}
	v2 := values[2]
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

func TestStreammySlice3StreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3Stream()
	strong := &mySlice3Stream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := (strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
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

func TestStreamMySlicePStreamMove(t *testing.T) {
	s := streams.New()
	values := valuesForMySlicePStream()
	strong := &MySlicePStream{Stream: s, Value: values[1]}
	v2 := *values[2]
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

func TestStreamMySlicePStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesForMySlicePStream()
	strong := &MySlicePStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (*values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := *values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := (*strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
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

func TestStreammySlice2PStreamMove(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2PStream()
	strong := &mySlice2PStream{Stream: s, Value: values[1]}
	v2 := *values[2]
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

func TestStreammySlice2PStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice2PStream()
	strong := &mySlice2PStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (*values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := *values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := (*strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
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

func TestStreammySlice3PStreamMove(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3PStream()
	strong := &mySlice3PStream{Stream: s, Value: values[1]}
	v2 := *values[2]
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

func TestStreammySlice3PStreamItem(t *testing.T) {
	s := streams.New()
	values := valuesFormySlice3PStream()
	strong := &mySlice3PStream{Stream: s, Value: values[1]}
	item0 := strong.Item(0)
	if !reflect.DeepEqual(item0.Value, (*values[1])[0]) {
		t.Error("Item() did the unexpected", item0.Value)
	}

	for kk := range values {
		item := *values[kk]
		l := len(item) - 1
		if l < 0 {
			continue
		}
		item0 = item0.Update(item[l])
		if !reflect.DeepEqual(item0.Value, item[l]) {
			t.Error("Update did not take effect", item0.Value, item[l])
		}
		strong = strong.Latest()
		v := (*strong.Value)[0]
		if !reflect.DeepEqual(v, item[l]) {
			t.Error("Update did not take effect", v, item[l])
		}
	}

	v := strong.Value.ApplyCollection(nil, changes.Splice{Before: strong.Value.Slice(0, 1), After: strong.Value.Slice(0, 0)})
	if !reflect.DeepEqual(v.Slice(0, 0), v) {
		t.Error("Could not slice away item", v)
	}
}
