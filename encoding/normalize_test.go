// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import (
	"github.com/dotchain/dot/encoding"
	"github.com/dotchain/dot/encoding/richtext"
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	text := map[string]interface{}{
		"dot:encoding": "RichText",
		"dot:encoded":  []interface{}{map[string]interface{}{"text": "hello"}},
	}

	expect := map[string]interface{}{
		"dot:encoding": "RichText",
		"dot:encoded":  []map[string]string{{"text": "hello"}},
	}

	tests := [][2]interface{}{
		{nil, nil},
		{5, 5},
		{"hello", "hello"},
		{[]int{1, 2}, []int{1, 2}},
		{[]interface{}{1, 2}, []interface{}{1, 2}},
		{map[string]interface{}{"hello": "world"}, map[string]interface{}{"hello": "world"}},
		{richtext.NewArray(encoding.Default, text), expect},
	}

	for _, test := range tests {
		x, expected := test[0], test[1]
		if y := encoding.Normalize(x); !reflect.DeepEqual(y, expected) {
			t.Error("Normalized", x, "=", y)
		}

		if _, ok := encoding.Default.TryGet(x); !ok {
			continue
		}

		if y := encoding.Normalize(encoding.Get(x)); !reflect.DeepEqual(y, expected) {
			t.Error("Normalized(Get())", encoding.Get(x), "=", y)
		}
	}
}
