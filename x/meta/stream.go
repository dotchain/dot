// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package meta

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

// DataStream is a stream of meta Data values
type DataStream struct {
	Stream streams.Stream
	Value  Data
}

// Next returns the next value
func (s *DataStream) Next() (*DataStream, changes.Change) {
	next, nextc := s.Stream.Next()
	if next == nil {
		return nil, nil
	}

	nextVal := s.Value.Apply(nil, nextc).(Data)
	return &DataStream{Stream: next, Value: nextVal}, nextc
}

// Latest returns the latest entry in the stream
func (s *DataStream) Latest() *DataStream {
	for n, _ := s.Next(); n != nil; n, _ = s.Next() {
		s = n
	}
	return s
}
