// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/run"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"

	"github.com/dotchain/dot/test/testops"
)

func newOp(opID, parentID interface{}, ver, basis int, c changes.Change) ops.Operation {
	return ops.Operation{OpID: opID, ParentID: parentID, VerID: ver, BasisID: basis, Change: c}
}

func TestClientErrors(t *testing.T) {
	c := &nw.Client{URL: "http://localhost:8183/nw?q=1"}

	// unknown type
	op1 := newOp("ID1", "", 100, -1, run.Run{})
	err := c.Append(getContext(), []ops.Op{op1})
	if err == nil {
		t.Fatal("Did not fail with unregistered change type")
	}

	c = &nw.Client{URL: "&amp;://ok"}
	err = c.Append(getContext(), []ops.Op{})
	if err == nil {
		t.Fatal("Did not fail with invalid url")
	}

	// simulate server error
	errout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boo", http.StatusMultipleChoices)
	})
	srv := httptest.NewServer(errout)
	defer srv.Close()

	c = &nw.Client{URL: srv.URL, Client: srv.Client()}

	err = c.Append(getContext(), []ops.Op{})
	if err == nil {
		t.Fatal("Unexpected success with client error")
	}

	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("gooopey")); err != nil {
			panic(err)
		}
	}))
	defer srv2.Close()

	c = &nw.Client{URL: srv2.URL, Client: srv2.Client()}
	err = c.Append(getContext(), []ops.Op{})
	if err == nil {
		t.Fatal("Unexpected success with client error")
	}
}

func TestHeaderPasssing(t *testing.T) {
	foundHeader := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundHeader = r.Header.Get("Zug") == "Zug"
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	h := http.Header{}
	h.Set("Zug", "Zug")
	c := &nw.Client{URL: srv.URL, Client: srv.Client(), Header: h}

	// ignore error because it is a fake server
	ignore(c.Append(getContext(), []ops.Op{}))
	if !foundHeader {
		t.Fatal("Missed Zug Header")
	}
}

func TestPollerStore(t *testing.T) {
	store := nw.MemPoller(testops.MemStore(nil))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := store.Poll(ctx, 0)
	if err == nil || err != ctx.Err() {
		t.Error("unexpected poll result", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// now kick off a go routine to wakeup
	go func() {
		time.Sleep(10 * time.Millisecond)
		op := newOp("ID2", "", 100, -1, run.Run{})
		_ = store.Append(context.Background(), []ops.Op{op})
	}()

	err = store.Poll(ctx, 0)
	if err != nil {
		t.Error("unexpected poll result", err)
	}
}

func TestClosedPollerStore(t *testing.T) {
	store := nw.MemPoller(testops.MemStore(nil))
	store.Close()

	err := store.Poll(context.Background(), 0)
	if err != nil {
		t.Error("unexpected poll result", err)
	}

	store = nw.MemPoller(testops.MemStore(nil))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(10 * time.Millisecond)
		store.Close()
	}()

	err = store.Poll(ctx, 0)
	if err != nil {
		t.Error("unexpected poll result", err)
	}
}

func TestStoreGetSince(t *testing.T) {
	operations := []ops.Op{
		newOp("ID1", "", 100, -1, run.Run{}),
		newOp("ID2", "", 100, -1, run.Run{}),
		newOp("ID3", "", 100, -1, run.Run{}),
	}
	store := testops.MemStore(operations)
	defer store.Close()

	result, err := store.GetSince(context.Background(), 0, 1)
	if err != nil || len(result) != 1 {
		t.Error("Unexpected results", err, len(result))
	}
}

func TestServerErrors(t *testing.T) {
	store := nw.MemPoller(fakeStore{})
	defer store.Close()
	srv := httptest.NewServer(&nw.Handler{Store: store})
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil || resp.StatusCode != 400 {
		t.Fatal("Unexpected fetch without content-type", err, resp)
	}

	data := bytes.NewReader([]byte("Some garbage"))
	resp, err = srv.Client().Post(srv.URL, "application/x-gob", data)
	if err != nil || resp.StatusCode != 400 {
		t.Fatal("Unexpected fetch without content-type", err, resp)
	}

	c := &nw.Client{URL: srv.URL, Client: srv.Client()}
	err = c.Append(getContext(), nil)
	if err == nil || err.Error() != "Append error" {
		t.Fatal("Unexpected append behavior", err)
	}

	opx, err := c.GetSince(getContext(), 0, 0)
	if err == nil {
		t.Fatal("Unexpected success for GetSince", err, opx)
	}
}

type fakeStore struct{}

func (f fakeStore) Append(_ context.Context, opx []ops.Op) error {
	return errors.New("Append error")
}

func (f fakeStore) GetSince(_ context.Context, version, limit int) ([]ops.Op, error) {
	unencodeable := []ops.Op{unencodeableOp{}}
	return unencodeable, errors.New("GetSince error")
}

func (f fakeStore) Poll(_ context.Context, version int) error {
	return errors.New("Poll error")
}

func (f fakeStore) Close() {
}

type unencodeableOp struct{}

func (u unencodeableOp) ID() interface{}                   { return nil }
func (u unencodeableOp) Version() int                      { return 0 }
func (u unencodeableOp) WithVersion(int) ops.Op            { return nil }
func (u unencodeableOp) Parent() interface{}               { return nil }
func (u unencodeableOp) Basis() int                        { return 0 }
func (u unencodeableOp) Changes() changes.Change           { return nil }
func (u unencodeableOp) WithChanges(changes.Change) ops.Op { return nil }

func getContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// in this test, we let cancel and the associated channel leak.
	_ = cancel
	return ctx
}

func ignore(err error) {}
