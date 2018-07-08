// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package set implements the set encoding which is useful
// mainly to illustrate how set-like encodings are worthwhile.
// For example, sparse graph edge adjacency sets are probably
// best encoded this way as they are very sparse.
//
// This class is really not meant for modeling of types of
// objects represented in DOT Models.  It is strictly an
// encoder/decoder for use with the dot transformer package
// and as such real world clients.
package set

import (
	"encoding/json"
	"github.com/dotchain/dot/encoding"
	"github.com/pkg/errors"
)

func init() {
	encoding.Default.RegisterConstructor("Set", NewSet)
}

// Set represents the encoding type. The catalog is not really
// worth storing but it would be interesting for containerr
// objects that can store other encoded objects.  In that case,
// decoding them would require use of the catalog.
type Set struct {
	c encoding.Catalog
	v []interface{}
}

// NewSet is the constructor registered with encoding.Catalog
func NewSet(c encoding.Catalog, m map[string]interface{}) Set {
	values, _ := m["dot:encoded"].([]interface{})
	return Set{c, values}
}

// MarshalJSON is required to make sure the wire format matches
// DOT protocol expectations
func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"dot:encoding": "Set",
		"dot:encoded":  s.v,
	})
}

// Get panics if a key is not present or has a zero value. Since
// set membership is boolean, this leaves the only possible return
// value as true
func (s Set) Get(key string) interface{} {
	for _, val := range s.v {
		if key == val.(string) {
			return true
		}
	}
	panic(errors.Errorf("encoding/set missing key: %v", key))
}

// Set updates the key and returns a new data structure (does not modify
// the exiseting data structure).  There are only three legitimate values:
// nil or false (indicating the item is being removed) and true
// indicating the item is being inserted.
func (s Set) Set(key string, v interface{}) encoding.ObjectLike {
	exists, ok := v.(bool)
	if !ok && v != nil {
		panic(errors.Errorf("encoding/set allows boolean values only but got: %v", v))
	}

	result := make([]interface{}, 0, len(s.v)+1)
	for _, val := range s.v {
		if key != val.(string) {
			result = append(result, val)
		}
	}
	if exists {
		result = append(result, key)
	}
	return Set{s.c, result}
}

// ForKeys iterates over all the keys of the set
func (s Set) ForKeys(fn func(string, interface{})) {
	for _, key := range s.v {
		fn(key.(string), true)
	}
}
