// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package heap_test

import (
	"testing"

	"github.com/dotchain/dot/x/heap"
)

func TestUpdate(t *testing.T) {
	h := heap.Heap{}
	h1 := h.Update("first", 5)
	h2 := h1.Update("second", 4)
	h3 := h2.Update("third", 4)

	if x := h3.UpdateChange("first", 5); x != nil {
		t.Error("Unexpected shift", x)
	}

	if x := h3.DeleteChange("nonexistent"); x != nil {
		t.Error("Unexpected shift", x)
	}

	keys := []string{}
	h3.Iterate(func(key interface{}, rank int) bool {
		keys = append(keys, key.(string))
		return false
	})
	if keys[0] != "first" || len(keys) != 1 {
		t.Error("Unexpected keys", keys)
	}

	h4 := h3.Delete("second")
	keys = []string{}
	h4.Iterate(func(key interface{}, rank int) bool {
		keys = append(keys, key.(string))
		return true
	})
	if keys[0] != "first" || keys[1] != "third" {
		t.Error("Unexpected keys", keys)
	}
}
