// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

func UnknownArrayTest() Array {
	make := func(v interface{}) interface{} {
		return map[string]interface{}{
			"dot:encoding": "hello",
			"dot:generic":  true,
			"dot:encoded":  v,
		}
	}

	return Array{
		initial: make([]interface{}{0, 1, 2, 3, 4, "hello", 6}),
		empty:   make([]interface{}{}),
		insert:  make([]interface{}{13, 14}),
		other: map[string]interface{}{
			"dot:encoding": "SparseArray",
			"dot:encoded":  []interface{}{5, "q"},
		},
		objectInsert: "hello",
		offset:       "2",
	}
}

func UnknownDictTest() Dict {
	make := func(x map[string]interface{}) interface{} {
		return map[string]interface{}{
			"dot:encoding": "unknown dict",
			"dot:generic":  true,
			"dot:encoded":  x,
		}
	}

	return Dict{
		initial:        make(map[string]interface{}{"hello": true, "world": "a"}),
		empty:          make(map[string]interface{}{}),
		someValue:      true,
		existingKeys:   []string{"hello", "world"},
		existingValues: []interface{}{true, "a"},
		nonExistingKey: "poo",
	}
}
