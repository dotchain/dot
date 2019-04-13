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

type Entry struct {
	Key1, Key2 int64
	Rank       int
}

// Heap implements a heap of keys
type Heap struct {
	Items map[interface{}]Entry
}

// UpdateChange returns the change to apply for either inserting an
// item into the heap or updating its Rank.
func (h Heap) UpdateChange(key interface{}, rank int) changes.Change {
	if Entry, ok := h.Items[key]; ok && Entry.Rank == rank {
		return nil
	}

	v := changes.Atomic{Value: Entry{Key1: rand.Int63(), Key2: rand.Int63(), Rank: rank}}
	p := []interface{}{key}
	return changes.PathChange{Path: p, Change: changes.Replace{Before: h.get(key), After: v}}
}

// DeleteChange returns the change to apply for deleting an item from the heap
func (h Heap) DeleteChange(key interface{}) changes.Change {
	if _, ok := h.Items[key]; !ok {
		return nil
	}

	p := []interface{}{key}
	return changes.PathChange{Path: p, Change: changes.Replace{Before: h.get(key), After: changes.Nil}}
}

// Update inserts a new Entry wite provided Rank or updates the Rank
// for an existng Entry
func (h Heap) Update(key interface{}, rank int) Heap {
	return h.Apply(nil, h.UpdateChange(key, rank)).(Heap)
}

// Delete removes the Entry from the heap
func (h Heap) Delete(key interface{}) Heap {
	return h.Apply(nil, h.DeleteChange(key)).(Heap)
}

// Apply implments Changes.Value
func (h Heap) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Get: h.get, Set: h.set}).Apply(ctx, c, h)
}

// Iterate calls the provided function for each Entry in the heap in
// descending order of Rank value.
func (h Heap) Iterate(fn func(key interface{}, rank int) bool) {
	keys := make([]interface{}, 0, len(h.Items))
	for k := range h.Items {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		v1 := h.Items[keys[j]]
		v2 := h.Items[keys[i]]
		if v1.Rank != v2.Rank {
			return v1.Rank < v2.Rank
		}
		if v1.Key1 != v2.Key1 {
			return v1.Key1 < v2.Key1
		}
		return v1.Key2 < v2.Key2
	})
	for _, key := range keys {
		if !fn(key, h.Items[key].Rank) {
			break
		}
	}
}

func (h Heap) get(key interface{}) changes.Value {
	if Entry, ok := h.Items[key]; ok {
		return changes.Atomic{Value: Entry}
	}
	return changes.Nil
}

func (h Heap) set(key interface{}, v changes.Value) changes.Value {
	clone := Heap{Items: map[interface{}]Entry{}}
	for k, v := range h.Items {
		if k != key {
			clone.Items[k] = v
		}
	}
	if v != changes.Nil {
		clone.Items[key] = v.(changes.Atomic).Value.(Entry)
	}
	return clone
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
