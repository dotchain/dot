// Generated.  DO NOT EDIT.
package example

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (t Todo) get(key interface{}) changes.Value {
	switch key {

	case "complete":
		return changes.Atomic{t.Complete}
	case "desc":
		return types.S16(t.Description)
	}
	panic(key)
}

func (t Todo) set(key interface{}, v changes.Value) changes.Value {
	tClone := t
	switch key {
	case "complete":
		tClone.Complete = (v).(changes.Atomic).Value.(bool)
	case "desc":
		tClone.Description = string((v).(types.S16))
	}
	return tClone
}

func (t Todo) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: t.get, Set: t.set}).Apply(ctx, c, t)
}

func (t Todo) SetComplete(value bool) Todo {
	tReplace := changes.Replace{changes.Atomic{t.Complete}, changes.Atomic{value}}
	tChange := changes.PathChange{[]interface{}{"complete"}, tReplace}
	return t.Apply(nil, tChange).(Todo)
}

func (t Todo) SetDescription(value string) Todo {
	tReplace := changes.Replace{types.S16(t.Description), types.S16(value)}
	tChange := changes.PathChange{[]interface{}{"desc"}, tReplace}
	return t.Apply(nil, tChange).(Todo)
}

// TodoStream implements a stream of Todo values
type TodoStream struct {
	Stream streams.Stream
	Value  Todo
}

// Next returns the next entry in the stream if there is one
func (s *TodoStream) Next() (*TodoStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(Todo); ok {
		return &TodoStream{Stream: next, Value: nextVal}, nextc
	}
	return &TodoStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *TodoStream) Latest() *TodoStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *TodoStream) Update(val Todo) *TodoStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &TodoStream{Stream: nexts, Value: val}
	}
	return s
}

func (s *TodoStream) Complete() *streams.Bool {
	return &streams.Bool{Stream: streams.Substream(s.Stream, "complete"), Value: s.Value.Complete}
}
func (s *TodoStream) Description() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "desc"), Value: s.Value.Description}
}

func (t TodoList) get(key interface{}) changes.Value {
	return t[key.(int)]
}

func (t TodoList) set(key interface{}, v changes.Value) changes.Value {
	tClone := TodoList(append([]Todo(nil), (t)...))
	tClone[key.(int)] = (v).(Todo)
	return tClone
}

func (t TodoList) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	tVal := t
	afterVal := (after.(TodoList))
	tNew := append(append(tVal[:offset:offset], afterVal...), tVal[end:]...)
	return tNew
}

// Slice implements changes.Collection Slice() method
func (t TodoList) Slice(offset, count int) changes.Collection {
	tSlice := (t)[offset : offset+count]
	return tSlice
}

// Count implements changes.Collection Count() method
func (t TodoList) Count() int {
	return len(t)
}

func (t TodoList) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: t.get, Set: t.set, Splice: t.splice}).Apply(ctx, c, t)
}

func (t TodoList) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: t.get, Set: t.set, Splice: t.splice}).ApplyCollection(ctx, c, t)
}

// Splice replaces [offset:offset+count] with insert...
func (t TodoList) Splice(offset, count int, insert ...Todo) TodoList {
	tInsert := TodoList(insert)
	return t.splice(offset, count, tInsert).(TodoList)
}

// Move shuffles [offset:offset+count] by distance.
func (t TodoList) Move(offset, count, distance int) TodoList {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return t.ApplyCollection(nil, c).(TodoList)
}

// TodoListStream implements a stream of TodoList values
type TodoListStream struct {
	Stream streams.Stream
	Value  TodoList
}

// Next returns the next entry in the stream if there is one
func (s *TodoListStream) Next() (*TodoListStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(TodoList); ok {
		return &TodoListStream{Stream: next, Value: nextVal}, nextc
	}
	return &TodoListStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *TodoListStream) Latest() *TodoListStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *TodoListStream) Update(val TodoList) *TodoListStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &TodoListStream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *TodoListStream) Item(index int) *TodoStream {
	return &TodoStream{Stream: streams.Substream(s.Stream, index), Value: (s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *TodoListStream) Splice(offset, count int, replacement ...Todo) *TodoListStream {
	after := TodoList(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: after}
	str := s.Stream.Append(c)
	return &TodoListStream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *TodoListStream) Move(offset, count, distance int) *TodoListStream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &TodoListStream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}
