// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom_test

import (
	"fmt"
	"github.com/dotchain/dot/ui/dom"
	"golang.org/x/net/html"
	"strings"
)

func ExampleReconciler_append() {
	before := "<div>hello world</div>"
	after := "<div>hello world<span>heya</span></div>"
	fragment, _ := html.ParseFragment(strings.NewReader(before), nil)
	initial := node{fragment[0].FirstChild.NextSibling}
	fragment, _ = html.ParseFragment(strings.NewReader(after), nil)
	expected := node{fragment[0].FirstChild.NextSibling}

	r := dom.Reconciler(newHTMLNode)
	reconciled := r.Reconcile(initial, expected)
	if reconciled != initial {
		fmt.Println("Unexpected reconciled output", toHTML(reconciled))
	}

	if toHTML(reconciled) != toHTML(expected) {
		fmt.Println("Unexpected reconciled output", toHTML(reconciled))
	}

	// Output:
}

func ExampleReconciler_reorder() {
	before := `<div><span id="1">one</span><span id="2">two</span></div>`
	after := `<div><span id="2">two</span><span id="1">one</span></div>`
	fragment, _ := html.ParseFragment(strings.NewReader(before), nil)
	initial := node{fragment[0].FirstChild.NextSibling.FirstChild}
	firstChild := initial.Node.FirstChild
	secondChild := firstChild.NextSibling
	fragment, _ = html.ParseFragment(strings.NewReader(after), nil)
	expected := node{fragment[0].FirstChild.NextSibling.FirstChild}

	r := dom.Reconciler(newHTMLNode)
	reconciled := r.Reconcile(initial, expected).(node)
	if reconciled != initial {
		fmt.Println("Unexpected reconciled output", toHTML(reconciled))
	}

	if toHTML(reconciled) != toHTML(expected) {
		fmt.Println("Unexpected reconciled output", toHTML(reconciled))
	}

	// confirm that the two children are properly swapped around instead
	// of new nodes being created
	if reconciled.Node.FirstChild.NextSibling != firstChild {
		fmt.Println("First child is unexpected")
	}

	if reconciled.Node.FirstChild != secondChild {
		fmt.Println("Second child is unexpected")
	}

	// Output:
}
