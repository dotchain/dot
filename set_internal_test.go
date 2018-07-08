// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"reflect"
	"testing"
)

func dictKeys(obj interface{}) []string {
	val := reflect.ValueOf(obj)
	kind := val.Kind()
	for kind == reflect.Ptr {
		val = val.Elem()
		kind = val.Kind()
	}

	if kind == reflect.Map {
		result := []string{}
		for _, kk := range val.MapKeys() {
			result = append(result, kk.String())
		}
		return result
	}

	if kind == reflect.Struct {
		result := []string{}
		for ii := 0; ii < val.NumField(); ii++ {
			result = append(result, val.Type().Field(ii).Name)
		}
		return result
	}

	return nil
}

func toInterface(v reflect.Value) interface{} {
	if v.IsValid() {
		return v.Interface()
	}
	return nil
}

func dictGet(obj interface{}, key string) interface{} {
	val := reflect.ValueOf(obj)
	kind := val.Kind()
	for kind == reflect.Ptr {
		val = val.Elem()
		kind = val.Kind()
	}

	if kind == reflect.Map {
		return toInterface(val.MapIndex(reflect.ValueOf(key)))
	}

	if kind == reflect.Struct {
		return toInterface(val.FieldByName(key))
	}

	return nil
}

func generateSets(input, replace interface{}) []Change {
	ops := []Change{}
	keys := dictKeys(input)
	if reflect.ValueOf(input).Kind() == reflect.Map {
		// for maps we support adding a new non-existent key as an op
		keys = append(keys, "New Key")
	}

	for _, key := range keys {
		// first implement a delete op
		ops = append(ops, Change{
			Set: &SetInfo{
				Key:    key,
				Before: dictGet(input, key),
				After:  nil,
			},
		})

		ops = append(ops, Change{
			Set: &SetInfo{
				Key:    key,
				Before: dictGet(input, key),
				After:  replace,
			},
		})
	}
	return ops
}

func TestMergeSetSetSamePath(t *testing.T) {
	input := map[string]interface{}{"hello": "world", "good": "bye"}
	ops := generateSets(input, "yo")
	testMerge(t, input, ops, ops)
}

func TestMergeSetSetSubPath(t *testing.T) {
	input := map[string]interface{}{"hello": "world", "good": "bye"}
	outer := map[string]interface{}{
		"key1": input,
		"key2": input,
		"key3": input,
		"key4": input,
	}
	outerSets := generateSets(outer, map[string]interface{}{"hello": "bye", "obladi": "oblada"})
	innerSets := []Change{}
	for key := range outer {
		path := []string{key}
		setInfo := &SetInfo{Key: "hello", Before: "world", After: "New World"}
		innerSets = append(innerSets, Change{Path: path, Set: setInfo})
	}

	testMerge(t, outer, outerSets, innerSets)
}

func TestMergeSetSplice(t *testing.T) {
	input := "hello world"
	outer := map[string]interface{}{
		"key1": input,
		"key2": input,
		"key3": input,
		"key4": input,
	}
	outerSets := generateSets(outer, "obladi")
	innerSets := []Change{}
	for key := range outer {
		path := []string{key}
		spliceInfo := &SpliceInfo{Offset: 2, Before: "llo", After: "LLO"}
		innerSets = append(innerSets, Change{Path: path, Splice: spliceInfo})
	}

	testMerge(t, outer, outerSets, innerSets)
}

func TestMergeSetMove(t *testing.T) {
	input := "hello world"
	outer := map[string]interface{}{
		"key1": input,
		"key2": input,
		"key3": input,
		"key4": input,
	}
	outerSets := generateSets(outer, "obladi")
	innerSets := []Change{}
	for key := range outer {
		path := []string{key}
		moveInfo := &MoveInfo{Offset: 2, Count: 3, Distance: 1}
		innerSets = append(innerSets, Change{Path: path, Move: moveInfo})
	}

	testMerge(t, outer, outerSets, innerSets)
}
