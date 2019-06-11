// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich/html"
)

func TestFontWeightApply(t *testing.T) {
	if x := html.FontThin.Apply(nil, nil); x != html.FontThin {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{
		Before: html.FontThin,
		After:  html.FontBold,
	}
	if x := html.FontThin.Apply(nil, replace); x != replace.After {
		t.Error("Unexpected replace", x)
	}

	c := changes.ChangeSet{replace}
	if x := html.FontThin.Apply(nil, c); x != replace.After {
		t.Error("Unexpected changeset", x)
	}
}
