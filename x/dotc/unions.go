// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import "io"

type UnionStream Union

// SubFields lists all fields which support streams
func (s UnionStream) SubFields() []Field {
	result := []Field{}
	for _, f := range s.Fields {
		if !f.Atomic || streamTypes[f.Type] != "" {
			result = append(result, f)
		}
	}
	return result
}

// GenerateStream generates the stream implementation
func (s UnionStream) GenerateStream(w io.Writer) error {
	return structStreamImpl.Execute(w, s)
}

// GenerateStreamTests generates the stream tests
func (s UnionStream) GenerateStreamTests(w io.Writer) error {
	return structStreamTests.Execute(w, s)
}
