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
	myClone := append(MySlice(nil), my...)
	myClone[key.(int)] = v.(changes.Atomic).Value.(bool)
	return myClone
}

func (my MySlice) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	return append(append(my[:offset:offset], after.(MySlice)...), my[end:]...)
}

// Slice implements changes.Collection Slice() method
func (my MySlice) Slice(offset, count int) changes.Collection {
	return my[offset : offset+count]
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
	end := offset + count
	return append(append(my[:offset:offset], insert...), my[end:]...)
}

func (my mySlice2) get(key interface{}) changes.Value {
	return my[key.(int)]
}

func (my mySlice2) set(key interface{}, v changes.Value) changes.Value {
	myClone := append(mySlice2(nil), my...)
	myClone[key.(int)] = v.(MySlice)
	return myClone
}

func (my mySlice2) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	return append(append(my[:offset:offset], after.(mySlice2)...), my[end:]...)
}

// Slice implements changes.Collection Slice() method
func (my mySlice2) Slice(offset, count int) changes.Collection {
	return my[offset : offset+count]
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
	end := offset + count
	return append(append(my[:offset:offset], insert...), my[end:]...)
}

func (my mySlice3) get(key interface{}) changes.Value {
	return changes.Atomic{my[key.(int)]}
}

func (my mySlice3) set(key interface{}, v changes.Value) changes.Value {
	myClone := append(mySlice3(nil), my...)
	myClone[key.(int)] = v.(changes.Atomic).Value.(*bool)
	return myClone
}

func (my mySlice3) splice(offset, count int, after changes.Collection) changes.Collection {
	end := offset + count
	return append(append(my[:offset:offset], after.(mySlice3)...), my[end:]...)
}

// Slice implements changes.Collection Slice() method
func (my mySlice3) Slice(offset, count int) changes.Collection {
	return my[offset : offset+count]
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
	end := offset + count
	return append(append(my[:offset:offset], insert...), my[end:]...)
}
