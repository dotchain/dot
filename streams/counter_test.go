// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func TestCounterStream(t *testing.T) {
	s := streams.New()
	strong := &streams.Counter{Stream: s, Value: 5}

	strong = strong.Update(22)
	if strong.Value != 22 {
		t.Error("Update did not change value", strong.Value)
	}
	s, c := s.Next()

	if !reflect.DeepEqual(c, changes.Replace{Before: types.Counter(5), After: types.Counter(22)}) {
		t.Error("Unexpected change on main stream", c)
	}

	c = changes.Replace{Before: types.Counter(22), After: types.Counter(9)}
	s = s.Append(c)
	c = changes.Replace{Before: types.Counter(9), After: types.Counter(21)}
	s = s.Append(c)
	strong = strong.Latest()

	if strong.Value != 21 {
		t.Error("Unexpected change on counter stream", strong.Value)
	}

	if _, c := strong.Next(); c != nil {
		t.Error("Unexpected change on counter stream", c)
	}

	s = s.Append(changes.Replace{Before: types.Counter(21), After: changes.Nil})

	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on counter stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: types.Counter(10)})
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on counter stream", c, strong)
	}

}

func TestCounterStreamIncrement(t *testing.T) {
	s := streams.New()
	strong1 := &streams.Counter{Stream: s, Value: 5}
	strong2 := &streams.Counter{Stream: s, Value: 5}

	strong1 = strong1.Increment(9)
	strong2 = strong2.Latest()

	if strong1.Value != 14 || strong1.Value != strong2.Value {
		t.Error("Increment diverged", strong1.Value, strong2.Value)
	}

	strong1 = strong1.Increment(-11)
	strong2 = strong2.Latest()

	if strong1.Value != 3 || strong1.Value != strong2.Value {
		t.Error("Increment diverged", strong1.Value, strong2.Value)
	}

}
