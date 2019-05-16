// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func TestBool(t *testing.T) {
	t1 := fred.Bool(false)
	t2 := fred.Bool(false)
	if t1 != t2 {
		t.Error("Unexpected inequality")
	}

	t3 := t1.Apply(nil, changes.ChangeSet{
		changes.Replace{
			Before: t1,
			After:  changes.Nil,
		},
	})
	if t3 != changes.Nil {
		t.Error("Unexpected apply", t3)
	}

	if t1.Apply(nil, nil) != t1 {
		t.Error("Unexpected apply", t3)
	}
}
