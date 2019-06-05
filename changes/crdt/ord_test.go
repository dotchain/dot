// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt_test

import (
	"testing"

	"github.com/dotchain/dot/changes/crdt"
)

func TestIntNextPrevLess(t *testing.T) {
	expected := []string{"1,1", "1,2", "1,3", "1,4", "1,5"}
	last := ""
	for _, v := range expected {
		if n := crdt.NextOrd(last); n != v {
			t.Error("Unexpected next", last, n, v)
		}

		if crdt.LessOrd(crdt.NextOrd(last), last) {
			t.Error("Less failed", crdt.NextOrd(last), last)
		}

		if n := crdt.PrevOrd(v); n != last {
			t.Error("Unexpected prev", last, n, v)
		}
		last = v
	}
}

func TestNegIntNextPrevLess(t *testing.T) {
	expected := []string{"1,-1", "1,-2", "1,-3", "1,-4", "1,-5"}
	last := ""
	for _, v := range expected {
		if n := crdt.PrevOrd(last); n != v {
			t.Error("Unexpected prev", last, n, v)
		}

		if crdt.LessOrd(last, crdt.PrevOrd(last)) {
			t.Error("Less failed", crdt.NextOrd(last), last)
		}

		if n := crdt.NextOrd(v); n != last {
			t.Error("Unexpected next", last, n, v)
		}
		last = v
	}
}

func TestBetween(t *testing.T) {
	left, right := "", crdt.NextOrd("")
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			left, right = right, left
		}
		mid := crdt.BetweenOrd(left, right, 1)[0]
		if mid == left || mid == right {
			t.Fatal("Between failed", left, right, mid)
		}
		if crdt.LessOrd(right, left) {
			left, right = right, left
		}
		if crdt.LessOrd(mid, left) || crdt.LessOrd(right, mid) {
			t.Fatal("Between failed", left, right, mid)
		}
		left = mid
	}
}

func TestBetweenN(t *testing.T) {
	left, right := "", crdt.NextOrd("")
	for i := 0; i < 1000; i++ {
		mids := crdt.BetweenOrd(left, right, 3)
		last := left
		for _, mid := range mids {
			if mid == left || mid == right {
				t.Fatal("Between failed", left, right, mid)
			}
			if crdt.LessOrd(mid, last) || crdt.LessOrd(right, mid) {
				t.Fatal("Between failed <", left, right, mid)
			}
			last = mid
		}
	}
}

func TestBetweenNeg(t *testing.T) {
	left, right := crdt.PrevOrd(""), ""
	for i := 0; i < 1000; i++ {
		mid := crdt.BetweenOrd(left, right, 1)[0]
		if mid == left || mid == right {
			t.Fatal("Between failed", left, right, mid)
		}
		if crdt.LessOrd(mid, left) || crdt.LessOrd(right, mid) {
			t.Fatal("Between failed", left, right, mid)
		}
		left = mid
	}
}

func TestBetweenNegPositive(t *testing.T) {
	left, right := "1,-1000", "1,1"
	for i := 0; i < 1000; i++ {
		mid := crdt.BetweenOrd(left, right, 1)[0]
		if mid == left || mid == right {
			t.Fatal("Between failed", left, right, mid)
		}
		if crdt.LessOrd(mid, left) || crdt.LessOrd(right, mid) {
			t.Fatal("Between failed", left, right, mid)
		}
		left = mid
	}
}
