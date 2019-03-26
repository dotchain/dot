package myunion

import (
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

type myUnion struct {
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
}

type myUnionp struct {
	boo   bool
	boop  *bool
	str   string
	Str16 types.S16
}

type boolStream struct {
	Stream streams.Stream
	Value *bool
}
