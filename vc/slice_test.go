// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc_test

import (
	"fmt"
	"github.com/dotchain/dot/vc"
)

func ExampleSlice_SpliceSync_insertionOrder() {
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

func ExampleSlice_SpliceSync_slices() {
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

func ExampleSlice_SpliceSync_branches() {
	initial := []interface{}{1, 2, 3, 4, 5}
	slice := vc.Slice{Version: vc.New(initial), Value: initial}

	// branch has value [1, 2, 4, 5]
	branch := slice.SpliceSync(2, 1, nil)

	// update the parent directly to: [0.5, 1, 2, 3, 4, 5, 5.5]
	slice2 := slice.SpliceSync(0, 0, []interface{}{0.5})
	slice2.SpliceSync(6, 0, []interface{}{5.5})
	// now update the stale branch to [1, 1.5, 2, 4, 5]
	branch = branch.SpliceSync(1, 0, []interface{}{1.5})

	// now verify that latest is properly merged
	latest, _ := slice.Latest()
	fmt.Println(branch.Value, latest.Value)

	// Output:
	// [1 1.5 2 4 5] [0.5 1 1.5 2 4 5 5.5]
}

func ExampleSlice_SpliceAsync() {
	initial := []interface{}{1, 2, 3, 4, 5}
	slice := vc.Slice{Version: vc.New(initial), Value: initial}

	slice1 := slice.SpliceAsync(0, 0, []interface{}{0})
	// There are no guarantees at this point that slice.Latest()
	// has been updated.
	slice1.SpliceSync(0, 0, []interface{}{0.5})
	// But there is a guarantee that by the time sync returns
	// the effects of its own history are reflected
	l, ok := slice.Latest()
	fmt.Println(len(l.Value) > len(initial), ok)

	// Output:
	// true true
}

func ExampleSlice_Latest_nested() {
	// initial is a slice of slices
	initial := []interface{}{
		[]interface{}{1, 2, 3, 4, 5},
	}
	slice := vc.Slice{Version: vc.New(initial), Value: initial}

	inner := slice.Version.ChildAt(0)
	innerSlice := vc.Slice{Version: inner, Value: initial[0].([]interface{})}

	inner2 := slice.Version.ChildAt(0)
	inner2Slice := vc.Slice{Version: inner2, Value: initial[0].([]interface{})}

	// now modify inner and see it reflected on inner2's latest
	innerSlice.SpliceSync(0, 0, []interface{}{0})
	inner2Latest, _ := inner2Slice.Latest()

	fmt.Println(innerSlice.Value, inner2Slice.Value, inner2Latest.Value)

	// now delete the whole inner slice and see latest fail
	slice.SpliceSync(0, 1, []interface{}{})
	_, ok := inner2Slice.Latest()

	fmt.Println("Latest:", ok)

	// Output:
	// [1 2 3 4 5] [1 2 3 4 5] [0 1 2 3 4 5]
	// Latest: false
}
