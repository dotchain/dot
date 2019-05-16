// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"fmt"
	"testing"

	"github.com/dotchain/dot/fred"
)

func TestVisitor(t *testing.T) {
	v := visitor{}

	t1 := &fred.Vals{
		&fred.ValMap{
			"booya": fred.Text("world"),
		},
		fred.Error("myerr"),
		fred.Nil().Eval(nil),
		fred.Text("hello").Field(nil, fred.Text("concat")),
		fred.Bool(true),
		fred.Bool(false),
	}

	t1.Visit(&v)
	expected := `<list>{
  [0] = <map>{
    [booya] = world
  } // <map>
  [1] = err: myerr
  [2] = <nil>
  [3] = <method>
  [4] = true
  [5] = false
} // <list>`

	if v.text != expected {
		t.Error("unexpected", v.text)
	}
}

type visitor struct {
	text string
	tabs string
}

func (vx *visitor) VisitLeaf(v fred.Val) {
	vx.text += v.Text()
}

func (vx *visitor) VisitChildrenBegin(v fred.Val) {
	vx.tabs += "  "
	vx.text += v.Text() + "{\n"
}

func (vx *visitor) VisitChild(val fred.Val, key interface{}) {
	vx.text += vx.tabs + fmt.Sprintf("[%v] = ", key)
	val.Visit(vx)
	vx.text += "\n"
}

func (vx *visitor) VisitChildrenEnd(v fred.Val) {
	vx.tabs = vx.tabs[2:]
	vx.text += vx.tabs + "} // " + v.Text()
}
