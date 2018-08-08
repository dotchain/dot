// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package counter_test

import (
	"github.com/dotchain/dot"
	"github.com/dotchain/dot/encoding"
	"github.com/dotchain/dot/encoding/counter"
	"reflect"
	"testing"
)

func TestTransforms(t *testing.T) {
	x := dot.Transformer{}
	u := dot.Utils(x)

	initial := counter.Counter(20)
	v1, c1 := initial.Increment(5)
	v2, c2 := initial.Increment(9)

	cx1, cx2 := x.MergeChanges([]dot.Change{c1}, []dot.Change{c2})

	final1, ok1 := u.TryApply(v1, cx1)
	final2, ok2 := u.TryApply(v2, cx2)

	if !ok1 || !ok2 {
		t.Fatal("TryApply failed", ok1, ok2)
	}
	if !u.AreSame(final1, final2) || final1 != final2 {
		t.Fatal("Final results not same", final1, final2)
	}

	if int64(final1.(counter.Counter)) != int64(v1)+9 {
		t.Fatal("Merging didnt work as expected", final1)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	x := dot.Transformer{}
	u := dot.Utils(x)

	c := counter.Counter(42)
	_, change := c.Increment(21)

	// update path so that change reflects the "initial" map
	initial := map[string]interface{}{
		"hello": map[string]interface{}{
			"dot:encoding": "Counter",
			"dot:generic":  true,
			"dot:encoded":  []interface{}{42.0},
		},
	}
	change.Path = []string{"hello"}
	final := encoding.Normalize(u.Apply(initial, []dot.Change{change}))

	expected := map[string]interface{}{
		"hello": map[string]interface{}{
			"dot:encoding": "Counter",
			"dot:generic":  true,
			"dot:encoded":  []interface{}{63.0},
		},
	}

	if !reflect.DeepEqual(final, expected) {
		t.Fatal("Unexpected value after apply", final)
	}
}
