// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"fmt"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func Example_newStream() {
	s := streams.New()
	s.Append(changes.Splice{
		Offset: 0,
		Before: types.S8(""),
		After:  types.S8("OK "),
	})

	_, c := streams.Latest(s)
	fmt.Println("Changed:", types.S8("Hello World").Apply(nil, c))

	// Output:
	// Changed: OK Hello World
}

func Example_streamMerge() {
	s := streams.New()
	s1 := s.Append(changes.Splice{
		Offset: 0,
		Before: types.S8(""),
		After:  types.S8("OK "),
	})

	_, c := streams.Latest(s)
	fmt.Println("Changed:", types.S8("Hello World").Apply(nil, c))

	// note that this works on s, so the offset location is based
	// off "Hello World", rather than "OK Hello World"
	_ = s.Append(changes.Splice{
		Offset: len("Hello World"),
		Before: types.S8(""),
		After:  types.S8("!"),
	})

	_, c = streams.Latest(s)
	fmt.Println("Changed:", types.S8("Hello World").Apply(nil, c))

	// now modify s1 again which is based off of "OK Hello World"
	s1.Append(changes.Splice{
		Offset: len("OK Hello World"),
		Before: types.S8(""),
		After:  types.S8("*"),
	})

	_, c = streams.Latest(s)
	fmt.Println("Changed:", types.S8("Hello World").Apply(nil, c))

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
	child.Append(changes.Splice{
		Offset: 0,
		Before: types.S8(""),
		After:  types.S8("OK "),
	})

	_, c := streams.Latest(s)
	fmt.Println("Latest:", val.Apply(nil, c))

	// merge child and parent, change will get reflected
	if err := child.Push(); err != nil {
		fmt.Println("Error", err)
	}

	_, c = streams.Latest(s)
	fmt.Println("Latest:", val.Apply(nil, c))

	// Output:
	// Latest: Hello World
	// Latest: OK Hello World
}
