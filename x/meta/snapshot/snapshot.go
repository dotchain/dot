// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package snapshot manages session storage
//
// Snapshots are session meta data (version, pending), the actual app
// value and the transformed/merged ops cache.
//
// The Bolt{} type writes the snapshot to a bolt-db (taking care to to
// incremental writes).
//
// It can also be used to initiate the session (in which case the
// snapshot is automatically updated as the app changes the stream
// syncs up).  The example illustrates how this would look like.
package snapshot

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"sync"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/meta"

	bolt "github.com/etcd-io/bbolt"
)

// Bolt represents a local file (bolt db) storage
//
// Path is the path to the file to write snapshots
//
// Initial is the initial app state. This can be nil if the app starts
// from scratch (in which case it is mapped to changes.Nil)
//
// Codec is the code to use. If not specified, the default codec is
// used (which is "application/gob")
type Bolt struct {
	Path    string
	Initial changes.Value
	nw.Codec
	close func()
}

// Load loads the snapshot value
func (f *Bolt) Load() (meta.Data, changes.Value, error) {
	db, data, v, err := f.load()
	if db != nil {
		f.must(db.Close())
	}
	return data, v, err
}

// Save updates the snapshot with changes
func (f *Bolt) Save(m meta.Data, latest changes.Value) error {
	db, data, v, err := f.load()
	if db != nil {
		defer func() {
			f.must(db.Close())
		}()
	}

	if err == nil {
		err = f.save(db, data, m, v, latest)
	}
	return err
}

// NewSession creates a new session from the snapshot
func (f *Bolt) NewSession(url string) (updates, metas streams.Stream) {
	db, data, v, err := f.load()
	f.must(err)

	var closef func()
	closef, updates, metas = reconnect(url, data)

	var mu sync.Mutex
	save := func() {
		mu.Lock()
		defer mu.Unlock()

		var c changes.Change
		updates, c = streams.Latest(updates)
		nextv := v.Apply(nil, c)
		metas, c = streams.Latest(metas)
		nextData := data.Apply(nil, c).(meta.Data)
		f.must(f.save(db, data, nextData, v, nextv))
		data, v = nextData, nextv
	}

	metas.Nextf(f, save)
	updates.Nextf(f, save)

	f.close = func() {
		closef()
		metas.Nextf(f, nil)
		updates.Nextf(f, nil)
		save()
		f.must(db.Close())
	}

	return updates, metas
}

// Close closes the session and the db
func (f *Bolt) Close() {
	f.close()
}

func (f *Bolt) load() (*bolt.DB, meta.Data, changes.Value, error) {
	db, err := bolt.Open(f.Path, 0666, nil)
	if err != nil {
		return nil, meta.Data{}, nil, err
	}

	m := meta.Data{Version: -1}
	v := f.Initial
	if v == nil {
		v = changes.Nil
	}

	err = db.View(func(tx *bolt.Tx) (e error) {
		defer catch(&e)()

		root := tx.Bucket([]byte("root"))
		if root == nil {
			return nil
		}

		v = f.decode(root.Get([]byte("Value"))).(changes.Value)
		m.Version = f.decode(root.Get([]byte("Version"))).(int)
		m.Pending = f.decode(root.Get([]byte("Pending"))).([]ops.Op)

		if m.Version >= 0 {
			m.TransformedOp = meta.CachedOp{}
			m.MergeOps = meta.CachedOps{}
		}

		for kk := 0; kk <= m.Version; kk++ {
			key := strconv.FormatUint(uint64(kk), 16)
			if val := root.Get([]byte("X-" + key)); val != nil {
				m.TransformedOp[kk] = f.decode(val).(ops.Op)
				val = root.Get([]byte("Merge-" + key))
				m.MergeOps[kk] = f.decode(val).([]ops.Op)
			}
		}
		return nil
	})

	return db, m, v, err
}

func (f *Bolt) save(db *bolt.DB, before, now meta.Data, beforev, nowv changes.Value) error {
	return db.Update(func(tx *bolt.Tx) (e error) {
		defer catch(&e)()
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		f.must(err)

		changed := false
		if before.Version < 0 || before.Version != now.Version {
			changed = true
			f.must(root.Put([]byte("Version"), f.encode(now.Version)))
		}

		if changed || len(before.Pending) != len(now.Pending) {
			changed = true
			f.must(root.Put([]byte("Pending"), f.encode(now.Pending)))
		}

		if changed {
			f.must(root.Put([]byte("Value"), f.encode(nowv)))
		}

		for k, v := range now.TransformedOp {
			if _, ok := before.TransformedOp[k]; ok {
				continue
			}
			key := strconv.FormatUint(uint64(k), 16)
			f.must(root.Put([]byte("X-"+key), f.encode(v)))
			f.must(root.Put([]byte("Merge-"+key), f.encode(now.MergeOps[k])))
		}
		return nil
	})
}

func (f *Bolt) encode(v interface{}) []byte {
	codec := f.Codec
	if codec == nil {
		codec = nw.DefaultCodecs["application/x-gob"]
	}
	var data bytes.Buffer
	f.must(codec.Encode(changes.Atomic{Value: v}, &data))
	return data.Bytes()
}

func (f *Bolt) decode(data []byte) interface{} {
	codec := f.Codec
	if codec == nil {
		codec = nw.DefaultCodecs["application/x-gob"]
	}
	var temp changes.Atomic
	f.must(codec.Decode(&temp, bytes.NewReader(data)))
	return temp.Value
}

func (f *Bolt) must(err error) {
	if err != nil {
		panic(err)
	}
}

func catch(e *error) func() {
	return func() {
		if r := recover(); r != nil {
			*e = r.(error)
		}
	}
}

func init() {
	gob.Register([]ops.Op(nil))
}
