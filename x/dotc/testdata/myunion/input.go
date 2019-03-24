package myunion

import (
	"github.com/dotchain/dot/changes/types"
)

type myUnion struct {
	boo   bool       `dotc:b`
	boop  *bool      `dotc:bp,atomic`
	str   string     `dotc:s`
	Str16 types.S16  `dotc:s16`
}

type myUnionp struct {
	boo   bool       `dotc:b`
	boop  *bool      `dotc:bp,atomic`
	str   string     `dotc:s`
	Str16 types.S16  `dotc:s16`
}
