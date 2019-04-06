// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/sync"
	"github.com/dotchain/dot/streams"

	"github.com/dotchain/dot/test/testops"
)

func TestSyncFromScratch(t *testing.T) {
	store := testops.MemStore(nil)
	c1, close1 := stream(store, -1, nil)
	c2, close2 := stream(store, -1, nil)
	defer store.Close()
	defer close1()
	defer close2()

	c2.Append(changes.Move{Offset: 2, Count: 3, Distance: 4})
	_, c1ops := next(c1)

	expected := changes.Move{Offset: 2, Count: 3, Distance: 4}
	if !reflect.DeepEqual(c1ops, expected) {
		t.Fatal("Unexpected merge", c1ops)
	}
}

func TestSyncReconnect(t *testing.T) {
	store := testops.MemStore([]ops.Op{
		ops.Operation{
			OpID:    "one",
			VerID:   0,
			BasisID: -1,
			Change:  changes.Splice{Offset: 5, Before: types.S8(" "), After: types.S8("--")},
		},
	})

	pending := ops.Operation{OpID: "two", VerID: 0, BasisID: 0}
	pending.Change = changes.Splice{Offset: 15, Before: types.S8(" "), After: types.S8("")}

	c1, close1 := stream(store, 0, []ops.Op{pending})
	c2, close2 := stream(store, -1, nil)
	defer store.Close()
	defer close1()
	defer close2()

	last := changes.Splice{Offset: 10, Before: types.S8(""), After: types.S8("OK")}
	c2 = c2.Append(last)

	// expect c1 to receive "last" but with offset shifted to factor in op#one
	_, x := next(c1)
	expected := last
	expected.Offset += len("--") - len(" ")

	if !reflect.DeepEqual(x, expected) {
		t.Fatal("Unexpected merge", x, expected)
	}

	// now fetch from c2 and expect the #one unchanged but #two modified
	c2, x = next(c2)
	expected = changes.Splice{Offset: 5, Before: types.S8(" "), After: types.S8("--")}
	if !reflect.DeepEqual(x, expected) {
		t.Fatal("Unexpected merge", x, expected)
	}

	c2, x = next(c2)
	expected = changes.Splice{Offset: 15, Before: types.S8(" "), After: types.S8("")}
	expected.Offset += len("OK") - len("")
	if !reflect.DeepEqual(x, expected) {
		t.Fatal("Unexpected merge", x, expected, c2)
	}
}

func TestSyncMultipleInFlight(t *testing.T) {
	store := testops.MemStore([]ops.Op{
		ops.Operation{
			OpID:    "one",
			VerID:   0,
			BasisID: -1,
			Change:  changes.Move{Offset: 100, Count: 101, Distance: 102},
		},
	})

	flushed := make(chan struct{}, 100)
	waitForFlush := func(version int, pending []ops.Op) {
		if len(pending) == 0 {
			flushed <- struct{}{}
		}
	}
	auto := sync.WithAutoTransform(testops.NullCache())
	s, close := sync.Stream(store, sync.WithNotify(waitForFlush), auto)
	defer close()

	// append a couple of moves
	s = s.Append(changes.Move{Offset: 1, Count: 2, Distance: 3})
	s = s.Append(changes.Move{Offset: 10, Count: 11, Distance: 12})

	// receive the original Move
	s, c := next(s)
	if c != (changes.Move{Offset: 100, Count: 101, Distance: 102}) {
		t.Fatal("Unexpected change", c)
	}

	// now wait till pending is flushed and then append one more
	<-flushed
	s.Append(changes.Move{Offset: 1000, Count: 1000, Distance: 1000})

	// wait for this to be flushed as well
	<-flushed

	// now confirm with the store that these operations have the
	// right parent and bassisIDs

	saved, err := store.GetSince(context.Background(), 0, 1000)
	if err != nil || len(saved) != 4 {
		t.Fatal("Wrong number of ops in the store", err, len(saved))
	}

	if saved[1].Parent() != nil || saved[2].Basis() != -1 {
		t.Fatal("Unexpected first op", saved[1].Parent(), saved[1].Basis())
	}

	if saved[2].Parent() != saved[1].ID() || saved[2].Basis() != -1 {
		t.Fatal("Unexpected first op", saved[2].Parent(), saved[2].Basis())
	}

	if saved[3].Parent() != nil || saved[3].Basis() != 2 {
		t.Fatal("Unexpected first op", saved[3].Parent(), saved[3].Basis())
	}
}

func TestSyncMismatchedVersions(t *testing.T) {
	store := &fakeStore{[]ops.Op{
		ops.Operation{
			OpID:    "one",
			VerID:   1,
			BasisID: -1,
			Change:  changes.Move{Offset: 100, Count: 101, Distance: 102},
		},
	}}

	errors := make(chan error, 100)
	_, close := sync.Stream(store, sync.WithLog(expectFatalLog(errors)))
	defer close()
	err := <-errors
	if !strings.Contains(err.Error(), "mismatch") {
		t.Fatal("Did not get a version mismatch error", err)
	}
}

func stream(s ops.Store, version int, pending []ops.Op) (streams.Stream, func()) {
	xformed := ops.Transformed(s, testops.NullCache())
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	opts := []sync.Option{
		sync.WithLog(l),
		sync.WithNotify(func(version int, pending []ops.Op) {}),
		sync.WithSession(-1, nil),
	}
	if version != -1 {
		opts = append(opts, sync.WithSession(version, pending))
	}

	return sync.Stream(xformed, opts...)
}

// next blocks until there is a next and returns that value
func next(s streams.Stream) (streams.Stream, changes.Change) {
	if next, c := s.Next(); next != nil {
		return next, c
	}

	wait := make(chan struct{}, 1000)
	s.Nextf("wait", func() { wait <- struct{}{} })
	defer s.Nextf("wait", nil)
	<-wait
	return s.Next()
}

type expectFatalLog chan error

func (e expectFatalLog) Println(args ...interface{})            {}
func (e expectFatalLog) Printf(fmt string, args ...interface{}) {}
func (e expectFatalLog) Fatal(args ...interface{}) {
	e <- args[1].(error)
}

type fakeStore struct {
	entries []ops.Op
}

func (f fakeStore) Append(ctx context.Context, args []ops.Op) error {
	return nil
}

func (f fakeStore) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	if version == 0 {
		return f.entries, nil
	}
	return nil, nil
}

func (f fakeStore) Poll(ctx context.Context, version int) error {
	return nil
}

func (f fakeStore) Close() {
}
