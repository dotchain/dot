// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package tree

import "github.com/dotchain/dot/changes"

// Nodes represents a sequence of nodes
type Nodes []*Node

// Slice implements changes.Collection
func (n Nodes) Slice(offset, count int) changes.Collection {
	return n[offset : offset+count]
}

// Count implements changes.Collection
func (n Nodes) Count() int {
	return len(n)
}

// ApplyCollection implements changes.Collection
func (n Nodes) ApplyCollection(c changes.Change) changes.Collection {
	switch c := c.(type) {
	case changes.Splice:
		remove := c.Before.Count()
		after := c.After.(Nodes)
		start, end := c.Offset, c.Offset+remove
		return append(append(n[:start:start], after...), n[end:]...)
	case changes.Move:
		ox, cx, dx := c.Offset, c.Count, c.Distance
		if dx < 0 {
			ox, cx, dx = ox+dx, -dx, cx
		}
		x1, x2, x3 := ox, ox+cx, ox+cx+dx
		return append(append(append(n[:x1:x1], n[x2:x3]...), n[x1:x2]...), n[x3:]...)
	case changes.PathChange:
		idx := c.Path[0].(int)
		clone := append([]*Node(nil), n...)
		v := clone[idx].Apply(changes.PathChange{c.Path[1:], c.Change})
		clone[idx] = v.(*Node)

		return Nodes(clone)
	}
	panic("Unexpected change on Apply")
}

// Apply implements changes.Value
func (n Nodes) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return n
	case changes.Replace:
		if !c.IsCreate() {
			return c.After
		}
	case changes.PathChange:
		if len(c.Path) == 0 {
			return n.Apply(c.Change)
		}
	case changes.Custom:
		return c.ApplyTo(n)
	}
	return n.ApplyCollection(c)
}
