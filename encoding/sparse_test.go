// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

// sparse should register itself, so it is imported only for
// its sideeffect.
import _ "github.com/dotchain/dot/encoding/sparse"

func SparseTest() Array {
	make := func(runLengthEncoded []interface{}) interface{} {
		return map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  runLengthEncoded,
		}
	}
	return Array{
		initial: make([]interface{}{2, "a", 1, "b", 1, 42, 10, "zebra"}),
		empty:   make([]interface{}{}),
		insert:  make([]interface{}{3, "z"}),
		offset:  "3",
	}
}
