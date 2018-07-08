// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package encoding defines the general encoding/decoding
// mechanism for the DOT transformer.  DOT works with virtual
// JSON models but the logical wire representation of operations
// allows a lot of flexibility.  The encoding package helps
// bridge those.  For example, the encoding/sparse package
// has a sparse representation for arrays.
//
// Note that the actual internal layout of objects in a client
// may not match the encoding at all. These structs defined here
// are not very useful for clients -- they are instead focused
// on allowing DOT transformers to work with non-standard
// wire representations.
//
// Note also that all encodings should essentially work off
// persistent data structures as that is a core expectation of
// the DOT internal merge code.
package encoding

import "github.com/pkg/errors"

// Dict represents a simple interface for manipulation of JSON objects.
// This is installed by default and so any DOT-based code can
// use actual map[string]interface{} types and DOT will transform
// those properly.
type Dict map[string]interface{}

// Get returns the value of the key.
func (s Dict) Get(key string) interface{} {
	if v, ok := s[key]; ok {
		return v
	}
	panic(errors.Errorf("cannot access field %s", key))
}

// Set returns a new map with the key updated to the provided  value
func (s Dict) Set(key string, value interface{}) ObjectLike {
	result := map[string]interface{}{}
	for k, v := range s {
		if k != key {
			result[k] = v
		}
	}
	if value != nil {
		result[key] = value
	}
	return Dict(result)
}

// ForKeys iterates over all the keys and values of the dictionary
func (s Dict) ForKeys(fn func(string, interface{})) {
	for k, v := range s {
		fn(k, v)
	}
}
