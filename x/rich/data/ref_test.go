// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich/data"
)

func TestRefApply(t *testing.T) {
	r1 := &data.Ref{ID: "one"}
	r2 := &data.Ref{ID: "two"}

	if r2.Name() != "Embed" {
		t.Error("Unexpected name", r2.Name())
	}

	if x := r1.Apply(nil, nil); !reflect.DeepEqual(x, r1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: r1, After: r2}
	if x := r1.Apply(nil, changes.ChangeSet{replace}); x != r2 {
		t.Error("Unexpected replace", x)
	}
}
