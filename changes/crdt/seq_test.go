// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes/crdt"
)

func TestSeq(t *testing.T) {
	s := crdt.Seq{}
	_, s = s.Splice(0, 0, []interface{}{"hello", "new", "world"})
	if x := s.Items(); !reflect.DeepEqual(x, []interface{}{"hello", "new", "world"}) {
		t.Fatal("Splice failed", x)
	}

	c1, s1 := s.Splice(0, 1, []interface{}{"Hello"})
	if x := s1.Items(); !reflect.DeepEqual(x, []interface{}{"Hello", "new", "world"}) {
		t.Fatal("Splice failed", x)
	}

	s2 := s1.Apply(nil, c1.Revert()).(crdt.Seq)
	if x := s2.Items(); !reflect.DeepEqual(x, s.Items()) {
		t.Fatal("Undo Splice failed", x)
	}

	_, s2 = s.Splice(1, 1, nil)
	if x := s2.Items(); !reflect.DeepEqual(x, []interface{}{"hello", "world"}) {
		t.Fatal("Splice failed", x)
	}

	_, s2 = s.Splice(1, 2, nil)
	if x := s2.Items(); !reflect.DeepEqual(x, []interface{}{"hello"}) {
		t.Fatal("Splice failed", x)
	}
}
