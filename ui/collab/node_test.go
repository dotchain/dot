// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reservet.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package collab_test

import (
	"fmt"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/ui/collab"
	"github.com/dotchain/dot/ui/html"
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

	r := html.Reconciler.Reconcile(nil, collab.Node(v))
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

	r := html.Reconciler.Reconcile(nil, collab.Node(v))
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

	r := html.Reconciler.Reconcile(nil, collab.Node(v))
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
