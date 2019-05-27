package mystruct

import "github.com/dotchain/dot/changes/types"

func valuesFormyStructStream() []myStruct {
	vTrue, vFalse := true, false

	return []myStruct{
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

func valuesFormyStructpStream() []*myStructp {
	vTrue, vFalse := true, false

	return []*myStructp{
		&myStructp{
			boo:   true,
			boop:  &vTrue,
			str:   "one",
			Str16: types.S16("one"),
		},
		&myStructp{
			boo:   false,
			boop:  &vFalse,
			str:   "two",
			Str16: types.S16("two"),
		},
		&myStructp{
			boo:   true,
			boop:  &vTrue,
			str:   "three",
			Str16: types.S16("three"),
		},
		&myStructp{
			boo:   false,
			boop:  &vFalse,
			str:   "four",
			Str16: types.S16("four"),
		},
	}
}

func valuesForMyStructStream() []MyStruct {
	vTrue, vFalse := true, false

	return []MyStruct{
		{
			boo:   true,
			boop:  &vTrue,
			str:   "one",
			Str16: types.S16("one"),
			Count: 1,
		},
		{
			boo:   false,
			boop:  &vFalse,
			str:   "two",
			Str16: types.S16("two"),
			Count: 2,
		},
		{
			boo:   true,
			boop:  &vTrue,
			str:   "three",
			Str16: types.S16("three"),
			Count: 3,
		},
		{
			boo:   false,
			boop:  &vFalse,
			str:   "four",
			Str16: types.S16("four"),
			Count: 4,
		},
	}
}
