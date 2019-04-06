// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/sync"

	"github.com/dotchain/dot/test/testops"
)

func TestSync(t *testing.T) {
	store := testops.MemStore(nil)
	xformed := ops.Transformed(store, testops.NullCache())
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	opts := []sync.Option{
		sync.WithLog(l),
		sync.WithNotify(func(version int, pending []ops.Op) {}),
		sync.WithSession(-1, nil),
	}

	c1, close1 := sync.Stream(xformed, opts...)
	defer close1()

	c2, close2 := sync.Stream(xformed, opts...)
	defer close2()

	var c1ops changes.Change
	wait := make(chan struct{}, 1000)
	c1.Nextf("key", func() {
		_, c1ops = c1.Next()
		wait <- struct{}{}
	})

	c2.Append(changes.Move{Offset: 2, Count: 3, Distance: 4})

	<-wait
	expected := changes.Move{Offset: 2, Count: 3, Distance: 4}
	if !reflect.DeepEqual(c1ops, expected) {
		t.Fatal("Unexpected merge", c1ops)
	}
}
