// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich/data"
)

func TestFontWeightApply(t *testing.T) {
	if data.FontThin.Name() != "FontWeight" {
		t.Error("Unexpected name", data.FontThin.Name())
	}

	if x := data.FontThin.Apply(nil, nil); x != data.FontThin {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{
		Before: data.FontThin,
		After:  data.FontBold,
	}
	if x := data.FontThin.Apply(nil, replace); x != replace.After {
		t.Error("Unexpected replace", x)
	}

	c := changes.ChangeSet{replace}
	if x := data.FontThin.Apply(nil, c); x != replace.After {
		t.Error("Unexpected changeset", x)
	}
}
