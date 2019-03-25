// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
)

func TestIntStream(t *testing.T) {
	s := streams.New()
	strong := &streams.Int{Stream: s, Value: 10}
	if strong.Stream != s {
		t.Fatal("Unexpected Stream()", strong.Stream)
	}

	strong = strong.Update(15)
	if strong.Value != 15 {
		t.Error("Update did not change value", strong.Value)
	}
	s, c := s.Next()

	before, after := changes.Atomic{Value: 10}, changes.Atomic{Value: 15}
	if !reflect.DeepEqual(c, changes.Replace{Before: before, After: after}) {
		t.Error("Unexpected change on main stream", c)
	}

	c = changes.Replace{Before: after, After: changes.Atomic{Value: 2}}
	s = s.Append(c)
	c = changes.Replace{Before: changes.Atomic{Value: 2}, After: changes.Atomic{Value: 99}}
	s = s.Append(c)
	strong = strong.Latest()

	if strong.Value != 99 {
		t.Error("Unexpected change on int stream", strong.Value)
	}

	if _, c := strong.Next(); c != nil {
		t.Error("Unexpected change on int stream", c)
	}

	s = s.Append(changes.Replace{Before: changes.Atomic{Value: 99}, After: changes.Nil})

	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on int stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: changes.Atomic{Value: 99}})
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on int stream", c, strong)
	}

}
