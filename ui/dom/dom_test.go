// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom_test

import (
	"github.com/dotchain/dot/ui/dom"
	"github.com/dotchain/dot/ui/html"
	"testing"
)

func Test(t *testing.T) {
	tests := map[string]string{
		``:                          `<x>booya</x>`,
		`<y></y>`:                   `<x>boo</x>`,
		`<x/>`:                      `<x><y>ok</y></x>`,
		`<hello>boo</hello>`:        `<hello>booya</hello>`,
		`<hello>booya</hello>`:      `<hello x="a">booya</hello>`,
		`<hello x="a">boo</hello>`:  `<hello x="b" y="c">booya</hello>`,
		`<hello id="a">boo</hello>`: `<hello id="b">booya</hello>`,
		`<x><y id="a">ok</y></x>`:   `<x><z>boo</z><y id="a">ok</y></x>`,
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
	b, a := parse(t, before), parse(t, after)
	result := html.Reconciler(nil, nil).Reconcile(b, a)
	if toHTML(result) != toHTML(a) {
		t.Error("Mismatched", toHTML(a), toHTML(result))
	}
}

func toHTML(n dom.MutableNode) string {
	if n == nil {
		return ""
	}

	return n.(stringer).String()
}

type stringer interface {
	String() string
}

func parse(t *testing.T, s string) dom.MutableNode {
	if s == "" {
		return nil
	}

	result, err := html.Parse(s)
	if err != nil {
		t.Fatal("invalid HTML", err)
	}
	return result
}
