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
package snapshot

import (
	"bytes"

	"github.com/dotchain/dot"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops/nw"

	bolt "github.com/etcd-io/bbolt"
)

// Bolt represents a local file (bolt db) storage
//
// Path is the path to the file to write snapshots
//
// Initial is the initial app state. This can be nil if the app starts
// from scratch (in which case it is mapped to changes.Nil)
//
// Codec is the codec to use. If not specified, the default codec is
// used
type Bolt struct {
	Path    string
	Initial changes.Value
	nw.Codec
}

// Load loads the snapshot value
func (f *Bolt) Load() (*dot.Session, changes.Value, error) {
	db, data, v, err := f.load()
	if db != nil {
		f.must(db.Close())
	}
	return data, v, err
}

// Save updates the snapshot with changes
func (f *Bolt) Save(s *dot.Session, latest changes.Value) error {
	db, data, v, err := f.load()
	if db != nil {
		defer func() {
			f.must(db.Close())
		}()
	}

	if err == nil {
		err = f.save(db, data, s, v, latest)
	}
	return err
}

func (f *Bolt) load() (*bolt.DB, *dot.Session, changes.Value, error) {
	db, err := bolt.Open(f.Path, 0666, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	s := dot.NewSession()
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
		s = f.decode(root.Get([]byte("Session"))).(*dot.Session)
		return nil
	})

	return db, s, v, err
}

func (f *Bolt) save(db *bolt.DB, before, now *dot.Session, beforev, nowv changes.Value) error {
	return db.Update(func(tx *bolt.Tx) (e error) {
		defer catch(&e)()
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		f.must(err)

		changed := before.Version < 0 || before.Version != now.Version ||
			len(before.Pending) != len(now.Pending) ||
			len(before.Merge) != len(now.Merge)

		if changed {
			f.must(root.Put([]byte("Value"), f.encode(nowv)))
			// TODO: the cached ops and such can be saved
			// incrementally as all old versions are guaranteed
			// to not change
			f.must(root.Put([]byte("Session"), f.encode(now)))
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
	nw.Register(&dot.Session{})
}
