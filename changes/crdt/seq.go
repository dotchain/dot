// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt

import (
	"sort"

	"github.com/dotchain/dot/changes"
)

// Seq implements a CRDT-style sequence
type Seq struct {
	Values Dict
	Ords   Dict
}

// Items returns the seq as an array
func (s Seq) Items() []interface{} {
	_, keys := s.items()
	result := make([]interface{}, len(keys))
	for kk, key := range keys {
		_, result[kk] = s.Values.Get(key)
	}
	return result
}

// Splice replaces the sub-sequence [offset:offset+remove]
func (s Seq) Splice(offset, remove int, replacement []interface{}) (changes.Change, Seq) {
	ords, keys := s.items()
	result := wrapper{}
	for kk := 0; kk < remove; kk++ {
		inner, _ := s.Values.Delete(keys[kk+offset])
		result = append(result, updValueSeq{inner})
		inner, _ = s.Ords.Delete(keys[kk+offset])
		result = append(result, updOrdSeq{inner})
	}

	newOrds := s.between(ords, offset, remove, len(replacement))
	for kk, ord := range newOrds {
		key := NewRank()
		inner, _ := s.Ords.Set(key, ord)
		result = append(result, updOrdSeq{inner})
		inner, _ = s.Values.Set(key, replacement[kk])
		result = append(result, updValueSeq{inner})
	}

	return result, result.ApplyTo(nil, s).(Seq)
}

// Move shifts the sub sequence (offset, offset +count) by distance
func (s Seq) Move(offset, count, distance int) (changes.Change, Seq) {
	ords, keys := s.items()
	if distance > 0 {
		ords = s.between(ords, offset+count+distance, 0, count)
	} else {
		ords = s.between(ords, offset+distance, 0, count)
	}
	result := wrapper{}
	for kk := 0; kk < count; kk++ {
		inner, _ := s.Ords.Set(keys[offset+kk], ords[kk])
		result = append(result, updOrdSeq{inner})
	}
	return result, result.ApplyTo(nil, s).(Seq)
}

// Update takes a change meant for the value at a specific index
// and wraps it so that it can applied on the Seq
func (s Seq) Update(idx int, inner changes.Change) (changes.Change, Seq) {
	_, keys := s.items()
	c, _ := s.Values.Update(keys[idx], inner)
	c = wrapper{updValueSeq{c}}
	return c, c.(changes.Custom).ApplyTo(nil, s).(Seq)
}

// Apply implements changes.Value
func (s Seq) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return c.(changes.Custom).ApplyTo(ctx, s)
}

func (s Seq) between(ords []string, offset, remove, count int) []string {
	var start, end string
	if offset > 0 {
		start = ords[offset-1]
	} else if len(ords) > 0 {
		start = PrevOrd(ords[0])
	}

	switch {
	case offset+remove < len(ords):
		end = ords[offset+remove]
	case len(ords) > 0:
		end = NextOrd(ords[len(ords)-1])
	default:
		end = NextOrd(start)
	}

	return BetweenOrd(start, end, count)
}

func (s Seq) items() ([]string, []interface{}) {
	ords := []string{}
	keysmap := map[string]interface{}{}
	for key, container := range s.Ords.Entries {
		if x, _ := s.Values.Get(key); x == nil {
			continue
		}

		if r, val := container.Get(); r != nil {
			ords = append(ords, val.(string))
			keysmap[val.(string)] = key
		}
	}
	sort.Slice(ords, func(i, j int) bool {
		return LessOrd(ords[i], ords[j])
	})
	keys := make([]interface{}, len(ords))
	for kk, o := range ords {
		keys[kk] = keysmap[o]
	}
	return ords, keys
}

type updValueSeq struct {
	changes.Change
}

func (u updValueSeq) Revert() crdtChange {
	return updValueSeq{u.Change.Revert()}
}

func (u updValueSeq) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Seq)
	result.Values = result.Values.Apply(ctx, u.Change).(Dict)
	return result
}

type updOrdSeq struct {
	changes.Change
}

func (u updOrdSeq) Revert() crdtChange {
	return updOrdSeq{u.Change.Revert()}
}

func (u updOrdSeq) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	result := v.(Seq)
	result.Ords = result.Ords.Apply(ctx, u.Change).(Dict)
	return result
}
