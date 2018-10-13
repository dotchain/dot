// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"fmt"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/types"
)

func Example_newStream() {
	latest := changes.Value(types.S8("Hello World"))
	s := streams.New()
	sLatest := s
	s.Nextf("apply", func() {
		var c changes.Change
		sLatest, c = sLatest.Next()
		latest = latest.Apply(c)
		fmt.Println("Changed:", latest)
	})

	s.Append(changes.Splice{0, types.S8(""), types.S8("OK ")})

	// Output:
	// Changed: OK Hello World
}

func Example_streamMergeUsingNextfAndApply() {
	latest := changes.Value(types.S8("Hello World"))
	s := streams.New()
	sLatest := s
	s.Nextf("apply", func() {
		var c changes.Change
		sLatest, c = sLatest.Next()
		latest = latest.Apply(c)
		fmt.Println("Changed:", latest)
	})

	s1 := s.Append(changes.Splice{0, types.S8(""), types.S8("OK ")})
	// note that this works on s, so the offset location is based
	// off "Hello World", rather than "OK Hello World"
	_ = s.Append(changes.Splice{len("Hello World"), types.S8(""), types.S8("!")})
	// now modify s1 again which is based off of "OK Hello World"
	s1.Append(changes.Splice{len("OK Hello World"), types.S8(""), types.S8("*")})

	// Output:
	// Changed: OK Hello World
	// Changed: OK Hello World!
	// Changed: OK Hello World!*
}

func Example_streamBranching() {
	val := changes.Value(types.S8("Hello World"))
	s := streams.New()
	child := streams.Branch(s)

	// update child, the changes won't be reflected on latest
	child.Append(changes.Splice{0, types.S8(""), types.S8("OK ")})
	_, c := streams.Latest(s)
	fmt.Println("Latest:", val.Apply(c))

	// merge child and parent, change will get reflected
	streams.Push(child)
	_, c = streams.Latest(s)
	fmt.Println("Latest:", val.Apply(c))

	// Output:
	// Latest: Hello World
	// Latest: OK Hello World
}
