// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import "fmt"

func Example_composition() {
	inner := []interface{}{1, 2, 3}
	outer := map[string]interface{}{"inner": inner}
	initial := []interface{}{outer}

	root := Slice{Version: New(initial), Value: initial}
	innerSlice := Slice{Version: root.Version.ChildAt(0).Child("inner"), Value: inner}
	window := innerSlice.Slice(1, 3)

	// now create a second window
	innerSlice2 := Slice{Version: root.Version.ChildAt(0).Child("inner"), Value: inner}
	window2 := innerSlice2.Slice(1, 3)

	fmt.Println("Before", window.Value, window2.Value)

	// append in the second window and see it appear in the first
	window2 = window2.SpliceSync(2, 0, []interface{}{4})

	latest, _ := window.Latest()
	fmt.Println("After", latest.Value, window2.Value)

	// Output:
	// Before [2 3] [2 3]
	// After [2 3 4] [2 3 4]
}
