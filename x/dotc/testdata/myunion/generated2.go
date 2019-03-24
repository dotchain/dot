// Generated.  DO NOT EDIT.
package myunion

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func (my *myUnionp) get(key interface{}) changes.Value {
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

func (my *myUnionp) set(key interface{}, v changes.Value) changes.Value {
	myClone := *my
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
	return &myClone
}

func (my *myUnionp) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my *myUnionp) setBoo(value bool) *myUnionp {
	return &myUnionp{boo: value}
}

func (my *myUnionp) setBoop(value *bool) *myUnionp {
	return &myUnionp{boop: value}
}

func (my *myUnionp) setStr(value string) *myUnionp {
	return &myUnionp{str: value}
}

func (my *myUnionp) SetStr16(value types.S16) *myUnionp {
	return &myUnionp{Str16: value}
}
