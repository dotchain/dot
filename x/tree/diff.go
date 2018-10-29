// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package tree

import "github.com/dotchain/dot/changes"

// Diff compares before and after and returns  a set of changes.
// Applying the changes to before is guaranteed to produce an
// equivalent value as after, though pointers may not match.
func Diff(before, after *Node) changes.Change {
	if before == after {
		return nil
	}

	if after == nil {
		return changes.Replace{before, changes.Nil}
	}

	if before == nil {
		return changes.Replace{changes.Nil, after}
	}

	c := changes.ChangeSet(nil)
	for k, v := range *before {
		pc := changes.PathChange{[]interface{}{k}, nil}
		v2, ok := (*after)[k]
		if !ok {
			pc.Change = changes.Replace{toValue(v), changes.Nil}
			c = append(c, pc)
			continue
		}
		if k == "Children" {
			if cx := DiffNodes(v.(Nodes), v2.(Nodes)); cx != nil {
				pc.Change = cx
				c = append(c, pc)
			}
			continue
		}
		if v != v2 {
			pc.Change = changes.Replace{toValue(v), toValue(v2)}
			c = append(c, pc)
		}
	}

	for k, v := range *after {
		if _, ok := (*before)[k]; ok {
			continue
		}
		cx := changes.Replace{changes.Nil, toValue(v)}
		pc := changes.PathChange{[]interface{}{k}, cx}
		c = append(c, pc)
	}

	if c == nil {
		return nil
	}

	return c
}

func toValue(v interface{}) changes.Value {
	switch v := v.(type) {
	case changes.Value:
		return v
	}
	return changes.Atomic{v}
}
