// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"fmt"
	"github.com/dotchain/dot/changes"
)

// A implements the Value interface over a slice of values
type A []changes.Value

func (a A) Slice(offset, count int) changes.Value {
	return a[offset : offset+count]
}

func (a A) Count() int {
	return len(a)
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
		ox, cx, dx := c.Offset, c.Count, c.Distance
		if dx < 0 {
			ox, cx, dx = ox+dx, -dx, cx
		}
		l, m1, m2, r := a[:ox:ox], a[ox:ox+cx], a[ox+cx:ox+cx+dx], a[ox+cx+dx:]
		return A(append(append(append(l, m2...), m1...), r...))
	case changes.ChangeSet:
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
