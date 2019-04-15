// Generated.  DO NOT EDIT.
package myslice

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func (my MySlice) get(key interface{}) changes.Value {
	return changes.Atomic{my[key.(int)]}
}

func (my MySlice) set(key interface{}, v changes.Value) changes.Value {
	myClone := MySlice(append([]bool(nil), (my)...))
	myClone[key.(int)] = (v).(changes.Atomic).Value.(bool)
	return myClone
}

func (my MySlice) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := my
	afterVal := (after.(MySlice))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return myNew
}

// Slice implements changes.Collection Slice() method
func (my MySlice) Slice(offset, count int) changes.Collection {
	mySlice := (my)[offset : offset+count]
	return mySlice
}

// Count implements changes.Collection Count() method
func (my MySlice) Count() int {
	return len(my)
}

func (my MySlice) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my MySlice) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

// Splice replaces [offset:offset+count] with insert...
func (my MySlice) Splice(offset, count int, insert ...bool) MySlice {
	myInsert := MySlice(insert)
	return my.splice(offset, count, myInsert).(MySlice)
}

// Move shuffles [offset:offset+count] by distance.
func (my MySlice) Move(offset, count, distance int) MySlice {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return my.ApplyCollection(nil, c).(MySlice)
}

func (my mySlice2) get(key interface{}) changes.Value {
	return my[key.(int)]
}

func (my mySlice2) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice2(append([]MySlice(nil), (my)...))
	myClone[key.(int)] = (v).(MySlice)
	return myClone
}

func (my mySlice2) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := my
	afterVal := (after.(mySlice2))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return myNew
}

// Slice implements changes.Collection Slice() method
func (my mySlice2) Slice(offset, count int) changes.Collection {
	mySlice := (my)[offset : offset+count]
	return mySlice
}

// Count implements changes.Collection Count() method
func (my mySlice2) Count() int {
	return len(my)
}

func (my mySlice2) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my mySlice2) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

// Splice replaces [offset:offset+count] with insert...
func (my mySlice2) Splice(offset, count int, insert ...MySlice) mySlice2 {
	myInsert := mySlice2(insert)
	return my.splice(offset, count, myInsert).(mySlice2)
}

// Move shuffles [offset:offset+count] by distance.
func (my mySlice2) Move(offset, count, distance int) mySlice2 {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return my.ApplyCollection(nil, c).(mySlice2)
}

func (my mySlice3) get(key interface{}) changes.Value {
	return changes.Atomic{my[key.(int)]}
}

func (my mySlice3) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice3(append([]*bool(nil), (my)...))
	myClone[key.(int)] = (v).(changes.Atomic).Value.(*bool)
	return myClone
}

func (my mySlice3) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := my
	afterVal := (after.(mySlice3))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return myNew
}

// Slice implements changes.Collection Slice() method
func (my mySlice3) Slice(offset, count int) changes.Collection {
	mySlice := (my)[offset : offset+count]
	return mySlice
}

// Count implements changes.Collection Count() method
func (my mySlice3) Count() int {
	return len(my)
}

func (my mySlice3) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my mySlice3) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

// Splice replaces [offset:offset+count] with insert...
func (my mySlice3) Splice(offset, count int, insert ...*bool) mySlice3 {
	myInsert := mySlice3(insert)
	return my.splice(offset, count, myInsert).(mySlice3)
}

// Move shuffles [offset:offset+count] by distance.
func (my mySlice3) Move(offset, count, distance int) mySlice3 {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return my.ApplyCollection(nil, c).(mySlice3)
}

func (my *MySliceP) get(key interface{}) changes.Value {
	return changes.Atomic{(*my)[key.(int)]}
}

func (my *MySliceP) set(key interface{}, v changes.Value) changes.Value {
	myClone := MySliceP(append([]bool(nil), (*my)...))
	myClone[key.(int)] = (v).(changes.Atomic).Value.(bool)
	return &myClone
}

func (my *MySliceP) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := *my
	afterVal := *(after.(*MySliceP))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return &myNew
}

// Slice implements changes.Collection Slice() method
func (my *MySliceP) Slice(offset, count int) changes.Collection {
	mySlice := (*my)[offset : offset+count]
	return &mySlice
}

// Count implements changes.Collection Count() method
func (my *MySliceP) Count() int {
	return len(*my)
}

func (my *MySliceP) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my *MySliceP) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

// Splice replaces [offset:offset+count] with insert...
func (my *MySliceP) Splice(offset, count int, insert ...bool) *MySliceP {
	myInsert := MySliceP(insert)
	return my.splice(offset, count, &myInsert).(*MySliceP)
}

// Move shuffles [offset:offset+count] by distance.
func (my *MySliceP) Move(offset, count, distance int) *MySliceP {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return my.ApplyCollection(nil, c).(*MySliceP)
}

