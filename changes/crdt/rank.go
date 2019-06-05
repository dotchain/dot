// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt

import (
	"math/rand"
	"time"
)

// Rank implements a unique ID that can be used for sorting.
//
// The Epoch field is meant to track the time when the rank was
// created, so this can be used for garbage collection and such
//
// Use Ord if it is important to generate specific sortable IDs with
// the ability to confine them between any two such values.
type Rank struct {
	Epoch int64
	Nonce [4]int64
}

// Less returns true iff r < o
func (r *Rank) Less(o *Rank) bool {
	if r.Epoch != o.Epoch {
		return r.Epoch < o.Epoch
	}
	for kk, v := range r.Nonce {
		if v != o.Nonce[kk] {
			return v < o.Nonce[kk]
		}
	}
	return false
}

// NewRank returns a new rank with current Epoch and some random bits
// to ensure uniqueness
func NewRank() *Rank {
	n := [4]int64{rand.Int63(), rand.Int63(), rand.Int63(), rand.Int63()}
	return &Rank{Epoch: time.Now().UnixNano(), Nonce: n}
}
