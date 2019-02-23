// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package simple

import "github.com/dotchain/dot/ux/core"

// Element is an simple element. It is not a component itself but
// is expected to be embedded within simple declarative components.
//
// Usage:
//
//      type MyComponent struct {
//           simple.Element,
//           ...
//      }
//
//      func NewMyComponent(....) *MyComponent {
//           result := &MyComponent{...}
//           result.Declare(props, children...)
//      }
//
//      func (c *MyComponent) Update(...) {
//           c.Declare(updatedProps, updatedChildren...)
//      }
//
// Note that each component still has to create all its children and
// reuse these from previous runs.
type Element struct {
	Root  core.Element
	props core.Props
}

// Declare creates the root element on first use or updates the
// element on subsequent use
func (e *Element) Declare(props core.Props, children ...core.Element) {
	children = e.filterNil(children)

	if e.Root == nil {
		e.Root = core.NewElement(props, children...)
		e.props = props
		return
	}

	if e.props != props {
		before, after := e.props.ToMap(), props.ToMap()
		e.props = props
		for k, v := range after {
			if before[k] != v {
				e.Root.SetProp(k, v)
			}
		}
	}
	e.updateChildren(children)
}

func (e *Element) filterNil(children []core.Element) []core.Element {
	result := children[:0]
	for _, elt := range children {
		if elt != nil {
			result = append(result, elt)
		}
	}
	return result
}

func (e *Element) updateChildren(after []core.Element) {
	for _, op := range bestDiff(e.Root.Children(), after, 0, nil) {
		if op.insert {
			e.Root.InsertChild(op.index, op.elt)
		} else {
			e.Root.RemoveChild(op.index)
		}
	}
}

type diff struct {
	insert bool
	elt    core.Element
	index  int
}

func bestDiff(before, after []core.Element, offset int, ops []diff) []diff {
	for len(before) > 0 && len(after) > 0 && before[0] == after[0] {
		offset++
		before, after = before[1:], after[1:]
	}

	switch {
	case len(before) == 0:
		for _, e := range after {
			ops = append(ops, diff{true, e, offset})
			offset++
		}
	case len(after) == 0:
		for range before {
			ops = append(ops, diff{false, nil, offset})
		}
	default:
		ops = chooseDiff(before, after, offset, ops)
	}

	return ops
}

func chooseDiff(before, after []core.Element, offset int, ops []diff) []diff {
	// choice1 = clone of ops + delete first before elt
	choice1 := append(ops, diff{false, nil, offset})
	choice1 = bestDiff(before[1:], after, offset, choice1)

	index := indexOf(before[0], after)
	if index == -1 {
		return choice1
	}

	// choice2 = clone of ops + insert index after elts
	choice2 := append([]diff(nil), ops...)
	for kk := 0; kk < index+1; kk++ {
		choice2 = append(choice2, diff{true, after[kk], offset + kk})
	}
	choice2 = bestDiff(before, after[index+1:], offset+index+1, choice2)
	if len(choice1) < len(choice2) {
		return choice1
	}
	return choice2
}

func indexOf(elt core.Element, elts []core.Element) int {
	for kk, e := range elts {
		if e == elt {
			return kk
		}
	}
	return -1
}
