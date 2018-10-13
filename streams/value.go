// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import "github.com/dotchain/dot/changes"

// ValueStream implements the Stream interface also caching the value
// with it.
type ValueStream struct {
	changes.Value
	Stream
}

// Append implements Stream.Append
func (vs *ValueStream) Append(c changes.Change) Stream {
	return &ValueStream{vs.Value.Apply(c), vs.Stream.Append(c)}
}

// ReverseAppend implements Stream.ReverseAppend
func (vs *ValueStream) ReverseAppend(c changes.Change) Stream {
	return &ValueStream{vs.Value.Apply(c), vs.Stream.ReverseAppend(c)}
}

// Next implements Stream.Next.  The stream returned is either nil or
// a ValueStream, so the updated value can be obtained by casting it.
func (vs *ValueStream) Next() (Stream, changes.Change) {
	if s, c := vs.Stream.Next(); s != nil {
		return &ValueStream{vs.Value.Apply(c), s}, c
	}
	return nil, nil
}
