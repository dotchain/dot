// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package collab_test

import (
	"fmt"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/ui/collab"
	"github.com/dotchain/dot/x/ui/html"
	"github.com/yosssi/gohtml"
)

var p = refs.Path{"Value"}

func ExampleNode_overlappingRegions() {
	v := collab.Text{
		Text:      "Hello World",
		SessionID: "me",
		Refs: map[interface{}]refs.Ref{
			"me":    refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 4, true}},
			"alpha": refs.Range{refs.Caret{p, 0, true}, refs.Caret{p, 2, true}},
			"beta":  refs.Range{refs.Caret{p, 4, true}, refs.Caret{p, 5, true}},
			"gamma": refs.Range{refs.Caret{p, 5, true}, refs.Caret{p, 5, true}},
		},
		Stream: nil,
	}

	r := html.Reconciler(nil, nil).Reconcile(nil, collab.Node(v))
	condense := gohtml.Condense
	defer func() {
		gohtml.Condense = condense
	}()
	gohtml.Condense = true
	fmt.Printf("%s", gohtml.Format(fmt.Sprintf("%v", r)))

	// Output:
	// <div contenteditable="true">
	//   <span class="range other">H</span>
	//   <span class="range both">e</span>
	//   <span class="range own">ll</span>
	//   <span class="range other">o</span>
	//   <span class="caret other"></span>
	//   <span> World</span>
	// </div>
}

func ExampleNode_carets() {
	v := collab.Text{
		Text:      "Hello World",
		SessionID: "me",
		Refs: map[interface{}]refs.Ref{
			"me":    refs.Range{refs.Caret{p, 2, true}, refs.Caret{p, 2, true}},
			"alpha": refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
			"beta":  refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
			"gamma": refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
		},
		Stream: nil,
	}

	r := html.Reconciler(nil, nil).Reconcile(nil, collab.Node(v))
	condense := gohtml.Condense
	defer func() {
		gohtml.Condense = condense
	}()
	gohtml.Condense = true
	fmt.Printf("%s", gohtml.Format(fmt.Sprintf("%v", r)))

	// Output:
	// <div contenteditable="true">
	//   <span>H</span>
	//   <span class="caret other"></span>
	//   <span>e</span>
	//   <span class="caret own"></span>
	//   <span>llo World</span>
	// </div>
}

func ExampleNode_sharedCarets() {
	v := collab.Text{
		Text:      "Hello World",
		SessionID: "me",
		Refs: map[interface{}]refs.Ref{
			"me":    refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
			"alpha": refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
			"beta":  refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
			"gamma": refs.Range{refs.Caret{p, 1, true}, refs.Caret{p, 1, true}},
		},
		Stream: nil,
	}

	r := html.Reconciler(nil, nil).Reconcile(nil, collab.Node(v))
	condense := gohtml.Condense
	defer func() {
		gohtml.Condense = condense
	}()
	gohtml.Condense = true
	fmt.Printf("%s", gohtml.Format(fmt.Sprintf("%v", r)))

	// Output:
	// <div contenteditable="true">
	//   <span>H</span>
	//   <span class="caret both"></span>
	//   <span>ello World</span>
	// </div>
}

func ExampleNode_insert() {
	v := collab.Text{
		Text:      "Hullo World",
		SessionID: "me",
		Refs: map[interface{}]refs.Ref{
			"me": refs.Range{refs.Caret{p, 2, true}, refs.Caret{p, 2, true}},
		},
		Stream: streams.New(),
	}

	kbd := html.Keyboard{}
	r := html.Reconciler(nil, kbd).Reconcile(nil, collab.Node(v)).(html.Node)

	f := kbd.Focus()
	// a couple of random events
	f.ArrowLeft()
	f.ArrowRight()

	// convert Hullo World to Hey World via a sequence of events
	f.ArrowLeft()
	f.ArrowRight()
	f.ShiftArrowLeft()
	f.Insert("e")
	f.ShiftArrowRight()
	f.ShiftArrowRight()
	f.ShiftArrowRight()
	f.ShiftArrowRight()
	f.Remove()
	f.Insert("y")
	f.Insert(" ")

	// render the latest
	r = html.Reconciler(nil, kbd).Reconcile(r, collab.Node(v.Latest())).(html.Node)

	// display
	condense := gohtml.Condense
	defer func() {
		gohtml.Condense = condense
	}()
	gohtml.Condense = true
	fmt.Printf("%s", gohtml.Format(fmt.Sprintf("%v", r)))

	// Output:
	// <div contenteditable="true">
	//   <span>Hey </span>
	//   <span class="caret own"></span>
	//   <span>World</span>
	// </div>
}
