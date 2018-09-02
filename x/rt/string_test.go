// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rt_test

import (
	"fmt"
	"github.com/dotchain/dot/changes"
)

// S implements the Value interface. Note that this implementation is
// only safe when using a byte encoding and so is not entirely
// portable across other languages.  It is meant for test purposes
// only.
type S string

func (s S) Slice(offset, count int) changes.Value {
	return S(s[offset : offset+count])
}

func (s S) Count() int {
	return len(s)
}

func (s S) toString(v changes.Value) string {
	if v == changes.Nil {
		return ""
	}
	return string(v.(S))
}

func (s S) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return s
	case changes.Replace:
		if c.IsDelete {
			return changes.Nil
		}
		return c.After
	case changes.Splice:
		remove := len(s.toString(c.Before))
		after := S(s.toString(c.After))
		return S(s[:c.Offset] + after + s[c.Offset+remove:])
	case changes.Move:
		ox, cx, dx := c.Offset, c.Count, c.Distance
		if dx < 0 {
			ox, cx, dx = ox+dx, -dx, cx
		}
		l, m1, m2, r := s[:ox], s[ox:ox+cx], s[ox+cx:ox+cx+dx], s[ox+cx+dx:]
		return S(l + m2 + m1 + r)
	case changes.ChangeSet:
		v := changes.Value(s)
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
			return s.Apply(c.Change)
		}
	}
	panic(fmt.Sprintf("Unexpected Apply %#v", c))
}
