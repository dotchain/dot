// Generated.  DO NOT EDIT.
package mystruct

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

func TestStreamMyStructStream(t *testing.T) {
	s := streams.New()
	values := valuesForMyStructStream()
	strong := &MyStructStream{Stream: s, Value: values[0]}

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

func TestStreamMyStructStreamboo(t *testing.T) {
	s := streams.New()
	values := valuesForMyStructStream()
	strong := &MyStructStream{Stream: s, Value: values[0]}
	expected := strong.Value.boo
	if !reflect.DeepEqual(expected, strong.boo().Value) {
		t.Error("Substream returned unexpected value", strong.boo().Value)
	}

	child := strong.boo()
	for kk := range values {
		child = child.Update(values[kk].boo)
		strong = strong.Latest()
		if !reflect.DeepEqual(child.Value, values[kk].boo) {
			t.Error("updating child didn't  take effect", child.Value)
		}
		if !reflect.DeepEqual(child.Value, strong.Value.boo) {
			t.Error("updating child didn't  take effect", child.Value)
		}
	}

	v := strong.Value.setBoo(values[0].boo)
	if !reflect.DeepEqual(v.boo, values[0].boo) {
		t.Error("Could not update", "setBoo")
	}
}
func TestStreamMyStructStreamboop(t *testing.T) {
	s := streams.New()
	values := valuesForMyStructStream()
	strong := &MyStructStream{Stream: s, Value: values[0]}
	expected := strong.Value.boop
	if !reflect.DeepEqual(expected, strong.boop().Value) {
		t.Error("Substream returned unexpected value", strong.boop().Value)
	}

	child := strong.boop()
	for kk := range values {
		child = child.Update(values[kk].boop)
		strong = strong.Latest()
		if !reflect.DeepEqual(child.Value, values[kk].boop) {
			t.Error("updating child didn't  take effect", child.Value)
		}
		if !reflect.DeepEqual(child.Value, strong.Value.boop) {
			t.Error("updating child didn't  take effect", child.Value)
		}
	}

	v := strong.Value.setBoop(values[0].boop)
	if !reflect.DeepEqual(v.boop, values[0].boop) {
		t.Error("Could not update", "setBoop")
	}
}
func TestStreamMyStructStreamstr(t *testing.T) {
	s := streams.New()
	values := valuesForMyStructStream()
	strong := &MyStructStream{Stream: s, Value: values[0]}
	expected := strong.Value.str
	if !reflect.DeepEqual(expected, strong.str().Value) {
		t.Error("Substream returned unexpected value", strong.str().Value)
	}

	child := strong.str()
	for kk := range values {
		child = child.Update(values[kk].str)
		strong = strong.Latest()
		if !reflect.DeepEqual(child.Value, values[kk].str) {
			t.Error("updating child didn't  take effect", child.Value)
		}
		if !reflect.DeepEqual(child.Value, strong.Value.str) {
			t.Error("updating child didn't  take effect", child.Value)
		}
	}

	v := strong.Value.setStr(values[0].str)
	if !reflect.DeepEqual(v.str, values[0].str) {
		t.Error("Could not update", "setStr")
	}
}
func TestStreamMyStructStreamCount(t *testing.T) {
	s := streams.New()
	values := valuesForMyStructStream()
	strong := &MyStructStream{Stream: s, Value: values[0]}
	expected := strong.Value.Count
	if !reflect.DeepEqual(expected, strong.Count().Value) {
		t.Error("Substream returned unexpected value", strong.Count().Value)
	}

	child := strong.Count()
	for kk := range values {
		child = child.Update(values[kk].Count)
		strong = strong.Latest()
		if !reflect.DeepEqual(child.Value, values[kk].Count) {
			t.Error("updating child didn't  take effect", child.Value)
		}
		if !reflect.DeepEqual(child.Value, strong.Value.Count) {
			t.Error("updating child didn't  take effect", child.Value)
		}
	}

	v := strong.Value.SetCount(values[0].Count)
	if !reflect.DeepEqual(v.Count, values[0].Count) {
		t.Error("Could not update", "SetCount")
	}
}
