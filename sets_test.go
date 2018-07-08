// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import (
	"github.com/dotchain/dot"
	"testing"
)

func TestSetSet(t *testing.T) {
	input := map[string]interface{}{"hello": "world"}
	t.Run("NewKey NewKey", func(t *testing.T) {
		left := []dot.Change{{Set: &dot.SetInfo{Key: "key1", After: "value1"}}}
		right := []dot.Change{{Set: &dot.SetInfo{Key: "key2", After: "value2"}}}
		output := map[string]interface{}{"hello": "world", "key1": "value1", "key2": "value2"}
		testOps(t, input, output, left, right)
	})
	t.Run("Update NewKey", func(t *testing.T) {
		left := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world", After: "world2"}}}
		right := []dot.Change{{Set: &dot.SetInfo{Key: "key2", After: "value2"}}}
		output := map[string]interface{}{"hello": "world2", "key2": "value2"}
		testOps(t, input, output, left, right)
	})
	t.Run("Update Update", func(t *testing.T) {
		left := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world", After: "world2"}}}
		right := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world", After: "world3"}}}
		// either op can win here, doesn't matter which
		output := map[string]interface{}{"hello": "world2"}
		testOps(t, input, output, left, right)
	})
	t.Run("Update Delete", func(t *testing.T) {
		left := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world", After: "world2"}}}
		right := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world"}}}
		// either op can win here, doesn't matter which
		output := map[string]interface{}{"hello": "world2"}
		testOps(t, input, output, left, right)
	})
	t.Run("Delete Update", func(t *testing.T) {
		left := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world"}}}
		right := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world", After: "world2"}}}
		// either op can win here, doesn't matter which
		output := map[string]interface{}{}
		testOps(t, input, output, left, right)
	})
	t.Run("Delete Delete", func(t *testing.T) {
		left := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world"}}}
		right := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world"}}}
		// either op can win here, doesn't matter which
		output := map[string]interface{}{}
		testOps(t, input, output, left, right)
	})
}

func TestSetSpliceMove(t *testing.T) {
	input := map[string]interface{}{"hello": "world"}
	setOp := []dot.Change{{Set: &dot.SetInfo{Key: "hello", Before: "world", After: "jimmy"}}}
	spliceOp := []dot.Change{{
		Path:   []string{"hello"},
		Splice: &dot.SpliceInfo{Offset: 2, Before: "rld", After: "RLD"},
	}}
	moveOp := []dot.Change{{
		Path: []string{"hello"},
		Move: &dot.MoveInfo{Offset: 2, Count: 2, Distance: 1},
	}}
	output := map[string]interface{}{"hello": "jimmy"}

	// set wins over the rest
	testOps(t, input, output, setOp, spliceOp)
	testOps(t, input, output, spliceOp, setOp)

	// set wins over move
	testOps(t, input, output, setOp, moveOp)
	testOps(t, input, output, moveOp, setOp)
}
