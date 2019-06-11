// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"

	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func ExampleFormat() {
	s := rich.NewText("hello", html.FontBold)
	fmt.Println("html =", html.Format(s))
	// Output:html = <b>hello</b>
}
