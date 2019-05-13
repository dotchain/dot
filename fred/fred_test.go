// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred_test

import (
	"fmt"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/fred"
)

func ExampleDir_basic() {
	dir := fred.Dir{
		"heya": fred.Error("custom error"),
		"suya": fred.Error("suya is suya"),
		"ptr":  fred.Ref("heya"),
	}
	s := fred.NewDirStream(dir, nil)
	ptr, s1 := s.Eval(fred.Ref("heya"))
	ptr2, s2 := s.Eval(fred.Ref("ptr"))
	fmt.Println("Got", ptr, ptr2)

	// Now replace heya with heya2 and ptr from &heya to &suya
	s.Stream.Append(changes.ChangeSet{
		changes.PathChange{
			Path: []interface{}{"heya"},
			Change: changes.Replace{
				Before: fred.Error("custom error"),
				After:  fred.Error("new error"),
			},
		},
		changes.PathChange{
			Path: []interface{}{"ptr"},
			Change: changes.Replace{
				Before: fred.Ref("heya"),
				After:  fred.Ref("suya"),
			},
		},
	})

	s1x, c1 := s1.Next()
	s2x, c2 := s2.Next()

	ptr = s1x.Eval()
	ptr2 = s2x.Eval()

	fmt.Println("Got", ptr, ptr2)
	fmt.Println("Changes", c1, c2)

	// Output:
	// Got custom error custom error
	// Got new error suya is suya
	// Changes {custom error new error} {custom error suya is suya}
}

func ExampleDir_addNewRef() {
	dir := fred.Dir{
		"heya": fred.Error("custom error"),
		"suya": fred.Error("suya is suya"),
		"ptr":  fred.Ref("heya"),
	}
	s := fred.NewDirStream(dir, nil)
	ptr, s1 := s.Eval(fred.Ref("heya"))
	ptr2, s2 := s.Eval(fred.Ref("ptr"))
	fmt.Println("Got", ptr, ptr2)

	// Now replace heya with heya2 and ptr from &heya to &suya
	s.Stream.Append(changes.ChangeSet{
		changes.PathChange{
			Path: []interface{}{"oomi"},
			Change: changes.Replace{
				Before: changes.Nil,
				After:  fred.Error("gumi"),
			},
		},
		changes.PathChange{
			Path: []interface{}{"ptr"},
			Change: changes.Replace{
				Before: fred.Ref("heya"),
				After:  fred.Ref("oomi"),
			},
		},
	})

	s1x, c1 := s1.Next()
	s2x, c2 := s2.Next()

	ptr = s1x.Eval()
	ptr2 = s2x.Eval()

	fmt.Println("Got", ptr, ptr2)
	fmt.Println("Changes", c1, c2)

	// Output:
	// Got custom error custom error
	// Got custom error gumi
	// Changes <nil> {custom error gumi}
}

func ExampleDir_ensureCacheLeak() {
	dir := fred.Dir{
		"heya": fred.Error("custom error"),
		"suya": fred.Error("suya is suya"),
		"ptr":  fred.Ref("heya"),
	}
	s := fred.NewDirStream(dir, nil)
	_, s1 := s.Eval(fred.Ref("ptr"))
	s.Eval(fred.Ref("suya"))

	s.Stream.Append(nil)
	_, _ = s1.Next()
	snext, _ := s.Next()

	fmt.Println(s.Cache[fred.Ref("suya")], snext.Cache[fred.Ref("suya")])

	// Output: suya is suya <nil>
}