func (my *mySlice2P) get(key interface{}) changes.Value {
	return (*my)[key.(int)]
}

func (my *mySlice2P) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice2P(append([]*MySliceP(nil), (*my)...))
	myClone[key.(int)] = (v).(*MySliceP)
	return &myClone
}

func (my *mySlice2P) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := *my
	afterVal := *(after.(*mySlice2P))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return &myNew
}

// Slice implements changes.Collection Slice() method
func (my *mySlice2P) Slice(offset, count int) changes.Collection {
	mySlice := (*my)[offset : offset+count]
	return &mySlice
}

// Count implements changes.Collection Count() method
func (my *mySlice2P) Count() int {
	return len(*my)
}

func (my *mySlice2P) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my *mySlice2P) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

// Splice replaces [offset:offset+count] with insert...
func (my *mySlice2P) Splice(offset, count int, insert ...*MySliceP) *mySlice2P {
	myInsert := mySlice2P(insert)
	return my.splice(offset, count, &myInsert).(*mySlice2P)
}

// Move shuffles [offset:offset+count] by distance.
func (my *mySlice2P) Move(offset, count, distance int) *mySlice2P {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return my.ApplyCollection(nil, c).(*mySlice2P)
}

func (my *mySlice3P) get(key interface{}) changes.Value {
	return changes.Atomic{(*my)[key.(int)]}
}

func (my *mySlice3P) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice3P(append([]*bool(nil), (*my)...))
	myClone[key.(int)] = (v).(changes.Atomic).Value.(*bool)
	return &myClone
}

func (my *mySlice3P) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := *my
	afterVal := *(after.(*mySlice3P))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return &myNew
}

// Slice implements changes.Collection Slice() method
func (my *mySlice3P) Slice(offset, count int) changes.Collection {
	mySlice := (*my)[offset : offset+count]
	return &mySlice
}

// Count implements changes.Collection Count() method
func (my *mySlice3P) Count() int {
	return len(*my)
}

func (my *mySlice3P) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my *mySlice3P) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

// Splice replaces [offset:offset+count] with insert...
func (my *mySlice3P) Splice(offset, count int, insert ...*bool) *mySlice3P {
	myInsert := mySlice3P(insert)
	return my.splice(offset, count, &myInsert).(*mySlice3P)
}

// Move shuffles [offset:offset+count] by distance.
func (my *mySlice3P) Move(offset, count, distance int) *mySlice3P {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	return my.ApplyCollection(nil, c).(*mySlice3P)
}

// MySliceStream implements a stream of MySlice values
type MySliceStream struct {
	Stream streams.Stream
	Value  MySlice
}

