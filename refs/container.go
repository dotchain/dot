// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs

import "github.com/dotchain/dot/changes"

// Container wraps a Value and a set of references into the value.
//
// It implements the changes.Value interface. It acts like a
// "map-like" object with the wrapped value taking up the key
// "Value".  The references are not accessed vai PathChange. Instead
// they are accessed via Update changes only. This allows
// all changes to be properly merged.
type Container struct {
	changes.Value
	refs map[interface{}]Ref
}

// NewContainer wraps the value and refs into a Container
func NewContainer(v changes.Value, refs map[interface{}]Ref) Container {
	return Container{v, refs}
}

// Refs returns the internal set of refs associated with the container
func (con Container) Refs() map[interface{}]Ref {
	return con.refs
}

// UpdateRef updates a ref. If r is nil, the reference is deleted.
// All references should have "Value" as a prefix to the path.
func (con Container) UpdateRef(key interface{}, r Ref) (Container, changes.Change) {
	old := con.refs[key]
	if old == nil && r == nil {
		return con, nil
	}
	c := Update{key, old, r}
	return con.Apply(c).(Container), c
}

// GetRef returns the current ref for the provided key
func (con Container) GetRef(key interface{}) Ref {
	return con.refs[key]
}

// Apply implements changes.Value:Apply
func (con Container) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return con
	case changes.Replace:
		if !c.IsCreate() {
			return c.After
		}
	case Update:
		refs := map[interface{}]Ref{}
		for k, v := range con.refs {
			if k != c.Key {
				refs[k] = v
			}
		}
		if c.After != nil {
			refs[c.Key] = c.After
		}
		if len(refs) == 0 {
			refs = nil
		}
		return Container{con.Value, refs}
	case changes.PathChange:
		if len(c.Path) == 0 {
			return con.Apply(c.Change)
		}
		return con.applyPathChange(c)
	case changes.Custom:
		return c.ApplyTo(con)
	}
	panic("Unknown change type")
}

func (con Container) applyPathChange(c changes.PathChange) changes.Value {
	updated := map[interface{}]Ref{}
	for k, ref := range con.refs {
		ref, _ = ref.Merge(c)
		if ref != InvalidRef {
			updated[k] = ref
		}
	}
	if len(updated) == 0 {
		updated = nil
	}
	if c.Path[0] != "Value" {
		panic("Unexpected path")
	}
	val := con.Value.Apply(changes.PathChange{c.Path[1:], c.Change})
	return Container{val, updated}
}
