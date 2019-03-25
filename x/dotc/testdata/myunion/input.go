package myunion

import (
	"github.com/dotchain/dot/changes/types"
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
