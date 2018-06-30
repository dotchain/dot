// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import "github.com/dotchain/dot"

// Map represents a dictionary of string keys and any type of value
type Map struct {
	// Control is the version control metadata
	Control

	// Value is the actual underlying map
	Value map[string]interface{}
}

// SetKey updates the Key with the new value and returns a new
// map with the new values.  A value of nil is the same as removing
// the key.
//
// The Sync suffix refers to the fact that an immediate call to Latest
// will show the current update reflected in it.  Note that the order
// guarantees are weak for multiple parallel updates (i.e. updates on
// the same initial Map version) but a linear update (where the result
// of the previous update is used in the next makes strong guarantees)
func (m Map) SetKey(key string, value interface{}) Map {
	before := m.Value[key]
	set := &dot.SetInfo{Key: key, Before: before, After: value}
	changes := []dot.Change{{Set: set}}
	result, _ := unwrap(utils.Apply(m.Value, changes)).(map[string]interface{})
	version := m.Control.UpdateSync(changes)
	return Map{Control: version, Value: result}
}

// SetKeyAsync updates the Key with the new value and returns a new
// map with the new values.  A value of nil is the same as removing
// the key.
//
// The are no guarantees on when the update will be reflected in a
// call to Latest() but there is an order guarantee that all updates
// that were involved in the current  Map will be applied before this
// update is applied.
func (m Map) SetKeyAsync(key string, value interface{}) Map {
	before := m.Value[key]
	set := &dot.SetInfo{Key: key, Before: before, After: value}
	changes := []dot.Change{{Set: set}}
	result, _ := unwrap(utils.Apply(m.Value, changes)).(map[string]interface{})
	version := m.Control.UpdateAsync(changes)
	return Map{Control: version, Value: result}
}

// Latest return the latest version of the current map
func (m Map) Latest() (Map, bool) {
	v, ver := m.Control.Latest()
	if ver == nil {
		return Map{}, false
	}

	val, _ := v.(map[string]interface{})
	return Map{Control: ver, Value: val}, true
}