package myunion

import (
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/heap"	
	"github.com/dotchain/dot/streams"
)

type myUnion struct {
	activeKeyHeap heap.Heap
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
}

type myUnionp struct {
	activeKeyHeap heap.Heap
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
}

type boolStream struct {
	Stream streams.Stream
	Value *bool
}
