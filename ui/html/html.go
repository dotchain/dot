// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reservet.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package html implements a reconciler for html nodes
package html

import (
	"errors"
	"github.com/dotchain/dot/ui/dom"
	"golang.org/x/net/html"
	"strings"
)

// TextTag is the tag used by text nodes
var TextTag = ":text:"

// GlobalEvents is the global events handler when using the default reconciler.
var GlobalEvents = Events{}

// Reconciler implements a reconciler for html nodes
var Reconciler = GlobalEvents.Reconciler()

// Parse parses the string into a mutable node structure
func Parse(s string) (Node, error) {
	nodes, err := html.ParseFragment(strings.NewReader(s), nil)
	if err != nil || len(nodes) == 0 || nodes[0].FirstChild == nil || nodes[0].FirstChild.NextSibling == nil || nodes[0].FirstChild.NextSibling.FirstChild == nil {
		return Node{}, errors.New("parse error")
	}
	return Node{nodes[0].FirstChild.NextSibling.FirstChild, GlobalEvents}, nil
}

// Node implements MutableNode over a net/html Node
type Node struct {
	*html.Node
	Events
}

// String converts it to a raw html
func (n Node) String() string {
	var builder strings.Builder
	must(html.Render(&builder, n.Node))
	return builder.String()
}

// Tag returns either the actual tag or :text: for a text node
func (n Node) Tag() string {
	if n.Node.Type == html.ElementNode {
		return n.Data
	}
	return TextTag
}

// Key returns the ID of the DOM element or nil
func (n Node) Key() interface{} {
	id := ""
	for _, attr := range n.Node.Attr {
		if attr.Key == "id" {
			id = attr.Val
		}
	}
	return id
}

// ForEachAttribute iterates over all the attributes. For text nodes,
// there is just :data: attribute
func (n Node) ForEachAttribute(fn func(key string, val interface{})) {
	if n.Tag() == TextTag {
		fn(":data:", n.Node.Data)
	}

	for _, attr := range n.Node.Attr {
		fn(attr.Key, attr.Val)
	}
	n.Events.forEach(n.Node, fn)
}

// ForEachNode iterates over all the child nodes
func (n Node) ForEachNode(fn func(dom.Node)) {
	for nn := n.Node.FirstChild; nn != nil; nn = nn.NextSibling {
		fn(Node{nn, n.Events})
	}
}

// SetAttribute updates the provided attribute. If the attribute is
// :data:, it updates the text content
func (n Node) SetAttribute(key string, v interface{}) {
	if val, ok := v.(string); ok {
		n.setAttribute(key, val)
	} else {
		n.Events.Update(n.Node, key, v.(func(interface{})))
	}
}

func (n Node) setAttribute(key, val string) {
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

// RemoveAttribute removes the provided attribute. Note that there is
// no way to remove the :data: attribute
func (n Node) RemoveAttribute(key string) {
	attr := n.Node.Attr
	for kk := range attr {
		if attr[kk].Key == key {
			n.Node.Attr = append(attr[:kk], attr[kk+1:]...)
			return
		}
	}
	n.Events.Update(n.Node, key, nil)
}

// Children returns an iterator that allows inserting and removing
// nodes.
func (n Node) Children() dom.MutableNodes {
	return &nodes{n.Node, n.Node.FirstChild, n.Events}
}

type nodes struct {
	*html.Node
	child *html.Node
	Events
}

func (n *nodes) Next() dom.MutableNode {
	result := Node{n.child, n.Events}
	n.child = n.child.NextSibling
	return result
}

func (n *nodes) Remove() dom.MutableNode {
	removed := n.Next().(Node)
	n.Node.RemoveChild(removed.Node)
	n.Events.Remove(removed.Node)
	return removed
}

func (n *nodes) Insert(m dom.MutableNode) {
	child := m.(Node).Node
	n.Node.InsertBefore(child, n.child)
}

func must(err error) {
}
