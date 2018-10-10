// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs

import "github.com/dotchain/dot/changes"

// List manages a value which may contain references.
//
// DOT values are expected to be pure trees but with references, it is
// possible to implement generic graph-like values.
//
// References are stored in a two-step fashion: every site that wants
// to point to some other part of the value "tree" only stores an "ID"
// to that part.  Separately, List maintains a map of such IDs to the
// location where the ID points to.
//
// This two step process is needed because, as changes are applied to
// the value, the location of references can change.  For instance, if
// one of the references (say rowX) points to an element in an array
// (say rows/42) and a change comes in inserting a new element at the
// start of the array.  Now, the reference needs to be updated (to
// rows/43).  For every change, all the call sites where there are
// possible references would need to be tracked and updated.  This
// gives rise to much code complexity.
//
// The simpler approach taken here is to restrict storage of actual
// reference paths to a single map (which is updated on every
// change).  What is stored acorss the tree is only the key within
// this map.
//
// List implements the changes.Value interface as if it were a
// map-like object: the "Value" key refers to the inner value and the
// "Refs" key refers to the refs map.  Note that all the ref paths
// are expected to include the "Value" key as the first key.
//
// This is an immutable type. All mutations return the new list as
// well as a change which captures the effect of that mutation.
type List struct {
	V changes.Value
	R map[interface{}]Ref
}

// Add adds an entry in the ref map. It returns the new list and the
// equivalent change
func (l List) Add(key interface{}, ref Ref) (changes.Change, List) {
	if _, ok := l.R[key]; ok {
		panic("refs.List: duplicate key add")
	}

	p := []interface{}{"Refs", key}
	c := changes.PathChange{p, changes.Replace{changes.Nil, changes.Atomic{ref}}}
	return c, l.Apply(c).(List)
}

// Remove removes an entry in the ref map. It returns the new list and
// the equivalent change,
func (l List) Remove(key interface{}) (changes.Change, List) {
	r, ok := l.R[key]
	if !ok {
		panic("refs.List: key does not exist")
	}

	p := []interface{}{"Refs", key}
	c := changes.PathChange{p, changes.Replace{changes.Atomic{r}, changes.Nil}}
	return c, l.Apply(c).(List)
}

// Update modifies the ref for the given key. It returns the new list
// and the equivalent change.
func (l List) Update(key interface{}, ref Ref) (changes.Change, List) {
	r, ok := l.R[key]
	if !ok {
		panic("refs.List: key does not exist")
	}

	before := changes.Atomic{r}
	after := changes.Atomic{ref}
	p := []interface{}{"Refs", key}
	c := changes.PathChange{p, changes.Replace{before, after}}
	return c, l.Apply(c).(List)
}

// Apply implements changes.Value:Apply. It accepts changes to the
// value (these will have path of "Value") and also automatically
// updates all the refs that are affected by the change.
func (l List) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return l
	case changes.Replace:
		if !c.IsCreate() {
			return c.After
		}
	case changes.PathChange:
		if len(c.Path) == 0 {
			return l.Apply(c.Change)
		}
		switch c.Path[0] {
		case "Value":
			v := l.V.Apply(changes.PathChange{c.Path[1:], c.Change})
			return List{v, l.updateRefs(c)}
		case "Refs":
			return l.applyRef(c.Path[1:], c.Change)
		}
	}
	return c.(changes.Custom).ApplyTo(l)
}

func (l List) updateRefs(c changes.Change) map[interface{}]Ref {
	refs := map[interface{}]Ref{}
	for k, v := range l.R {
		refs[k], _ = v.Merge(c)
	}
	return refs
}

func (l List) cloneRefs() map[interface{}]Ref {
	refs := map[interface{}]Ref{}
	for k, v := range l.R {
		refs[k] = v
	}
	return refs
}

func (l List) applyRef(path []interface{}, c changes.Change) changes.Value {
	if c == nil {
		return l
	}

	refs := l.cloneRefs()
	after := c.(changes.Replace).After
	if after == changes.Nil {
		delete(refs, path[0])
	} else {
		refs[path[0]] = after.(changes.Atomic).Value.(Ref)
	}

	return List{l.V, refs}
}

// Count implements changes.Value:Count
func (l List) Count() int {
	panic("refs:List does not implement Count")
}

// Slice implements changes.Value:Slice
func (l List) Slice(offset, count int) changes.Value {
	panic("refs:List does not implement Slice")
}
