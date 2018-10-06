// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fold_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/fold"
	"reflect"
	"testing"
)

func TestFold(t *testing.T) {
	// much of the folding test is covered by the examples.
	// The only thing left is to validate the Scheduler() and
	// WithScheduler() methods.

	async := &streams.AsyncScheduler{}
	upstream := streams.New()

	// move [0 - 5] to the right by 10
	folded := fold.New(changes.Move{0, 5, 10}, upstream).WithScheduler(async)

	var next changes.Change
	upstream.Nextf("key", func(c changes.Change, _ streams.Stream) {
		next = c
	})

	// move [1 - 2] to the right by 20 and see it on upstream
	folded = folded.Append(changes.Move{1, 1, 20})
	expected := changes.ChangeSet{changes.Change(nil), changes.Move{6, 1, 15}}
	if next != nil {
		t.Error("Unexpected sync change", next)
	}
	async.Run(1)

	if !reflect.DeepEqual(next, expected) {
		t.Error("Unexpected Next() behavior", next)
	}

	// validate scheduler
	if x := folded.Scheduler(); x != async {
		t.Error("folded.Scheduler()", x)
	}
}
