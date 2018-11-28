// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package diff_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/diff"
	"github.com/dotchain/dot/changes/types"
	"math/rand"
	"testing"
)

func TestS8Diff(t *testing.T) {
	d := diff.Std{}
	s1, s2, s3, s4 := randS(), randS(), randS(), randS()

	old := types.S8(s1 + s2 + s3)
	new := types.S8(s1 + s4 + s3)
	if x := d.Diff(d, old, old); x != nil {
		t.Error("Failed identity", x)
	}

	if x := old.Apply(nil, d.Diff(d, old, new)); x != new {
		t.Error("Failed update", old, new, x)
	}

	if x := d.Diff(d, old, types.S8(s1+s3)); len(x.(changes.ChangeSet)) != 1 {
		t.Error("Failed delete", x)
	}

	if x := d.Diff(d, types.S8(s1+s3), old); len(x.(changes.ChangeSet)) != 1 {
		t.Error("Failed insert", x)
	}
}

func TestS16Diff(t *testing.T) {
	d := diff.Std{}
	s1, s2, s3, s4 := randS(), randS(), randS(), randS()

	old := types.S16(s1 + s2 + s3)
	new := types.S16(s1 + s4 + s3)
	if x := d.Diff(d, old, old); x != nil {
		t.Error("Failed identity", x)
	}

	if x := old.Apply(nil, d.Diff(d, old, new)); x != new {
		t.Error("Failed update", old, new, x)
	}

	if x := d.Diff(d, old, types.S16(s1+s3)); len(x.(changes.ChangeSet)) != 1 {
		t.Error("Failed delete", x)
	}

	if x := d.Diff(d, types.S16(s1+s3), old); len(x.(changes.ChangeSet)) != 1 {
		t.Error("Failed insert", x)
	}
}

func randS() string {
	n := rand.Intn(10)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(42)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
