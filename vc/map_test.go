// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import "fmt"

func ExampleMap_SetKey() {
	initial := map[string]interface{}{"x": 1}
	m := Map{Control: New(initial), Value: initial}

	for kk := 5; kk < 10; kk++ {
		v := m.SetKey("x", kk)
		fmt.Println("Inserted", v.Value)
	}

	latest, _ := m.Latest()
	fmt.Println(m.Value, "=>", latest.Value)

	// Output:
	// Inserted map[x:5]
	// Inserted map[x:6]
	// Inserted map[x:7]
	// Inserted map[x:8]
	// Inserted map[x:9]
	// map[x:1] => map[x:5]
}

func ExampleMap_SetKey_branches() {
	initial := map[string]interface{}{"x": 1, "y": 5}
	m := Map{Control: New(initial), Value: initial}

	// branch sets y to 25
	branch := m.SetKey("y", 25)

	// update m directly by deleting x and setting y to 40
	m = m.SetKey("x", nil)
	m = m.SetKey("y", 40)

	// now update branch once again by setting y to 300
	branch = branch.SetKey("y", 300)

	// now verify that latest is propery merged
	latest, _ := m.Latest()
	x, y := branch.Value["x"], branch.Value["y"]
	fmt.Println(m.Value, x, y, "=>", latest.Value)

	// Output:
	// map[y:40] 1 300 => map[y:300]
}

func ExampleMap_SetKeyAsync() {
	initial := map[string]interface{}{"x": 1, "y": 5}
	m := Map{Control: New(initial), Value: initial}

	m1 := m.SetKeyAsync("y", 50)
	// There are no guarantees that at this point m.Latest()
	// would have been updated
	m1.SetKey("z", 100)
	// But there is  a guarantee that when a sync call finishes,
	// the latest  operation and any direct parent would be
	// reflected in laatest.
	latest, _ := m.Latest()

	fmt.Println("y", latest.Value["y"])

	// Output:
	// y 50
}

func ExampleMap_Latest_nested() {
	// initial is a slice of slices
	innerval := map[string]interface{}{"x": 1}
	outerval := map[string]interface{}{"inner": innerval}
	initial := map[string]interface{}{"outer": outerval}

	m := Map{Control: New(initial), Value: initial}

	inner := m.Control.Child("outer").Child("inner")
	innerMap := Map{Control: inner, Value: innerval}

	inner2 := m.Control.Child("outer").Child("inner")
	inner2Map := Map{Control: inner2, Value: innerval}

	// now modify inner and see it reflected on inner2's latest
	innerMap = innerMap.SetKey("x", 200)
	inner2Latest, _ := inner2Map.Latest()

	fmt.Println(innerMap.Value, inner2Map.Value, inner2Latest.Value)

	// now delete the whole inner map and see latest fail
	m.SetKey("outer", nil)
	_, ok := inner2Map.Latest()

	fmt.Println("Latest:", ok)

	// Output:
	// map[x:200] map[x:1] map[x:200]
	// Latest: false
}
