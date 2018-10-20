// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	store := MemStore(nil)
	xformed := ops.Transformed(store)
	c1 := ops.NewConnector(-1, nil, xformed, rand.Float64)
	c2 := ops.NewConnector(-1, nil, xformed, rand.Float64)
	c1.Connect()
	defer c1.Disconnect()

	c2.Connect()
	defer c2.Disconnect()

	var c1ops changes.Change
	var wg sync.WaitGroup
	c1.Async.Run(func() {
		s := c1.Stream
		c1.Stream.Nextf("key", func() {
			_, c1ops = s.Next()
			wg.Done()
		})
	})

	wg.Add(1)
	s := c2.Stream

	c2.Async.Run(func() {
		s = s.Append(changes.Move{2, 3, 4})
	})

	go func() {
		for {
			c1.Async.Loop(-1)
			c2.Async.Loop(-1)
			time.Sleep(time.Millisecond * 100)
		}
	}()
	wg.Wait()
	expected := changes.ChangeSet{changes.Move{2, 3, 4}}
	if !reflect.DeepEqual(c1ops, expected) {
		t.Fatal("Unexpected merge", c1ops)
	}
}
