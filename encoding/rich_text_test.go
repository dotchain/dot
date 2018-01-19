// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

// sparse should register itself, so it is imported only for
// its sideeffect.
import _ "github.com/dotchain/dot/encoding/rich_text"

func RichTextTest() Array {
	make := func(encoded []interface{}) interface{} {
		return map[string]interface{}{
			"dot:encoding": "RichText",
			"dot:encoded":  encoded,
		}
	}
	return Array{
		initial: make([]interface{}{
			map[string]interface{}{
				"text":   "hello",
				"bold":   "true",
				"strike": "true",
			},
			map[string]interface{}{
				"text": " ",
			},
			map[string]interface{}{
				"text": "world",
				"bold": "true",
			},
		}),
		empty: make([]interface{}{}),
		insert: make([]interface{}{
			map[string]interface{}{
				"text": "!",
				"bold": "true",
			},
			map[string]interface{}{
				"text":    "!",
				"italics": "true",
			},
		}),
		objectInsert: map[string]interface{}{
			"text": "?",
			"bold": "true",
		},
		other: map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  []interface{}{2, map[string]interface{}{"text": "?"}},
		},
		offset: "3",
	}
}
