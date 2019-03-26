package myunion

import (
	"github.com/dotchain/dot/changes/types"
)

func valuesFormyUnionStream() []myUnion {
	vTrue, vFalse := true, false

	return []myUnion{
		{
			boo:   true,
			boop:  &vTrue,
			str:   "one",
			Str16: types.S16("one"),
		},
		{
			boo:   false,
			boop:  &vFalse,
			str:   "two",
			Str16: types.S16("two"),
		},
		{
			boo:   true,
			boop:  &vTrue,
			str:   "three",
			Str16: types.S16("three"),
		},
		{
			boo:   false,
			boop:  &vFalse,
			str:   "four",
			Str16: types.S16("four"),
		},
	}
}

func valuesFormyUnionpStream() []*myUnionp {
	vTrue, vFalse := true, false

	return []*myUnionp{
		&myUnionp{
			boo:   true,
			boop:  &vTrue,
			str:   "one",
			Str16: types.S16("one"),
		},
		&myUnionp{
			boo:   false,
			boop:  &vFalse,
			str:   "two",
			Str16: types.S16("two"),
		},
		&myUnionp{
			boo:   true,
			boop:  &vTrue,
			str:   "three",
			Str16: types.S16("three"),
		},
		&myUnionp{
			boo:   false,
			boop:  &vFalse,
			str:   "four",
			Str16: types.S16("four"),
		},
	}
}
