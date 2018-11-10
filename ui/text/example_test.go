// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package text_test

import (
	"fmt"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/ui/text"
)

func Example_stream_confluence() {
	s := text.StreamFromString("Hello", true)
	s = s.SetSelection(3, 3, true)

	// try three separate edits on s
	s.Insert("A")
	s.Insert("B")
	s.Insert("C")

	// now validate that the latest has HelABC^lo
	l, _ := streams.Latest(s)
	latest := l.(*text.Stream)
	start, _ := latest.Editable.Start(false)
	end, _ := latest.Editable.End(false)
	fmt.Println("Text:", latest.Editable.Text, start, end)

	// Output:
	// Text: HelABClo 6 6
}

func Example_stream_confluenceWithCursors() {
	s := text.StreamFromString("Hello", true)

	// setting a selection at the same time as insert
	s.Insert("A")
	s.SetSelection(3, 3, true)

	// now validate that the latest has AHel^lo
	l, _ := streams.Latest(s)
	latest := l.(*text.Stream)
	start, _ := latest.Editable.Start(false)
	end, _ := latest.Editable.End(false)
	fmt.Println("Text:", latest.Value(), start, end)

	// Output:
	// Text: AHello 4 4
}
