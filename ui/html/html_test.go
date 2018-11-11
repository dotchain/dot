// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"fmt"
	"github.com/dotchain/dot/ui/html"
	"testing"
)

func Test(t *testing.T) {
	if n, err := html.Parse("<html><html>"); err == nil {
		t.Error("Unexpected successful parse", n)
	}

	tests := map[string]string{
		`<hello x="a">boo</hello>`:                 `<hello x="b" y="c">booya</hello>`,
		`<hello id="a">boo</hello>`:                `<hello id="b">booya</hello>`,
		`<x><z id="b">boo</z><y id="a">ok</y></x>`: `<x><y id="a">ok</y><z id="b">boo</z></x>`,
	}

	for before, after := range tests {
		t.Run(before+"=>"+after, func(t *testing.T) {
			validate(t, before, after)
			validate(t, after, before)
		})
	}
}

func validate(t *testing.T, before, after string) {
	b, _ := html.Parse(before)
	a, _ := html.Parse(after)
	result := html.Reconciler.Reconcile(b, a)
	if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", a) {
		t.Error("Mismatched", a, result)
	}
}
