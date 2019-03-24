// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import "io"

// Struct has the type information of a struct for code generation of
// the Apply() and SetField(..) methods
type Struct struct {
	Recv, Type string
	Fields     []Field
}

// Pointer specifies if the struct type is itself a pointer
func (s Struct) Pointer() bool {
	return s.Type[0] == '*'
}

// GenerateApply generates the code for the changes.Value Apply() method
func (s Struct) GenerateApply(w io.Writer) error {
	return structApply.Execute(w, s)
}

// GenerateSetters generates the code for the field setters
func (s Struct) GenerateSetters(w io.Writer) error {
	return structSetters.Execute(w, s)
}
