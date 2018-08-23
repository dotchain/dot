// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fold_test

import (
	"fmt"
	"github.com/dotchain/dot"
	"github.com/dotchain/dot/fold"
)

func Example_Foldable_simple() {
	s := "OK: Hello World"
	f0 := fold.Foldable{Local: s, Remote: s}
	fmt.Println("Fresh:", f0.LocalValue(), f0.RemoteValue())
	f1 := f0.Fold(splice(0, "OK: ", ""))
	fmt.Println("After fold1:", f1.LocalValue(), f1.RemoteValue())
	f2 := f1.Fold(splice(0, "Hell", "L"))
	fmt.Println("After fold2:", f2.LocalValue(), f2.RemoteValue())
	f3 := fold.Folded(f2).Apply(splice(3, "W", "w"))
	fmt.Println("After W->w:", f3.LocalValue(), f3.RemoteValue())
	f4 := fold.Unfolded(f3).Apply(splice(15, "", "!"))
	fmt.Println("After adding !:", f4.LocalValue(), f4.RemoteValue())
	f5 := fold.Unfolded(f4).Apply(splice(5, "e", "E"))
	fmt.Println("After e=>E:", f5.LocalValue(), f5.RemoteValue())
	f6 := f5.Unfold(0, 1)
	fmt.Println("After unfold:", f6.LocalValue(), f6.RemoteValue())

	// Output:
	// Fresh: OK: Hello World OK: Hello World
	// After fold1: Hello World OK: Hello World
	// After fold2: Lo World OK: Hello World
	// After W->w: Lo world OK: Hello world
	// After adding !: Lo world! OK: Hello world!
	// After e=>E: Lo world! OK: HEllo world!
	// After unfold: OK: Lo world! OK: HEllo world!
}

func splice(offset int, before, after string) []dot.Change {
	return []dot.Change{{Splice: &dot.SpliceInfo{offset, before, after}}}
}
