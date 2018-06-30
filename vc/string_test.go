// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import "fmt"

func ExampleString_Splice_insertionOrder() {
	initial := "123"
	str := String{Control: New(initial), Value: initial}

	// SpliceSync behaves like an immutable Splice
	for kk := 5; kk < 10; kk++ {
		v := str.Splice(1, 0, fmt.Sprintf("%d", kk))
		fmt.Println("Inserted", v.Value)
	}

	// Splice makes a weak guarantee that insertions at the same
	// will be ordered in the order of calls to SpliceSync
	latest, _ := str.Latest()
	fmt.Println("Latest", latest.Value)

	// Output:
	// Inserted 1523
	// Inserted 1623
	// Inserted 1723
	// Inserted 1823
	// Inserted 1923
	// Latest 15678923
}

func ExampleString_Splice_strs() {
	initial := "12345"
	str := String{Control: New(initial), Value: initial}

	// we can create window into this str "234" like so:
	window := str.Slice(1, 4)
	fmt.Println("Window", window.Value)

	// we can edit the original str like so:
	str.Splice(3, 0, "E")

	// and update just the window like so:
	wlatest, _ := window.Latest()
	latest, _ := str.Latest()
	fmt.Println("New Window, Latest", wlatest.Value, latest.Value)

	// Further more, we can edit the window separately
	// and see things merge cleanly as well
	window = window.Splice(1, 0, "T")
	wlatest, _ = window.Latest()
	latest, _ = str.Latest()

	// Basically the guarantee is that all splices preserve the
	// weak order guarantee: when concurrent operations get
	// merged, the indices of items change but an item which was
	// than another in a particular version of the array will
	// never become later because of splice operations.
	fmt.Println("New Window, Latest", wlatest.Value, latest.Value)

	// Output:
	// Window 234
	// New Window, Latest 23E4 123E45
	// New Window, Latest 2T3E4 12T3E45
}

func ExampleString_Splice_branches() {
	initial := "12345"
	str := String{Control: New(initial), Value: initial}

	// branch has value 1245
	branch := str.Splice(2, 1, "")

	// update the parent directly to: 0123456
	str2 := str.Splice(0, 0, "0")
	str2.Splice(6, 0, "6")
	// now update the stale branch to 1X245
	branch = branch.Splice(1, 0, "X")

	// now verify that latest is properly merged
	latest, _ := str.Latest()
	fmt.Println(branch.Value, latest.Value)

	// Output:
	// 1X245 01X2456
}

func ExampleString_SpliceAsync() {
	initial := "12345"
	str := String{Control: New(initial), Value: initial}

	str1 := str.SpliceAsync(0, 0, "0")
	// There are no guarantees at this point that str.Latest()
	// has been updated.
	str1.Splice(0, 0, "a")
	// But there is a guarantee that by the time sync returns
	// the effects of its own history are reflected
	l, ok := str.Latest()
	fmt.Println(len(l.Value) > len(initial), ok)

	// Output:
	// true true
}

func ExampleString_Branch() {
	initial := "Hello"
	main := String{Control: New(initial), Value: initial}

	// branch off "Hell"
	b, child := main.Slice(1, 4).Branch()

	// update child to H3ll
	child.Splice(0, 1, "3")

	// update parent to "Hello world!"
	main.Splice(5, 0, " world!")

	// print both to validate child changes have not been merged
	l1, _ := main.Latest()
	l2, _ := child.Latest()
	fmt.Println(l1.Value, l2.Value)

	// push the branch and repeat
	b.Push()
	l1, _ = main.Latest()
	l2, _ = child.Latest()
	fmt.Println(l1.Value, l2.Value)

	// Output:
	// Hello world! 3ll
	// H3llo world! 3ll
}
