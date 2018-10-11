// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package text_test

import (
	"fmt"
	"github.com/dotchain/dot/streams/text"
)

func ExampleStream_confluence() {
	s := text.StreamFromString("Hello", true)
	s = s.SetSelection(3, 3, true)

	// try three separate edits on s
	s.Insert("A")
	s.Insert("B")
	s.Insert("C")

	// now validate that the latest has HelABC^lo
	latest := latestValue(s)
	start, _ := latest.E.Start()
	end, _ := latest.E.End()
	fmt.Println("Text:", latest.E.Text, start, end)

	// Output:
	// Text: HelABClo 6 6
}

func ExampleStream_confluenceWithCursors() {
	s := text.StreamFromString("Hello", true)

	// setting a selection at the same time as insert
	s.Insert("A")
	s.SetSelection(3, 3, true)

	// now validate that the latest has AHel^lo
	latest := latestValue(s)
	start, _ := latest.E.Start()
	end, _ := latest.E.End()
	fmt.Println("Text:", latest.E.Text, start, end)

	// Output:
	// Text: AHello 4 4
}

func latestValue(s *text.Stream) *text.Stream {
	for _, v := s.Next(); v != nil; _, v = s.Next() {
		s = v.(*text.Stream)
	}
	return s
}
