// +build !js,!tiny

// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

import (
	"encoding/json"
	"unicode/utf16"
)

func (u unknown) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.NormalizeDOT())
}

// MarshalJSON is a custom json marshaler, needed because Array wraps
// the regular array.
func (s Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.v)
}

func (e enrichArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ArrayLike)
}

func (o enrichObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ObjectLike)
}

// MarshalJSON is implemented to convert back to the required
// string format
func (s String16) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal("")
	}

	return json.Marshal(string(utf16.Decode(s)))
}
