// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func TestHeadingApply(t *testing.T) {
	s1, s2 := rich.NewText("heading1"), rich.NewText("heading2")
	h1 := html.Heading{Level: 1, Text: &s1}
	h2 := html.Heading{Level: 2, Text: &s2}

	if x := h1.Apply(nil, nil); !reflect.DeepEqual(x, h1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: h2, After: h2}
	if x := h1.Apply(nil, replace); !reflect.DeepEqual(x, h2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path: []interface{}{"Level"},
		Change: changes.Replace{
			Before: changes.Atomic{Value: 1},
			After:  changes.Atomic{Value: 2},
		},
	}
	if x := h1.Apply(nil, c).(html.Heading); x.Level != 2 {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"Text"},
		Change: changes.Replace{
			Before: s1,
			After:  s2,
		},
	}
	if x := h1.Apply(nil, c).(html.Heading); !reflect.DeepEqual(*x.Text, s2) {
		t.Error("Unexpected change", x.Text)
	}
}

func TestHeading(t *testing.T) {
	levels := []string{"h1", "h1", "h2", "h3", "h4", "h5", "h6", "h1"}
	for l, str := range levels {
		test := fmt.Sprintf("%s-%d", str, l)
		t.Run(test, func(t *testing.T) {
			h := html.NewHeading(l, rich.NewText("x", html.FontBold))
			expected := fmt.Sprintf("<%s><b>x</b></%s>", str, str)
			if x := html.Format(h, nil); x != expected {
				t.Error("Unexpected", x, expected)
			}
		})
	}
}
