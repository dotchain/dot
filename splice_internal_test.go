// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"github.com/dotchain/dot/encoding"
	"strconv"
	"testing"
)

func generateSplices(input, insert interface{}) []Change {
	ops := []Change{}
	for offset := 0; offset < arraySize(input); offset++ {
		for count := 0; count < arraySize(input)-offset; count++ {
			// first append a deletion op
			ops = append(ops, Change{
				Splice: &SpliceInfo{
					Offset: offset,
					Before: encoding.Get(input).Slice(offset, count),
					After:  encoding.Get(input).Slice(0, 0),
				},
			})
			ops = append(ops, Change{
				Splice: &SpliceInfo{
					Offset: offset,
					Before: encoding.Get(input).Slice(offset, count),
					After:  insert,
				},
			})
		}
	}
	return ops
}

func TestMergeSpliceSpliceSamePath(t *testing.T) {
	input := "Hello World"
	ops := generateSplices(input, "yo-")
	testMerge(t, input, ops, ops)
}

func TestMergeSpliceSpliceSubPath(t *testing.T) {
	input := "Hello World"
	inputOuter := []interface{}{input, input, input, input}

	outerSplices := generateSplices(inputOuter, []interface{}{"yo", "yo"})
	innerSplices := []Change{}
	for ii := 0; ii < len(inputOuter); ii++ {
		innerSplices = append(innerSplices, Change{
			Path: []string{strconv.Itoa(ii)},
			Splice: &SpliceInfo{
				Offset: 2,
				Before: "llo",
				After:  "LLO",
			},
		})
	}

	testMerge(t, inputOuter, outerSplices, innerSplices)
}

func TestMergeSpliceSet(t *testing.T) {
	input := map[string]interface{}{"hello": "world", "good": "bye"}
	outer := []interface{}{input, input, input, input}

	replace := map[string]interface{}{"hello": "bye"}
	outerSplices := generateSplices(outer, []interface{}{replace, replace})

	innerSets := []Change{}
	for ii := range outer {
		path := []string{strconv.Itoa(ii)}
		setInfo := &SetInfo{Key: "hello", Before: "world", After: "New World"}
		innerSets = append(innerSets, Change{Path: path, Set: setInfo})
	}

	testMerge(t, outer, outerSplices, innerSets)
}
