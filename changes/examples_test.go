// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"fmt"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/types"
)

func ExampleStream_create() {
	latest := changes.Value(types.S8("Hello World"))
	s := changes.NewStream().On("apply", func(c changes.Change, _ interface{}, _ *changes.Stream) {
		latest = latest.Apply(c)
		fmt.Println("Changed:", latest)
	})

	s.Apply(changes.Splice{0, types.S8(""), types.S8("OK ")}, nil)

	// Output:
	// Changed: OK Hello World
}

func ExampleStream_mergeUsingOnAndApply() {
	latest := changes.Value(types.S8("Hello World"))
	s := changes.NewStream().On("apply", func(c changes.Change, _ interface{}, _ *changes.Stream) {
		latest = latest.Apply(c)
		fmt.Println("Changed:", latest)
	})

	s1 := s.Apply(changes.Splice{0, types.S8(""), types.S8("OK ")}, nil)
	// note that this works on s, so the offset location is based
	// off "Hello World", rather than "OK Hello World"
	_ = s.Apply(changes.Splice{len("Hello World"), types.S8(""), types.S8("!")}, nil)
	// now modify s1 again which is based off of "OK Hello World"
	s1.Apply(changes.Splice{len("OK Hello World"), types.S8(""), types.S8("*")}, nil)

	// Output:
	// Changed: OK Hello World
	// Changed: OK Hello World!
	// Changed: OK Hello World!*
}

func ExampleBranch_create() {
	latest := changes.Value(types.S8("Hello World"))
	s := changes.NewStream().On("apply", func(c changes.Change, _ interface{}, _ *changes.Stream) {
		latest = latest.Apply(c)
	})

	// create a new stream for the "child"
	child := changes.NewStream()
	branch := &changes.Branch{s, child}

	// update child, the changes won't be reflected on latest
	child.Apply(changes.Splice{0, types.S8(""), types.S8("OK ")}, nil)
	fmt.Println("Latest:", latest)

	// merge child and parent, change will get reflected
	branch.Merge()
	fmt.Println("Latest:", latest)

	// Output:
	// Latest: Hello World
	// Latest: OK Hello World
}
