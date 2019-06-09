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

func TestFontStyleApply(t *testing.T) {
	normal := html.FontStyleNormal
	italic := html.FontStyleItalic

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

func TestFontStyle(t *testing.T) {
	styles := []html.FontStyle{
		html.FontStyleNormal,
		html.FontStyleItalic,
		html.FontStyleOblique,
	}

	for _, style := range styles {
		t.Run(string(style), func(t *testing.T) {
			s := rich.NewText("hello", style)
			expected := fmt.Sprintf(
				"<span style=\"font-style: %s\">hello</span>",
				style,
			)
			if style == html.FontStyleItalic {
				expected = "<i>hello</i>"
			}
			if x := html.Format(s, nil); x != expected {
				t.Error("Unexpected", x)
			}
		})
	}
}
