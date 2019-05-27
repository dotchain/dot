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
	folded := fold.New(changes.Move{Offset: 0, Count: 5, Distance: 10}, upstream)

	// move [1 - 2] to the right by 20 and see it on upstream
	folded = folded.Append(changes.Move{Offset: 1, Count: 1, Distance: 20})
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
	folded := fold.New(changes.Move{Offset: 0, Count: 5, Distance: 10}, upstream)

	// move [1 - 2] to the right by 1 and see it on the folded
	upstream.Append(changes.Move{Offset: 1, Count: 1, Distance: 1})
	_, c := folded.Next()

	fmt.Printf("%#v\n", c)

	//  Output:
	// changes.ChangeSet{changes.Change(nil), changes.Move{Offset:11, Count:1, Distance:1}}
}

func Example_nilFold() {
	upstream := streams.New()

	folded := fold.New(changes.Splice{Offset: 0, Before: types.S8(""), After: types.S8("hello")}, upstream)
	folded2 := folded.Append(changes.Splice{Offset: 1, Before: types.S8("e"), After: types.S8("u")})

	folded.Append(changes.Splice{Offset: 10, Before: types.S8(""), After: types.S8("insert")})

	_, c := upstream.Next()
	fmt.Println("Got Change:", c)

	c, _ = fold.Unfold(folded2)
	fmt.Println("Unfolded:", c)

	// Output:
	// Got Change: {5  insert}
	// Unfolded: {0  hullo}
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic but didn't happen")
		}
	}()

	upstream := streams.New()
	folded := fold.New(changes.Move{Offset: 1, Count: 2, Distance: 3}, upstream)
	folded.ReverseAppend(changes.Move{Offset: 3, Count: 4, Distance: 5})
}
