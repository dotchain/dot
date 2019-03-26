// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import "io"

type UnionStream Union

// GenerateStream generates the stream implementation
func (s UnionStream) GenerateStream(w io.Writer) error {
	return structStreamImpl.Execute(w, StructStream(Union(s)))
}

// GenerateStreamTests generates the stream tests
func (s UnionStream) GenerateStreamTests(w io.Writer) error {
	return structStreamTests.Execute(w, StructStream(Union(s)))
}
