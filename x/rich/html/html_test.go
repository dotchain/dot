// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"testing"

	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func ExampleFormat() {
	s := rich.NewText("hello", html.FontBold)
	fmt.Println("html =", html.Format(s, nil))
	// Output:html = <b>hello</b>
}

func TestFormatBold(t *testing.T) {
	s := rich.NewText("Hello ").
		Concat(rich.NewText("beautiful", html.FontBold)).
		Concat(rich.NewText(" world"))

	result := html.Format(s, nil)
	if result != "Hello <b>beautiful</b> world" {
		t.Error("Unexpected", result, s)
	}
}

func TestFormatBoldAndItalic(t *testing.T) {
	s := rich.NewText("Hello ").
		Concat(rich.NewText("bold", html.FontBold)).
		Concat(rich.NewText("and", html.FontBold, html.FontStyleItalic)).
		Concat(rich.NewText("italic", html.FontStyleItalic)).
		Concat(rich.NewText(" world"))

	result := html.Format(s, nil)
	if result != "Hello <b>bold</b><i><b>and</b></i><i>italic</i> world" {
		t.Error("Unexpected", result, s)
	}
}
