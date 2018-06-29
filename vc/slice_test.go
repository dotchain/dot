// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc_test

import (
	"fmt"
	"github.com/dotchain/dot/vc"
)

func ExampleSliceSpliceSync_insertionOrder() {
	initial := []interface{}{1, 2, 3}
	slice := vc.Slice{Version: vc.New(initial), Value: initial}

	// SpliceSync behaves like an immutable Splice
	for kk := 5; kk < 10; kk++ {
		v := slice.SpliceSync(1, 0, []interface{}{kk})
		fmt.Println("Inserted", v.Value)
	}

	// Splice makes a weak guarantee that insertions at the same
	// will be ordered in the order of calls to SpliceSync
	latest, _ := slice.Latest()
	fmt.Println("Latest", latest.Value)

	// Output:
	// Inserted [1 5 2 3]
	// Inserted [1 6 2 3]
	// Inserted [1 7 2 3]
	// Inserted [1 8 2 3]
	// Inserted [1 9 2 3]
	// Latest [1 5 6 7 8 9 2 3]
}

func ExampleSliceSpliceSync_slices() {
	initial := []interface{}{1, 2, 3, 4, 5}
	slice := vc.Slice{Version: vc.New(initial), Value: initial}

	// we can create window into this slice ([2 3 4]) like so:
	window := slice.Slice(1, 4)
	fmt.Println("Window", window.Value)

	// we can edit the original slice like so:
	slice.SpliceSync(3, 0, []interface{}{3.5})

	// and update just the window like so:
	wlatest, _ := window.Latest()
	latest, _ := slice.Latest()
	fmt.Println("New Window, Latest", wlatest.Value, latest.Value)

	// Further more, we can edit the window separately
	// and see things merge cleanly as well
	window = window.SpliceSync(1, 0, []interface{}{2.5})
	wlatest, _ = window.Latest()
	latest, _ = slice.Latest()

	// Basically the guarantee is that all splices preserve the
	// weak order guarantee: when concurrent operations get
	// merged, the indices of items change but an item which was
	// than another in a particular version of the array will
	// never become later because of splice operations.
	fmt.Println("New Window, Latest", wlatest.Value, latest.Value)

	// Output:
	// Window [2 3 4]
	// New Window, Latest [2 3 3.5 4] [1 2 3 3.5 4 5]
	// New Window, Latest [2 2.5 3 3.5 4] [1 2 2.5 3 3.5 4 5]
}
