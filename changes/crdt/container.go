// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt

import "github.com/dotchain/dot/changes"

// Container holds a single CRDT value.
//
// This is implemented as a map of unique rank to values. Every update
// creates a new entry in this map.  Undo of an update causes the undo
// count for that rank to be incremented (and undo of the undo causes
// the count to be decremented).
//
// Tho whole container itself can be deleted (with Deleted tracking
// the count) and undeleted.
type Container struct {
	Entries map[*Rank]interface{}
	Undos   map[*Rank]int
	Deleted int
}

// Set updates the value using the provided rank. The effective value
// is given by the largest undeleted rank.  Reverting the returned
// change will simply cause the undo count of the rank to be modified.
func (c Container) Set(r *Rank, v interface{}) (changes.Change, Container) {
	cx := wrapper{setContainer{r, v}}
	return cx, cx.ApplyTo(nil, c).(Container)
}

// Delete deletes the whole container
func (c Container) Delete() (changes.Change, Container) {
	cx := wrapper{delContainer(1)}
	return cx, cx.ApplyTo(nil, c).(Container)
}

// Update wraps a change meant for the value within the container.
func (c Container) Update(inner changes.Change) (changes.Change, Container) {
	r, _ := c.Get()
	cx := wrapper{updContainer{r, inner}}
	return cx, cx.ApplyTo(nil, c).(Container)
}

// Get returns the latest undeleted value or nil if no such value exists.
func (c Container) Get() (*Rank, interface{}) {
	if c.Deleted > 0 {
		return nil, nil
	}

	var r *Rank
	var result interface{}
	for rank, val := range c.Entries {
		if c.Undos[rank] > 0 {
			continue
		}
		if r == nil || r.Less(rank) {
			r, result = rank, val
		}
	}
	return r, result
}

// Clone duplicates the whole container
func (c Container) Clone() Container {
	entries := map[*Rank]interface{}{}
	for r, v := range c.Entries {
		entries[r] = v
	}
	undos := map[*Rank]int{}
	for r, u := range c.Undos {
		undos[r] = u
	}
	return Container{Entries: entries, Undos: undos, Deleted: c.Deleted}
}

// Apply implements changes.Value
func (c Container) Apply(ctx changes.Context, cx changes.Change) changes.Value {
	return cx.(changes.Custom).ApplyTo(ctx, c)
}

type setContainer struct {
	*Rank
	Value interface{}
}

func (sc setContainer) Revert() crdtChange {
	return unsetContainer{sc.Rank, 1}
}

func (sc setContainer) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Container).Clone()
	result.Entries[sc.Rank] = sc.Value
	return result
}

type unsetContainer struct {
	*Rank
	UndoCount int
}

func (uc unsetContainer) Revert() crdtChange {
	return unsetContainer{uc.Rank, -1}
}

func (uc unsetContainer) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Container).Clone()
	result.Undos[uc.Rank] += uc.UndoCount
	return result
}

type delContainer int

func (dc delContainer) Revert() crdtChange {
	return delContainer(-1)
}

func (dc delContainer) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Container).Clone()
	result.Deleted += int(dc)
	return result
}

type updContainer struct {
	*Rank
	Change changes.Change
}

func (uc updContainer) Revert() crdtChange {
	return updContainer{uc.Rank, uc.Change.Revert()}
}

func (uc updContainer) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Container).Clone()
	result.Entries[uc.Rank] = result.Entries[uc.Rank].(changes.Value).Apply(ctx, uc.Change)
	return result
}
