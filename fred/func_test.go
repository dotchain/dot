// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

type Func49 struct{}

func (f Func49) Eval(dir *fred.DirStream, args []fred.Object) fred.Object {
	if len(args) == 0 {
		return fred.Error("No args")
	}
	return args[0].Eval(dir)
}

func (f Func49) Diff(old, next *fred.DirStream, c changes.Change, args []fred.Object) changes.Change {
	if len(args) == 0 {
		return nil
	}
	return args[0].Diff(old, next, c)
}

func TestFunc_noArgs(t *testing.T) {
	dir := fred.Dir{
		"heya":   fred.Error("custom error"),
		"suya":   fred.Error("suya is suya"),
		"foo":    fred.Func{Func49{}, nil},
		"fooptr": fred.Ref("foo"),
	}
	s := fred.NewDirStream(dir, nil)
	ptr, s1 := s.Eval(fred.Ref("fooptr"))
	if ptr != fred.Error("No args") {
		t.Error("Unexpected eval", ptr)
	}
	s.Stream.Append(nil)
	if s1, c := s1.Next(); c != nil || s1 == nil || s1.Eval() != ptr {
		t.Error("Unexpected", c, s1 == nil)
	}

}

func TestFunc_refref(t *testing.T) {
	dir := fred.Dir{
		"heya": fred.Error("custom error"),
		"suya": fred.Error("suya is suya"),
		// the second ref is recursive and will crash the test if it is evaluated.
		"foo": fred.Func{
			Functor: Func49{},
			Args:    fred.ToTuple([]fred.Object{fred.Ref("suya"), fred.Ref("foo")}),
		},
		"fooptr": fred.Ref("foo"),
	}
	s := fred.NewDirStream(dir, nil)
	ptr, s1 := s.Eval(fred.Ref("fooptr"))
	if ptr != dir["suya"] {
		t.Error("Unexpected eval", ptr)
	}
	s.Stream.Append(changes.PathChange{
		Path:   []interface{}{"suya"},
		Change: changes.Replace{Before: dir["suya"], After: fred.Ref("heya")},
	})

	s1, c := s1.Next()
	if s1.Eval() != dir["heya"] {
		t.Error("unexpected val", s1.Eval())
	}
	if c != (changes.Replace{Before: ptr, After: dir["heya"]}) {
		t.Error("unexpected change", c)
	}
}
