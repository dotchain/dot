// Generated.  DO NOT EDIT.
package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/heap"
)

func (my *myUnionp) get(key interface{}) changes.Value {
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

func (my *myUnionp) set(key interface{}, v changes.Value) changes.Value {
	myClone := *my
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
	return &myClone
}

func (my *myUnionp) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my *myUnionp) activeKey() string {
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

func (my *myUnionp) maxRank() int {
	rank := -1
	// fetch the largest rank
	my.activeKeyHeap.Iterate(func(_ interface{}, r int) bool {
		rank = r
		return false
	})
	return rank
}
func (my *myUnionp) setBoo(value bool) *myUnionp {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("b", rank)
	return &myUnionp{activeKeyHeap: h, boo: value}
}

func (my *myUnionp) getBoo() (bool, bool) {
	var result bool
	if my.activeKey() != "b" {
		return result, false
	}
	return my.boo, true
}

func (my *myUnionp) setBoop(value *bool) *myUnionp {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("bp", rank)
	return &myUnionp{activeKeyHeap: h, boop: value}
}

func (my *myUnionp) getBoop() (*bool, bool) {
	var result *bool
	if my.activeKey() != "bp" {
		return result, false
	}
	return my.boop, true
}

func (my *myUnionp) setStr(value string) *myUnionp {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("s", rank)
	return &myUnionp{activeKeyHeap: h, str: value}
}

func (my *myUnionp) getStr() (string, bool) {
	var result string
	if my.activeKey() != "s" {
		return result, false
	}
	return my.str, true
}

func (my *myUnionp) SetStr16(value types.S16) *myUnionp {
	rank := my.maxRank() + 1
	h := my.activeKeyHeap.Update("s16", rank)
	return &myUnionp{activeKeyHeap: h, Str16: value}
}

func (my *myUnionp) GetStr16() (types.S16, bool) {
	var result types.S16
	if my.activeKey() != "s16" {
		return result, false
	}
	return my.Str16, true
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

func (s *myUnionpStream) transformer() func(changes.Change) changes.Change {
	h := (*s.Value).activeKeyHeap
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

func (s *myUnionpStream) boo() *streams.Bool {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &streams.Bool{Stream: streams.Substream(stream, "b"), Value: s.Value.boo}
}
func (s *myUnionpStream) boop() *boolStream {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &boolStream{Stream: streams.Substream(stream, "bp"), Value: s.Value.boop}
}
func (s *myUnionpStream) str() *streams.S16 {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &streams.S16{Stream: streams.Substream(stream, "s"), Value: s.Value.str}
}
func (s *myUnionpStream) Str16() *streams.S16 {
	stream := streams.Transform(s.Stream, s.transformer(), nil)
	return &streams.S16{Stream: streams.Substream(stream, "s16"), Value: string(s.Value.Str16)}
}