// Next returns the next entry in the stream if there is one
func (s *MySliceStream) Next() (*MySliceStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(MySlice); ok {
		return &MySliceStream{Stream: next, Value: nextVal}, nextc
	}
	return &MySliceStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *MySliceStream) Latest() *MySliceStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *MySliceStream) Update(val MySlice) *MySliceStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &MySliceStream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *MySliceStream) Item(index int) *streams.Bool {
	return &streams.Bool{Stream: streams.Substream(s.Stream, index), Value: (s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *MySliceStream) Splice(offset, count int, replacement ...bool) *MySliceStream {
	after := MySlice(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: after}
	str := s.Stream.Append(c)
	return &MySliceStream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *MySliceStream) Move(offset, count, distance int) *MySliceStream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &MySliceStream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}

// mySlice2Stream implements a stream of mySlice2 values
type mySlice2Stream struct {
	Stream streams.Stream
	Value  mySlice2
}

// Next returns the next entry in the stream if there is one
func (s *mySlice2Stream) Next() (*mySlice2Stream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(mySlice2); ok {
		return &mySlice2Stream{Stream: next, Value: nextVal}, nextc
	}
	return &mySlice2Stream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *mySlice2Stream) Latest() *mySlice2Stream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *mySlice2Stream) Update(val mySlice2) *mySlice2Stream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &mySlice2Stream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *mySlice2Stream) Item(index int) *MySliceStream {
	return &MySliceStream{Stream: streams.Substream(s.Stream, index), Value: (s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *mySlice2Stream) Splice(offset, count int, replacement ...MySlice) *mySlice2Stream {
	after := mySlice2(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: after}
	str := s.Stream.Append(c)
	return &mySlice2Stream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *mySlice2Stream) Move(offset, count, distance int) *mySlice2Stream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &mySlice2Stream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}

// mySlice3Stream implements a stream of mySlice3 values
type mySlice3Stream struct {
	Stream streams.Stream
	Value  mySlice3
}

// Next returns the next entry in the stream if there is one
func (s *mySlice3Stream) Next() (*mySlice3Stream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(mySlice3); ok {
		return &mySlice3Stream{Stream: next, Value: nextVal}, nextc
	}
	return &mySlice3Stream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *mySlice3Stream) Latest() *mySlice3Stream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *mySlice3Stream) Update(val mySlice3) *mySlice3Stream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &mySlice3Stream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *mySlice3Stream) Item(index int) *boolStream {
	return &boolStream{Stream: streams.Substream(s.Stream, index), Value: (s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *mySlice3Stream) Splice(offset, count int, replacement ...*bool) *mySlice3Stream {
	after := mySlice3(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: after}
	str := s.Stream.Append(c)
	return &mySlice3Stream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *mySlice3Stream) Move(offset, count, distance int) *mySlice3Stream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &mySlice3Stream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}

// MySlicePStream implements a stream of *MySliceP values
type MySlicePStream struct {
	Stream streams.Stream
	Value  *MySliceP
}

// Next returns the next entry in the stream if there is one
func (s *MySlicePStream) Next() (*MySlicePStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(*MySliceP); ok {
		return &MySlicePStream{Stream: next, Value: nextVal}, nextc
	}
	return &MySlicePStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *MySlicePStream) Latest() *MySlicePStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *MySlicePStream) Update(val *MySliceP) *MySlicePStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &MySlicePStream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *MySlicePStream) Item(index int) *streams.Bool {
	return &streams.Bool{Stream: streams.Substream(s.Stream, index), Value: (*s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *MySlicePStream) Splice(offset, count int, replacement ...bool) *MySlicePStream {
	after := MySliceP(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: &after}
	str := s.Stream.Append(c)
	return &MySlicePStream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *MySlicePStream) Move(offset, count, distance int) *MySlicePStream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &MySlicePStream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}

// mySlice2PStream implements a stream of *mySlice2P values
type mySlice2PStream struct {
	Stream streams.Stream
	Value  *mySlice2P
}

// Next returns the next entry in the stream if there is one
func (s *mySlice2PStream) Next() (*mySlice2PStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(*mySlice2P); ok {
		return &mySlice2PStream{Stream: next, Value: nextVal}, nextc
	}
	return &mySlice2PStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *mySlice2PStream) Latest() *mySlice2PStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *mySlice2PStream) Update(val *mySlice2P) *mySlice2PStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &mySlice2PStream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *mySlice2PStream) Item(index int) *MySlicePStream {
	return &MySlicePStream{Stream: streams.Substream(s.Stream, index), Value: (*s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *mySlice2PStream) Splice(offset, count int, replacement ...*MySliceP) *mySlice2PStream {
	after := mySlice2P(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: &after}
	str := s.Stream.Append(c)
	return &mySlice2PStream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *mySlice2PStream) Move(offset, count, distance int) *mySlice2PStream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &mySlice2PStream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}

// mySlice3PStream implements a stream of *mySlice3P values
type mySlice3PStream struct {
	Stream streams.Stream
	Value  *mySlice3P
}

// Next returns the next entry in the stream if there is one
func (s *mySlice3PStream) Next() (*mySlice3PStream, changes.Change) {
	if s.Stream == nil {
		return nil, nil
	}

	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	if nextVal, ok := s.Value.Apply(nil, nextc).(*mySlice3P); ok {
		return &mySlice3PStream{Stream: next, Value: nextVal}, nextc
	}
	return &mySlice3PStream{Value: s.Value}, nil
}

// Latest returns the latest entry in the stream
func (s *mySlice3PStream) Latest() *mySlice3PStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}

// Update replaces the current value with the new value
func (s *mySlice3PStream) Update(val *mySlice3P) *mySlice3PStream {
	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: s.Value, After: val})
		s = &mySlice3PStream{Stream: nexts, Value: val}
	}
	return s
}

// Item returns the sub item stream
func (s *mySlice3PStream) Item(index int) *boolStream {
	return &boolStream{Stream: streams.Substream(s.Stream, index), Value: (*s.Value)[index]}
}

// Splice splices the items replacing Value[offset:offset+count] with replacement
func (s *mySlice3PStream) Splice(offset, count int, replacement ...*bool) *mySlice3PStream {
	after := mySlice3P(replacement)
	c := changes.Splice{Offset: offset, Before: s.Value.Slice(offset, count), After: &after}
	str := s.Stream.Append(c)
	return &mySlice3PStream{Stream: str, Value: s.Value.Splice(offset, count, replacement...)}
}

// Move shuffles Value[offset:offset+count] over by distance
func (s *mySlice3PStream) Move(offset, count, distance int) *mySlice3PStream {
	c := changes.Move{Offset: offset, Count: count, Distance: distance}
	str := s.Stream.Append(c)
	return &mySlice3PStream{Stream: str, Value: s.Value.Move(offset, count, distance)}
}
