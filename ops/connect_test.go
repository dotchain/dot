// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"math/rand"
	"reflect"
	"sync"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
)

func TestConnect(t *testing.T) {
	store := MemStore(nil)
	xformed := ops.Transformed(store, nullCache{})
	c1 := ops.NewConnector(-1, nil, xformed, rand.Float64)
	c2 := ops.NewConnector(-1, nil, xformed, rand.Float64)
	c1.Connect()
	defer c1.Disconnect()

	c2.Connect()
	defer c2.Disconnect()

	var c1ops changes.Change
	var wg sync.WaitGroup
	s := c1.Stream
	c1.Stream.Nextf("key", func() {
		_, c1ops = s.Next()
		wg.Done()
	})

	wg.Add(1)
	s2 := c2.Stream
	s2.Append(changes.Move{Offset: 2, Count: 3, Distance: 4})

	wg.Wait()
	expected := changes.Move{Offset: 2, Count: 3, Distance: 4}
	if !reflect.DeepEqual(c1ops, expected) {
		t.Fatal("Unexpected merge", c1ops)
	}

	c1.Async.Close()
	c2.Async.Close()
}
