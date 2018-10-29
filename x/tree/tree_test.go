// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package tree_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/tree"
	"github.com/dotchain/dot/x/types"
	"golang.org/x/net/html"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestEdgeCases(t *testing.T) {
	n := &tree.Node{"Tag": "div", "X": types.S8("hello")}
	r := changes.Replace{types.S8("hello"), types.S8("boo")}
	result := n.Apply(changes.PathChange{[]interface{}{"X"}, r})
	expected := &tree.Node{"Tag": "div", "X": types.S8("boo")}

	if !reflect.DeepEqual(result, expected) {
		t.Error("setting attributes doesn't work", result)
	}

	if !reflect.DeepEqual(n.Apply(nil), n) {
		t.Error("Unexpected nil apply")
	}

	nn := tree.Nodes{&tree.Node{"Tag": "a"}, &tree.Node{"Tag": "b"}}
	if !reflect.DeepEqual(nn.Slice(0, 1), nn[:1]) {
		t.Error("nodes:slice failed", nn.Slice(0, 1))
	}
	if nn.Count() != 2 {
		t.Error("nodes:count failed", nn.Count())
	}
	if !reflect.DeepEqual(nn.Apply(nil), nn) {
		t.Error("nil apply")
	}

	expected2 := tree.Nodes{nn[1], nn[0]}
	move := changes.Move{1, 1, -1}
	if !reflect.DeepEqual(nn.Apply(move), expected2) {
		t.Error("Move failed", nn.Apply(move))
	}
}

func TestPanics(t *testing.T) {
	catch := func(message string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic", message)
			}
		}()
		fn()
	}

	catch("attributeApply", func() {
		mm := types.M{"Y": types.S8("hello")}
		(&tree.Node{"X": mm}).Apply(
			changes.PathChange{
				[]interface{}{"X", "Y"},
				changes.Replace{changes.Nil, types.S8("boo")},
			},
		)
	})
	catch("badChange", func() { (&tree.Node{}).Apply(badChange{}) })
	catch("badChange", func() { (tree.Nodes(nil)).Apply(badChange{}) })
}

func Test(t *testing.T) {
	// test = map[before]after
	tests := map[string]string{
		``:                          `<hello>booya</hello>`,
		`<x></x>`:                   `<x><y>ok</y></x>`,
		`<hello>boo</hello>`:        `<hello>booya</hello>`,
		`<hello x="a">boo</hello>`:  `<hello x="b" y="c">booya</hello>`,
		`<hello id="a">boo</hello>`: `<hello id="b">booya</hello>`,
		`<x><y id="a">ok</y></x>`:   `<x><z>boo</z><y id="a">ok</y></x>`,
		`<x><z id="b">boo</z><y id="a">ok</y></x>`: `<x><y id="a">ok</y><z id="b">boo</z></x>`,
	}

	for before, after := range tests {
		t.Run(before+"=>"+after, func(t *testing.T) {
			if x := toHTML(fromHTML(t, before)); !strings.EqualFold(before, x) {
				t.Fatal("html", before, x)
			}

			if x := toHTML(fromHTML(t, after)); !strings.EqualFold(after, x) {
				t.Fatal("html", after, x)
			}

			if x := tree.Diff(fromHTML(t, before), fromHTML(t, before)); x != nil {
				t.Fatal("self check before", x)
			}

			if x := tree.Diff(fromHTML(t, after), fromHTML(t, after)); x != nil {
				t.Fatal("self check after", x)
			}

			diff := tree.Diff(fromHTML(t, before), fromHTML(t, after))
			applied, _ := fromHTML(t, before).Apply(diff).(*tree.Node)
			if !strings.EqualFold(after, toHTML(applied)) {
				t.Error("Expected", after, "got", toHTML(applied))
			}

			diff = tree.Diff(fromHTML(t, after), fromHTML(t, before))
			applied, _ = fromHTML(t, after).Apply(diff).(*tree.Node)
			if !strings.EqualFold(before, toHTML(applied)) {
				t.Error("Expected", after, "got", toHTML(applied))
			}
		})
	}
}

func fromHTML(t *testing.T, s string) *tree.Node {
	nodes, err := html.ParseFragment(strings.NewReader(s), nil)
	if err != nil {
		t.Fatal("invalid HTML", err)
	}
	body := nodes[0].FirstChild.NextSibling
	return fromHTMLNode(body.FirstChild)
}

func fromHTMLNode(n *html.Node) *tree.Node {
	if n == nil {
		return nil
	}

	if n.Type != html.ElementNode {
		if n.Data == "" {
			return nil
		}

		return &tree.Node{"Text": n.Data}
	}

	children := tree.Nodes(nil)
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if nn := fromHTMLNode(child); nn != nil {
			children = append(children, nn)
		}
	}
	result := tree.Node{"Tag": n.Data}
	if len(children) > 0 {
		result["Children"] = children
	}

	for _, attr := range n.Attr {
		if attr.Key == "id" {
			result = result.WithKey(attr.Val)
		} else {
			result[attr.Key] = attr.Val
		}
	}

	return &result
}

func toHTML(np *tree.Node) string {
	if np == nil {
		return ""
	}

	n := *np

	if n["Tag"] == nil {
		return n["Text"].(string)
	}

	attributes := ""
	if v := n.Key(); v != nil {
		attributes = ` id="` + v.(string) + `"`
	}

	keys := []string{}
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if k == "Key" || k == "Children" || k == "Tag" {
			continue
		}
		attributes += " " + k + `="` + n[k].(string) + `"`
	}

	children := ""
	for _, child := range n.Children() {
		children += toHTML(child)
	}

	tag := n["Tag"].(string)
	return "<" + tag + attributes + ">" + children + "</" + tag + ">"
}

type badChange struct{}

func (b badChange) Merge(o changes.Change) (ox, bx changes.Change) {
	return nil, nil
}

func (b badChange) Revert() changes.Change {
	return nil
}
