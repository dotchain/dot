// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package compiler_test

import (
	"github.com/andreyvit/diff"
	"github.com/dotchain/dot/compiler"
	"strings"
	"testing"
)

func TestTask(t *testing.T) {
	info := compiler.Info{
		Package: "task",
		Imports: [][2]string{
			{"hello", "github.com/boo/hello"},
			{"hello2", "github.com/boo/hello2"},
		},
		Streams: []compiler.StreamInfo{
			{
				StreamType: "TaskStream",
				ValueType:  "Tasks",
				Fields: []compiler.FieldInfo{{
					Field:           "Done",
					FieldType:       "bool",
					FieldStreamType: "BoolStream",
				}},
				EntryStreamType: "",
			},
		},
	}
	got := strings.TrimSpace(compiler.Compile(info))
	want := strings.TrimSpace(taskExpected)
	if got != want {
		t.Errorf("Diff:\n%v", diff.LineDiff(want, got))
	}

	// Output:
}

var taskExpected = `
// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
//

// task is generated streams
package task

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/streams"
)

// TaskStream is a stream of Tasks values.
type TaskStream struct {
	// Notifier provides On/Off/Notify support. New instances of
	// TaskStream created via the AppendLocal or AppendRemote
	// share the same Notifier value.
	*streams.Notifier

	// Value holds the current value. The latest value may be
	// fetched via the Latest() method.
	Value Tasks

	// Change tracks the change that leads to the next value.
	Change changes.Change

	// Next tracks the next value in the stream.
	Next *TaskStream
}

// NewTaskStream creates a new Tasks stream
func NewTaskStream(value Tasks) *TaskStream {
	return &TaskStream{&streams.Notifier{}, value, nil, nil}
}

// Append appends a local change. isLocal identifies if the caller is
// local or remote. It returns the updated stream whose value matches
// the provided value and whose Latest() converges to the latest of
// the stream.
func (s *TaskStream) Append(c changes.Change, value Tasks, isLocal bool) *TaskStream {
	if c == nil {
		c = changes.Replace{Before: s.wrapValue(s.Value), After: s.wrapValue(value)}
	}

	// return value: after is correctly set to provided value
	result := &TaskStream{Notifier: s.Notifier, Value: value}

	// before tracks s, after tracks result, v tracks latest value
	// of after chain
	before := s
	v := changes.Atomic{value}

	// walk the chain of Next and find corresponding values to
	// add to after so that both s annd after converge
	for after := result; before.Next != nil; before = before.Next {
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
		v = v.Apply(afterChange)
		after.Change = afterChange
		after.Next = &TaskStream{Notifier: s.Notifier, Value: s.unwrapValue(v)}
		after = after.Next
	}

	// append the residual change (c) to converge to wherever
	// after has landed. Notify since s.Latest() has now changed
	before.Change, before.Next = c, after
	s.Notify()
	return result
}

func (s *TaskStream) wrapValue(value Tasks) changes.Value {
	if x, ok := value.(changes.Value); ok {
		return x
	}

	return changes.Atomic{value}
}

func (s *TaskStream) unwrapValue(v changes.Value) Tasks {
	if x, ok := v.(Tasks); ok {
		return x
	}
	return v.(changes.Atomic).Value
}

// SetDone updates the field with a new value
func (s *TaskStream) SetDone(v bool) *TaskStream {
	c := changes.Replace{s.wrapValue(s.Value.Done), s.wrapValue(v)}
	value := s.Value
	value.Done = v
	key := []interface{}{"Done"}
	return s.Append(changes.PathChange{key, c}, value)
}

// FieldSubstream returns a stream for Done that is automatically
// connected to the current TaskStream instance.  Updates to one
// automatically update the other.
func (s *TaskStream) FieldSubstream(cache SubstreamCache) (field *DoneTaskStream) {
	n := s.Notifier
	handler := &streams.Handler{nil}
	if f, h, ok := cache.GetSubstream(n, "Done"); ok {
		field, handler = f.(*DoneTaskStream), h.(*streams.Handler)
	} else {
		field = NewDoneTaskStream(s.Value.Done)
		parent, merging, path := s, false, []interface{}{"Done"}
		handler.Handle = func() {
			if merging {
				return
			}

			merging := true
			for ; field.Next != nil; field = field.Next {
				v := parent.Value
				v.Done = field.Value
				c := changes.PathChange{path, field.Change}
				parent = parent.Append(c, v, true)
			}

			for ; parent.Next != nil; parent = parent.Next {
				result := refs.Merge(path, parent.Change)
				if result == nil {
					field = field.Append(nil, parent.Value.Done, true)
				} else {
					field = field.Append(result.Affected, parent.Value.Done, true)
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
	close := func() { n.Off(handler); n2.Off(handler) }
	cache.SetSubstream(n, "Done", field, handler, close)
	return field
}
`
