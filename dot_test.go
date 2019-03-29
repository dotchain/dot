// Code generated by github.com/tvastar/test/cmd/testmd/testmd.go. DO NOT EDIT.

package dot_test

import (
	"fmt"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/undo"
	"github.com/dotchain/dot/x/fold"
)

func Example_applying_changes() {
	// import fmt
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types

	// S8 is DOT-compatible string type with UTF8 string indices
	initial := types.S8("hello")

	append := changes.Splice{
		Offset: len("hello"),       // end of "hello"
		Before: types.S8(""),       // nothing to remove
		After:  types.S8(" world"), // insert " world"
	}

	// apply the change
	updated := initial.Apply(nil, append)

	fmt.Println(updated)
	// Output: hello world

}

func Example_apply_stream() {
	// import fmt
	// import github.com/dotchain/dot/streams

	initial := &streams.S8{Stream: streams.New(), Value: "hello"}
	updated := initial.Splice(5, 0, " world")

	fmt.Println(updated.Value)
	// Output: hello world

}

func Example_changeset_composition() {
	// import fmt
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types

	initial := types.S8("hello")

	// append " world" => "hello world"
	append1 := changes.Splice{
		Offset: len("hello"),
		Before: types.S8(""),
		After:  types.S8(" world"),
	}

	// append "." => "hello world."
	append2 := changes.Splice{
		Offset: len("hello world"),
		Before: types.S8(""),
		After:  types.S8("."),
	}

	// now combine the two appends and apply
	both := changes.ChangeSet{append1, append2}
	updated := initial.Apply(nil, both)
	fmt.Println(updated)

	// Output: hello world.

}

func Example_path_composition() {
	// import fmt
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types

	// types.A is a generic array type and types.M is a map type
	initial := types.A{types.M{"hello": types.S8("world")}}

	// replace "world" with "world!"
	replace := changes.Replace{Before: types.S8("world"), After: types.S8("world!")}

	// replace "world" with "world!" of initial[0]["hello"]
	path := []interface{}{0, "hello"}
	c := changes.PathChange{Path: path, Change: replace}
	updated := initial.Apply(nil, c)
	fmt.Println(updated)

	// Output: [map[hello:world!]]

}

func Example_convergence() {
	// import fmt
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types

	initial := types.S8("hello")

	// two changes: append " world" and delete "lo"
	insert := changes.Splice{Offset: 5, Before: types.S8(""), After: types.S8(" world")}
	remove := changes.Splice{Offset: 3, Before: types.S8("lo"), After: types.S8("")}

	// two versions derived from initial
	inserted := initial.Apply(nil, insert)
	removed := initial.Apply(nil, remove)

	// merge the changes
	removex, insertx := insert.Merge(remove)

	// converge by applying the above
	final1 := inserted.Apply(nil, removex)
	final2 := removed.Apply(nil, insertx)

	fmt.Println(final1, final1 == final2)
	// Output: hel world true

}

func Example_convergence_streams() {
	// import fmt
	// import github.com/dotchain/dot/streams

	initial := streams.S8{Stream: streams.New(), Value: "hello"}

	// two changes: append " world" and delete "lo"
	s1 := initial.Splice(5, 0, " world")
	s2 := initial.Splice(3, len("lo"), "")

	// streams automatically merge because they are both
	// based on initial
	s1 = s1.Latest()
	s2 = s2.Latest()

	fmt.Println(s1.Value, s1.Value == s2.Value)
	// Output: hel world true

}

func Example_undo_streams() {
	// import fmt
	// import github.com/dotchain/dot/streams
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types
	// import github.com/dotchain/dot/streams/undo

	// create master, undoable child and the undo stack itself
	master := &streams.S16{Stream: streams.New(), Value: "hello"}
	s, stack := undo.New(master.Stream)
	undoableChild := &streams.S16{Stream: s, Value: master.Value}

	// change hello => Hello
	undoableChild = undoableChild.Splice(0, len("h"), "H")
	fmt.Println(undoableChild.Value)

	// for kicks, update master hello => hello$ as if it came
	// from the server
	master.Splice(len("hello"), 0, "$")

	// now undo this via the stack
	stack.Undo()

	// now undoableChild should be hello$
	undoableChild = undoableChild.Latest()
	fmt.Println(undoableChild.Value)

	// now redo the last operation to get Hello$
	stack.Redo()
	undoableChild = undoableChild.Latest()
	fmt.Println(undoableChild.Value)

	// Output:
	// Hello
	// hello$
	// Hello$

}

func Example_folding() {
	// import fmt
	// import github.com/dotchain/dot/streams
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types
	// import github.com/dotchain/dot/x/fold

	// create master, folded child and the folding itself
	master := &streams.S16{Stream: streams.New(), Value: "hello world!"}
	foldChange := changes.Splice{
		Offset: len("hello"),
		Before: types.S16(" world"),
		After:  types.S16("..."),
	}
	foldedStream := fold.New(foldChange, master.Stream)
	folded := &streams.S16{Stream: foldedStream, Value: "hello...!"}

	// folded:  hello...! => Hello...!!!
	folded = folded.Splice(0, len("h"), "H")
	folded = folded.Splice(len("Hello...!"), 0, "!!")
	fmt.Println(folded.Value)

	// master: hello world => hullo world
	master = master.Splice(len("h"), len("e"), "u")
	fmt.Println(master.Value)

	// now folded = Hullo...!!!
	fmt.Println(folded.Latest().Value)

	// master = Hullo world!!!
	fmt.Println(master.Latest().Value)

	// Output:
	// Hello...!!!
	// hullo world!
	// Hullo...!!!
	// Hullo world!!!

}

func Example_branching() {
	// import fmt
	// import github.com/dotchain/dot/streams
	// import github.com/dotchain/dot/changes
	// import github.com/dotchain/dot/changes/types

	// local is a branch of master
	master := &streams.S16{Stream: streams.New(), Value: "hello"}
	local := &streams.S16{Stream: streams.Branch(master.Stream), Value: master.Value}

	// edit locally: hello => hallo
	local.Splice(len("h"), len("e"), "a")

	// changes will not be reflected on master yet
	fmt.Println(master.Latest().Value)

	// push local changes up to master now
	streams.Push(local.Stream)

	// now master = hallo
	fmt.Println(master.Latest().Value)

	// Output:
	// hello
	// hallo

}
