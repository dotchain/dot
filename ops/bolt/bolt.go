// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package bolt implements the dot storage for files using boltdb
//
// A http server can be implemented like so:
//      import "github.com/dotchain/dot/ops/bolt"
//      import "github.com/dotchain/dot/ops/nw"
//      store, _ := bolt.New("file.bolt", "instance", nil)
//      defer  store.Close()
//      handler := &nw.Handler{Store: store}
//      h := func(w http.ResponseWriter, req  *http.Request) {
//              // Enable CORS
//              w.Header().Set("Access-Control-Allow-Origin", "*")
//              w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//              if req.Method == "OPTIONS" {
//                    return
//              }
//              handler.ServeHTTP(w, req)
//      }
//      http.HandleFunc("/api/", h)
//      http.ListenAndServe()
//
// Concurrency
//
// A single store instance is safe for concurrent access but the
// provided file is locked until the store is closed.
package bolt

import (
	"bytes"
	"context"
	"strconv"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/nw"
	bolt "github.com/etcd-io/bbolt"
)

// New returns a store with the file store backed by the provided filename
func New(fileName, id string, codec nw.Codec) (ops.Store, error) {
	db, err := bolt.Open(fileName, 0666, nil)
	if err != nil {
		return nil, err
	}
	if codec == nil {
		codec = nw.DefaultCodecs["application/x-gob"]
	}

	return &store{db, []byte(id), codec}, nil
}

type store struct {
	db    *bolt.DB
	id    []byte
	codec nw.Codec
}

func (s *store) Close() {
	must(s.db.Close())
}

func (s *store) Append(ctx context.Context, ops []ops.Op) error {
	if len(ops) == 0 {
		return nil
	}

	ids := make([][]byte, 0, len(ops))
	datas := make([][]byte, 0, len(ops))
	seen := map[interface{}]bool{}
	for _, op := range ops {
		if !seen[op.ID()] {
			seen[op.ID()] = true
			id, data, err := s.encode(op)
			if err != nil {
				return err
			}
			ids = append(ids, id)
			datas = append(datas, data)
		}
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists(s.id)
		if err == nil {
			for kk := range ids {
				v := root.Get(ids[kk])
				if v != nil {
					continue
				}
				seq, err := root.NextSequence()
				must(err)
				must(root.Put(ids[kk], []byte{0}))
				must(root.Put([]byte(strconv.FormatUint(seq-1, 16)), datas[kk]))
			}
		}
		return err
	})
}

func (s *store) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	datas := [][]byte(nil)
	err := s.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket(s.id)
		if root == nil {
			return nil
		}

		count := root.Sequence()
		if count > uint64(version+limit) {
			count = uint64(version + limit)
		}

		for kk := uint64(version); kk < count; kk++ {
			d := root.Get([]byte(strconv.FormatUint(kk, 16)))
			d = append([]byte(nil), d...)
			datas = append(datas, d)
		}
		return nil
	})

	var opx []ops.Op
	if err == nil {
		opx = make([]ops.Op, 0, len(datas))
		for kk, data := range datas {
			var op ops.Op
			op, err = s.decode(data)
			if err != nil {
				opx = nil
				break
			}
			opx = append(opx, op.WithVersion(version+kk))
		}
	}
	return opx, err
}

type opdata struct {
	ops.Op
}

func (s *store) encode(op ops.Op) ([]byte, []byte, error) {
	codec := nw.DefaultCodecs["application/x-gob"]
	if s.codec != nil {
		codec = s.codec
	}
	var data, id bytes.Buffer
	err := codec.Encode(opdata{op}, &data)
	if err == nil {
		err = codec.Encode(op.ID(), &id)
	}
	if err != nil {
		return nil, nil, err
	}

	return id.Bytes(), data.Bytes(), nil
}

func (s *store) decode(data []byte) (ops.Op, error) {
	codec := nw.DefaultCodecs["application/x-gob"]
	if s.codec != nil {
		codec = s.codec
	}
	var opd opdata
	if err := codec.Decode(&opd, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return opd.Op, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
