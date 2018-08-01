// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"github.com/dotchain/dot/conv"
	"github.com/dotchain/dot/encoding"
)

// Utils implements a bunch of utlities on top of transformer.
type Utils Transformer

// AreSame compares if two interfaces can be treated as being
// equivalent.
func (u Utils) AreSame(i1, i2 interface{}) bool {
	if encoding.IsString(i1) && encoding.IsString(i2) {
		return encoding.ToString(i1) == encoding.ToString(i2)
	}

	x1, ok1 := u.C.TryGet(i1)
	x2, ok2 := u.C.TryGet(i2)

	if ok1 != ok2 {
		return false
	}

	if !ok1 || !ok2 {
		return deepEqual(i1, i2)
	}

	if x1 == nil {
		return u.isEmpty(x2)
	}

	if x2 == nil {
		return u.isEmpty(x1)
	}

	if x1.IsArray() != x2.IsArray() {
		return false
	}

	if x1.IsArray() {
		c := x1.Count()
		if c != x2.Count() {
			return false
		}
		failed := false
		x1.ForEach(func(kk int, val interface{}) {
			failed = failed || !u.AreSame(val, x2.Get(conv.FromIndex(kk)))
		})
		return !failed
	}

	failed, seen, count := false, map[string]bool{}, 0
	x1.ForKeys(func(key string, val interface{}) { seen[key] = true })
	x2.ForKeys(func(key string, val interface{}) {
		if failed || !seen[key] {
			failed = true
			return
		}
		failed = !u.AreSame(val, x1.Get(key))
		count++
	})
	return !failed && count == len(seen)
}

// TryApply tries to apply the input method and returns the applied
// value if it succeeded.  If it failed, it returns nil and sets ok
// to false.
func (u Utils) TryApply(obj interface{}, changes []Change) (result interface{}, ok bool) {
	// we do a generic catch because the inner tryApplyXYZ only checks for the
	// set of activities that are not covered by the actual apply mechanism --
	// mainly the "before" activities. Egregious type mismatches will happen
	// via panics still as the encoding interface currently does not provide
	// a way to return false.  There is no point dealing with that because
	// those are not supposed to happen.  OTOH, the other checks that return
	// bool can happen with good clients if there is an odd merge with
	// RangeApply.  This just ensures that panic only happens if there is an
	// actual client bug.
	defer func() {
		if r := recover(); r != nil {
			// log.Println("TryApply panic'ed", r)
			result = nil
			ok = false
		}
	}()

	r := obj
	ok = true
	for _, change := range changes {
		r, ok = u.tryApplyChange(r, change.Path, func(input interface{}) (interface{}, bool) {
			switch {
			case change.Splice != nil:
				splice := change.Splice
				return u.tryApplySplice(input, splice.Offset, splice.Before, splice.After)
			case change.Move != nil:
				move := change.Move
				return u.applyMove(input, move.Offset, move.Count, move.Distance), true
			case change.Set != nil:
				set := change.Set
				return u.tryApplySet(input, set.Key, set.Before, set.After)
			case change.Range != nil:
				r := change.Range
				return u.tryApplyRange(input, r.Offset, r.Count, r.Changes)
			}
			// log.Println("Ignoring empty", input)
			return input, true
		})
		if !ok {
			break
		}
	}
	return r, ok
}

// Apply applies a sequence of changes to the provided input
// and returns the outpuu. This is not really meant for client
// models to use.
func (u Utils) Apply(obj interface{}, changes []Change) interface{} {
	return u.check(u.TryApply(obj, changes))
}

func (u Utils) tryApplyChange(obj interface{}, path []string, fn func(interface{}) (interface{}, bool)) (interface{}, bool) {
	if len(path) == 0 {
		return fn(obj)
	}

	if data, ok := u.C.TryGet(obj); ok {
		inner := data.Get(path[0])
		if updated, ok := u.tryApplyChange(inner, path[1:], fn); ok {
			return data.Set(path[0], updated), true
		}
	}
	return nil, false
}

func (u Utils) tryApplySplice(input interface{}, offset int, before, after interface{}) (interface{}, bool) {
	data, ok := u.C.TryGet(input)
	if !ok {
		return nil, false
	}

	if !u.AreSame(before, nil) {
		count := u.C.Get(before).Count()
		if !u.AreSame(data.Slice(offset, count), before) {
			return nil, false
		}
	}

	return data.Splice(offset, before, after), true
}

func (u Utils) applySplice(input interface{}, offset int, before, after interface{}) interface{} {
	return u.check(u.tryApplySplice(input, offset, before, after))
}

func (u Utils) applyMove(input interface{}, offset, count, distance int) interface{} {
	data := u.C.Get(input)
	before := data.Slice(offset, count)
	return data.Splice(offset, before, nil).Splice(offset+distance, nil, before)
}

func (u Utils) tryApplySet(input interface{}, key string, before, after interface{}) (interface{}, bool) {
	i := u.C.Get(input)
	if u.AreSame(before, nil) {
		found := false
		i.ForKeys(func(k string, v interface{}) {
			found = found || k == key
		})
		if found {
			return nil, false
		}
	} else if !u.AreSame(i.Get(key), before) {
		return nil, false
	}

	return i.Set(key, after), true
}

func (u Utils) tryApplyRange(input interface{}, offset, count int, changes interface{}) (res interface{}, ok bool) {
	defer func() {
		if r := recover(); r != nil || !ok {
			ok = false
			res = nil
		}
	}()

	c := changes.([]Change)
	ok = true
	res = u.C.Get(input).RangeApply(offset, count, func(before interface{}) interface{} {
		if ok {
			if inner, innerOK := u.TryApply(before, c); innerOK {
				return inner
			}
		}
		ok = false
		return before
	})
	return res, ok
}

func (u Utils) applyRange(input interface{}, offset, count int, changes interface{}) interface{} {
	return u.check(u.tryApplyRange(input, offset, count, changes))
}

func (u Utils) isEmpty(e encoding.UniversalEncoding) bool {
	if e.IsArray() {
		return e.Count() == 0
	}
	count := 0
	e.ForKeys(func(string, interface{}) { count++ })
	return count == 0
}

func (u Utils) check(i interface{}, ok bool) interface{} {
	if ok {
		return i
	}
	return nil
}
