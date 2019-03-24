// Generated.  DO NOT EDIT.
package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func (my myUnion) get(key interface{}) changes.Value {
	switch key {

	case "b":
		return changes.Atomic{my.boo}
	case "bp":
		return changes.Atomic{my.boop}
	case "s":
		return types.S16(my.str)
	case "s16":
		return my.Str16
	}
	panic(key)
}

func (my myUnion) set(key interface{}, v changes.Value) changes.Value {
	myClone := my
	switch key {
	case "b":
		myClone.boo = v.(changes.Atomic).Value.(bool)
	case "bp":
		myClone.boop = v.(changes.Atomic).Value.(*bool)
	case "s":
		myClone.str = string(v.(types.S16))
	case "s16":
		myClone.Str16 = v.(types.S16)
	default:
		panic(key)
	}
	return myClone
}

func (my myUnion) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my myUnion) setBoo(value bool) myUnion {
	return myUnion{boo: value}
}

func (my myUnion) setBoop(value *bool) myUnion {
	return myUnion{boop: value}
}

func (my myUnion) setStr(value string) myUnion {
	return myUnion{str: value}
}

func (my myUnion) SetStr16(value types.S16) myUnion {
	return myUnion{Str16: value}
}
