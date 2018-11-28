// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package diff_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/diff"
	"github.com/dotchain/dot/changes/types"
	"testing"
)

func TestChangingType(t *testing.T) {
	d := diff.Std{}
	old := types.S8("hello")
	new := types.S16("World")
	c := d.Diff(d, old, new)
	expected := changes.Replace{old, new}
	if c != expected {
		t.Error("Unexpected", c)
	}
}
