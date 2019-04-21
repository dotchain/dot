// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package sjson implements a strongly-typed json that can be used to
// preserve type information across languages.
//
// This is accomplished by representing very value as a map with the
// key being the type name and the value being the actual JSON
// encoding.
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
