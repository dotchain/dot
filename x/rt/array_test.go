// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rt_test

import (
	"fmt"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rt"
)

// A implements the Value interface over a slice of values
type A []changes.Value

func (a A) Slice(offset, count int) changes.Value {
	return a[offset : offset+count]
}

func (a A) Count() int {
	return len(a)
}

func (a A) move(ox, cx, dx int) changes.Value {
	if dx < 0 {
		ox, cx, dx = ox+dx, -dx, cx
	}
	l, m1, m2, r := a[:ox:ox], a[ox:ox+cx], a[ox+cx:ox+cx+dx], a[ox+cx+dx:]
	return A(append(append(append(l, m2...), m1...), r...))
}

func (a A) applyChangeSet(c changes.ChangeSet) changes.Value {
	v := changes.Value(a)
	for _, cx := range c {
		if cx != nil {
			if v == nil {
				fmt.Println("Unexpected nil v", cx)
			}
			v = v.Apply(cx)
		}
	}
	return v
}

func (a A) run(offset, count int, c changes.Change) changes.Value {
	copy := append([]changes.Value(nil), a...)
	for kk := offset; kk < offset+count; kk++ {
		copy[kk] = copy[kk].Apply(c)
	}
	return A(copy)
}

func (a A) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return a
	case changes.Replace:
		if c.IsDelete {
			return changes.Nil
		}
		return c.After
	case changes.Splice:
		after := c.After.(A)
		result := append(a[:c.Offset:c.Offset], after...)
		return A(append(result, a[c.Offset+c.Before.Count():]...))
	case changes.Move:
		return a.move(c.Offset, c.Count, c.Distance)
	case rt.Run:
		return a.run(c.Offset, c.Count, c.Change)
	case changes.ChangeSet:
		return a.applyChangeSet(c)
	case changes.PathChange:
		if len(c.Path) == 0 {
			return a.Apply(c.Change)
		}
		idx := c.Path[0].(int)
		clone := append([]changes.Value(nil), a...)
		clone[idx] = clone[idx].Apply(changes.PathChange{c.Path[1:], c.Change})
		return A(clone)
	default:
		panic(fmt.Sprintf("Unexpected Apply %#v", c))
	}
}
