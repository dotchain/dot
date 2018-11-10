// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom_test

import (
	"github.com/dotchain/dot/ui/dom"
	"golang.org/x/net/html"
	"strings"
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

	r := dom.Reconciler(newHTMLNode)

	for before, after := range tests {
		t.Run(before+"=>"+after, func(t *testing.T) {
			validate(t, r, before, after)
			validate(t, r, after, before)
		})
	}
}

func newHTMLNode(tag string, key interface{}) dom.MutableNode {
	n := &html.Node{Type: html.ElementNode, Data: tag}
	if tag == ":text:" {
		n.Type = html.TextNode
	}
	return node{n}
}

func validate(t *testing.T, r dom.Reconciler, before, after string) {
	b, a := parse(t, before), parse(t, after)
	result := r.Reconcile(b, a)
	if toHTML(result) != toHTML(a) {
		t.Error("Mismatched", toHTML(a), toHTML(result))
	}
}

func toHTML(n dom.MutableNode) string {
	if n == nil {
		return ""
	}

	var builder strings.Builder
	html.Render(&builder, n.(node).Node)
	return builder.String()
}

func parse(t *testing.T, s string) dom.MutableNode {
	if s == "" {
		return nil
	}

	nodes, err := html.ParseFragment(strings.NewReader(s), nil)
	if err != nil {
		t.Fatal("invalid HTML", err)
	}
	body := nodes[0].FirstChild.NextSibling
	return node{body.FirstChild}
}

type node struct {
	*html.Node
}

func (n node) Tag() string {
	if n.Node.Type == html.ElementNode {
		return n.Data
	}
	return ":text:"
}

func (n node) Key() interface{} {
	id := ""
	for _, attr := range n.Node.Attr {
		if attr.Key == "id" {
			id = attr.Val
		}
	}
	return id
}

func (n node) ForEachAttribute(fn func(key string, val interface{})) {
	if n.Tag() == ":text:" {
		fn(":data:", n.Node.Data)
	}

	for _, attr := range n.Node.Attr {
		fn(attr.Key, attr.Val)
	}
}

func (n node) ForEachNode(fn func(dom.Node)) {
	for nn := n.Node.FirstChild; nn != nil; nn = nn.NextSibling {
		fn(node{nn})
	}
}

func (n node) SetAttribute(key string, v interface{}) {
	val := v.(string)
	if key == ":data:" {
		n.Node.Data = val
	}

	for kk := range n.Node.Attr {
		if n.Node.Attr[kk].Key == key {
			n.Node.Attr[kk].Val = val
			return
		}
	}
	n.Node.Attr = append(n.Node.Attr, html.Attribute{"", key, val})
}

func (n node) RemoveAttribute(key string) {
	attr := n.Node.Attr
	for kk := range attr {
		if attr[kk].Key == key {
			n.Node.Attr = append(attr[:kk], attr[kk+1:]...)
			return
		}
	}
}

func (n node) Children() dom.MutableNodes {
	return &nodes{n.Node, n.Node.FirstChild}
}

type nodes struct {
	*html.Node
	child *html.Node
}

func (n *nodes) Next() dom.MutableNode {
	result := node{n.child}
	n.child = n.child.NextSibling
	return result
}

func (n *nodes) Remove() dom.MutableNode {
	removed := n.Next().(node)
	n.Node.RemoveChild(removed.Node)
	return removed
}

func (n *nodes) Insert(m dom.MutableNode) {
	child := m.(node).Node
	n.Node.InsertBefore(child, n.child)
}
