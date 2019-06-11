// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich/data"
)

func TestFontStyleApply(t *testing.T) {
	normal := data.FontStyleNormal
	italic := data.FontStyleItalic

	if normal.Name() != "FontStyle" {
		t.Error("Unexpected name", normal.Name())
	}

	if x := normal.Apply(nil, nil); x != normal {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: normal, After: italic}
	if x := normal.Apply(nil, replace); x != replace.After {
		t.Error("Unexpected replace", x)
	}

	c := changes.ChangeSet{replace}
	if x := normal.Apply(nil, c); x != replace.After {
		t.Error("Unexpected changeset", x)
	}
}
