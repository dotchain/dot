// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
)

func TestLinkApply(t *testing.T) {
	s1, s2 := rich.NewText("link1"), rich.NewText("link2")
	l1 := data.Link{URL: "url1", Value: s1}
	l2 := data.Link{URL: "url2", Value: s2}

	if l2.Name() != "Embed" {
		t.Error("Unexpected name", l2.Name())
	}

	if x := l1.Apply(nil, nil); !reflect.DeepEqual(x, l1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: l2, After: l2}
	if x := l1.Apply(nil, replace); !reflect.DeepEqual(x, l2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path: []interface{}{"URL"},
		Change: changes.Splice{
			Before: types.S16("u"),
			After:  types.S16("U"),
		},
	}
	if x := l1.Apply(nil, c).(data.Link); x.URL != "Url1" {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"Value"},
		Change: changes.Replace{
			Before: s1,
			After:  s2,
		},
	}
	if x := l1.Apply(nil, c).(data.Link); !reflect.DeepEqual(x.Value, s2) {
		t.Error("Unexpected change", x.Value)
	}
}
