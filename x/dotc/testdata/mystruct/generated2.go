// Generated.  DO NOT EDIT.
package mystruct

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func (my *myStructp) get(key interface{}) changes.Value {
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

func (my *myStructp) set(key interface{}, v changes.Value) changes.Value {
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

func (my *myStructp) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set}).Apply(ctx, c, my)
}

func (my *myStructp) setBoo(value bool) *myStructp {
	myReplace := changes.Replace{changes.Atomic{my.boo}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"b"}, myReplace}
	return my.Apply(nil, myChange).(*myStructp)
}

func (my *myStructp) setBoop(value *bool) *myStructp {
	myReplace := changes.Replace{changes.Atomic{my.boop}, changes.Atomic{value}}
	myChange := changes.PathChange{[]interface{}{"bp"}, myReplace}
	return my.Apply(nil, myChange).(*myStructp)
}

func (my *myStructp) setStr(value string) *myStructp {
	myReplace := changes.Replace{types.S16(my.str), types.S16(value)}
	myChange := changes.PathChange{[]interface{}{"s"}, myReplace}
	return my.Apply(nil, myChange).(*myStructp)
}

func (my *myStructp) SetStr16(value types.S16) *myStructp {
	myReplace := changes.Replace{my.Str16, value}
	myChange := changes.PathChange{[]interface{}{"s16"}, myReplace}
	return my.Apply(nil, myChange).(*myStructp)
}
