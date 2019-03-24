// Generated.  DO NOT EDIT.
package myslice

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

func (my MySlice) get(key interface{}) changes.Value {
	return changes.Atomic{my[key.(int)]}
}

func (my MySlice) set(key interface{}, v changes.Value) changes.Value {
	myClone := MySlice(append([]bool(nil), (my)...))
	myClone[key.(int)] = v.(changes.Atomic).Value.(bool)
	return myClone
}

func (my MySlice) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := my
	afterVal := (after.(MySlice))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return myNew
}

// Slice implements changes.Collection Slice() method
func (my MySlice) Slice(offset, count int) changes.Collection {
	mySlice := (my)[offset : offset+count]
	return mySlice
}

// Count implements changes.Collection Count() method
func (my MySlice) Count() int {
	return len(my)
}

func (my MySlice) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my MySlice) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

func (my MySlice) Splice(offset, count int, insert ...bool) MySlice {
	myInsert := MySlice(insert)
	return my.splice(offset, count, myInsert).(MySlice)
}

func (my mySlice2) get(key interface{}) changes.Value {
	return my[key.(int)]
}

func (my mySlice2) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice2(append([]MySlice(nil), (my)...))
	myClone[key.(int)] = v.(MySlice)
	return myClone
}

func (my mySlice2) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := my
	afterVal := (after.(mySlice2))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return myNew
}

// Slice implements changes.Collection Slice() method
func (my mySlice2) Slice(offset, count int) changes.Collection {
	mySlice := (my)[offset : offset+count]
	return mySlice
}

// Count implements changes.Collection Count() method
func (my mySlice2) Count() int {
	return len(my)
}

func (my mySlice2) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my mySlice2) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

func (my mySlice2) Splice(offset, count int, insert ...MySlice) mySlice2 {
	myInsert := mySlice2(insert)
	return my.splice(offset, count, myInsert).(mySlice2)
}

func (my mySlice3) get(key interface{}) changes.Value {
	return changes.Atomic{my[key.(int)]}
}

func (my mySlice3) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice3(append([]*bool(nil), (my)...))
	myClone[key.(int)] = v.(changes.Atomic).Value.(*bool)
	return myClone
}

func (my mySlice3) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := my
	afterVal := (after.(mySlice3))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return myNew
}

// Slice implements changes.Collection Slice() method
func (my mySlice3) Slice(offset, count int) changes.Collection {
	mySlice := (my)[offset : offset+count]
	return mySlice
}

// Count implements changes.Collection Count() method
func (my mySlice3) Count() int {
	return len(my)
}

func (my mySlice3) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my mySlice3) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

func (my mySlice3) Splice(offset, count int, insert ...*bool) mySlice3 {
	myInsert := mySlice3(insert)
	return my.splice(offset, count, myInsert).(mySlice3)
}

func (my *MySliceP) get(key interface{}) changes.Value {
	return changes.Atomic{(*my)[key.(int)]}
}

func (my *MySliceP) set(key interface{}, v changes.Value) changes.Value {
	myClone := MySliceP(append([]bool(nil), (*my)...))
	myClone[key.(int)] = v.(changes.Atomic).Value.(bool)
	return &myClone
}

func (my *MySliceP) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := *my
	afterVal := *(after.(*MySliceP))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return &myNew
}

// Slice implements changes.Collection Slice() method
func (my *MySliceP) Slice(offset, count int) changes.Collection {
	mySlice := (*my)[offset : offset+count]
	return &mySlice
}

// Count implements changes.Collection Count() method
func (my *MySliceP) Count() int {
	return len(*my)
}

func (my *MySliceP) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my *MySliceP) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

func (my *MySliceP) Splice(offset, count int, insert ...bool) *MySliceP {
	myInsert := MySliceP(insert)
	return my.splice(offset, count, &myInsert).(*MySliceP)
}

func (my *mySlice2P) get(key interface{}) changes.Value {
	return (*my)[key.(int)]
}

func (my *mySlice2P) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice2P(append([]*MySliceP(nil), (*my)...))
	myClone[key.(int)] = v.(*MySliceP)
	return &myClone
}

func (my *mySlice2P) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := *my
	afterVal := *(after.(*mySlice2P))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return &myNew
}

// Slice implements changes.Collection Slice() method
func (my *mySlice2P) Slice(offset, count int) changes.Collection {
	mySlice := (*my)[offset : offset+count]
	return &mySlice
}

// Count implements changes.Collection Count() method
func (my *mySlice2P) Count() int {
	return len(*my)
}

func (my *mySlice2P) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my *mySlice2P) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

func (my *mySlice2P) Splice(offset, count int, insert ...*MySliceP) *mySlice2P {
	myInsert := mySlice2P(insert)
	return my.splice(offset, count, &myInsert).(*mySlice2P)
}

func (my *mySlice3P) get(key interface{}) changes.Value {
	return changes.Atomic{(*my)[key.(int)]}
}

func (my *mySlice3P) set(key interface{}, v changes.Value) changes.Value {
	myClone := mySlice3P(append([]*bool(nil), (*my)...))
	myClone[key.(int)] = v.(changes.Atomic).Value.(*bool)
	return &myClone
}

func (my *mySlice3P) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	myVal := *my
	afterVal := *(after.(*mySlice3P))
	myNew := append(append(myVal[:offset:offset], afterVal...), myVal[end:]...)
	return &myNew
}

// Slice implements changes.Collection Slice() method
func (my *mySlice3P) Slice(offset, count int) changes.Collection {
	mySlice := (*my)[offset : offset+count]
	return &mySlice
}

// Count implements changes.Collection Count() method
func (my *mySlice3P) Count() int {
	return len(*my)
}

func (my *mySlice3P) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).Apply(ctx, c, my)
}

func (my *mySlice3P) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Get: my.get, Set: my.set, Splice: my.splice}).ApplyCollection(ctx, c, my)
}

func (my *mySlice3P) Splice(offset, count int, insert ...*bool) *mySlice3P {
	myInsert := mySlice3P(insert)
	return my.splice(offset, count, &myInsert).(*mySlice3P)
}
