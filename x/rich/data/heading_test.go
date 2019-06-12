// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

func TestHeadingApply(t *testing.T) {
	s1, s2 := rich.NewText("heading1"), rich.NewText("heading2")
	h1 := data.Heading{Level: 1, Text: s1}
	h2 := data.Heading{Level: 2, Text: s2}

	if h1.Name() != "Embed" {
		t.Error("Unexpected name", h1.Name())
	}

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
	if x := h1.Apply(nil, c).(data.Heading); x.Level != 2 {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"Text"},
		Change: changes.Replace{
			Before: s1,
			After:  s2,
		},
	}
	if x := h1.Apply(nil, c).(data.Heading); !reflect.DeepEqual(x.Text, s2) {
		t.Error("Unexpected change", x.Text)
	}
}
