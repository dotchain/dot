//+build stress

// Generated.  DO NOT EDIT.
package stress

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (s State) get(key interface{}) changes.Value {
	switch key {

	case "text":
		return types.S8(s.Text)
	case "count":
		return s.Count
	}
	panic(key)
}

func (s State) set(key interface{}, v changes.Value) changes.Value {
	sClone := s
	switch key {
	case "text":
		sClone.Text = string(v.(types.S8))
	case "count":
		sClone.Count = v.(types.Counter)
	}
	return sClone
}

func (s State) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: s.get, Set: s.set}).Apply(ctx, c, s)
}

func (s State) SetText(value string) State {
	sReplace := changes.Replace{types.S8(s.Text), types.S8(value)}
	sChange := changes.PathChange{[]interface{}{"text"}, sReplace}
	return s.Apply(nil, sChange).(State)
}

func (s State) SetCount(value types.Counter) State {
	sReplace := changes.Replace{s.Count, value}
	sChange := changes.PathChange{[]interface{}{"count"}, sReplace}
	return s.Apply(nil, sChange).(State)
}

// StateStream implements a stream of State values
type StateStream struct {
	Stream streams.Stream
	Value  State
}

// Next returns the next entry in the stream if there is one
func (s *StateStream) Next() (*StateStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(State); ok {
		return &StateStream{Stream: next, Value: nextVal}, nextc
	}
	return &StateStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *StateStream) Latest() *StateStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *StateStream) Update(val State) *StateStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &StateStream{Stream: nexts, Value: val}
	}
	return s
}

func (s *StateStream) Text() *streams.S8 {
	return &streams.S8{Stream: streams.Substream(s.Stream, "text"), Value: s.Value.Text}
}
func (s *StateStream) Count() *streams.Counter {
	return &streams.Counter{Stream: streams.Substream(s.Stream, "count"), Value: int32(s.Value.Count)}
}
