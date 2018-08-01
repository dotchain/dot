// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package conv_test

import (
	"github.com/dotchain/dot/conv"
	"testing"
)

func TestFromIndex(t *testing.T) {
	tests := []int{0, 5, 10, 15, 100, 1001}
	expected := []string{"0", "5", "10", "15", "100", "1001"}

	for kk, x := range tests {
		if y := conv.FromIndex(x); y != expected[kk] {
			t.Error("FromIndex(", x, ") =", y)
		}
	}
}

func TestToIndex(t *testing.T) {
	tests := []string{"0", "5", "10", "15", "100", "1001"}
	expected := []int{0, 5, 10, 15, 100, 1001}

	for kk, x := range tests {
		if y := conv.ToIndex(x); y != expected[kk] {
			t.Error("FromIndex(", x, ") =", y)
		}
	}
}

func TestIsIndex(t *testing.T) {
	yes := []string{"0", "5", "10", "15", "100", "1001"}
	no := []string{"", "-1", "1.5"}

	for _, x := range yes {
		if y := conv.IsIndex(x); !y {
			t.Error("IsIndex(", x, ") =", y)
		}
	}
	for _, x := range no {
		if y := conv.IsIndex(x); y {
			t.Error("IsIndex(", x, ") =", y)
		}
	}
}
