// Generated.  DO NOT EDIT.
package mystruct

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func (my myStruct) get(key interface{}) changes.Value {
	switch key {
	case "b":
		return changes.Atomic{my.boo}
	case "bp":
		if my.boop == nil {
			return changes.Nil
		}
		return changes.Atomic{my.boop}
	case "s":
		return types.S16(my.str)
	case "s16":
		return my.str16
	case "x":
		if my.xtr == nil {
			return changes.Nil
		}
		return types.S16(*my.xtr)
	case "x16":
		if my.xtr16 == nil {
			return changes.Nil
		}
		return *my.xtr16
	}
	panic(key)
}

func (my myStruct) set(key interface{}, v changes.Value) changes.Value {
	myClone := my
	switch key {
	case "b":
		myClone.boo = v.(changes.Atomic).Value.(bool)
	case "bp":
		if v == changes.Nil {
			myClone.boop = nil
			break
		}
		myClone.boop = v.(changes.Atomic).Value.(*bool)
	case "s":
		myClone.str = string(v.(types.S16))
	case "s16":
		myClone.str16 = v.(types.S16)
	case "x":
		if v == changes.Nil {
			myClone.xtr = nil
			break
		}
		myClone.xtr = func(x string) *string { return &x }(string(v.(types.S16)))
	case "x16":
		if v == changes.Nil {
			myClone.xtr16 = nil
			break
		}
		myClone.xtr16 = func(x types.S16) *types.S16 { return &x }(v.(types.S16))
	default:
		panic(key)
	}
	return myClone
}

func (my myStruct) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}
