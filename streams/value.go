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
func (vs *ValueStream) Next() (changes.Change, Stream) {
	if c, s := vs.Stream.Next(); s != nil {
		return c, &ValueStream{vs.Value.Apply(c), s}
	}
	return nil, nil
}

// Nextf implements Stream.Nextf. The stream provided in the callback
// is of type ValueStream
func (vs *ValueStream) Nextf(key interface{}, fn func(changes.Change, Stream)) {
	if fn == nil {
		vs.Stream.Nextf(key, nil)
	} else {
		vs.Stream.Nextf(key, func(c changes.Change, s Stream) {
			fn(c, &ValueStream{vs.Value.Apply(c), s})
		})
	}
}
