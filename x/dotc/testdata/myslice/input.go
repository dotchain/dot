package myslice

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// MySlice is public
type MySlice []bool
type mySlice2 []MySlice
type mySlice3 []*bool

// MySliceP is public
type MySliceP []bool
type mySlice2P []*MySliceP
type mySlice3P []*bool

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
