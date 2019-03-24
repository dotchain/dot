// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import "strings"

// Field holds info for a struct field
type Field struct {
	Name, Key, Type string
	Atomic          bool
}

// WrapR returns a string form of the field that maps to a changes.Value type
func (f Field) WrapR(recv string) string {
	return f.Wrap(recv + "." + f.Name)
}

// Wrap returns a string form of the field that maps to a changes.Value type
func (f Field) Wrap(val string) string {
	return wrapValue(val, f.Type, f.Atomic)
}

// Unwrap converts "v" to the type of the field
func (f Field) Unwrap() string {
	return unwrapValue("v", f.Type, f.Atomic)
}

// Setter returns the method name of the field setter
func (f Field) Setter() string {
	title := strings.Title(f.Name)
	if title == f.Name {
		return "Set" + title
	}
	return "set" + title
}
