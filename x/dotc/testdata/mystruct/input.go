package mystruct

import (
	"github.com/dotchain/dot/changes/types"
)

type myStruct struct {
	boo   bool       `dotc:b`
	boop  *bool      `dotc:bp,atomic`
	str   string     `dotc:s`
	Str16 types.S16  `dotc:s16`
}

type myStructp struct {
	boo   bool       `dotc:b`
	boop  *bool      `dotc:bp,atomic`
	str   string     `dotc:s`
	Str16 types.S16  `dotc:s16`
}
