// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt

import "github.com/dotchain/dot/changes"

// Dict implements a CRDT-style dictionary
type Dict struct {
	Entries map[interface{}]Container
}

// Set creates or updates the value for a key
func (d Dict) Set(key, value interface{}) (changes.Change, Dict) {
	inner, _ := d.Entries[key].Set(NewRank(), value)
	c := wrapper{updateDict{key, inner}}
	return c, c.ApplyTo(nil, d).(Dict)
}

// Delete removes the value for the key.
func (d Dict) Delete(key interface{}) (changes.Change, Dict) {
	inner, _ := d.Entries[key].Delete()
	c := wrapper{updateDict{key, inner}}
	return c, c.ApplyTo(nil, d).(Dict)
}

// Get returns the current value or nil if there is no value available.
func (d Dict) Get(key interface{}) (*Rank, interface{}) {
	return d.Entries[key].Get()
}

// Update takes a change meant for the value at the provided key
// and wraps it so that it can applied on the dict
func (d Dict) Update(key interface{}, inner changes.Change) (changes.Change, Dict) {
	c, _ := d.Entries[key].Update(inner)
	c = wrapper{updateDict{key, c}}
	return c, c.(changes.Custom).ApplyTo(nil, d).(Dict)
}

// Clone duplicates the whole state
func (d Dict) Clone() Dict {
	result := Dict{Entries: map[interface{}]Container{}}
	for k, c := range d.Entries {
		result.Entries[k] = c
	}
	return result
}

// Apply implments changes.Value
func (d Dict) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return c.(changes.Custom).ApplyTo(ctx, d)
}

type updateDict struct {
	Key interface{}
	changes.Change
}

func (ud updateDict) Revert() crdtChange {
	return updateDict{ud.Key, ud.Change.Revert()}
}

func (ud updateDict) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Dict).Clone()
	result.Entries[ud.Key] = result.Entries[ud.Key].Apply(ctx, ud.Change).(Container)
	return result
}
