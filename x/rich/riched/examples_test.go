// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched_test

import (
	"fmt"

	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/data"
	"github.com/dotchain/dot/x/rich/html"
	"github.com/dotchain/dot/x/rich/riched"
)

func ExampleStream() {
	s := riched.NewStream(rich.NewText("Hello world", data.FontBold))
	s = s.SetSelection([]interface{}{5}, []interface{}{5})
	s = s.InsertString(" beautiful")
	fmt.Println(html.Format(s.Editor.Text))

	// Output: <b>Hello beautiful world</b>
}
