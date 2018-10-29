// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package tree

import (
	"github.com/dotchain/dot/changes"
	"github.com/tvastar/hamming"
)

// DiffNodes compares before and after and returns  a set of changes.
// Applying the changes to before is guaranteed to produce an
// equivalent value as after, though pointers may not match.
func DiffNodes(before, after Nodes) changes.Change {
	c := changes.ChangeSet(innerChanges(before, after))
	splice := func(offset int, b []hamming.Item, a []hamming.Item) {
		c = append(c, changes.Splice{offset, fromItems(b), fromItems(a)})
	}
	move := func(offset, count, distance int) {
		c = append(c, changes.Move{offset, count, distance})
	}
	hamming.Edits(toItems(before), toItems(after), splice, move)

	if len(c) == 0 {
		return nil
	}

	return c
}

func innerChanges(before, after Nodes) []changes.Change {
	result := []changes.Change(nil)

	counter := map[interface{}]int{}
	indices := map[interface{}]int{}
	for kk, n := range before {
		key := n.Key()
		pair := [2]interface{}{key, counter[key]}
		counter[key]++
		indices[pair] = kk
	}

	counter = map[interface{}]int{}
	for _, n := range after {
		key := n.Key()
		pair := [2]interface{}{key, counter[key]}
		counter[key]++

		idx, ok := indices[pair]
		if !ok {
			continue
		}

		c := Diff(before[idx], n)
		if c == nil {
			continue
		}
		c = changes.PathChange{[]interface{}{idx}, c}
		result = append(result, c)
	}

	return result
}

func toItems(n Nodes) []hamming.Item {
	items := make([]hamming.Item, len(n))
	for kk, nn := range n {
		items[kk] = nn
	}
	return items
}

func fromItems(items []hamming.Item) Nodes {
	result := make(Nodes, len(items))
	for kk, ii := range items {
		result[kk] = ii.(*Node)
	}
	return result
}
