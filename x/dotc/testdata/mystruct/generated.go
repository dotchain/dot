// Generated.  DO NOT EDIT.
package mystruct

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (my myStruct) get(key interface{}) changes.Value {
	switch key {

	case "b":
		return changes.Atomic{my.boo}
	case "bp":
		return changes.Atomic{my.boop}
	case "s":
		return types.S16(my.str)
	case "s16":
		return my.Str16
	}
	panic(key)
}

func (my myStruct) set(key interface{}, v changes.Value) changes.Value {
	myClone := my
	switch key {
	case "b":
		myClone.boo = (v).(changes.Atomic).Value.(bool)
	case "bp":
		myClone.boop = (v).(changes.Atomic).Value.(*bool)
	case "s":
		myClone.str = string((v).(types.S16))
	case "s16":
		myClone.Str16 = (v).(types.S16)
	default:
		panic(key)
	}
	return myClone
}

func (my myStruct) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my myStruct) setBoo(value bool) myStruct {
	myReplace := changes.Replace{changes.Atomic{my.boo}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"b"}, myReplace}
	return my.Apply(nil, myChange).(myStruct)
}

func (my myStruct) setBoop(value *bool) myStruct {
	myReplace := changes.Replace{changes.Atomic{my.boop}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"bp"}, myReplace}
	return my.Apply(nil, myChange).(myStruct)
}

func (my myStruct) setStr(value string) myStruct {
	myReplace := changes.Replace{types.S16(my.str), types.S16(value)}
	myChange := changes.PathChange{[]interface{}{"s"}, myReplace}
	return my.Apply(nil, myChange).(myStruct)
}

func (my myStruct) SetStr16(value types.S16) myStruct {
	myReplace := changes.Replace{my.Str16, value}
	myChange := changes.PathChange{[]interface{}{"s16"}, myReplace}
	return my.Apply(nil, myChange).(myStruct)
}

// myStructStream implements a stream of myStruct values
type myStructStream struct {
	Stream streams.Stream
	Value  myStruct
}

// Next returns the next entry in the stream if there is one
func (s *myStructStream) Next() (*myStructStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(myStruct); ok {
		return &myStructStream{Stream: next, Value: nextVal}, nextc
	}
	return &myStructStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *myStructStream) Latest() *myStructStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *myStructStream) Update(val myStruct) *myStructStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &myStructStream{Stream: nexts, Value: val}
	}
	return s
}

func (s *myStructStream) boo() *streams.Bool {
	return &streams.Bool{Stream: streams.Substream(s.Stream, "b"), Value: s.Value.boo}
}
func (s *myStructStream) boop() *boolStream {
	return &boolStream{Stream: streams.Substream(s.Stream, "bp"), Value: s.Value.boop}
}
func (s *myStructStream) str() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "s"), Value: s.Value.str}
}
func (s *myStructStream) Str16() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "s16"), Value: string(s.Value.Str16)}
}
