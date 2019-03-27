// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package heap implements a heap value type
package heap

import (
	"math/rand"
	"sort"
	"time"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

type entry struct {
	key1, key2 int64
	rank       int
}

// Heap implements a heap of keys
type Heap struct {
	items map[interface{}]entry
}

// UpdateChange returns the change to apply for either inserting an
// item into the heap or updating its rank.
func (h Heap) UpdateChange(key interface{}, rank int) changes.Change {
	if entry, ok := h.items[key]; ok && entry.rank == rank {
		return nil
	}

	v := changes.Atomic{Value: entry{key1: rand.Int63(), key2: rand.Int63(), rank: rank}}
	p := []interface{}{key}
	return changes.PathChange{Path: p, Change: changes.Replace{Before: h.get(key), After: v}}
}

// DeleteChange returns the change to apply for deleting an item from the heap
func (h Heap) DeleteChange(key interface{}) changes.Change {
	if _, ok := h.items[key]; !ok {
		return nil
	}

	p := []interface{}{key}
	return changes.PathChange{Path: p, Change: changes.Replace{Before: h.get(key), After: changes.Nil}}
}

// Update inserts a new entry wite provided rank or updates the rank
// for an existng entry
func (h Heap) Update(key interface{}, rank int) Heap {
	return h.Apply(nil, h.UpdateChange(key, rank)).(Heap)
}

// Delete removes the entry from the heap
func (h Heap) Delete(key interface{}) Heap {
	return h.Apply(nil, h.DeleteChange(key)).(Heap)
}

// Apply implments Changes.Value
func (h Heap) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: h.get, Set: h.set}).Apply(ctx, c, h)
}

// Iterate calls the provided function for each entry in the heap in
// descending order of rank value.
func (h Heap) Iterate(fn func(key interface{}, rank int) bool) {
	keys := make([]interface{}, 0, len(h.items))
	for k := range h.items {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		v1 := h.items[keys[j]]
		v2 := h.items[keys[i]]
		if v1.rank != v2.rank {
			return v1.rank < v2.rank
		}
		if v1.key1 != v2.key1 {
			return v1.key1 < v2.key1
		}
		return v1.key2 < v2.key2
	})
	for _, key := range keys {
		if !fn(key, h.items[key].rank) {
			break
		}
	}
}

func (h Heap) get(key interface{}) changes.Value {
	if entry, ok := h.items[key]; ok {
		return changes.Atomic{Value: entry}
	}
	return changes.Nil
}

func (h Heap) set(key interface{}, v changes.Value) changes.Value {
	clone := Heap{items: map[interface{}]entry{}}
	for k, v := range h.items {
		if k != key {
			clone.items[k] = v
		}
	}
	if v != changes.Nil {
		clone.items[key] = v.(changes.Atomic).Value.(entry)
	}
	return clone
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
