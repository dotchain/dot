// Generated.  DO NOT EDIT.
package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (my myUnion) get(key interface{}) changes.Value {
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

func (my myUnion) set(key interface{}, v changes.Value) changes.Value {
	myClone := my
	switch key {
	case "b":
		myClone.boo = v.(changes.Atomic).Value.(bool)
	case "bp":
		myClone.boop = v.(changes.Atomic).Value.(*bool)
	case "s":
		myClone.str = string(v.(types.S16))
	case "s16":
		myClone.Str16 = v.(types.S16)
	default:
		panic(key)
	}
	return myClone
}

func (my myUnion) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my myUnion) setBoo(value bool) myUnion {
	return myUnion{boo: value}
}

func (my myUnion) setBoop(value *bool) myUnion {
	return myUnion{boop: value}
}

func (my myUnion) setStr(value string) myUnion {
	return myUnion{str: value}
}

func (my myUnion) SetStr16(value types.S16) myUnion {
	return myUnion{Str16: value}
}

// myUnionStream implements a stream of myUnion values
type myUnionStream struct {
	Stream streams.Stream
	Value  myUnion
}

// Next returns the next entry in the stream if there is one
func (s *myUnionStream) Next() (*myUnionStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(myUnion); ok {
		return &myUnionStream{Stream: next, Value: nextVal}, nextc
	}
	return &myUnionStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *myUnionStream) Latest() *myUnionStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *myUnionStream) Update(val myUnion) *myUnionStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &myUnionStream{Stream: nexts, Value: val}
	}
	return s
}

func (s *myUnionStream) boo() *streams.Bool {
	return &streams.Bool{Stream: streams.Substream(s.Stream, "b"), Value: (s.Value.boo)}
}
func (s *myUnionStream) str() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "s"), Value: (s.Value.str)}
}
func (s *myUnionStream) Str16() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "s16"), Value: string(s.Value.Str16)}
}
