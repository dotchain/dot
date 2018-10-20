// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
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
			_, c1ops = streams.Latest(s)
			c1ops = simplify(c1ops)
			wg.Done()
		})
	})

	wg.Add(2)
	c2.Async.Run(func() {
		s := c2.Stream.Append(changes.Move{2, 3, 4})
		s.Append(changes.Move{10, 11, 2})
	})

	c2.Async.Loop(-1)
	store.Poll(context.Background(), 0)
	required := 2
	for required > 0 {
		required -= c1.Async.Loop(-1)
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
	expected := changes.ChangeSet{changes.Move{2, 3, 4}, changes.Move{10, 11, 2}}
	if !reflect.DeepEqual(c1ops, expected) {
		t.Fatal("Unexpected merge", c1ops)
	}
}

func simplify(c changes.Change) changes.Change {
	switch c := c.(type) {
	case changes.ChangeSet:
		result := []changes.Change{}
		for _, cx := range c {
			cx = simplify(cx)
			if cx != nil {
				result = append(result, cx)
			}
		}

		switch len(result) {
		case 0:
			return nil
		case 1:
			return result[0]
		default:
			return changes.ChangeSet(result)
		}
	}
	return c
}
