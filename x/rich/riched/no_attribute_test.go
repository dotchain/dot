// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/riched"
)

func TestNoAttribute(t *testing.T) {
	n := riched.NoAttribute("something")
	if n.Name() != "something" {
		t.Error("Unexpected name", string(n))
	}

	if x := n.Apply(nil, nil); x != n {
		t.Error("Unexpected nil apply", x)
	}

	replace := changes.Replace{Before: n, After: types.S16("ok")}
	if x := n.Apply(nil, replace); x != replace.After {
		t.Error("Unexpected replace result", x)
	}

	if x := n.Apply(nil, changes.ChangeSet{replace}); x != replace.After {
		t.Error("Unexpected replace result", x)
	}
}
