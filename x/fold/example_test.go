// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fold_test

import (
	"fmt"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/fold"
)

func Example_appendFolded() {
	upstream := streams.New()

	// move [0 - 5] to the right by 10
	folded := fold.New(changes.Move{0, 5, 10}, upstream)

	// move [1 - 2] to the right by 20 and see it on upstream
	folded = folded.Append(changes.Move{1, 1, 20})
	if x, _ := folded.Next(); x != nil {
		fmt.Println("Unexpected Next() behavior", x)
	}

	_, c := upstream.Next()

	cx, _ := fold.Unfold(folded)
	fmt.Printf("%#v\n%#v\n", c, cx)

	//  Output:
	// changes.ChangeSet{changes.Change(nil), changes.Move{Offset:6, Count:1, Distance:15}}
	// changes.Move{Offset:0, Count:5, Distance:9}
}

func Example_appendUpstream() {
	upstream := streams.New()

	// move [0 - 5] to the right by 10
	folded := fold.New(changes.Move{0, 5, 10}, upstream)

	// move [1 - 2] to the right by 1 and see it on the folded
	upstream.Append(changes.Move{1, 1, 1})
	_, c := folded.Next()

	fmt.Printf("%#v\n", c)

	//  Output:
	// changes.ChangeSet{changes.Change(nil), changes.Move{Offset:11, Count:1, Distance:1}}
}

func Example_nilFold() {
	upstream := streams.New()
	latest := upstream
	upstream.Nextf("mykey", func() {
		var c changes.Change
		latest, c = latest.Next()
		fmt.Println("Got Change:", c)
	})
	defer upstream.Nextf("mykey", nil)

	folded := fold.New(changes.Splice{0, types.S8(""), types.S8("hello")}, upstream)
	folded2 := folded.Append(changes.Splice{0, types.S8("j"), types.S8("j")})
	folded.Append(changes.Splice{10, types.S8(""), types.S8("insert")})
	c, _ := fold.Unfold(folded2)

	fmt.Println("Unfolded:", c)

	// Output:
	// Got Change: {5  insert}
	// Unfolded: {0  jello}
}

func Example_nextf() {
	upstream := streams.New()
	folded := fold.New(changes.Splice{0, types.S8(""), types.S8("hello")}, upstream)
	var latest streams.Stream = folded
	folded.Nextf("mykey", func() {
		var c changes.Change
		latest, c = latest.Next()
		fmt.Println("Got Change:", c)
	})
	defer folded.Nextf("mykey", nil)

	// because of the folded splicee, offset 5 should get transformed to offset  10
	upstream.Append(changes.Move{5, 6, 7})

	// Output:
	// Got Change: {10 6 7}
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic but didn't happen")
		}
	}()

	upstream := streams.New()
	folded := fold.New(changes.Move{1, 2, 3}, upstream)
	folded.ReverseAppend(changes.Move{3, 4, 5})
}
