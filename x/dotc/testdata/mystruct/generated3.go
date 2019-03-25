// Generated.  DO NOT EDIT.
package mystruct

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (my MyStruct) get(key interface{}) changes.Value {
	switch key {

	case "b":
		return changes.Atomic{my.boo}
	case "bp":
		return changes.Atomic{my.boop}
	case "s":
		return types.S16(my.str)
	case "count":
		return changes.Atomic{my.Count}
	}
	panic(key)
}

func (my MyStruct) set(key interface{}, v changes.Value) changes.Value {
	myClone := my
	switch key {
	case "b":
		myClone.boo = v.(changes.Atomic).Value.(bool)
	case "bp":
		myClone.boop = v.(changes.Atomic).Value.(*bool)
	case "s":
		myClone.str = string(v.(types.S16))
	case "count":
		myClone.Count = v.(changes.Atomic).Value.(int)
	default:
		panic(key)
	}
	return myClone
}

func (my MyStruct) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my MyStruct) setBoo(value bool) MyStruct {
	myReplace := changes.Replace{changes.Atomic{my.boo}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"b"}, myReplace}
	return my.Apply(nil, myChange).(MyStruct)
}

func (my MyStruct) setBoop(value *bool) MyStruct {
	myReplace := changes.Replace{changes.Atomic{my.boop}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"bp"}, myReplace}
	return my.Apply(nil, myChange).(MyStruct)
}

func (my MyStruct) setStr(value string) MyStruct {
	myReplace := changes.Replace{types.S16(my.str), types.S16(value)}
	myChange := changes.PathChange{[]interface{}{"s"}, myReplace}
	return my.Apply(nil, myChange).(MyStruct)
}

func (my MyStruct) SetCount(value int) MyStruct {
	myReplace := changes.Replace{changes.Atomic{my.Count}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"count"}, myReplace}
	return my.Apply(nil, myChange).(MyStruct)
}

// MyStructStream implements a stream of MyStruct values
type MyStructStream struct {
	Stream streams.Stream
	Value  MyStruct
}

// Next returns the next entry in the stream if there is one
func (s *MyStructStream) Next() (*MyStructStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(MyStruct); ok {
		return &MyStructStream{Stream: next, Value: nextVal}, nextc
	}
	return &MyStructStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *MyStructStream) Latest() *MyStructStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *MyStructStream) Update(val MyStruct) *MyStructStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &MyStructStream{Stream: nexts, Value: val}
	}
	return s
}

func (s *MyStructStream) Count() *streams.Int {
	return &streams.Int{Stream: streams.Substream(s.Stream, "count"), Value: s.Value.Count}
}
