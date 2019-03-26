package mystruct

import (
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
	Count int
}

type boolStream struct {
	Stream streams.Stream
	Value *bool
}
