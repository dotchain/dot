// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/data"
)

func TestImageApply(t *testing.T) {
	i1 := data.Image{Src: "uri1", AltText: "link1"}
	i2 := data.Image{Src: "uri2", AltText: "link2"}

	if i1.Name() != "Embed" {
		t.Error("Unexpected name", i1.Name())
	}

	if x := i1.Apply(nil, nil); !reflect.DeepEqual(x, i1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: i2, After: i2}
	if x := i1.Apply(nil, replace); !reflect.DeepEqual(x, i2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path: []interface{}{"Src"},
		Change: changes.Splice{
			Before: types.S16("u"),
			After:  types.S16("U"),
		},
	}
	if x := i1.Apply(nil, c).(data.Image); x.Src != "Uri1" {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"AltText"},
		Change: changes.Splice{
			Before: types.S16("text1"),
			After:  types.S16("text2"),
		},
	}
	if x := i1.Apply(nil, c).(data.Image); x.AltText != "text2" {
		t.Error("Unexpected change", x.AltText)
	}
}
