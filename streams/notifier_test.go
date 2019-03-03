// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"github.com/dotchain/dot/streams"
	"testing"
)

func TestNotifier(t *testing.T) {
	count := 0
	h := &streams.Handler{func() { count++ }}
	var n streams.Notifier

	// add a dummy handler
	n.On(&streams.Handler{func() {}})

	// add a real one and test
	n.On(h)
	n.Notify()
	if count != 1 {
		t.Error("Unexpected", count)
	}

	// add yet another dummy handler and test
	n.On(&streams.Handler{func() {}})
	n.Notify()
	if count != 2 {
		t.Error("Unexpected", count)
	}

	// remove and test
	n.Off(h)
	n.Notify()
	if count != 2 {
		t.Error("Unexpected", count)
	}
}
