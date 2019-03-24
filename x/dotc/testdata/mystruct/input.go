package mystruct

import (
	"github.com/dotchain/dot/changes/types"
)

type myStruct struct {
	boo   bool       `dotc:b`
	boop  *bool      `dotc:bp,atomic`
	str   string     `dotc:s`
	str16 types.S16  `dotc:s16`
	xtr   *string    `dotc:x,nullable`
	xtr16 *types.S16 `dotc:x16,nullable`
}
