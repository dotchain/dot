// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import (
	"github.com/dotchain/dot/encoding"
	"reflect"
	"testing"
)

type Dict struct {
	initial, empty interface{}
	zeroValue      interface{}
	someValue      interface{}
	existingKeys   []string
	existingValues []interface{}
	nonExistingKey string
}

func DictTest() Dict {
	return Dict{
		initial:        map[string]interface{}{"hello": "world", "q": 42.1},
		empty:          map[string]interface{}{},
		someValue:      "doozy",
		existingKeys:   []string{"hello", "q"},
		existingValues: []interface{}{"world", 42.1},
		nonExistingKey: "world",
	}
}

func (d Dict) TestAll(t *testing.T) {
	t.Run("Empty", d.TestEmpty)
	t.Run("MarshalUnmarshal", d.TestMarshalUnmarshal)
	t.Run("ArrayBehaviors", d.TestArrayBehaviors)
	t.Run("NonEmpty", d.TestNonEmpty)
	t.Run("ForKeys", d.TestForKeys)
}

func (d Dict) TestMarshalUnmarshal(t *testing.T) {
	u := encoding.Get(d.empty)
	if !reflect.DeepEqual(unmarshal(marshal(u)), d.empty) {
		t.Error("Unmarshal/Marshal changed stuff", unmarshal(marshal(u)), u)
	}

	u = encoding.Get(d.initial)
	if !reflect.DeepEqual(unmarshal(marshal(u)), d.initial) {
		t.Errorf("Unmarshal/Marshal changed stuff %#v %#v\n", unmarshal(marshal(u)), d.initial)
	}
}

func (d Dict) TestEmpty(t *testing.T) {
	s := encoding.Get(d.empty)
	ensureEqual(t, s.IsArray(), false)

	shouldPanic(t, "fetching empty", func() { s.Get(d.nonExistingKey) })

	u := encoding.Get(s.Set(d.nonExistingKey, d.someValue))

	ensureEqual(t, u.Get(d.nonExistingKey), d.someValue)

	e1 := u.Set(d.nonExistingKey, nil)
	e2 := u.Set(d.nonExistingKey, d.zeroValue)
	ensureEqual(t, e1, d.empty)
	ensureEqual(t, e2, d.empty)
}

func (d Dict) TestNonEmpty(t *testing.T) {
	s := encoding.Get(d.initial)
	ensureEqual(t, s.IsArray(), false)

	for kk := range d.existingKeys {
		key, val := d.existingKeys[kk], d.existingValues[kk]
		ensureEqual(t, s.Get(key), val)
		ensureEqual(t, s.Set(key, d.someValue).Set(key, val).Get(key), val)
	}
	ensureEqual(t, s, s.Set(d.nonExistingKey, d.someValue).Set(d.nonExistingKey, nil))
	ensureEqual(t, s, s.Set(d.nonExistingKey, d.someValue).Set(d.nonExistingKey, d.zeroValue))

	updated := s.Set(d.nonExistingKey, d.someValue).Set(d.existingKeys[0], nil)
	ensureEqual(t, updated.Get(d.nonExistingKey), d.someValue)
	ensureEqual(t, updated.Get(d.existingKeys[1]), d.existingValues[1])

	shouldPanic(t, "Get deleted key", func() { updated.Get(d.existingKeys[0]) })
	shouldPanic(t, "Get non-existent key", func() { encoding.Get(d.initial).Get(d.nonExistingKey) })
}

func (d Dict) TestForKeys(t *testing.T) {
	x := encoding.Get(d.empty).(encoding.ObjectLike)
	expected, actual := map[string]interface{}{}, map[string]interface{}{}
	for kk := range d.existingKeys {
		x = x.Set(d.existingKeys[kk], d.existingValues[kk])
		expected[d.existingKeys[kk]] = d.existingValues[kk]
	}
	x.ForKeys(func(key string, val interface{}) {
		actual[key] = val
	})
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Mismatched", expected, actual)
	}
}

func (d Dict) TestArrayBehaviors(t *testing.T) {
	s := encoding.Get(d.initial)
	shouldPanic(t, "Count on dict", func() { s.Count() })
	shouldPanic(t, "Slice on dict", func() { s.Slice(0, 0) })
	shouldPanic(t, "Splice on dict", func() { s.Splice(0, nil, nil) })
	shouldPanic(t, "ForEach on dict", func() { s.ForEach(func(int, interface{}) {}) })
	shouldPanic(t, "RangeApply on dict", func() { s.RangeApply(0, 0, func(interface{}) interface{} { return nil }) })
}
