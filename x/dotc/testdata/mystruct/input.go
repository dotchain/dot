package mystruct

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

type myStruct struct {
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
}

type myStructp struct {
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
}

// MyStruct is public
type MyStruct struct {
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
	Count int32
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
