// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
	"github.com/dotchain/dot/x/rich/riched"
)

func TestEditorApply(t *testing.T) {
	s := riched.NewEditor(rich.NewText("Hello world", data.FontBold))

	if x := s.Apply(nil, nil); x != s {
		t.Error("Unexpected nil apply", x)
	}

	replace := changes.Replace{Before: s, After: types.S16("boo")}
	if x := s.Apply(nil, replace); x != replace.After {
		t.Error("Unexpected replace", x)
	}

	if x := s.Apply(nil, changes.PathChange{Path: nil, Change: replace}); x != replace.After {
		t.Error("Unexpected replace", x)
	}
}
