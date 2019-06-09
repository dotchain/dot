// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
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

func TestFontWeight(t *testing.T) {
	weights := []html.FontWeight{
		html.FontThin,
		html.FontExtraLight,
		html.FontLight,
		html.FontNormal,
		html.FontMedium,
		html.FontSemibold,
		html.FontBold,
		html.FontExtraBold,
		html.FontBlack,
	}
	for _, weight := range weights {
		str := fmt.Sprintf("%d", weight)
		t.Run(str, func(t *testing.T) {
			s := rich.NewText("hello", weight)
			expected := fmt.Sprintf(
				"<span style=\"font-weight: %d\">hello</span>",
				weight,
			)
			if weight == html.FontBold {
				expected = "<b>hello</b>"
			}
			if x := html.Format(s, nil); x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}
