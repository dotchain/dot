// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"github.com/dotchain/dot/ui/html"
	"testing"
)

func TestMutableNode(t *testing.T) {
	if n, err := html.Parse("<html><html>"); err == nil {
		t.Error("Unexpected successful parse", n)
	}

	tests := map[string]string{
		`<hello x="a">boo</hello>`:                 `<hello x="b" y="c">booya</hello>`,
		`<hello id="a">boo</hello>`:                `<hello id="b">booya</hello>`,
		`<x><z id="b">boo</z><y id="a">ok</y></x>`: `<x><y id="a">ok</y><z id="b">boo</z></x>`,
	}

	for before, after := range tests {
		t.Run(before+"=>"+after, func(t *testing.T) {
			validate(t, before, after)
			validate(t, after, before)
		})
	}
}

func validate(t *testing.T, before, after string) {
	b, _ := html.Parse(before)
	a, _ := html.Parse(after)
	result := html.Reconciler(nil, nil).Reconcile(b, a)
	if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", a) {
		t.Error("Mismatched", a, result)
	}
}

func TestEventHandlers(t *testing.T) {
	e := html.Events{}
	r := html.Reconciler(e, html.Keyboard{})
	div, _ := html.Parse("<div><div></div></div>")

	clicked := ""
	div.Children().Next().SetAttribute("onclick", func(arg interface{}) {
		clicked = arg.(string)
	})

	node := div.Children().Next().(html.Node).Node
	e.Fire(node, "onclick", "hello")
	if clicked != "" {
		t.Fatal("Unexpected firing yo")
	}

	root := r.Reconcile(nil, div)
	node = root.Children().Next().(html.Node).Node
	e.Fire(node, "onclick", "boo")
	if clicked != "boo" {
		t.Fatal("Firing failed", clicked)
	}

	div.Children().Next().RemoveAttribute("onclick")
	root = r.Reconcile(root, div)

	e.Fire(node, "onclick", "boohoo")
	if clicked != "boo" {
		t.Fatal("Firing failed", clicked)
	}

	div.Children().Next().SetAttribute("onclick", func(arg interface{}) {
		clicked = arg.(string)
	})
	r.Reconcile(root, div)
	div.Children().Remove()
	r.Reconcile(root, div)

	e.Fire(node, "onclick", "boohoo")
	if clicked != "boo" {
		t.Fatal("Firing failed", clicked)
	}
}

func TestFocusHandling(t *testing.T) {
	k := html.Keyboard{}
	r := html.Reconciler(html.Events{}, k)
	div, _ := html.Parse("<div><div></div></div>")
	event := ""
	div.Children().Next().SetAttribute("keyboard", kbd(func(s string) {
		event = s
	}))

	root := r.Reconcile(nil, div)
	if k.Focus() != nil {
		t.Fatal("Unexpected keyboard")
	}

	div.Children().Next().SetAttribute("focus", -2)
	root = r.Reconcile(root, div)
	if k.Focus() != nil {
		t.Fatal("Unexpected keyboard")
	}

	div.Children().Next().SetAttribute("focus", 1)
	root = r.Reconcile(root, div)
	if k.Focus() == nil {
		t.Fatal("Unexpected keyboard")
	}
	k.Focus().Insert("x")
	if event != "Insert: x" {
		t.Error("Unexpected event", event)
	}

	div.Children().Next().RemoveAttribute("focus")
	root = r.Reconcile(root, div)
	if k.Focus() != nil {
		t.Fatal("Unexpected keyboard")
	}

	div.Children().Next().RemoveAttribute("keyboard")
	r.Reconcile(root, div)
	if k.Focus() != nil {
		t.Fatal("Unexpected keyboard")
	}
}

type kbd func(s string)

func (k kbd) Insert(ch string) {
	k("Insert: " + ch)
}
func (k kbd) Remove() {
	k("Remove")
}
func (k kbd) ArrowRight() {
	k("ArrowRight")
}

func (k kbd) ArrowLeft() {
	k("ArrowLeft")
}

func (k kbd) ShiftArrowRight() {
	k("ShiftArrowRight")
}

func (k kbd) ShiftArrowLeft() {
	k("ShiftArrowLeft")
}
