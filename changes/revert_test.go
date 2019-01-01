// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes_test

import (
	"github.com/dotchain/dot/changes"
	"testing"
)

func TestReverts(t *testing.T) {
	initial := S("hello")
	cx := []changes.Change{
		changes.Replace{initial, S("World")},
		changes.Replace{initial, changes.Nil},
		changes.Splice{1, S(""), S("OK!")},
		changes.Splice{1, S("el"), S("")},
		changes.Splice{1, S("el"), S("goo")},

		changes.Move{1, 2, -1},
		changes.Move{1, 2, 1},
		changes.Move{0, 2, 1},
	}
	for _, c := range cx {
		changed := initial.Apply(nil, c).Apply(nil, c.Revert())
		if changed != initial {
			t.Error("Revert failed to properly revert", c)
		}
		if c.Revert().Revert() != c {
			t.Error("Revert failed to properly revert", c, c.Revert(), c.Revert().Revert())
		}

	}
}

func TestRevertsChangeSet(t *testing.T) {
	initial := S("hello")
	cx := changes.ChangeSet{
		changes.Splice{0, S(""), S("OK")},
		changes.Splice{3, S("el"), S("je")},
	}
	changed := initial.Apply(nil, cx).Apply(nil, cx.Revert())
	if changed != initial {
		t.Error("Unexpected revert", changed)
	}
}
