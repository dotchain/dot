// Generated.  DO NOT EDIT.
package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (my *myUnionp) get(key interface{}) changes.Value {
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

func (my *myUnionp) set(key interface{}, v changes.Value) changes.Value {
	myClone := *my
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
	return &myClone
}

func (my *myUnionp) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my *myUnionp) setBoo(value bool) *myUnionp {
	return &myUnionp{boo: value}
}

func (my *myUnionp) setBoop(value *bool) *myUnionp {
	return &myUnionp{boop: value}
}

func (my *myUnionp) setStr(value string) *myUnionp {
	return &myUnionp{str: value}
}

func (my *myUnionp) SetStr16(value types.S16) *myUnionp {
	return &myUnionp{Str16: value}
}

// myUnionpStream implements a stream of *myUnionp values
type myUnionpStream struct {
	Stream streams.Stream
	Value  *myUnionp
}

// Next returns the next entry in the stream if there is one
func (s *myUnionpStream) Next() (*myUnionpStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(*myUnionp); ok {
		return &myUnionpStream{Stream: next, Value: nextVal}, nextc
	}
	return &myUnionpStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *myUnionpStream) Latest() *myUnionpStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *myUnionpStream) Update(val *myUnionp) *myUnionpStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &myUnionpStream{Stream: nexts, Value: val}
	}
	return s
}

func (s *myUnionpStream) boo() *streams.Bool {
	return &streams.Bool{Stream: streams.Substream(s.Stream, "b"), Value: (s.Value.boo)}
}
func (s *myUnionpStream) str() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "s"), Value: (s.Value.str)}
}
func (s *myUnionpStream) Str16() *streams.S16 {
	return &streams.S16{Stream: streams.Substream(s.Stream, "s16"), Value: string(s.Value.Str16)}
}
