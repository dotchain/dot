// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

var nilVal = fred.Nil().Eval(nil)

func TestFixedEval(t *testing.T) {
	t1 := fred.Fixed(nilVal)
	if x := t1.Eval(nil); x != (nilVal) {
		t.Error("Unexpected Eval", x)
	}
}

func TestFixedApply(t *testing.T) {
	t1 := fred.Fixed(nilVal)
	t2 := fred.Fixed(fred.Error("boo"))
	t3 := t1.Apply(nil, changes.Replace{Before: t1, After: t2})
	if t3 != t2 {
		t.Error("replace failed", t3)
	}

	t4 := t1.Apply(nil, changes.PathChange{
		Path:   []interface{}{"Val"},
		Change: changes.Replace{Before: t1.Eval(nil), After: t2.Eval(nil)},
	})
	if !reflect.DeepEqual(t4, t3) {
		t.Error("Unexpected eval", t4)
	}
}

func TestFixedBadKey(t *testing.T) {
	t1 := fred.Fixed(nilVal)

	// must panic
	var r interface{}
	defer func() { r = recover() }()
	t1.Apply(nil, changes.PathChange{Path: []interface{}{"boo"}})

	t.Error("Unexpected success", r)
}
