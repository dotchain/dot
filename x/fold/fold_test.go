// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fold_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/fold"
	"testing"
)

func TestFold(t *testing.T) {
	// much of the folding test is covered by the examples.
	// The only thing left is to validate the Scheduler() and
	// WithScheduler() methods.

	async := &streams.AsyncScheduler{}
	upstream := streams.New()

	folded := fold.New(changes.Move{0, 5, 10}, upstream).WithScheduler(async)
	if x := folded.Scheduler(); x != async {
		t.Error("Unexpected scheduler", x)
	}
}
