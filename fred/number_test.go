// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"math/big"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

type NumSum struct{}

func (f *NumSum) Eval(dir *fred.DirStream, args []fred.Object) fred.Object {
	sum := big.NewRat(0, 1)
	for _, arg := range args {
		v := arg.Eval(dir)
		n, ok := v.(fred.Number)
		if !ok {
			return fred.Error("non-numeric arg provided")
		}
		var r big.Rat
		if err := r.UnmarshalText([]byte(string(n))); err != nil {
			return fred.Error(err.Error())
		}
		sum.Add(sum, &r)
	}
	s, err := sum.MarshalText()
	if err != nil {
		return fred.Error(err.Error())
	}
	return fred.Number(string(s))
}

func (f *NumSum) Diff(old, next *fred.DirStream, c changes.Change, args []fred.Object) changes.Change {
	before := f.Eval(old, args)
	after := f.Eval(next, args)
	if before == after {
		return nil
	}
	return changes.Replace{Before: before, After: after}
}

func TestNum_sum(t *testing.T) {
	dir := fred.Dir{
		"one": fred.Number("1"),
		"two": fred.Number("2"),
		"three": fred.Func{
			Functor: &NumSum{},
			Args:    [2]fred.Object{fred.Ref("one"), fred.Ref("two")},
		},
		"six": fred.Func{
			Functor: &NumSum{},
			Args:    [2]fred.Object{fred.Ref("three"), fred.Ref("three")},
		},
		"sixptr": fred.Ref("six"),
	}
	s := fred.NewDirStream(dir, nil)
	ptr, _ := s.Eval(fred.Ref("sixptr"))
	if ptr != fred.Number("6") {
		t.Error("Unexpected eval", ptr)
	}
}
