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

func TestDirApply(t *testing.T) {
	s1, s2 := rich.NewText("dir1"), rich.NewText("dir2")
	d1 := &data.Dir{Root: s1, Objects: types.M{"boo": types.S16("hoo1")}}
	d2 := &data.Dir{Root: s2, Objects: types.M{"boo": types.S16("hoo2")}}

	if d2.Name() != "Embed" {
		t.Error("Unexpected name", d2.Name())
	}

	if x := d1.Apply(nil, nil); !reflect.DeepEqual(x, d1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: d1, After: d2}
	if x := d1.Apply(nil, replace); !reflect.DeepEqual(x, d2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path: []interface{}{"Root"},
		Change: changes.Replace{
			Before: s1,
			After:  s2,
		},
	}
	if x := d1.Apply(nil, c).(*data.Dir); x.Root != s2 {
		t.Error("Unexpected change", x)
	}

	c = changes.PathChange{
		Path: []interface{}{"Objects"},
		Change: changes.Replace{
			Before: d1.Objects,
			After:  d2.Objects,
		},
	}
	if x := d1.Apply(nil, c).(*data.Dir); !reflect.DeepEqual(x.Objects, d2.Objects) {
		t.Error("Unexpected change", x.Objects)
	}
}
