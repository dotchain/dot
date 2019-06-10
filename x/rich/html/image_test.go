// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/html"
)

func TestImageApply(t *testing.T) {
	i1 := html.Image{Src: "uri1", AltText: "link1"}
	i2 := html.Image{Src: "uri2", AltText: "link2"}

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
	if x := i1.Apply(nil, c).(html.Image); x.Src != "Uri1" {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"AltText"},
		Change: changes.Splice{
			Before: types.S16("text1"),
			After:  types.S16("text2"),
		},
	}
	if x := i1.Apply(nil, c).(html.Image); x.AltText != "text2" {
		t.Error("Unexpected change", x.AltText)
	}
}

func TestImageEncodings(t *testing.T) {
	i := html.NewImage("quote\"d", "a < b")

	if x := html.Format(i, nil); x != "<img src=\"quote&#34;d\" alt=\"a &lt; b\"></img>" {
		t.Error("Unexpected", x)
	}
}
