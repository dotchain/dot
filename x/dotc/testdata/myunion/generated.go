// Generated.  DO NOT EDIT.
package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/heap"
)

func (my myUnion) get(key interface{}) changes.Value {
	switch key {
	case "_heap_":
		return my.activeKeyHeap

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
	case "_heap_":
		myClone.activeKeyHeap = v.(heap.Heap)
	case "b":
		myClone.boo = (v).(changes.Atomic).Value.(bool)
	case "bp":
		myClone.boop = (v).(changes.Atomic).Value.(*bool)
	case "s":
		myClone.str = string((v).(types.S16))
	case "s16":
		myClone.Str16 = (v).(types.S16)
	}
	return myClone
}

func (my myUnion) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my myUnion) activeKey() string {
	result := ""
	// fetch the largest ranked key => latest
	my.activeKeyHeap.Iterate(func(key interface{}, _ int) bool {
		if s, ok := key.(string); ok {
			result = s
		}
		return false
	})
	return result
}

func (my myUnion) maxRank() int {
	rank := -1
	// fetch the largest rank
	my.activeKeyHeap.Iterate(func(_ interface{}, r int) bool {
		rank = r
		return false
	})
	return rank
}
func (my myUnion) setBoo(value bool) myUnion {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("b", rank)
	return myUnion{activeKeyHeap: h, boo: value}
}

func (my myUnion) getBoo() (bool, bool) {
	var result bool
	if my.activeKey() != "b" {
		return result, false
	}
	return my.boo, true
}

func (my myUnion) setBoop(value *bool) myUnion {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("bp", rank)
	return myUnion{activeKeyHeap: h, boop: value}
}

func (my myUnion) getBoop() (*bool, bool) {
	var result *bool
	if my.activeKey() != "bp" {
		return result, false
	}
	return my.boop, true
}

func (my myUnion) setStr(value string) myUnion {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("s", rank)
	return myUnion{activeKeyHeap: h, str: value}
}

func (my myUnion) getStr() (string, bool) {
	var result string
	if my.activeKey() != "s" {
		return result, false
	}
	return my.str, true
}

func (my myUnion) SetStr16(value types.S16) myUnion {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("s16", rank)
	return myUnion{activeKeyHeap: h, Str16: value}
}

func (my myUnion) GetStr16() (types.S16, bool) {
	var result types.S16
	if my.activeKey() != "s16" {
		return result, false
	}
	return my.Str16, true
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

func (s *myUnionStream) transformer() func(changes.Change) changes.Change {
	h := (s.Value).activeKeyHeap
	var xform func(changes.Change) changes.Change
	p := []interface{}{"_heap_"}

	maxRank := func() int {
		result := -1
		h.Iterate(func(_ interface{}, r int) bool {
			result = r
			return false
		})
		return result
	}

	xform = func(c changes.Change) changes.Change {
		switch c := c.(type) {
		case changes.ChangeSet:
			result := make(changes.ChangeSet, 0, len(c))
			for _, cx := range c {
				if cx = xform(cx); cx != nil {
					result = append(result, cx)
				}
			}
			return result
		case changes.PathChange:
			if len(c.Path) == 0 {
				return xform(c.Change)
			}
			if c.Path[0] != p[0] {
				cx := h.UpdateChange(c.Path[0], maxRank()+1)
				h = h.Update(c.Path[0], maxRank()+1)
				return changes.ChangeSet{changes.PathChange{Path: p, Change: cx}, c}
			}
		}
		return c
	}
	return xform
}

func (s *myUnionStream) boo() *streams.Bool {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &streams.Bool{Stream: streams.Substream(stream, "b"), Value: s.Value.boo}
}
func (s *myUnionStream) boop() *boolStream {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &boolStream{Stream: streams.Substream(stream, "bp"), Value: s.Value.boop}
}
func (s *myUnionStream) str() *streams.S16 {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &streams.S16{Stream: streams.Substream(stream, "s"), Value: s.Value.str}
}
func (s *myUnionStream) Str16() *streams.S16 {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &streams.S16{Stream: streams.Substream(stream, "s16"), Value: string(s.Value.Str16)}
}
