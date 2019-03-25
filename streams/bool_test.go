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

func TestBoolStream(t *testing.T) {
	vTrue := changes.Atomic{Value: true}
	vFalse := changes.Atomic{Value: false}

	s := streams.New()
	strong := &streams.Bool{Stream: s, Value: true}
	if strong.Stream != s {
		t.Fatal("Unexpected Stream()", strong.Stream)
	}

	strong = strong.Update(false)
	if strong.Value {
		t.Error("Update did not change value", strong.Value)
	}
	s, c := s.Next()

	if !reflect.DeepEqual(c, changes.Replace{Before: vTrue, After: vFalse}) {
		t.Error("Unexpected change on main stream", c)
	}

	c = changes.Replace{Before: vFalse, After: vTrue}
	s = s.Append(c)
	c = changes.Replace{Before: vTrue, After: vFalse}
	s = s.Append(c)
	strong = strong.Latest()

	if strong.Value {
		t.Error("Unexpected change on bool stream", strong.Value)
	}

	if _, c := strong.Next(); c != nil {
		t.Error("Unexpected change on bool stream", c)
	}

	s = s.Append(changes.Replace{Before: vFalse, After: changes.Nil})

	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on bool stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: vTrue})
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on bool stream", c, strong)
	}

}
