// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reservet.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package collab

import (
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/x/ui/dom"
	"sort"
)

// Node implements dom.Node on top of a Text value
func Node(t Text) dom.Node {
	return node(t)
}

type node Text

func (n node) Tag() string {
	return "div"
}

func (n node) Key() interface{} {
	return nil
}

func (n node) ForEachAttribute(fn func(key string, val interface{})) {
	fn("contenteditable", "true")

	// TODO: get proper focus ID
	fn("focus", 1)
	kbd := Keyboard(Text(n))
	fn("keyboard", &kbd)
}

func (n node) ForEachNode(fn func(n dom.Node)) {
	last := 0
	n.forEachCaret(last, fn)
	for _, next := range n.indices() {
		if next > 0 {
			n.forEachRange(last, next, fn)
			last = next
			n.forEachCaret(last, fn)
		}
	}
	if size := len(n.Text); size > last {
		n.forEachRange(last, size, fn)
		n.forEachCaret(size, fn)
	}
}

func (n node) indices() []int {
	indices := map[int]bool{}
	result := []int(nil)
	for _, ref := range n.Refs {
		indices[ref.(refs.Range).Start.Index] = true
		indices[ref.(refs.Range).End.Index] = true
	}
	for idx := range indices {
		result = append(result, idx)
	}
	sort.Ints(result)
	return result
}

func (n node) forEachCaret(idx int, fn func(dom.Node)) {
	var own, other bool
	for id, ref := range n.Refs {
		s, e := ref.(refs.Range).Start.Index, ref.(refs.Range).End.Index
		if s == e && s == idx {
			own = own || (id == n.SessionID)
			other = other || (id != n.SessionID)
		}
	}
	switch {
	case own && other:
		fn(region{"caret both", ""})
	case own:
		fn(region{"caret own", ""})
	case other:
		fn(region{"caret other", ""})
	}
}

func (n node) forEachRange(start, end int, fn func(dom.Node)) {
	var own, other bool
	for id, ref := range n.Refs {
		s, e := ref.(refs.Range).Start.Index, ref.(refs.Range).End.Index
		if s <= start && e >= end || e <= start && s >= end {
			own = own || (id == n.SessionID)
			other = other || (id != n.SessionID)
		}
	}

	text := n.Text[start:end]
	switch {
	case own && other:
		fn(region{"range both", text})
	case own:
		fn(region{"range own", text})
	case other:
		fn(region{"range other", text})
	default:
		fn(region{"", text})
	}
}

type region struct {
	class, text string
}

func (r region) Tag() string {
	return "span"
}

func (r region) Key() interface{} {
	return nil
}

func (r region) ForEachAttribute(fn func(key string, val interface{})) {
	if r.class != "" {
		fn("class", r.class)
	}
}

func (r region) ForEachNode(fn func(n dom.Node)) {
	if r.text != "" {
		fn(textnode(r.text))
	}
}

type textnode string

func (t textnode) Tag() string {
	return ":text:"
}

func (t textnode) Key() interface{} {
	return nil
}

func (t textnode) ForEachAttribute(fn func(key string, val interface{})) {
	fn(":data:", string(t))
}

func (t textnode) ForEachNode(fn func(n dom.Node)) {
}
