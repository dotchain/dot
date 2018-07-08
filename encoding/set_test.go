// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import _ "github.com/dotchain/dot/encoding/set"

func SetTest() Dict {
	make := func(x []interface{}) interface{} {
		return map[string]interface{}{
			"dot:encoding": "Set",
			"dot:encoded":  x,
		}
	}

	return Dict{
		initial:        make([]interface{}{"hello", "world"}),
		empty:          make([]interface{}{}),
		someValue:      true,
		existingKeys:   []string{"hello", "world"},
		existingValues: []interface{}{true, true},
		nonExistingKey: "poo",
	}
}
