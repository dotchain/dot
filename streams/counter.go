// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Counter implements a counter stream.
type Counter struct {
	Stream Stream
	Value  int32
}

// Next returns the next if there is one.
func (c *Counter) Next() (*Counter, changes.Change) {
	if c.Stream == nil {
		return nil, nil
	}

	next, nextc := c.Stream.Next()
	if next == nil {
		return nil, nil
	}

	v := c.Value
	val, ok := (types.Counter(v)).Apply(nil, nextc).(types.Counter)
	if ok {
		v = int32(val)
	} else {
		next = nil
		nextc = nil
	}
	return &Counter{Stream: next, Value: v}, nextc
}

// Latest returns the latest non-nil entry in the stream
func (c *Counter) Latest() *Counter {
	for next, _ := c.Next(); next != nil; next, _ = c.Next() {
		c = next
	}
	return c
}

// Update replaces the current value with the new value
func (c *Counter) Update(val int32) *Counter {
	before, after := types.Counter(c.Value), types.Counter(val)

	if c.Stream != nil {
		nexts := c.Stream.Append(changes.Replace{Before: before, After: after})
		c = &Counter{Stream: nexts, Value: val}
	}
	return c
}

// Increment by specified amount
func (c *Counter) Increment(by int32) *Counter {
	if c.Stream != nil {
		nexts := c.Stream.Append(types.Counter(c.Value).Increment(by))
		c = &Counter{Stream: nexts, Value: c.Value + by}
	}
	return c
}
