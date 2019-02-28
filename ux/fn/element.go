// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fn

import "github.com/dotchain/dot/ux/core"

//go:generate go run codegen.go - $GOFILE

// Element represents a DOM element
//
// codegen: pure
func Element(c *elementCtx, props core.Props, children ...core.Element) core.Element {
	return c.rootElement("root").reconcile(props, children)
}

// codegen: pure
func rootElement(c *rootElementCtx) *element {
	return &element{}
}

type element struct {
	root  core.Element
	props core.Props
}

func (e *element) reconcile(props core.Props, children []core.Element) core.Element {
	children = e.filterNil(children)

	if e.root == nil {
		e.root = core.NewElement(props, children...)
		e.props = props
		return e.root
	}

	if e.props != props {
		before, after := e.props.ToMap(), props.ToMap()
		e.props = props
		for k, v := range after {
			if before[k] != v {
				e.root.SetProp(k, v)
			}
		}
	}
	e.updateChildren(children)
	return e.root
}

func (e *element) filterNil(children []core.Element) []core.Element {
	result := children[:0]
	for _, elt := range children {
		if elt != nil {
			result = append(result, elt)
		}
	}
	return result
}

func (e *element) updateChildren(after []core.Element) {
	for _, op := range e.bestDiff(e.root.Children(), after, 0, nil) {
		if op.insert {
			e.root.InsertChild(op.index, op.elt)
		} else {
			e.root.RemoveChild(op.index)
		}
	}
}

type diff struct {
	insert bool
	elt    core.Element
	index  int
}

func (e *element) bestDiff(before, after []core.Element, offset int, ops []diff) []diff {
	for len(before) > 0 && len(after) > 0 && before[0] == after[0] {
		offset++
		before, after = before[1:], after[1:]
	}

	switch {
	case len(before) == 0:
		for _, elt := range after {
			ops = append(ops, diff{true, elt, offset})
			offset++
		}
	case len(after) == 0:
		for range before {
			ops = append(ops, diff{false, nil, offset})
		}
	default:
		ops = e.chooseDiff(before, after, offset, ops)
	}

	return ops
}

func (e *element) chooseDiff(before, after []core.Element, offset int, ops []diff) []diff {
	// choice1 = clone of ops + delete first before elt
	choice1 := append(ops, diff{false, nil, offset})
	choice1 = e.bestDiff(before[1:], after, offset, choice1)

	index := e.indexOf(before[0], after)
	if index == -1 {
		return choice1
	}

	// choice2 = clone of ops + insert index after elts
	choice2 := append([]diff(nil), ops...)
	for kk := 0; kk < index+1; kk++ {
		choice2 = append(choice2, diff{true, after[kk], offset + kk})
	}
	choice2 = e.bestDiff(before, after[index+1:], offset+index+1, choice2)
	if len(choice1) < len(choice2) {
		return choice1
	}
	return choice2
}

func (e *element) indexOf(elt core.Element, elts []core.Element) int {
	for kk, elt1 := range elts {
		if elt1 == elt {
			return kk
		}
	}
	return -1
}
