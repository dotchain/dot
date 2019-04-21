// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package sjson implements a portable strongly-typed json-like codec.
//
// The serialization format is json-like in that the output is
// always valid JSON but it differs from JSON in a few ways. It is
// somewhat readable but it is not meant for manual editing.  It is
// meant to be relatively easy to create serializers that work in any
// language without needing intermediate DSLs to specify the layout.
//
// The source of truth of the types is assumed to be Golang and all
// type names referred to in code effectively uses golang conventions.
//
// The serialization always includes top-level type information via:
//
//     {"typeName": <value>}
//
// For example, a simple int32 would get serialized to:
//
//     {"int32": 42}
//
// The <value> part matches JSON for integer and string types. Float32
// and Float64 are serialized as strings using Golang's `g` format
// specifier.
//
// Pointers are serialized to their underlying value (or `null` if
// they are nil):
//
//      {"*int32": null} or {"*int32": 42}
//
// Slices are styled similar to JSON:
//
//      {"[]int32": null} or {"[]int32": [42]}
//
// Maps do not use the JSON-notation since it does not adequately
// capture the type of the key for non-string keys.  Instead, maps are
// represented as a sequence with keys and values alternating:
//
//     {"map[string]string": ["hello", "world", "goodbye", "cruel sky"]}
//     // map[string]string{"hello": "world", "goodbye": "cruel sky"}
//
// Structs *could* use the JSON approach but since the type
// information is available, they instead simply encode the fields in
// sequence:
//
//      // type myStruct { Field1 string }
//      {"myStruct": ["hello"]}
//
// Structs only look at exported fields, though they do include
// embedded fields.
//
// Interfaces are encoded with the underlying type:
//
//      // type myStruct { Field1 Stringer }
//      {"myStruct": [{"string": "hello"}]}
//
// Note that interfaces can occur in slices, map keys or values. In
// all cases, the associated type is encoded.
//
// All named types should be registered via Codec.Register. This is
// needed for properly decoding types
package sjson

import "reflect"

// Codec exposes both the encoder and decoder
type Codec struct {
	Encoder
	Decoder
}

// Register registers a value and its associated type
func (c *Codec) Register(v interface{}) {
	c.Decoder.register(reflect.TypeOf(v))
}

// Std is a global encoder/decoder
var Std = Codec{}
