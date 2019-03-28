package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/heap"
)

type myUnion struct {
	activeKeyHeap heap.Heap
	boo           bool
	boop          *bool
	str           string
	Str16         types.S16
}

type myUnionp struct {
	activeKeyHeap heap.Heap
	boo           bool
	boop          *bool
	str           string
	Str16         types.S16
}

type boolStream struct {
	Stream streams.Stream
	Value  *bool
}

func (s *boolStream) Update(val *bool) *boolStream {
	before, after := changes.Atomic{Value: s.Value}, changes.Atomic{Value: val}

	if s.Stream != nil {
		nexts := s.Stream.Append(changes.Replace{Before: before, After: after})
		s = &boolStream{Stream: nexts, Value: val}
	}
	return s
}
